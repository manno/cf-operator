package boshdeployment_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"code.cloudfoundry.org/cf-operator/pkg/bosh/manifest/fakes"
	bdv1 "code.cloudfoundry.org/cf-operator/pkg/kube/apis/boshdeployment/v1alpha1"
	esv1 "code.cloudfoundry.org/cf-operator/pkg/kube/apis/extendedsecret/v1alpha1"
	"code.cloudfoundry.org/cf-operator/pkg/kube/controllers"
	cfd "code.cloudfoundry.org/cf-operator/pkg/kube/controllers/boshdeployment"
	cfakes "code.cloudfoundry.org/cf-operator/pkg/kube/controllers/fakes"
	cfcfg "code.cloudfoundry.org/cf-operator/pkg/kube/util/config"
	"code.cloudfoundry.org/cf-operator/pkg/kube/util/ctxlog"
	helper "code.cloudfoundry.org/cf-operator/pkg/testhelper"
)

var _ = Describe("ReconcileGeneratedVariable", func() {
	var (
		manager               *cfakes.FakeManager
		reconciler            reconcile.Reconciler
		recorder              *record.FakeRecorder
		request               reconcile.Request
		ctx                   context.Context
		resolver              fakes.FakeResolver
		log                   *zap.SugaredLogger
		config                *cfcfg.Config
		client                *cfakes.FakeClient
		manifestWithOpsSecret *corev1.Secret
	)

	BeforeEach(func() {
		controllers.AddToScheme(scheme.Scheme)
		recorder = record.NewFakeRecorder(20)
		manager = &cfakes.FakeManager{}
		manager.GetSchemeReturns(scheme.Scheme)
		manager.GetRecorderReturns(recorder)
		resolver = fakes.FakeResolver{}

		request = reconcile.Request{NamespacedName: types.NamespacedName{Name: "foo-with-ops", Namespace: "default"}}

		manifestWithOpsSecret = &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo-with-ops",
				Namespace: "default",
			},
			StringData: map[string]string{
				"manifest.yaml": `---
name: fake-manifest
releases:
- name: bar
  url: docker.io/cfcontainerization
  version: 1.0
  stemcell:
    os: opensuse
    version: 42.3
instance_groups:
- name: fakepod
  jobs:
  - name: foo
    release: bar
    properties:
      password: ((foo_password))
      bosh_containerization:
        ports:
        - name: foo
          protocol: TCP
          internal: 8080
variables:
- name: foo_password
  type: password
`,
			},
		}
		config = &cfcfg.Config{CtxTimeOut: 10 * time.Second}
		_, log = helper.NewTestLogger()
		ctx = ctxlog.NewParentContext(log)
		ctx = ctxlog.NewContextWithRecorder(ctx, "TestRecorder", recorder)

		client = &cfakes.FakeClient{}
		client.GetCalls(func(context context.Context, nn types.NamespacedName, object runtime.Object) error {
			switch object.(type) {
			case *corev1.Secret:
				if nn.Name == "foo-with-ops" {
					manifestWithOpsSecret.DeepCopyInto(object.(*corev1.Secret))
				}
			}

			return nil
		})

		manager.GetClientReturns(client)
	})

	JustBeforeEach(func() {
		reconciler = cfd.NewGeneratedVariableReconciler(ctx, config, manager, &resolver, controllerutil.SetControllerReference)
	})

	Describe("Reconcile", func() {
		Context("when manifest with ops is created", func() {
			It("handles an error when generating variables", func() {
				client.GetCalls(func(context context.Context, nn types.NamespacedName, object runtime.Object) error {
					switch object.(type) {
					case *corev1.Secret:
						if nn.Name == "foo-with-ops" {
							manifestWithOpsSecret.DeepCopyInto(object.(*corev1.Secret))
						}
					case *esv1.ExtendedSecret:
						return apierrors.NewNotFound(schema.GroupResource{}, nn.Name)
					}

					return nil
				})
				client.CreateCalls(func(context context.Context, object runtime.Object) error {
					switch object.(type) {
					case *esv1.ExtendedSecret:
						return errors.New("fake-error")
					}
					return nil
				})

				By("From ops applied state to variable interpolated state")
				_, err := reconciler.Reconcile(request)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed to generate variables"))

			})

			It("creates generated variables", func() {
				result, err := reconciler.Reconcile(request)
				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(Equal(reconcile.Result{}))

				newInstance := &bdv1.BOSHDeployment{}
				err = client.Get(context.Background(), types.NamespacedName{Name: "foo", Namespace: "default"}, newInstance)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})
})
