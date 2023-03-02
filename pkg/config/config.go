package config

import (
	"github.com/jinzhu/configor"
)

const DEFAULT_CONFIG_PATH string = "./config.yml"

type Configuration struct {
	Owner         string `required:"true" env:"OWNER"`
	BlueprintRepo string `required:"true" env:"BLUEPRINT_REPO"`
	ClusterRepo   string `required:"true" env:"CLUSTER_REPO"`
	GitHubToken   string `required:"true" env:"GITHUB_TOKEN"`
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
