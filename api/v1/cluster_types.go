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

// Blueprint Spec of Cluster Resource referred
type BlueprintInfoSpec struct {
	// Name of kind Blueprint
	Name string `json:"name,omitempty"`
	// Type Blueprint
	Type string `json:"type,omitempty"`
	// Revision of Blueprint
	Revision string `json:"revision,omitempty"`
	// Published Version of Blueprint
	Version string `json:"version,omitempty"`
}

// Content of Blueprint Packages
type BlueprintInfo struct {
	// Name of Blueprint
	Name string `json:"name,omitempty"`
	// Spec
	Spec BlueprintSpec `json:"spec,omitempty"`
	// Override field of blueprint
	Override map[string]string `json:"override,omitempty"`
}

type BlueprintInfoList struct {
	Items []Blueprint `json:"items,omitempty"`
}

// ClusterSpec defines the desired state of Cluster
type ClusterSpec struct {
	Infrastructure BlueprintInfoList `json:"infrastructure,omitempty"`
	Software       BlueprintInfoList `json:"software,omitempty"`
}

// ClusterStatus defines the observed state of Cluster
type ClusterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Cluster is the Schema for the clusters API
type Cluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterSpec   `json:"spec,omitempty"`
	Status ClusterStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ClusterList contains a list of Cluster
type ClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Cluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Cluster{}, &ClusterList{})
}
