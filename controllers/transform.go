package controllers

import (
	"context"

	intentv1 "github.com/ntnguyencse/intent-kaas/api/v1"
	"github.com/ntnguyencse/intent-kaas/pkg/git"
)

func (r *ClusterReconciler) TransformClusterToClusterDescription(ctx context.Context, clusterCR intentv1.Cluster, listBluePrint []intentv1.Blueprint, gitClient *git.GitClient) (intentv1.ClusterDescription, error) {
	logger1.Info("Starting transform cluster to cluster description")
	// Transform cluster
	var clusterDescription intentv1.ClusterDescription
	// Get list blueprint Infra
	blueprintInfraInfos := clusterCR.Spec.Infrastructure
	for _, bpInfra := range blueprintInfraInfos {
		// Find all value of blueprint and nested blueprint
		valu, _ := findInfoOfBluePrint(bpInfra, listBluePrint)
		logger1.Info("Print value infra blueprint", "value", valu)
	}
	// Get list Blueprint Software
	blueprintSoftware := clusterCR.Spec.Software

	for _, bpSoftware := range blueprintSoftware {
		valu, _ := findInfoOfBluePrint(bpSoftware, listBluePrint)
		logger1.Info("Print value software blueprint", "value", valu)
	}
	// r.Client.Get()
	return clusterDescription, nil
}
func findInfoOfBluePrint(info intentv1.BlueprintInfo, listBP []intentv1.Blueprint) (map[string]string, error) {
	var infoBP map[string]string
	// Recursive find info of  nested blueprint
	// Name string `json:"name,omitempty"`
	// Spec BlueprintInfoSpec `json:"spec,omitempty"`
	// Override map[string]string `json:"override,omitempty"`
	// Layer 1 blueprint
	for _, bp := range listBP {
		if bp.Name == info.Name {
			// infoBP = bp.Spec.Values
			// Get all data from blueprint
			infoBP = merge2map(infoBP, bp.Spec.Values)
			// Layer 2 of blueprint
			if len(bp.Spec.Blueprints) > 0 {
				for _, subBP := range bp.Spec.Blueprints {
					infoSubBP, _ := findInforOfBlueprintSpec(subBP, listBP)
					infoBP = merge2map(infoBP, infoSubBP)
				}
			}
		}
	}
	return infoBP, nil
}
func findInforOfBlueprintSpec(inforSpec intentv1.BlueprintInfoSpec, listBP []intentv1.Blueprint) (map[string]string, error) {

	var infoBP map[string]string
	for _, bp := range listBP {
		if bp.Name == inforSpec.Name {
			infoBP = merge2map(infoBP, bp.Spec.Values)
			return infoBP, nil
		}
	}
	return infoBP, nil
}
func merge2map(map1, map2 map[string]string) map[string]string {
	if len(map1) < 1 {
		return map2
	}
	if len(map2) < 1 {
		return map1
	}

	for key, value := range map2 {
		map1[key] = value
	}
	return map1
}
