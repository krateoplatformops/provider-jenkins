package folderrolebinding

import (
	"context"
	"errors"
	"fmt"
	"strings"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
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

	v1alpha1 "github.com/krateoplatformops/provider-jenkins/apis/folderrolebinding/v1alpha1"

	"github.com/krateoplatformops/provider-jenkins/internal/clients/jenkins"
	"github.com/krateoplatformops/provider-jenkins/internal/helpers"
)

const (
	errInvalidCRD = "managed resource is not an FolderRoleBinding custom resource"

	externalNameFmt = "folderrolebinding/%s/%s"

	reasonCannotCreate = "CannotCreateExternalResource"
	reasonCreated      = "CreatedExternalResource"
	reasonDeleted      = "DeletedExternalResource"
	reasonCannotDelete = "CannotDeleteExternalResource"
)

func Setup(mgr ctrl.Manager, o controller.Options) error {
	name := managed.ControllerName(v1alpha1.FolderRoleBindingGroupKind)

	log := o.Logger.WithValues("controller", name)

	recorder := mgr.GetEventRecorderFor(name)

	r := managed.NewReconciler(mgr,
		resource.ManagedKind(v1alpha1.FolderRoleBindingGroupVersionKind),
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
		For(&v1alpha1.FolderRoleBinding{}).
		Complete(ratelimiter.NewReconciler(name, r, o.GlobalRateLimiter))
}

type connector struct {
	kube     client.Client
	log      logging.Logger
	recorder record.EventRecorder
	clientFn func(opts *jenkins.ClientOpts) *jenkins.Client
}

func (c *connector) Connect(ctx context.Context, mg resource.Managed) (managed.ExternalClient, error) {
	cr, ok := mg.(*v1alpha1.FolderRoleBinding)
	if !ok {
		return nil, errors.New(errInvalidCRD)
	}

	opts, err := jenkins.ClientFromProviderConfig(ctx, c.kube, cr)
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
	cr, ok := mg.(*v1alpha1.FolderRoleBinding)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errInvalidCRD)
	}

	parts := strings.Split(meta.GetExternalName(cr), "/")
	if len(parts) < 3 {
		return managed.ExternalObservation{
			ResourceExists:   false,
			ResourceUpToDate: true,
		}, nil
	}

	// TODO: we don't have any API ref to check if resource exists
	/*
		_, err := e.cli.GetJobConfig(ctx, parts[1])
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
	*/
	spec := cr.Spec.ForProvider.DeepCopy()

	cr.Status.AtProvider = generateObservation(spec)
	cr.Status.SetConditions(xpv1.Available())

	return managed.ExternalObservation{
		ResourceExists:   true,
		ResourceUpToDate: true,
	}, nil
}

func (e *external) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*v1alpha1.FolderRoleBinding)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errInvalidCRD)
	}

	cr.SetConditions(xpv1.Creating())

	spec := cr.Spec.ForProvider.DeepCopy()

	sid := strings.TrimSpace(spec.SID)
	if len(sid) == 0 {
		err := fmt.Errorf("SID not specified in FolderRole (name: %s)", spec.Name)
		e.log.Debug("SID not specified in FolderRole", "name", spec.Name)
		e.rec.Eventf(cr, corev1.EventTypeNormal, "MissingSID", err.Error())
		return managed.ExternalCreation{}, err
	}

	err := e.cli.AddFolderRole(ctx, jenkins.AddFolderRoleOpts{
		Name:        spec.Name,
		Permissions: spec.Permissions,
		FolderNames: spec.FolderNames,
	})
	if err != nil {
		e.log.Info("FolderRole NOT created", "name", spec.Name, "error", err.Error())
		e.rec.Eventf(cr, corev1.EventTypeWarning, reasonCannotCreate, "FolderRole NOT created (name: %s, error: %s)", spec.Name, err.Error())
		return managed.ExternalCreation{}, err
	}

	err = e.cli.AssignSidToFolderRole(ctx, spec.SID, spec.Name)
	if err != nil {
		return managed.ExternalCreation{}, err
	}

	meta.SetExternalName(cr, fmt.Sprintf(externalNameFmt, spec.Name, spec.SID))
	return managed.ExternalCreation{}, nil
}

func (e *external) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	e.log.Debug("Update requested but not implemented (yet?)")
	return managed.ExternalUpdate{}, nil
}

func (e *external) Delete(ctx context.Context, mg resource.Managed) error {
	cr, ok := mg.(*v1alpha1.FolderRoleBinding)
	if !ok {
		return errors.New(errInvalidCRD)
	}

	cr.SetConditions(xpv1.Deleting())

	spec := cr.Spec.ForProvider.DeepCopy()

	err := e.cli.DeleteFolderRole(ctx, spec.Name)
	if err != nil {
		e.log.Info("Error deleting FolderRole", "name", spec.Name, "error", err.Error())
		//e.rec.Eventf(cr, corev1.EventTypeWarning, reasonCannotDelete, "Error deleting FolderRole (name: %s, error: %s)", spec.Name, err.Error())
	}

	return nil
}

func generateObservation(e *v1alpha1.FolderRoleBindingParams) v1alpha1.FolderRoleBindingObservation {
	return v1alpha1.FolderRoleBindingObservation{
		Name: helpers.StringPtr(e.Name),
		SID:  helpers.StringPtr(e.SID),
	}
}
