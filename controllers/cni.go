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
	"fmt"
	_ "fmt"
	_ "os"
	"time"

	kubernetesclient "github.com/ntnguyencse/L-KaaS/pkg/kubernetes-client"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type InstallComponent struct {
	Version         string
	Name            string
	URL             string
	TargetNamespace string
	KubeConfigPath  string
}
type Installer struct {
	Client     client.Client
	Components []InstallComponent
	Func       func() error
}
type InstallOptions struct {
	WaitProviders       bool
	WaitProviderTimeout time.Duration
}

// Install middle ware
func SetUpInstaller(client client.Client) Installer {
	var installer Installer
	installer.Client = client
	return installer
}
func CreateInstallerComponent(Name, Version, URL, TargetNamespace, KubeconfigPath string) InstallComponent {
	return InstallComponent{
		Name:            Name,
		Version:         Version,
		URL:             URL,
		TargetNamespace: TargetNamespace,
		KubeConfigPath:  KubeconfigPath,
	}
}

func (i *Installer) AddInstallComponent(item InstallComponent) {
	i.Components = append(i.Components, item)
}
func (i *Installer) Install(clusterName string) error {
	// folder := "/tmp/" + clusterName + "/"
	//
	for _, item := range i.Components {
		// Download file with URL
		// r, err := cloudfile.Open(item.URL)
		// if err != nil {
		// 	loggerLKP.Error(err, "Error read file from remote url: "+item.URL)
		// 	return err
		// }

		// defer r.Close()
		// strBinary, err := io.ReadAll(r)
		// if err != nil {
		// 	loggerLKP.Error(err, "Error read file GetTemplateFile io.ReadAll: "+item.URL)
		// 	return err
		// }
		// result := string(strBinary)
		// componentYamlFilePath, err := emcoctl.SaveValueFile(item.Name+".yaml", folder, &result)
		// if err != nil {
		// 	fmt.Println(err, "Error when store yaml fiel for install: "+item.URL)
		// 	return err
		// }
		//  Install yaml file
		fmt.Println("Applying Resource....")
		fmt.Println("Resource Name: ", item.Name, "Resource URL:", item.URL)
		err := kubernetesclient.ApplyResourceKubernetesWithKubeConfig(item.KubeConfigPath, item.URL)
		if err != nil {
			// fmt.Println("Resource Name: ", item.Name, "Resource URL:", item.URL)
			fmt.Println("Error when apply resources", err)

			return err
		}
	}
	return nil
}

// Install Calico
func InstallCalioCNI(installer *Installer, clusterName string, components ...InstallComponent) {
	// Download Calico Operator
	// Components for install Calico Operator
	for _, item := range components {
		installer.AddInstallComponent(item)
	}
	installer.Install(clusterName)

}
