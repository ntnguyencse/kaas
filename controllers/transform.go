package controllers

import (
	"context"
	"os"

	intentv1 "github.com/ntnguyencse/L-KaaS/api/v1"
	// "github.com/ntnguyencse/L-KaaS/pkg/git"
	CAPIClient "github.com/ntnguyencse/L-KaaS/pkg/client"
	config "github.com/ntnguyencse/L-KaaS/pkg/config"
)

const OPENSTACK_PROVIDER_URL string = "https://github.com/kubernetes-sigs/cluster-api-provider-openstack/releases/download/v0.7.1/infrastructure-components.yaml"
const DEFAULT_CAPI_CONFIG_PATH string = "config/capi/clusterctl-config.yaml"

var KUBECONFIG string
var configs = map[string]string{
	"OPENSTACK_IMAGE_NAME":                   "OPENSTACK_IMAGE_NAME",
	"OPENSTACK_EXTERNAL_NETWORK_ID":          "OPENSTACK_EXTERNAL_NETWORK_ID",
	"OPENSTACK_DNS_NAMESERVERS":              "OPENSTACK_DNS_NAMESERVERS",
	"OPENSTACK_SSH_KEY_NAME":                 "OPENSTACK_SSH_KEY_NAME",
	"OPENSTACK_CLOUD_CACERT_B64":             "OPENSTACK_CLOUD_CACERT_B64",
	"OPENSTACK_CLOUD_PROVIDER_CONF_B64":      "OPENSTACK_CLOUD_PROVIDER_CONF_B64",
	"OPENSTACK_CLOUD_YAML_B64":               "OPENSTACK_CLOUD_YAML_B64",
	"OPENSTACK_FAILURE_DOMAIN":               "OPENSTACK_FAILURE_DOMAIN",
	"OPENSTACK_CLOUD":                        "OPENSTACK_CLOUD",
	"OPENSTACK_CONTROL_PLANE_MACHINE_FLAVOR": "OPENSTACK_CONTROL_PLANE_MACHINE_FLAVOR",
	"OPENSTACK_NODE_MACHINE_FLAVOR":          "OPENSTACK_NODE_MACHINE_FLAVOR",
	"KUBERNETES_VERSION":                     "1.24.5",
	"CONTROL_PLANE_MACHINE_COUNT":            "3",
	"WORKER_MACHINE_COUNT":                   "3",
}

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

// Transform to CAPI Resource
// Remember export KUBECONFIG
func TranslateFromClusterDescritionToCAPI(clusterDes *intentv1.ClusterDescription, configForProvider intentv1.ProviderConfig, configForCluster map[string]string) (string, error) {
	// Create Cluster API SDK Client
	// Get provider config
	// Currently, only use InfrastructureProviderType for boostraping cluster
	providerConfigs := CAPIClient.CreateProviderConfig(configForProvider.Name, configForProvider.URL, configForProvider.ProviderType)
	// Create client
	KUBECONFIG = os.Getenv("KUBECONFIG")
	clientctl, err := CAPIClient.CreateNewClient(KUBECONFIG, configForCluster, providerConfigs)
	if err != nil {
		loggerCL.Error(err, "Error when create CAPI client")
	}
	// Generate Cluster
	loggerCL.Info("Print URL", "URL", configForCluster["CAPI_TEMPLATE_URL"])
	// url := "https://github.com/kubernetes-sigs/cluster-api-provider-openstack/blob/main/templates/cluster-template.yaml"
	clusterString, err := clientctl.GetClusterTemplate(clusterDes.Name, clusterDes.Namespace, configForCluster["CAPI_TEMPLATE_URL"])

	return clusterString, err
}

func getCredentialsForOpenStackProvider(configPath string) (map[string]string, error) {
	// Current only support for OPENSTACK, Edit this function to support more provider
	if configPath == "" {
		configPath = DEFAULT_CAPI_CONFIG_PATH
	}

	providerConfigLoaded := config.LoadOpenStackCredentials(configPath)
	loggerCL.Info("Print LoadOpenStackCredentials", "Configs", providerConfigLoaded)
	secrets := map[string]string{
		"OPENSTACK_IMAGE_NAME":                   providerConfigLoaded.OpenstackImageName,
		"OPENSTACK_EXTERNAL_NETWORK_ID":          providerConfigLoaded.OpenstackExternalNetworkId,
		"OPENSTACK_DNS_NAMESERVERS":              providerConfigLoaded.OpenstackDNSNameservers,
		"OPENSTACK_SSH_KEY_NAME":                 providerConfigLoaded.OpenstackSshKeyName,
		"OPENSTACK_CLOUD_CACERT_B64":             providerConfigLoaded.OpenstackCloudCacertB64,
		"OPENSTACK_CLOUD_PROVIDER_CONF_B64":      providerConfigLoaded.OpenstackCloudProviderConfB64,
		"OPENSTACK_CLOUD_YAML_B64":               providerConfigLoaded.OpenstackCloudYamlB64,
		"OPENSTACK_FAILURE_DOMAIN":               providerConfigLoaded.OpenstackFailureDomain,
		"OPENSTACK_CLOUD":                        providerConfigLoaded.OpenstackCloud,
		"OPENSTACK_CONTROL_PLANE_MACHINE_FLAVOR": providerConfigLoaded.OpenstackControlPlaneMachineFlavor,
		"OPENSTACK_NODE_MACHINE_FLAVOR":          providerConfigLoaded.OpenstackNodeMachineFlavor,
		"KUBERNETES_VERSION":                     providerConfigLoaded.KubernetesVersion,
	}
	// providerConfig["OPENSTACK_DNS_NAMESERVERS"] = providerConfigLoaded.OPENSTACK_DNS_NAMESERVERS
	// providerConfig["OPENSTACK_IMAGE_NAME"] = providerConfigLoaded.OPENSTACK_IMAGE_NAME
	// providerConfig["OPENSTACK_EXTERNAL_NETWORK_ID"] = providerConfigLoaded.OPENSTACK_EXTERNAL_NETWORK_ID
	// providerConfig["OPENSTACK_SSH_KEY_NAME"] = providerConfigLoaded.OPENSTACK_SSH_KEY_NAME
	// providerConfig["OPENSTACK_CLOUD_CACERT_B64"] = providerConfigLoaded.OPENSTACK_CLOUD_CACERT_B64
	// providerConfig["OPENSTACK_CLOUD_PROVIDER_CONF_B64"] = providerConfigLoaded.OPENSTACK_CLOUD_PROVIDER_CONF_B64
	// providerConfig["OPENSTACK_CLOUD_YAML_B64"] = providerConfigLoaded.OPENSTACK_CLOUD_PROVIDER_CONF_B64
	// providerConfig["OPENSTACK_CLOUD"] = providerConfigLoaded.OPENSTACK_CLOUD_PROVIDER_CONF_B64
	// providerConfig["OPENSTACK_CONTROL_PLANE_MACHINE_FLAVOR"] = providerConfigLoaded.OPENSTACK_CLOUD_PROVIDER_CONF_B64
	// providerConfig["OPENSTACK_NODE_MACHINE_FLAVOR"] = providerConfigLoaded.OPENSTACK_CLOUD_PROVIDER_CONF_B64

	return secrets, nil
}
func GetConfigForOpenStack() intentv1.ProviderConfig {
	return intentv1.ProviderConfig{
		Name:         CAPIClient.OPENSTACK,
		URL:          CAPIClient.OPENSTACK_URL,
		ProviderType: CAPIClient.InfrastructureProviderType,
	}
}
