/*
Copyright 2023 Nguyen Thanh Nguyen.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DescriptionSpec struct {
	BlueprintInfo ProfileInfo `json:"info,omitempty"`

	Spec map[string]string `json:"spec,omitempty"`
}

// ClusterDescriptionSpec defines the desired state of ClusterDescription
type ClusterDescriptionSpec struct {
	Infrastructure []DescriptionSpec `json:"infrastructure,omitempty"`

	Software []DescriptionSpec `json:"software,omitempty"`
}

// ClusterDescriptionStatus defines the observed state of ClusterDescription
type ClusterDescriptionStatus struct {
	Status string `json:"status,omitempty"`
	// Revision of cluster description
	Revision int64 `json:"revision,omitempty"`
	// Number of master node
	MasterStatus int `json:"masterstatus,omitempty"`
	// Number of worker node
	WorkerStatus int `json:"workerstatus,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:pruning:PreserveUnknownFields
// ClusterDescription is the Schema for the clusterdescriptions API
type ClusterDescription struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterDescriptionSpec   `json:"spec,omitempty"`
	Status ClusterDescriptionStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:pruning:PreserveUnknownFields
// ClusterDescriptionList contains a list of ClusterDescription
type ClusterDescriptionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterDescription `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ClusterDescription{}, &ClusterDescriptionList{})
}
