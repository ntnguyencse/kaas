package config

import (
	"fmt"

	"github.com/jinzhu/configor"
)

const DEFAULT_CONFIG_PATH string = "./config.yml"
const DEFAULT_OPENSTACKCONFIG_PATH string = "./openstack-config.yml"
const DEFAULT_AWSCONFIG_PATH string = "./capactl.yml"

type Configuration struct {
	Owner                          string `required:"true" env:"OWNER"`
	BlueprintRepo                  string `required:"true" env:"BLUEPRINT_REPO"`
	ClusterRepo                    string `required:"true" env:"CLUSTER_REPO"`
	GitHubToken                    string `required:"true" env:"GITHUB_TOKEN"`
	OpenStackConfigurationFilePath string `required:"false" env:"OS_CONFIGPATH"`
	// AWSConfigurationFilePath       string `required:"false" env:"AWS_CONFIGPATH"`
}

func LoadConfig(path string) Configuration {
	var config Configuration
	if len(path) < 1 {
		configor.Load(&config, DEFAULT_CONFIG_PATH)
	} else {
		configor.Load(&config, path)
	}
	return config
}

// Variables of config
// In config file, all variables must not be in UPPER CASE, All names must in LOWER CASE
// Ex: OpenstackImageName => openstackimagename
type OpenStackConfiguration struct {
	OpenstackImageName                 string `required:"true" env:"OPENSTACK_IMAGE_NAME"`
	OpenstackExternalNetworkId         string `required:"true" env:"OPENSTACK_EXTERNAL_NETWORK_ID"`
	OpenstackDNSNameservers            string `required:"true" env:"OPENSTACK_DNS_NAMESERVERS"`
	OpenstackSshKeyName                string `required:"true" env:"OPENSTACK_SSH_KEY_NAME"`
	OpenstackCloudCacertB64            string `required:"true" env:"OPENSTACK_CLOUD_CACERT_B64"`
	OpenstackCloudProviderConfB64      string `required:"true" env:"OPENSTACK_CLOUD_PROVIDER_CONF_B64"`
	OpenstackCloudYamlB64              string `required:"true" env:"OPENSTACK_CLOUD_YAML_B64"`
	OpenstackFailureDomain             string `required:"true" env:"OPENSTACK_FAILURE_DOMAIN"`
	OpenstackCloud                     string `required:"true" env:"OPENSTACK_CLOUD"`
	OpenstackControlPlaneMachineFlavor string `required:"true" env:"OPENSTACK_CONTROL_PLANE_MACHINE_FLAVOR"`
	OpenstackNodeMachineFlavor         string `required:"true" env:"OPENSTACK_NODE_MACHINE_FLAVOR"`
	KubernetesVersion                  string `required:"true" env:"KUBERNETES_VERSION"`
}

type AWSConfiguration struct {
	AwsRegion                  string `required:"true" env:"AWS_REGION"`
	AwsSSHKeyName              string `required:"true" env:"AWS_SSH_KEY_NAME"`
	AwsControlPlaneMachineType string `required:"true" env:"AWS_CONTROL_PLANE_MACHINE_TYPE"`
	AwsNodeMachineType         string `required:"true" env:"AWS_NODE_MACHINE_TYPE"`
	KubernetesVersion          string `required:"true" env:"KUBERNETES_VERSION"`
	AwsB64EncodedCredentials   string `required:"true" env:"AWS_B64ENCODED_CREDENTIALS"`
}

func LoadAWSCredentials(path string) AWSConfiguration {
	var config AWSConfiguration
	if len(path) < 1 {
		configor.Load(&config, DEFAULT_AWSCONFIG_PATH)
	} else {
		configor.Load(&config, path)
	}
	return config
}

func LoadOpenStackCredentials(path string) OpenStackConfiguration {
	var config OpenStackConfiguration
	if len(path) < 1 {
		configor.Load(&config, DEFAULT_OPENSTACKCONFIG_PATH)
	} else {
		configor.Load(&config, path)
	}
	return config
}

// prerequisitestemplateurl: "https://raw.githubusercontent.com/ntnguyencse/L-KaaS/main/templates/emco/dcm/prerequisites.yaml"
// registrationtemplateurl: "https://raw.githubusercontent.com/ntnguyencse/L-KaaS/main/templates/emco/dcm/1stcluster.yaml"
// instantiatetemplatetrl: "https://raw.githubusercontent.com/ntnguyencse/L-KaaS/main/templates/emco/dcm/instantiate-lc.yaml"
// prerevaluestemplateurl: "https://raw.githubusercontent.com/ntnguyencse/L-KaaS/dev/templates/emco/dcm/values/prerequisites-values.yaml"
// updatetemplateurl: "https://raw.githubusercontent.com/ntnguyencse/L-KaaS/main/templates/emco/dcm/update-lc.yaml"

type EMCOResourceConfiguration struct {
	PrerequisiteTemplateUrl       string `required:"true" env:"PREREQUISITETEMPLATEURL"`
	RegistrationTemplateUrl       string `required:"true" env:"REGISTRATIONTEMPLATEURL"`
	InstantiateTemplateUrl        string `required:"true" env:"INSTANTIATETEMPLATEURL"`
	PrerequisiteValuesTemplateUrl string `required:"true" env:"PREREQUISITEVALUESTEMPLATEURL"`
	UpdateTemplateUrl             string `required:"true" env:"UPDATETEMPLATEURL"`
}

func LoadEMCOResourceConfiguration(path string) EMCOResourceConfiguration {
	var config EMCOResourceConfiguration

	er := configor.Load(&config, path)
	if er != nil {
		fmt.Println("Error when get config LoadEMCOResourceConfiguration ", er)
	}
	return config
}

//-------------------Example------------------//
// var Config = struct {
// 	APPName string `default:"app name"`

// 	DB struct {
// 		Name     string
// 		User     string `default:"root"`
// 		Password string `required:"true" env:"DBPassword"`
// 		Port     uint   `default:"3306"`
// 	}

// 	Contacts []struct {
// 		Name  string
// 		Email string `required:"true"`
// 	}
// }{}
//----------------------------------------//
// YAML File:
// appname: test

// db:
//     name:     test
//     user:     test
//     password: test
//     port:     1234

// contacts:
// - name: i test
//   email: test@test.com
