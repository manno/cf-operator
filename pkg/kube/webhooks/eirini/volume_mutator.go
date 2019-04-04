package eirini

import (
	"context"

	"code.cloudfoundry.org/cf-operator/pkg/kube/util/config"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission/types"
)

type setReferenceFunc func(owner, object metav1.Object, scheme *runtime.Scheme) error

// VolumeMutator changes pod definitions
type VolumeMutator struct {
	client       client.Client
	scheme       *runtime.Scheme
	setReference setReferenceFunc
	log          *zap.SugaredLogger
	config       *config.Config
	decoder      types.Decoder
}

// Implement admission.Handler so the controller can handle admission request.
var _ admission.Handler = &VolumeMutator{}

// NewPodMutator returns a new reconcile.Reconciler
func NewVolumeMutator(log *zap.SugaredLogger, config *config.Config, mgr manager.Manager, srf setReferenceFunc) admission.Handler {
	mutatorLog := log.Named("eirini-volume-mutator")
	mutatorLog.Info("Creating a Volume mutator")

	return &VolumeMutator{
		log:          mutatorLog,
		config:       config,
		client:       mgr.GetClient(),
		scheme:       mgr.GetScheme(),
		setReference: srf,
	}
}

// Handle manages volume claims for ExtendedStatefulSet pods
func (m *VolumeMutator) Handle(ctx context.Context, req types.Request) types.Response {

	return types.Response{}
}
