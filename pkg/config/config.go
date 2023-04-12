package config

import (
	"github.com/jinzhu/configor"
)

const DEFAULT_CONFIG_PATH string = "./config.yml"
const DEFAULT_OPENSTACKCONFIG_PATH string = "./openstack-config.yml"

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

type OpenStackConfiguration struct {
	OPENSTACK_IMAGE_NAME                   string `required:"true" env:"OPENSTACK_IMAGE_NAME"`
	OPENSTACK_EXTERNAL_NETWORK_ID          string `required:"true" env:"OPENSTACK_EXTERNAL_NETWORK_ID"`
	OPENSTACK_DNS_NAMESERVERS              string `required:"true" env:"OPENSTACK_DNS_NAMESERVERS"`
	OPENSTACK_SSH_KEY_NAME                 string `required:"true" env:"OPENSTACK_SSH_KEY_NAME"`
	OPENSTACK_CLOUD_CACERT_B64             string `required:"true" env:"OPENSTACK_CLOUD_CACERT_B64"`
	OPENSTACK_CLOUD_PROVIDER_CONF_B64      string `required:"true" env:"OPENSTACK_CLOUD_PROVIDER_CONF_B64"`
	OPENSTACK_CLOUD_YAML_B64               string `required:"true" env:"OPENSTACK_CLOUD_YAML_B64"`
	OPENSTACK_FAILURE_DOMAIN               string `required:"true" env:"OPENSTACK_FAILURE_DOMAIN"`
	OPENSTACK_CLOUD                        string `required:"true" env:"OPENSTACK_CLOUD"`
	OPENSTACK_CONTROL_PLANE_MACHINE_FLAVOR string `required:"true" env:"OPENSTACK_CONTROL_PLANE_MACHINE_FLAVOR"`
	OPENSTACK_NODE_MACHINE_FLAVOR          string `required:"true" env:"OPENSTACK_NODE_MACHINE_FLAVOR"`
}

func LoadOpenStackConfig(path string) OpenStackConfiguration {
	var config OpenStackConfiguration
	if len(path) < 1 {
		configor.Load(&config, DEFAULT_OPENSTACKCONFIG_PATH)
	} else {
		configor.Load(&config, path)
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
