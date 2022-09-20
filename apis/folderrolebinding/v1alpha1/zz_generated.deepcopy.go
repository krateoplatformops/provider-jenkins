//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2022 Kiratech S.p.A.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FolderRoleBinding) DeepCopyInto(out *FolderRoleBinding) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FolderRoleBinding.
func (in *FolderRoleBinding) DeepCopy() *FolderRoleBinding {
	if in == nil {
		return nil
	}
	out := new(FolderRoleBinding)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *FolderRoleBinding) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FolderRoleBindingList) DeepCopyInto(out *FolderRoleBindingList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]FolderRoleBinding, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FolderRoleBindingList.
func (in *FolderRoleBindingList) DeepCopy() *FolderRoleBindingList {
	if in == nil {
		return nil
	}
	out := new(FolderRoleBindingList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *FolderRoleBindingList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FolderRoleBindingObservation) DeepCopyInto(out *FolderRoleBindingObservation) {
	*out = *in
	if in.SID != nil {
		in, out := &in.SID, &out.SID
		*out = new(string)
		**out = **in
	}
	if in.Name != nil {
		in, out := &in.Name, &out.Name
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FolderRoleBindingObservation.
func (in *FolderRoleBindingObservation) DeepCopy() *FolderRoleBindingObservation {
	if in == nil {
		return nil
	}
	out := new(FolderRoleBindingObservation)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FolderRoleBindingParams) DeepCopyInto(out *FolderRoleBindingParams) {
	*out = *in
	if in.Permissions != nil {
		in, out := &in.Permissions, &out.Permissions
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.FolderNames != nil {
		in, out := &in.FolderNames, &out.FolderNames
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FolderRoleBindingParams.
func (in *FolderRoleBindingParams) DeepCopy() *FolderRoleBindingParams {
	if in == nil {
		return nil
	}
	out := new(FolderRoleBindingParams)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FolderRoleBindingSpec) DeepCopyInto(out *FolderRoleBindingSpec) {
	*out = *in
	in.ResourceSpec.DeepCopyInto(&out.ResourceSpec)
	in.ForProvider.DeepCopyInto(&out.ForProvider)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FolderRoleBindingSpec.
func (in *FolderRoleBindingSpec) DeepCopy() *FolderRoleBindingSpec {
	if in == nil {
		return nil
	}
	out := new(FolderRoleBindingSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FolderRoleBindingStatus) DeepCopyInto(out *FolderRoleBindingStatus) {
	*out = *in
	in.ResourceStatus.DeepCopyInto(&out.ResourceStatus)
	in.AtProvider.DeepCopyInto(&out.AtProvider)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FolderRoleBindingStatus.
func (in *FolderRoleBindingStatus) DeepCopy() *FolderRoleBindingStatus {
	if in == nil {
		return nil
	}
	out := new(FolderRoleBindingStatus)
	in.DeepCopyInto(out)
	return out
}
