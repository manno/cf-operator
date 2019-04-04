package eirini

import (
	"context"
	"fmt"
	"net/http"

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
	getPodFunc   GetPodFuncType
}

// Implement admission.Handler so the controller can handle admission request.
var _ admission.Handler = &VolumeMutator{}

// NewPodMutator returns a new reconcile.Reconciler
func NewVolumeMutator(log *zap.SugaredLogger, config *config.Config, mgr manager.Manager, srf setReferenceFunc, getPodFunc GetPodFuncType) admission.Handler {
	mutatorLog := log.Named("eirini-volume-mutator")
	mutatorLog.Info("Creating a Volume mutator")

	return &VolumeMutator{
		log:          mutatorLog,
		config:       config,
		client:       mgr.GetClient(),
		scheme:       mgr.GetScheme(),
		setReference: srf,
		getPodFunc:   getPodFunc,
	}
}

// Handle manages volume claims for ExtendedStatefulSet pods
func (m *VolumeMutator) Handle(ctx context.Context, req types.Request) types.Response {

	pod, err := m.getPodFunc(m.decoder, req)
	if err != nil {
		return admission.ErrorResponse(http.StatusBadRequest, err)
	}
	fmt.Println(pod)

	return types.Response{}
}
