package controllers

import (
	"context"

	intentv1 "github.com/ntnguyencse/L-KaaS/api/v1"
	// "github.com/ntnguyencse/L-KaaS/pkg/git"
)

func (r *ClusterReconciler) TransformClusterToClusterDescription(ctx context.Context, clusterCR intentv1.Cluster, listBluePrint []intentv1.Profile) (intentv1.ClusterDescription, error) {
	loggerCL.Info("Starting transform cluster to cluster description")
	// Transform cluster
	var clusterDescription intentv1.ClusterDescription
	// Form a clusterDescription from cluster
	clusterDescription.Name = clusterCR.Name
	clusterDescription.Labels = clusterCR.Labels
	clusterDescription.Annotations = clusterCR.Annotations
	clusterDescription.Namespace = "clusters"
	//
	// Get list blueprint Infra
	blueprintInfraInfos := clusterCR.Spec.Infrastructure
	for _, bpInfra := range blueprintInfraInfos {

		// Find all value of blueprint and nested blueprint
		// Get all value of blueprint
		valu, _ := findInfoOfBluePrint(bpInfra, listBluePrint)

		desSpec := intentv1.DescriptionSpec{
			BlueprintInfo: bpInfra,
			Spec:          valu,
		}
		clusterDescription.Spec.Infrastructure = append(clusterDescription.Spec.Infrastructure, desSpec)
		loggerCL.Info("Print value infra blueprint", "value", valu)
	}
	// Get list Blueprint Software
	blueprintSoftware := clusterCR.Spec.Software

	for _, bpSoftware := range blueprintSoftware {

		// Get all value of blueprint
		valu, _ := findInfoOfBluePrint(bpSoftware, listBluePrint)

		desSpec := intentv1.DescriptionSpec{
			BlueprintInfo: bpSoftware,
			Spec:          valu,
		}
		clusterDescription.Spec.Software = append(clusterDescription.Spec.Software, desSpec)
		loggerCL.Info("Print value software blueprint", "value", valu)
	}
	// r.Client.Get()
	// Get list blueprint Network
	blueprintNetwork := clusterCR.Spec.Network

	for _, bpNetwork := range blueprintNetwork {

		// Get all value of blueprint
		valu, _ := findInfoOfBluePrint(bpNetwork, listBluePrint)

		desSpec := intentv1.DescriptionSpec{
			BlueprintInfo: bpNetwork,
			Spec:          valu,
		}
		clusterDescription.Spec.Network = append(clusterDescription.Spec.Network, desSpec)
		loggerCL.Info("Print value network blueprint", "value", valu)
	}

	return clusterDescription, nil
}
func findInfoOfBluePrint(info intentv1.ProfileInfo, listBP []intentv1.Profile) (map[string]string, error) {
	var infoBP map[string]string
	// Recursive find info of  nested blueprint
	// Name string `json:"name,omitempty"`
	// Spec BlueprintInfoSpec `json:"spec,omitempty"`
	// Override map[string]string `json:"override,omitempty"`
	// Layer 1 blueprint
	for _, bp := range listBP {
		loggerCL.Info(info.Name, "findInfoOfBluePrint", bp.Name)

		if bp.Name == info.Spec.Name {

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
func findInforOfBlueprintSpec(inforSpec intentv1.ProfileInfoSpec, listBP []intentv1.Profile) (map[string]string, error) {

	var infoBP map[string]string
	for _, bp := range listBP {
		loggerCL.Info(inforSpec.Name, "findInforOfBlueprintSpec", bp.Name)
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
