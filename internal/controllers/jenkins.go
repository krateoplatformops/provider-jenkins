package controllers

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/crossplane-runtime/pkg/controller"

	"github.com/krateoplatformops/provider-jenkins/internal/controllers/config"
	"github.com/krateoplatformops/provider-jenkins/internal/controllers/folderrolebinding"
	"github.com/krateoplatformops/provider-jenkins/internal/controllers/pipeline"
)

// Setup creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		config.Setup,
		pipeline.Setup,
		folderrolebinding.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
