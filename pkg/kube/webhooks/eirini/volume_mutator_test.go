package eirini_test

import (
	"context"
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"

	admissionv1beta1 "k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
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
		request          types.Request
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

		request = types.Request{AdmissionRequest: &admissionv1beta1.AdmissionRequest{}}
	})

	// source_type: APP
	Describe("Handle", func() {

		It("passes on errors from the decoding step", func() {
			f := generateGetPodFunc(nil, fmt.Errorf("decode failed"))
			mutator := webhooks.NewVolumeMutator(log, config, manager, setReferenceFunc, f)

			res := mutator.Handle(ctx, request)
			Expect(res.Response.Result.Code).To(Equal(int32(http.StatusBadRequest)))
		})

		It("does not act if the source_type: APP label is not set", func() {
			pod := labeledPod("foo", map[string]string{})
			f := generateGetPodFunc(&pod, nil)

			mutator := webhooks.NewVolumeMutator(log, config, manager, setReferenceFunc, f)

			resp := mutator.Handle(ctx, request)
			Expect(resp).ToNot(BeNil())
		})
	})
})

// LabeledPod defines a pod with labels and a simple web server
func labeledPod(name string, labels map[string]string) corev1.Pod {
	return corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: labels,
		},
	}
}

func generateGetPodFunc(pod *corev1.Pod, err error) webhooks.GetPodFuncType {
	return func(_ types.Decoder, _ types.Request) (*corev1.Pod, error) {
		return pod, err
	}
}
