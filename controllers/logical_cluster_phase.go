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

package controllers

import (
	"context"

	intentv1 "github.com/ntnguyencse/L-KaaS/api/v1"
	capiv1alpha4 "sigs.k8s.io/cluster-api/api/v1alpha4"
)

var cluster *capiv1alpha4.Cluster

func (r *LogicalClusterControlPlaneProviderReconciler) ReconcileClusterPhase(ctx context.Context, logicalCluster *intentv1.LogicalCluster, capiCluster *capiv1alpha4.ClusterList) error {
	// Change phase of Logical Cluster
	// based on Status of CAPI Clusters
	// TODO: Defined Cluster Phase
	// logicalCluster.Status
	return nil
}
func SetLogicalClusterPhase()

