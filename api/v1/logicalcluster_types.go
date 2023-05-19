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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// LogicalClusterSpec defines the desired state of LogicalCluster
type LogicalClusterSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of LogicalCluster. Edit logicalcluster_types.go to remove/update
	Clusters []ClusterMember `json:"clusters,omitempty"`
}
type ClusterMemberSpec struct {
	// Cluster Catalog
	//+optional
	ClusterCatalog string `json:"clustercatalog,omitempty"`
	// Cluster Detail Spec
	//+optional
	ClusterSpec ClusterSpec `json:"spec,omitempty"`
}
type ClusterMember struct {
	Name string `json:"name,omitempty"`
	// ClusterCatalog    string `json:"clustercatalog,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	// Cluster member of logical cluster. Each cluster member associate with a physical cluster (CAPI)
	//+optional
	ClusterMemberSpec ClusterMemberSpec `json:"spec,omitempty"`
	// ClusterRef is a reference to a L-KaaS cluster that holds the details cluster
	// +optional
	ClusterRef *corev1.ObjectReference `json:"clusterref,omitempty"`
}

// LogicalClusterStatus defines the observed state of LogicalCluster
type LogicalClusterStatus struct {

	// Ready state of Logical cluster
	Ready bool `json:"ready,omitempty"`
	// State of Each Cluster Member
	ClusterMemberState string `json:"clusterMemberState,omitempty"`
	// Failure Message
	// +optional
	FailureMessage string `json:"failureMessage,omitempty"`
	// Failure Reason
	// +optional
	FailureReason string `json:"failureReason,omitempty"`
	// Logical Cluster Conditions
	// +optional
	Conditions ConditionType `json:"conditions,omitempty"`

	// Registration Status
	Registration bool `json:"registration,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// LogicalCluster is the Schema for the logicalclusters API
type LogicalCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LogicalClusterSpec   `json:"spec,omitempty"`
	Status LogicalClusterStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// LogicalClusterList contains a list of LogicalCluster
type LogicalClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []LogicalCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&LogicalCluster{}, &LogicalClusterList{})
}
