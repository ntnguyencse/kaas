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

// ClusterCatalogSpec defines the desired state of ClusterCatalog
type ClusterCatalogSpec struct {
	Infrastructure []ProfileInfo `json:"infrastructure,omitempty"`
	Network        []ProfileInfo `json:"network,omitempty"`
	Software       []ProfileInfo `json:"software,omitempty"`
}

// ClusterCatalogStatus defines the observed state of ClusterCatalog
type ClusterCatalogStatus struct {
	// Sync status of cluster
	Sync string `json:"sync,omitempty"`
	// SHA of cluster package
	SHA string `json:"sha,omitempty"`
	// Repo contains cluster package
	Repo string `json:"repo,omitempty"`
	// Version  of cluster package
	Version string `json:"version,omitempty"`
	// Revision of cluster package
	Revision int64 `json:"revision,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// ClusterCatalog is the Schema for the clustercatalogs API
type ClusterCatalog struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterCatalogSpec   `json:"spec,omitempty"`
	Status ClusterCatalogStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
// ClusterCatalogList contains a list of ClusterCatalog
type ClusterCatalogList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterCatalog `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ClusterCatalog{}, &ClusterCatalogList{})
}
