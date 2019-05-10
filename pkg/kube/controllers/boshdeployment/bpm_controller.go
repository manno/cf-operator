package boshdeployment

import (
	"code.cloudfoundry.org/cf-operator/pkg/kube/apis"
	"code.cloudfoundry.org/cf-operator/pkg/kube/util/names"
	"code.cloudfoundry.org/cf-operator/pkg/kube/util/owner"
	"code.cloudfoundry.org/cf-operator/pkg/kube/util/versionedsecretstore"
	"context"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/source"

	bdm "code.cloudfoundry.org/cf-operator/pkg/bosh/manifest"
	"code.cloudfoundry.org/cf-operator/pkg/kube/util/config"
	"code.cloudfoundry.org/cf-operator/pkg/kube/util/ctxlog"
)

// AddBPM creates a new BPM Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func AddBPM(ctx context.Context, config *config.Config, mgr manager.Manager) error {
	ctx = ctxlog.NewContextWithRecorder(ctx, "bpm-reconciler", mgr.GetRecorder("bpm-recorder"))
	r := NewBPMReconciler(ctx, config, mgr, bdm.NewResolver(mgr.GetClient(), func() bdm.Interpolator { return bdm.NewInterpolator() }), controllerutil.SetControllerReference)

	// Create a new controller
	c, err := controller.New("bpm-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	mapSecrets := handler.ToRequestsFunc(func(a handler.MapObject) []reconcile.Request {
		secret := a.Object.(*corev1.Secret)
		return reconcilesForVersionedSecret(ctx, mgr, *secret)
	})

	// Watch Secrets owned by resource ExtendedStatefulSet or referenced by resource ExtendedStatefulSet
	err = c.Watch(&source.Kind{Type: &corev1.Secret{}}, &handler.EnqueueRequestsFromMapFunc{ToRequests: mapSecrets})
	if err != nil {
		return err
	}

	return nil
}

// reconcilesForVersionedSecret
func reconcilesForVersionedSecret(ctx context.Context, mgr manager.Manager, secret corev1.Secret) []reconcile.Request {
	reconciles := []reconcile.Request{}

	// add requests for the ExtendedStatefulSet referencing the versioned secret
	secretLabels := secret.GetLabels()
	if secretLabels == nil {
		return reconciles
	}

	secretKind, ok := secretLabels[versionedsecretstore.LabelSecretKind]
	if !ok {
		return reconciles
	}
	if secretKind != versionedsecretstore.VersionSecretKind {
		return reconciles
	}

	referencedSecretName := names.GetPrefixFromVersionedSecretName(secret.GetName())
	if referencedSecretName == "" {
		return reconciles
	}

	exStatefulSets := &essv1.ExtendedStatefulSetList{}
	err := mgr.GetClient().List(ctx, &client.ListOptions{}, exStatefulSets)
	if err != nil || len(exStatefulSets.Items) < 1 {
		return reconciles
	}

	for _, exStatefulSet := range exStatefulSets.Items {
		_, referencedSecrets := owner.GetConfigNamesFromSpec(exStatefulSet.Spec.Template.Spec.Template.Spec)
		if _, ok := referencedSecrets[referencedSecretName]; ok {
			reconciles = append(reconciles, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      exStatefulSet.GetName(),
					Namespace: exStatefulSet.GetNamespace(),
				},
			})
		}

		// add requests for the ExtendedStatefulSet referencing the versioned secret when a new ExtendedStatefulSet template updated
		for secretName := range referencedSecrets {
			if strings.HasPrefix(secretName, referencedSecretName) {
				reconciles = append(reconciles, reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      exStatefulSet.GetName(),
						Namespace: exStatefulSet.GetNamespace(),
					},
				})
			}
		}
	}

	return reconciles
}
