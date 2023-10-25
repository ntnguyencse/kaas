package kubernetesclient

import (
	"fmt"

	"k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/tools/clientcmd"
)

func GetKubernetesServerVersion(kubeConfigPath string) (*version.Info, error) {

	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		fmt.Println("Error :", err)
		return nil, err
	}

	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		fmt.Println("error in discoveryClient", err)
		return nil, err
	}

	information, err := discoveryClient.ServerVersion()
	if err != nil {
		fmt.Println("Error while fetching server version information", err)
		return nil, err
	}

	fmt.Println("Version", information)
	return information, err

}
