package jenkins

import (
	"context"
	"fmt"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/krateoplatformops/provider-jenkins/apis/v1alpha1"
	"github.com/krateoplatformops/provider-jenkins/internal/helpers"
	httphelper "github.com/krateoplatformops/provider-jenkins/internal/helpers/http"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ClientFromProviderConfig(ctx context.Context, kc client.Client, mg resource.Managed) (*ClientOpts, error) {
	if mg.GetProviderConfigReference() == nil {
		return nil, errors.New("providerConfigRef is not given")
	}

	pc := &v1alpha1.ProviderConfig{}
	err := kc.Get(ctx, types.NamespacedName{Name: mg.GetProviderConfigReference().Name}, pc)
	if err != nil {
		return nil, errors.Wrap(err, "cannot get referenced Provider")
	}

	t := resource.NewProviderConfigUsageTracker(kc, &v1alpha1.ProviderConfigUsage{})
	err = t.Track(ctx, mg)
	if err != nil {
		return nil, errors.Wrap(err, "cannot track ProviderConfig usage")
	}

	return initJenkinsClientOpts(ctx, kc, pc)
}

func initJenkinsClientOpts(ctx context.Context, kc client.Client, pc *v1alpha1.ProviderConfig) (*ClientOpts, error) {
	opts := httphelper.ClientOpts{
		Verbose:  helpers.BoolValue(pc.Spec.Verbose),
		Insecure: helpers.BoolValue(pc.Spec.Insecure),
	}

	if s := pc.Spec.Credentials.Source; s != xpv1.CredentialsSourceSecret {
		return nil, fmt.Errorf("credentials source %s is not currently supported", s)
	}

	csr := pc.Spec.Credentials.SecretRef
	if csr == nil {
		return nil, fmt.Errorf("no credentials secret referenced")
	}

	var auth *Auth
	apiToken, err := helpers.GetSecretValue(ctx, kc, helpers.SecretKeySelector{
		Name: csr.Name, Namespace: csr.Namespace, Key: csr.Key,
	})
	if err != nil {
		return nil, err
	}

	username := helpers.StringValue(pc.Spec.Username)
	if len(username) > 0 && len(apiToken) > 0 {
		auth = &Auth{
			Username: username, ApiToken: apiToken,
		}
	}

	return &ClientOpts{
		BaseUrl:    pc.Spec.BaseUrl,
		Controller: helpers.StringValue(pc.Spec.Controller),
		Auth:       auth,
		HttpClient: httphelper.ClientFromOpts(opts),
	}, nil
}
