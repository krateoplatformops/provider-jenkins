package pipeline

import (
	"context"
	"errors"
	"fmt"
	"strings"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/google/go-cmp/cmp"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/record"

	"github.com/crossplane/crossplane-runtime/pkg/controller"
	"github.com/crossplane/crossplane-runtime/pkg/meta"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/crossplane/crossplane-runtime/pkg/event"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"github.com/crossplane/crossplane-runtime/pkg/ratelimiter"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"

	v1alpha1 "github.com/krateoplatformops/provider-jenkins/apis/pipeline/v1alpha1"

	"github.com/krateoplatformops/provider-jenkins/internal/clients/jenkins"
	"github.com/krateoplatformops/provider-jenkins/internal/helpers"
)

const (
	errInvalidCRD = "managed resource is not an Pipeline custom resource"

	reasonCannotCreate = "CannotCreateExternalResource"
	reasonCreated      = "CreatedExternalResource"
	reasonDeleted      = "DeletedExternalResource"
)

func Setup(mgr ctrl.Manager, o controller.Options) error {
	name := managed.ControllerName(v1alpha1.PipelineGroupKind)

	log := o.Logger.WithValues("controller", name)

	recorder := mgr.GetEventRecorderFor(name)

	r := managed.NewReconciler(mgr,
		resource.ManagedKind(v1alpha1.PipelineGroupVersionKind),
		managed.WithExternalConnecter(&connector{
			kube:     mgr.GetClient(),
			log:      log,
			recorder: recorder,
			clientFn: jenkins.NewClient,
		}),
		managed.WithLogger(log),
		managed.WithRecorder(event.NewAPIRecorder(recorder)))

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		WithOptions(o.ForControllerRuntime()).
		For(&v1alpha1.Pipeline{}).
		Complete(ratelimiter.NewReconciler(name, r, o.GlobalRateLimiter))
}

type connector struct {
	kube     client.Client
	log      logging.Logger
	recorder record.EventRecorder
	clientFn func(opts *jenkins.ClientOpts) *jenkins.Client
}

func (c *connector) Connect(ctx context.Context, mg resource.Managed) (managed.ExternalClient, error) {
	cr, ok := mg.(*v1alpha1.Pipeline)
	if !ok {
		return nil, errors.New(errInvalidCRD)
	}

	opts, err := jenkinsClientFromProviderConfig(ctx, c.kube, cr)
	if err != nil {
		c.log.Info(fmt.Sprintf("%s: initializing Jenkins API client", err.Error()))
		return nil, err
	}

	return &external{
		kube: c.kube,
		log:  c.log,
		rec:  c.recorder,
		cli:  c.clientFn(opts),
	}, nil
}

type external struct {
	kube client.Client
	log  logging.Logger
	rec  record.EventRecorder
	cli  *jenkins.Client
}

func (e *external) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	cr, ok := mg.(*v1alpha1.Pipeline)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errInvalidCRD)
	}

	parts := strings.Split(meta.GetExternalName(cr), "/")
	if len(parts) < 2 {
		return managed.ExternalObservation{
			ResourceExists:   false,
			ResourceUpToDate: true,
		}, nil
	}

	observed, err := e.cli.GetJobConfig(ctx, parts[1])
	if err != nil {
		var notFound *jenkins.HTTPStatusError
		if errors.As(err, &notFound) {
			return managed.ExternalObservation{
				ResourceExists:   false,
				ResourceUpToDate: true,
			}, nil
		}
		return managed.ExternalObservation{}, err
	}
	e.log.Debug("Observed job configuration fetched", "name", meta.GetExternalName(cr))

	spec := cr.Spec.ForProvider.DeepCopy()

	desired, err := e.getJobConfig(ctx, spec)
	if err != nil {
		return managed.ExternalObservation{}, err
	}

	if !cmp.Equal(desired, string(observed)) {
		e.log.Debug("Configuration drift detected!")
		return managed.ExternalObservation{
			ResourceExists:   true,
			ResourceUpToDate: false,
		}, nil
	}

	cr.Status.AtProvider = generateObservation(spec)
	cr.Status.SetConditions(xpv1.Available())

	return managed.ExternalObservation{
		ResourceExists:   true,
		ResourceUpToDate: true,
	}, nil
}

func (e *external) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*v1alpha1.Pipeline)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errInvalidCRD)
	}

	cr.SetConditions(xpv1.Creating())

	spec := cr.Spec.ForProvider.DeepCopy()

	data, err := e.getJobConfig(ctx, spec)
	if err != nil {
		return managed.ExternalCreation{}, err
	}

	err = e.cli.CreateJob(ctx, spec.JobName, []byte(data))
	if err != nil {
		return managed.ExternalCreation{}, err
	}
	e.log.Debug("Job created", "name", spec.JobName)
	e.rec.Eventf(cr, corev1.EventTypeNormal, reasonCreated, "Job created (name: %s)", spec.JobName)

	meta.SetExternalName(cr, fmt.Sprintf("jenkins/%s", spec.JobName))

	return managed.ExternalCreation{}, nil
}

func (e *external) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	e.log.Debug("Update requested but not implemented (yet?)")
	return managed.ExternalUpdate{}, nil
}

func (e *external) Delete(ctx context.Context, mg resource.Managed) error {
	cr, ok := mg.(*v1alpha1.Pipeline)
	if !ok {
		return errors.New(errInvalidCRD)
	}

	cr.SetConditions(xpv1.Deleting())

	spec := cr.Spec.ForProvider.DeepCopy()

	return e.cli.DeleteJob(ctx, spec.JobName)
}

func (e *external) getJobConfig(ctx context.Context, spec *v1alpha1.PipelineParams) (string, error) {
	ref := helpers.ConfigMapKeySelector{
		Name:      spec.JobConfigRef.Name,
		Namespace: spec.JobConfigRef.Namespace,
		Key:       spec.JobConfigRef.Key,
	}
	res, err := helpers.GetConfigMapValue(ctx, e.kube, ref)
	if err == nil {
		e.log.Debug("Desired configuration retrieved", "configMapRef", ref.String())
	}

	return res, err
}

func generateObservation(e *v1alpha1.PipelineParams) v1alpha1.PipelineObservation {
	return v1alpha1.PipelineObservation{
		JobName: helpers.StringPtr(e.JobName),
	}
}
