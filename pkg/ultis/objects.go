package ultis

import (
	"context"
	"fmt"
	"os"

	intentv1 "github.com/ntnguyencse/L-KaaS/api/v1"
	// apierrors "k8s.io/apimachinery/pkg/api/errors"
	// "k8s.io/apimachinery/pkg/runtime"
	// "k8s.io/apimachinery/pkg/types"
	// capiulti "sigs.k8s.io/cluster-api/util"
	// ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	// "sigs.k8s.io/controller-runtime/pkg/log"
)

func GetCatalogByNameAndNameSpace(ctx context.Context, c client.Client, catalogName, namespaceName string) (intentv1.ClusterCatalog, error) {
	catalog := intentv1.ClusterCatalog{}
	err := c.Get(ctx, client.ObjectKey{Name: catalogName, Namespace: namespaceName}, &catalog)
	if err != nil {
		return catalog, err
	}
	return catalog, nil
}

func GetProfileByNameAndNameSpace(ctx context.Context, c client.Client, profileName, namespaceName string) (intentv1.Blueprint, error) {
	blueprint := intentv1.Blueprint{}
	err := c.Get(ctx, client.ObjectKey{Name: profileName, Namespace: namespaceName}, &blueprint)
	if err != nil {
		return blueprint, err
	}
	return blueprint, nil
}

func GetClusterByNameAndNameSpace(ctx context.Context, c client.Client, clusterName, namespaceName string) (intentv1.Cluster, error) {
	cluster := intentv1.Cluster{}
	err := c.Get(ctx, client.ObjectKey{Name: clusterName, Namespace: namespaceName}, &cluster)
	if err != nil {
		return cluster, err
	}
	return cluster, nil
}

func SaveYamlStringToFile(fileName string, folder string, yamlContent *string) (string, error) {
	// If folder is empty, save file to current folder
	if folder == "" {
		folder, _ = os.Getwd()
	}
	filePath := folder + "/" + fileName
	// Remove file before open
	e := os.Remove(filePath)
	if e != nil {
		fmt.Println(e, "Error when remove file", filePath)
	}
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	defer file.Close()
	if err != nil {
		fmt.Println(err, "Error when open file")
		return "error Open", err
	}
	_, err = file.WriteString(*yamlContent)
	if err != nil {
		fmt.Println(err, "Error when write file")
		return "error Write", err
	}
	return filePath, nil
}
