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

// SVNUserSpec defines the desired state of SVNUser
type SVNUserSpec struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	// The name of the SVNServer
	SVNServer string `json:"svnServer,omitempty"`

	// Groups is a list of SVNGroups that the user belongs to.
	Groups []string `json:"groups,omitempty"`

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern="^[a-zA-Z0-9+/=.${}]+$"
	// EncryptedPassword is a password encrypted by `htpasswd`.
	// This must be computed elsewhere in order to avoid additional complexity of
	// letting controllers manage sensitive values.
	//
	// This field can be computed by the following command:
	//   $ htpasswd -nB USERNAME | cut -d : -f 2-
	//   New password: (TYPE YOUR PASSWORD HERE)
	//   Re-type new password: (TYPE YOUR PASSWORD HERE)
	//   $2y$05$Z9loUIkf0DynjbD0UMEpneKCSKYfkTCaE/pwY8wt7MtKQILxKRwjG (example output)
	//
	// See https://httpd.apache.org/docs/2.4/misc/password_encryptions.html for more information.
	EncryptedPassword string `json:"encryptedPassword,omitempty"`
}

// SVNUserStatus defines the observed state of SVNUser
type SVNUserStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// SVNUser is the Schema for the svnusers API
type SVNUser struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SVNUserSpec   `json:"spec,omitempty"`
	Status SVNUserStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// SVNUserList contains a list of SVNUser
type SVNUserList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SVNUser `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SVNUser{}, &SVNUserList{})
}
