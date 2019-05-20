package v1alpha1

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"

	"code.cloudfoundry.org/cf-operator/pkg/kube/apis"
)

// This file is safe to edit
// It's used as input for the Kube code generator
// Run "make generate" after modifying this file

var (
	// LabelReferencedJobName is the name key for dependent job
	LabelReferencedJobName = fmt.Sprintf("%s/referenced-job-name", apis.GroupName)
)

const (
	// LabelExtendedJob key for label used to identify extendedjob. Value
	// is set to true if the batchv1.Job is from an ExtendedJob
	LabelExtendedJob = "extendedjob"
	// LabelEJobName key for label on a batchv1.Job's pod, which is set to the ExtendedJob's name
	LabelEJobName = "ejob-name"
	// LabelTriggeringPod key for label, which is set to the UID of the pod that triggered an ExtendedJob
	LabelTriggeringPod = "triggering-pod"
)

// ExtendedJobSpec defines the desired state of ExtendedJob
type ExtendedJobSpec struct {
	Output               *Output                `json:"output,omitempty"`
	Trigger              Trigger                `json:"trigger"`
	Template             corev1.PodTemplateSpec `json:"template"`
	UpdateOnConfigChange bool                   `json:"updateOnConfigChange"`
}

// Strategy describes the trigger strategy
type Strategy string

const (
	// TriggerManual is the default for errand jobs, change to TriggerNow to run them
	TriggerManual Strategy = "manual"
	// TriggerNow instructs the controller to run the job now,
	// resets to TriggerManual after starting the job
	TriggerNow Strategy = "now"
	// TriggerOnce jobs run only once, when created, then switches to TriggerDone
	TriggerOnce Strategy = "once"
	// TriggerDone jobs are no longer triggered. It's the final state for TriggerOnce strategies
	TriggerDone Strategy = "done"
)

// Trigger decides how to trigger the ExtendedJob
type Trigger struct {
	Strategy Strategy         `json:"strategy"`
	PodState *PodStateTrigger `json:"podstate,omitempty"`
}

// PodState is our abstraction of the pods state with regards to triggered
// extended jobs
type PodState string

const (
	// PodStateUnknown means we could not identify the state
	PodStateUnknown PodState = ""

	// PodStateReady means the pod is in phase=running with condition=ready
	PodStateReady PodState = "ready"

	// PodStateCreated means the pod did not exist before and is ready
	PodStateCreated PodState = "created"

	// PodStateNotReady means the pod is in phase pending
	PodStateNotReady PodState = "notready"

	// PodStateDeleted means the pod is in phase=succeeded or disappeared or phase=''
	PodStateDeleted PodState = "deleted"
)

// PodStateTrigger specifies how to trigger depending on a Job
type PodStateTrigger struct {
	When     PodState  `json:"when"`
	Selector *Selector `json:"selector,omitempty"`
}

// Selector filter objects
type Selector struct {
	MatchLabels      *labels.Set    `json:"matchLabels,omitempty"`
	MatchExpressions []*Requirement `json:"matchExpressions,omitempty"`
}

// Requirement describes a label requirement
type Requirement struct {
	Key      string             `json:"key"`
	Operator selection.Operator `json:"operator"`
	Values   []string           `json:"values"`
}

// Output contains options to persist job output
type Output struct {
	NamePrefix     string            `json:"namePrefix"`           // the secret name will be <NamePrefix><container name>
	OutputType     string            `json:"outputType,omitempty"` // only json is supported for now
	SecretLabels   map[string]string `json:"secretLabels,omitempty"`
	WriteOnFailure bool              `json:"writeOnFailure,omitempty"`
	Versioned      bool              `json:"versioned,omitempty"`
}

// ExtendedJobStatus defines the observed state of ExtendedJob
type ExtendedJobStatus struct {
	// Specifies whether the Job has reached phase Succeeded
	Succeeded bool `json:"succeeded"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ExtendedJob is the Schema for the extendedstatefulsetcontroller API
// +k8s:openapi-gen=true
type ExtendedJob struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ExtendedJobSpec   `json:"spec,omitempty"`
	Status ExtendedJobStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ExtendedJobList contains a list of ExtendedJob
type ExtendedJobList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ExtendedJob `json:"items"`
}

// ToBeDeleted checks whether this ExtendedJob has been marked for deletion
func (e *ExtendedJob) ToBeDeleted() bool {
	// IsZero means that the object hasn't been marked for deletion
	return !e.GetDeletionTimestamp().IsZero()
}

// IsAutoErrand returns true if this ext job is an auto errand
func (e *ExtendedJob) IsAutoErrand() bool {
	return e.Spec.Trigger.Strategy == TriggerOnce || e.Spec.Trigger.Strategy == TriggerDone
}
