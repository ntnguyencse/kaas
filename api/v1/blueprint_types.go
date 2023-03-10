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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// BlueprintSpec defines the desired state of Blueprint
type BlueprintSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Blueprint refer to another blueprint
	Blueprints []BlueprintInfoSpec `json:"blueprint,omitempty"`
	Values     map[string]string   `json:"values,omitempty"`
}

// +kubebuilder:subresource:status
// BlueprintStatus defines the observed state of Blueprint
type BlueprintStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// Status of Blueprint
	Status string `json:"status,omitempty"`
	// Sync status of Blueprint
	Sync string `json:"sync,omitempty"`
	// SHA of blueprint
	SHA string `json:"sha,omitempty"`
	// Repo contains blueprint
	Repo string `json:"repo,omitempty"`
	// Version of blueprint
	Version string `json:"version,omitempty"`
	// Revision of blueprint
	Revision int64 `json:"revision,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:pruning:PreserveUnknownFields
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="Time duration since creation of Cluster"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.status",description="Cluster status"
// +kubebuilder:printcolumn:name="SHA",type="string",JSONPath=".status.sha",description="SHA"
// +kubebuilder:printcolumn:name="Repo",type="string",JSONPath=".status.repo",description="Repo"
// +kubebuilder:printcolumn:name="Sync",type="string",JSONPath=".status.sync",description="Sync"
// +kubebuilder:printcolumn:name="Revision",type="integer",JSONPath=".status.revision",description="Revision"
// Blueprint is the Schema for the blueprints API
type Blueprint struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BlueprintSpec   `json:"spec,omitempty"`
	Status BlueprintStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// BlueprintList contains a list of Blueprint
type BlueprintList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Blueprint `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Blueprint{}, &BlueprintList{})
}
