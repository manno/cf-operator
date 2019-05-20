package reference

import (
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"

	bdv1 "code.cloudfoundry.org/cf-operator/pkg/kube/apis/boshdeployment/v1alpha1"
	ejobv1 "code.cloudfoundry.org/cf-operator/pkg/kube/apis/extendedjob/v1alpha1"
	estsv1 "code.cloudfoundry.org/cf-operator/pkg/kube/apis/extendedstatefulset/v1alpha1"
)

// GetConfigMapsReferencedBy returns a list of all names for ConfigMaps referenced by the object
// The object can be an ExtendedStatefulSet, an ExtendedeJob or a BOSHDeployment
func GetConfigMapsReferencedBy(object interface{}) (map[string]bool, error) {
	// Figure out the type of object
	switch object.(type) {
	case bdv1.BOSHDeployment:
		return getConfMapRefFromBdpl(object.(bdv1.BOSHDeployment)), nil
	case ejobv1.ExtendedJob:
		return getConfMapRefFromEJob(object.(ejobv1.ExtendedJob)), nil
	case estsv1.ExtendedStatefulSet:
		return getConfMapRefFromESts(object.(estsv1.ExtendedStatefulSet)), nil
	default:
		return nil, errors.New("can't get config map references for unkown type; supported types are BOSHDeployment, ExtendedJob and ExtendedStatefulSet")
	}
}

func getConfMapRefFromBdpl(object bdv1.BOSHDeployment) map[string]bool {
	result := map[string]bool{}

	if object.Spec.Manifest.Type == bdv1.ConfigMapType {
		result[object.Spec.Manifest.Ref] = true
	}

	for _, ops := range object.Spec.Ops {
		if ops.Type == bdv1.ConfigMapType {
			result[ops.Ref] = true
		}
	}

	return result
}

func getConfMapRefFromESts(object estsv1.ExtendedStatefulSet) map[string]bool {
	return getConfMapRefFromPod(object.Spec.Template.Spec.Template.Spec)
}

func getConfMapRefFromEJob(object ejobv1.ExtendedJob) map[string]bool {
	return getConfMapRefFromPod(object.Spec.Template.Spec)
}

func getConfMapRefFromPod(object corev1.PodSpec) map[string]bool {
	result := map[string]bool{}

	// Look at all volumes
	for _, volume := range object.Volumes {
		if volume.VolumeSource.ConfigMap != nil {
			result[volume.VolumeSource.ConfigMap.Name] = true
		}
	}

	// Look at all init containers
	for _, container := range object.Containers {
		for _, envFrom := range container.EnvFrom {
			if envFrom.ConfigMapRef != nil {
				result[envFrom.ConfigMapRef.Name] = true
			}
		}

		for _, envVar := range container.Env {
			if envVar.ValueFrom != nil && envVar.ValueFrom.ConfigMapKeyRef != nil {
				result[envVar.ValueFrom.ConfigMapKeyRef.Name] = true
			}
		}
	}

	// Look at all containers
	for _, container := range object.Containers {
		for _, envFrom := range container.EnvFrom {
			if envFrom.ConfigMapRef != nil {
				result[envFrom.ConfigMapRef.Name] = true
			}
		}

		for _, envVar := range container.Env {
			if envVar.ValueFrom != nil && envVar.ValueFrom.ConfigMapKeyRef != nil {
				result[envVar.ValueFrom.ConfigMapKeyRef.Name] = true
			}
		}
	}

	return result
}
