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

// SVNRepositorySpec defines the desired state of SVNRepository
type SVNRepositorySpec struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern="^[a-zA-Z0-9][a-zA-Z0-9.-]*$"
	// The name of the SVNServer
	SVNServer string `json:"svnServer,omitempty"`
}

// SVNRepositoryStatus defines the observed state of SVNRepository
type SVNRepositoryStatus struct {
	// +Kubebuilder:validation:Optional
	Conditions []Condition `json:"conditions"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// SVNRepository is the Schema for the svnrepositories API
//
// TODO: what if SVNRepository is deleted?
type SVNRepository struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SVNRepositorySpec   `json:"spec,omitempty"`
	Status SVNRepositoryStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// SVNRepositoryList contains a list of SVNRepository
type SVNRepositoryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SVNRepository `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SVNRepository{}, &SVNRepositoryList{})
}
