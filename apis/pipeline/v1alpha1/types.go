package v1alpha1

import (
	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// A ValueSelector is a selector for a configMap or a secret in an arbitrary namespace.
type ValueSelector struct {
	// Name of the secret.
	Name string `json:"name"`

	// Namespace of the secret.
	Namespace string `json:"namespace"`

	// The key to select.
	Key string `json:"key"`
}

type PipelineParams struct {
	// JobName: the name you would like to give the new job name.
	JobName string `json:"jobName"`

	// JobConfigRef: configMap containing the config that
	// can be used to create the new job.
	JobConfigRef ValueSelector `json:"jobConfigRef"`
}

type PipelineObservation struct {
	// JobName: the name of the created job.
	JobName *string `json:"jobName,omitempty"`
}

// A PipelineSpec defines the desired state of a Pipeline.
type PipelineSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       PipelineParams `json:"forProvider"`
}

// A PipelineStatus represents the observed state of a Pipeline.
type PipelineStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          PipelineObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A Pipeline is a managed resource that represents an Jenkins Pipeline
// +kubebuilder:printcolumn:name="JOB_NAME",type="string",JSONPath=".status.atProvider.jobName"
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status",priority=1
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status",priority=1
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,krateo,jenkins}
type Pipeline struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PipelineSpec   `json:"spec"`
	Status PipelineStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// PipelineList contains a list of Pipeline.
type PipelineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Pipeline `json:"items"`
}
