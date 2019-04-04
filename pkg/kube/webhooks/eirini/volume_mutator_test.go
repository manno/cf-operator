package eirini_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"

	admissionv1beta1 "k8s.io/api/admission/v1beta1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission/types"

	"code.cloudfoundry.org/cf-operator/pkg/kube/client/clientset/versioned/scheme"
	"code.cloudfoundry.org/cf-operator/pkg/kube/controllers"
	cfakes "code.cloudfoundry.org/cf-operator/pkg/kube/controllers/fakes"
	"code.cloudfoundry.org/cf-operator/pkg/kube/util/config"
	webhooks "code.cloudfoundry.org/cf-operator/pkg/kube/webhooks/eirini"
	helper "code.cloudfoundry.org/cf-operator/pkg/testhelper"
	"code.cloudfoundry.org/cf-operator/testing"
)

var _ = Describe("Volume Mutator", func() {

	var (
		manager          *cfakes.FakeManager
		client           *cfakes.FakeClient
		ctx              context.Context
		config           *config.Config
		env              testing.Catalog
		log              *zap.SugaredLogger
		setReferenceFunc func(owner, object metav1.Object, scheme *runtime.Scheme) error = func(owner, object metav1.Object, scheme *runtime.Scheme) error { return nil }
	)

	BeforeEach(func() {
		controllers.AddToScheme(scheme.Scheme)
		client = &cfakes.FakeClient{}
		restMapper := meta.NewDefaultRESTMapper([]schema.GroupVersion{})
		restMapper.Add(schema.GroupVersionKind{Group: "", Kind: "Pod", Version: "v1"}, meta.RESTScopeNamespace)

		manager = &cfakes.FakeManager{}
		manager.GetSchemeReturns(scheme.Scheme)
		manager.GetClientReturns(client)
		manager.GetRESTMapperReturns(restMapper)

		config = env.DefaultConfig()
		ctx = testing.NewContext()
		_, log = helper.NewTestLogger()
	})

	Describe("Handle", func() {
		It("does something", func() {
			mutator := webhooks.NewVolumeMutator(log, config, manager, setReferenceFunc)
			adm_request := &admissionv1beta1.AdmissionRequest{}
			request := types.Request{AdmissionRequest: adm_request}

			resp := mutator.Handle(ctx, request)
			Expect(resp).ToNot(BeNil())
		})
	})
})
