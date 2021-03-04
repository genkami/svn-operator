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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SVNServerSpec defines the desired state of SVNServer
type SVNServerSpec struct {
	// +kubebuilder:validation:Required
	// PodTemplate is a template to create Pods.
	PodTemplate PodTemplate `json:"podTemplate,omitempty"`

	// +kubebuilder:validation:Required
	// VolumeClaimTemplate is a PVC to store SVN repositories and configuration files in.
	VolumeClaimTemplate corev1.PersistentVolumeClaimSpec `json:"volumeClaimTemplate,omitempty"`
}

// PodTemplate is an optional template to create SVN server pods.
type PodTemplate struct {
	// +kubebuilder:validation:Optional
	// Image specifies a container image of SVN server.
	// If not specified, the default value will be used.
	Image string `json:"image,omitempty"`

	// +kubebuilder:validation:Optional
	// NodeSelector is a selector which must be true for the pod to fit on a node.
	// Selector which must match a node's labels for the pod to be scheduled on that node.
	// More info: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	// ServiceAccountName is the name of the ServiceAccount to use to run this pod.
	// More info: https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/
	ServiceAccountName string `json:"serviceAccountName,omitempty"`

	// +kubebuilder:validation:Optional
	// ImagePullSecrets is an optional list of references to secrets in the same namespace to use for pulling any of the images used by this PodSpec.
	// If specified, these secrets will be passed to individual puller implementations for them to use. For example,
	// in the case of docker, only DockerConfig type secrets are honored.
	// More info: https://kubernetes.io/docs/concepts/containers/images#specifying-imagepullsecrets-on-a-pod
	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty"`

	// +kubebuilder:validation:Optional
	// If specified, the pod's scheduling constraints
	Affinity *corev1.Affinity `json:"affinity,omitempty"`

	// +kubebuilder:validation:Optional
	// If specified, the pod's tolerations.
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`
}

// SVNServerStatus defines the observed state of SVNServer
type SVNServerStatus struct {
	// +kubebuilder:validation:Optional
	Conditions []Condition `json:"conditions"`
}

type Condition struct {
	Type ConditionType `json:"type"`

	Reason string `json:"reason,omitempty"`

	// The time when the SVNServer's condition changed in RFC3339 format.
	TransitionTime string `json:"transitionTime"`
}

type ConditionType string

const (
	ConditionTypeNone   ConditionType = ""
	ConditionTypeSynced ConditionType = "Synced"
	ConditionTypeFailed ConditionType = "Failed"
)

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// SVNServer is the Schema for the svnservers API
type SVNServer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SVNServerSpec   `json:"spec,omitempty"`
	Status SVNServerStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// SVNServerList contains a list of SVNServer
type SVNServerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SVNServer `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SVNServer{}, &SVNServerList{})
}
