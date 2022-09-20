package v1alpha1

import (
	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type FolderRoleBindingParams struct {
	// SID:
	SID string `json:"sid"`

	// Name:
	Name string `json:"name"`

	// Permissions:
	Permissions []string `json:"permissions"`

	// FolderNames:
	FolderNames []string `json:"folderNames"`
}

type FolderRoleBindingObservation struct {
	// SID:
	SID *string `json:"sid,omitempty"`

	// Name:
	Name *string `json:"name,omitempty"`
}

// A FolderRoleBindingSpec defines the desired state of a FolderRoleBinding.
type FolderRoleBindingSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       FolderRoleBindingParams `json:"forProvider"`
}

// A FolderRoleBindingStatus represents the observed state of a FolderRoleBinding.
type FolderRoleBindingStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          FolderRoleBindingObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A FolderRoleBinding is a managed resource that represents an Jenkins FolderRoleBinding
// +kubebuilder:printcolumn:name="NAME",type="string",JSONPath=".status.atProvider.name"
// +kubebuilder:printcolumn:name="SID",type="string",JSONPath=".status.atProvider.sid"
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status",priority=1
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status",priority=1
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,krateo,jenkins}
type FolderRoleBinding struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FolderRoleBindingSpec   `json:"spec"`
	Status FolderRoleBindingStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// FolderRoleBindingList contains a list of FolderRoleBinding.
type FolderRoleBindingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []FolderRoleBinding `json:"items"`
}
