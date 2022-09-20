package v1alpha1

import (
	"reflect"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

const (
	Group   = "jenkins.krateo.io"
	Version = "v1alpha1"
)

var (
	// SchemeGroupVersion is group version used to register these objects
	SchemeGroupVersion = schema.GroupVersion{Group: Group, Version: Version}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilder = &scheme.Builder{GroupVersion: SchemeGroupVersion}
)

var (
	FolderRoleBindingKind             = reflect.TypeOf(FolderRoleBinding{}).Name()
	FolderRoleBindingGroupKind        = schema.GroupKind{Group: Group, Kind: FolderRoleBindingKind}.String()
	FolderRoleBindingKindAPIVersion   = FolderRoleBindingKind + "." + SchemeGroupVersion.String()
	FolderRoleBindingGroupVersionKind = SchemeGroupVersion.WithKind(FolderRoleBindingKind)
)

func init() {
	SchemeBuilder.Register(&FolderRoleBinding{}, &FolderRoleBindingList{})
}
