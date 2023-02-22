package controllers

import (
	intentv1 "github.com/ntnguyencse/intent-kaas/api/v1"
	"github.com/ntnguyencse/intent-kaas/pkg/git"
)

func (r *ClusterReconciler) TransformClusterToClusterDescription(clusterCR intentv1.Cluster, listBluePrint []intentv1.Blueprint, gitClient *git.GitClient) (intentv1.ClusterDescription, error) {
	logger1.Info("Starting transform cluster to cluster description")
	// Transform cluster
	var clusterDescription intentv1.ClusterDescription
	// Get list blueprint Infra
	blueprintInfraInfos := clusterCR.Spec.Infrastructure
	for _, bpInfra := range blueprintInfraInfos {
		// Find all value of blueprint and nested blueprint
		findInfoOfBluePrint(bpInfra, listBluePrint)
	}
	// Get list Blueprint Software

	return clusterDescription, nil
}
func findInfoOfBluePrint(info intentv1.BlueprintInfo, listBP []intentv1.Blueprint) (map[string]string, error) {
	var infoBP map[string]string
	// Recursive find info of  nested blueprint
	// Name string `json:"name,omitempty"`
	// Spec BlueprintInfoSpec `json:"spec,omitempty"`
	// Override map[string]string `json:"override,omitempty"`
	
	return infoBP, nil
}
