/*
Copyright 2025.

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

// MLModelDeploymentSpec defines the desired state of MLModelDeployment.
type MLModelDeploymentSpec struct {
	Image    string `json:"image"`
	Replicas int    `json:"replicas"`
	Host     string `json:"host"`
}

// MLModelDeploymentStatus defines the observed state of MLModelDeployment.
type MLModelDeploymentStatus struct {
	Phase string `json:"phase,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// MLModelDeployment is the Schema for the mlmodeldeployments API.
type MLModelDeployment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MLModelDeploymentSpec   `json:"spec,omitempty"`
	Status MLModelDeploymentStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MLModelDeploymentList contains a list of MLModelDeployment.
type MLModelDeploymentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MLModelDeployment `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MLModelDeployment{}, &MLModelDeploymentList{})
}
