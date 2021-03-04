/*
Copyright 2021 Genta Kamitani.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// SVNGroupSpec defines the desired state of SVNGroup
type SVNGroupSpec struct {
	// +kubebuilder:validation:Required
	// The name of the SVNServer
	SVNServer string `json:"server,omitempty"`

	// +kubebuilder:validation:Required
	// The permissions that the group have.
	Permissions []Permission `json:"permissions,omitempty"`
}

type Permission struct {
	// +kubebuilder:validation:Required
	// The name of the SVNRepository to give access to.
	// The SVNRepository must reside in the same namespace as the SVNGroup.
	Repository string `json:"repository,omitempty"`

	// +kubebuilder:validation:Pattern="r|rw|"
	// The permission to access to the repository.
	Permission string `json:"permission,omitempty"`
}

// Here is a list of allowed permissions.
// We do not define types like Permission because currently kubebuilder seems
// to not support kubebuilder:validation:Pattern with named types.
const (
	// PermissionNone means the group have no permission to the repository.
	PermissionNone = ""

	// PermissionR means the group only can read from the repository.
	PermissionR = "r"

	// PermissionRW means the group can both read from and write to the repository.
	PermissionRW = "rw"
)

// SVNGroupStatus defines the observed state of SVNGroup
type SVNGroupStatus struct {
	// +Kubebuilder:validation:Optional
	Conditions []Condition `json:"conditions"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// SVNGroup is the Schema for the svngroups API
type SVNGroup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SVNGroupSpec   `json:"spec,omitempty"`
	Status SVNGroupStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// SVNGroupList contains a list of SVNGroup
type SVNGroupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SVNGroup `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SVNGroup{}, &SVNGroupList{})
}
