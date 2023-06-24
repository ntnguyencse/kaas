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
	"bytes"
	"context"
	"io"
	"os"

	// jsonclassic "encoding/json"

	// config "github.com/ntnguyencse/L-KaaS/pkg/config"
	"github.com/go-logr/logr"
	goyaml "github.com/go-yaml/yaml"
	intentv1 "github.com/ntnguyencse/L-KaaS/api/v1"
	kubernetesclient "github.com/ntnguyencse/L-KaaS/pkg/kubernetes-client"
	"github.com/ntnguyencse/L-KaaS/pkg/ultis"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	capiv1alpha4 "sigs.k8s.io/cluster-api/api/v1alpha4"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ClusterReconciler reconciles a Cluster object
type ClusterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	l      logr.Logger
	s      *json.Serializer
}

var (
	loggerCL                        = ctrl.Log.WithName("Cluster Controller")
	CAPIConfigFilePath              string
	OpenStackProviderConfigFilePath string
)

const (
	ClusterFinalizer                   string = "cluster.intent.automation.dcn.ssu.ac.kr"
	NAMESPACE_DEFAULT                  string = "default"
	STATUS_GENERATED                   string = "Generated"
	OPENSTACK_DEFAULT_CONFIG_FILE_PATH string = "/.l-kaas/config/openstack/openstack-config.yml"
	// CAPIDefaultFilePath      string = "/.l-kaas/config/capi/capictl.yml"
	CNIFlannelType string = "flannel"
	CNICalicoType  string = "calico"
)

//+kubebuilder:rbac:groups=intent.automation.dcn.ssu.ac.kr,resources=clusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=intent.automation.dcn.ssu.ac.kr,resources=clusters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=intent.automation.dcn.ssu.ac.kr,resources=clusters/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Cluster object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *ClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	OpenStackProviderConfigFilePath := os.Getenv("OPENSTACKCONFIGFILEPATH")
	if len(OpenStackProviderConfigFilePath) < 1 {
		OpenStackProviderConfigFilePath = OPENSTACK_DEFAULT_CONFIG_FILE_PATH
	}
	CAPIConfigFilePath := os.Getenv("CAPICONFIGFILEPATH")
	if len(CAPIConfigFilePath) < 1 {
		CAPIConfigFilePath = DEFAULT_CAPI_CONFIG_PATH
	}
	// Load Configuration
	// configuration := config.LoadConfig(config.DEFAULT_CONFIG_PATH)

	// gitclient1, err := git.NewClient(configuration.ClusterRepo, configuration.Owner, configuration.GitHubToken, ctx)
	// if err != nil {
	// 	r.l.Error(err, "Error while create new Github client")
	// 	_ = gitclient1
	// 	return ctrl.Result{}, nil
	// }
	loggerCL.Info("Reconciling.... CLuster")
	loggerCL.Info("Start reconciling Cluster Resource")

	// CLuster Resource object get from Kubernetes API Server
	var cluster intentv1.Cluster
	err := r.Get(context.Background(), req.NamespacedName, &cluster)
	if err != nil {
		if errors.IsNotFound(err) {
			// The Cluster Resources has been deleted, so we need to delete the cluster resource description corresponding
			loggerCL.V(1).Info("The Cluster Resources has been deleted, so we need to delete the cluster resource description corresponding")
			/////
			// TO-DO: Delete the cluster resource description
			////
			////
			return ctrl.Result{}, nil
		}
		// There was an error getting the Deployment, so we'll retry later
		loggerCL.V(1).Info("There was an error getting the Deployment, so we'll retry later")
		return ctrl.Result{}, err
	}

	// Handle deletion reconciliation loop.
	if !cluster.ObjectMeta.DeletionTimestamp.IsZero() {
		return r.ReconcileDelete(ctx, &cluster)
	}
	clusterStatus := cluster.Status
	if clusterStatus.Ready && string(clusterStatus.Phase) == string(capiv1alpha4.ClusterPhaseProvisioned) && !clusterStatus.Registration {
		r.ReconcileInstallSoftware(ctx, req, &cluster)
	}
	return r.ReconcileNormal(ctx, &cluster)
}

// SetupWithManager sets up the controller with the Manager.
func (r *ClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&intentv1.Cluster{}).
		Complete(r)
}

/*
// Draft
// func test() {
// 	var deploy intentv1.Cluster
// 	gitclient1, err := git.NewClient("a", "a", "", context.TODO())
// 	contentLastAppliedConfigurationBytes := deploy.Annotations["kubectl.kubernetes.io/last-applied-configuration"]
// 	loggerCL.V(1).Info("Print Deploy:", "Cluster", deploy)
// 	var content intentv1.Cluster
// 	err = jsonclassic.Unmarshal([]byte(contentLastAppliedConfigurationBytes), &content)

// 	if err != nil {
// 		loggerCL.V(0).Info("Error when unmarshal cluster resource")
// 	} else {
// 		if len(deploy.ManagedFields) > 0 {
// 			loggerCL.V(0).Info("ManagedField", deploy.ManagedFields[0].Manager, deploy.ManagedFields[0].Operation, deploy.ManagedFields[0].Time.String(), deploy.ManagedFields[0].APIVersion)
// 		}
// 	}
// 	// Check status and sha, version, revision
// 	loggerCL.Info("Print status of object:", "Status", deploy.Status.Status, deploy.Status.SHA, deploy.Status.Sync)

// 	///////////////////////////================
// 	loggerCL.Info("Print Before status of object:", "Status", deploy.Status.Status, deploy.Status.SHA, deploy.Status.Sync)
// 	if deploy.Status.Status != "Reconciling" {
// 		deploy.Status.Status = "Reconciling"
// 		deploy.Status.Sync = "Not Synced"
// 		deploy.Status.SHA = "SHA Tests"
// 		// Update status of Kubernetes objects
// 		// err = r.Client.Status().Update(ctx, &deploy)
// 		if err != nil {
// 			loggerCL.V(0).Error(err, "Error while update status object Not Synced")
// 			// return ctrl.Result{}, err
// 		}
// 		var deploy1 intentv1.Cluster
// 		// err = r.Get(context.Background(), req.NamespacedName, &deploy1)
// 		if err != nil {
// 			if errors.IsNotFound(err) {
// 				// The Deployment has been deleted, so we don't need to do anything
// 				loggerCL.V(1).Info("The Deployment has been deleted, so we don't need to do anything")
// 				// return ctrl.Result{}, nil
// 			}
// 			// There was an error getting the Deployment, so we'll retry later
// 			loggerCL.V(1).Info("There was an error getting the Deployment, so we'll retry later")
// 			// return ctrl.Result{}, err
// 		}
// 		loggerCL.Info("Print status of object:", "Status", deploy1.Status.Status, deploy1.Status.SHA, deploy1.Status.Sync)

// 	}

// 	if false {
// 		// if deploy.CreationTimestamp == *deploy.ManagedFields[0].Time {
// 		if len(deploy.ObjectMeta.Labels["SHA"]) < 2 {
// 			isFileNotExist, err := gitclient1.IsFileNotExist(deploy.Name+".yaml", "test/")
// 			if err != nil {
// 				loggerCL.Error(err, "Error while check file existing..")
// 			} else {
// 				content1, _ := jsonclassic.MarshalIndent(content, " ", "    ")
// 				var sha string
// 				if !isFileNotExist {
// 					resp, err := gitclient1.UpdateFile(deploy.Name+".yaml", "test/", content1)
// 					if err != nil {
// 						loggerCL.V(0).Error(err, "Error whilde reconciling cluster resource. Update file")
// 					} else {
// 						sha = *resp.SHA
// 						loggerCL.V(0).Info("Updated SHA", "SHA", sha)
// 					}

// 				} else {
// 					resp, err := gitclient1.CommitNewFile(deploy.Name+".yaml", "main", "test/", content1)
// 					if err != nil {
// 						loggerCL.V(0).Error(err, "Error whilde reconciling cluster resource. Commit file")
// 					} else {
// 						sha = *resp.SHA
// 						loggerCL.V(0).Info("Updated SHA", "SHA", sha)
// 					}
// 				}

// 				deploy.ObjectMeta.Labels["SHA"] = sha
// 				deploy.ObjectMeta.Labels["Sync"] = "Synced"
// 				// deploy.Status.Status
// 				// err := r.Update(ctx, &deploy)
// 				if err != nil {
// 					loggerCL.V(0).Error(err, "Error while update status object Synced")
// 					// return ctrl.Result{}, err
// 				}
// 			}
// 		} else {
// 			// deploy.ObjectMeta.Labels["Sync"] = "UnSynced"
// 			deploy.Status.Sync = "Not Synced"
// 			// err := r.Update(ctx, &deploy)
// 			if err != nil {
// 				loggerCL.V(0).Error(err, "Error while update status object Not Synced")
// 				// return ctrl.Result{}, err
// 			}
// 			var deploy1 intentv1.Cluster
// 			// err = r.Get(context.Background(), req.NamespacedName, &deploy1)
// 			if err != nil {
// 				if errors.IsNotFound(err) {
// 					// The Deployment has been deleted, so we don't need to do anything
// 					loggerCL.V(1).Info("The Deployment has been deleted, so we don't need to do anything")
// 					// return ctrl.Result{}, nil
// 				}
// 				// There was an error getting the Deployment, so we'll retry later
// 				loggerCL.V(1).Info("There was an error getting the Deployment, so we'll retry later")
// 				// return ctrl.Result{}, err
// 			}
// 			loggerCL.Info("Print status of object:", "Status", deploy1.Status.Status, deploy1.Status.SHA, deploy1.Status.Sync)
// 		}
// 	}
// 	// r.Status().Update()
// }
// // Check the revision of package (equal the generation of kubernetes objects)
// 	// The generation of k8s object change means the metadata or the spec of k8s object change => detect change using the change of generation number of k8s object.
// 	// Store the generation in status of object
// 	var oldGeneration int64
// 	if deploy.Status.Revision < 1 {
// 		oldGeneration = 0
// 	} else {
// 		oldGeneration = deploy.Status.Revision
// 	}
// 	// Checking changes of cluster contents
// 	if oldGeneration != deploy.Generation {
// 		// Sync the change of cluster resource to github
// 		// TO-DO: Sync the change to github
// 		// Check is the new object or existing object
// 		// Each file will store in separate folder name same name as namespace
// 		fileName := deploy.Name + ".yaml"
// 		fileFolderName := deploy.Namespace
// 		if len(deploy.Namespace) < 1 {
// 			fileFolderName = NAMESPACE_DEFAULT
// 		}
// 		// Get current content of kubernetes object
// 		// That content was inside last-applied-configuration
// 		contentLastAppliedConfigurationBytes := deploy.Annotations["kubectl.kubernetes.io/last-applied-configuration"]
// 		var content intentv1.Cluster
// 		err = jsonclassic.Unmarshal([]byte(contentLastAppliedConfigurationBytes), &content)
// 		// Convert to pretty-print yaml file
// 		content1, _ := jsonclassic.MarshalIndent(content, " ", "    ")
// 		isFileExist, err := gitclient1.IsFileNotExist(fileName, fileFolderName)
// 		var shaFile string
// 		if err == nil {
// 			if isFileExist {

// 				// Temporaty disable git features
// 				// result, err := gitclient1.CommitNewFile(fileName, git.GIT_MAINBRANCH, fileFolderName, content1)
// 				// if err != nil {
// 				// 	loggerCL.V(1).Error(err, "Error while commit a new file", "File name and folder name", fileName, fileFolderName)
// 				// } else {
// 				// 	shaFile = *result.SHA
// 				// }
// 				_ = content1 // Remember to delete
// 			} else {
// 				// Temporaty disable git features
// 				// result, err := gitclient1.UpdateFile(fileName, fileFolderName, content1)
// 				// if err != nil {
// 				// 	loggerCL.V(1).Error(err, "Error while Update a new file", "File name and folder name", fileName, fileFolderName)
// 				// } else {
// 				// 	shaFile = *result.SHA
// 				// }
// 				_ = shaFile
// 			}
// 			//
// 			// deploy.Status.SHA = shaFile
// 			// deploy.Status.Sync = git.SYNCED_STATUS
// 			// deploy.Status.Repo = gitclient1.GetOwner() + "/" + gitclient1.GetRepoName()
// 		} else {
// 			loggerCL.V(1).Error(err, "Error while check file not existing", "File Name and folder:", fileFolderName+"/"+fileName)
// 			deploy.Status.Sync = git.NOT_SYNC_STATUS
// 		}

// 		//----------------------------------------------//-------------------------------------------------------------------//

// 		// Transform Cluster Resource to Cluster Description
// 		// Get list blueprints
// 		var listBP intentv1.ProfileList
// 		err = r.Client.List(ctx, &listBP)
// 		if err != nil {
// 			loggerCL.Error(err, "Error while list blueprints")
// 		}
// 		// Transform cluster Resource to Cluster Description Resources and apply to Kubernetes API Server

// 		clusterDes, err := r.TransformClusterToClusterDescription(ctx, deploy, listBP.Items)
// 		if err != nil {
// 			loggerCL.Error(err, "Transform Cluster Failed")
// 		} else {
// 			loggerCL.Info("Applying cluster description")
// 			// Apply the Cluster Resource Description to Kubernetes API Server
// 			err := r.Client.Create(ctx, &clusterDes)
// 			// r.Client.
// 			if err != nil {
// 				loggerCL.Error(err, "Error while applying cluster resource")
// 			} else {
// 				loggerCL.Info("Applying successful")
// 			}
// 		}
// 		loggerCL.Info("Cluster Description:", "value", clusterDes)

// 		// Update status of kubernetes object
// 		// Update Revision
// 		deploy.Status.Revision = deploy.Generation
// 		deploy.Status.Status = STATUS_GENERATED

// 		// Update the changes to Kubernetes Server
// 		r.Client.Update(ctx, &deploy)
// 		return ctrl.Result{}, nil
*/
func (r *ClusterReconciler) ReconcileNormal(ctx context.Context, cluster *intentv1.Cluster) (ctrl.Result, error) {
	// TODO:
	// Reconcile Normal
	loggerCL.Info("Print Cluster: ", cluster.Name, cluster)
	r.GetOrCreateCluster(ctx, "default", cluster)
	// Do not forget defer func (){}()
	return ctrl.Result{}, nil
}

func (r *ClusterReconciler) ReconcileDelete(ctx context.Context, cluster *intentv1.Cluster) (ctrl.Result, error) {
	// TODO
	// Reconcile Delete
	// Do not forget defer func

	return ctrl.Result{}, nil
}

// Update/ Reconcile Phase of cluster
func (r *ClusterReconciler) ReconcileClusterPhase(ctx context.Context, cluster *intentv1.Cluster) {

}

func (r *ClusterReconciler) GetOrCreateCluster(ctx context.Context, clusterNameSpace string, cluster *intentv1.Cluster) (intentv1.ClusterDescription, error) {
	// Get list Profiles
	listProfiles := intentv1.ProfileList{}
	err := r.Client.List(ctx, &listProfiles)
	if err != nil {
		loggerCL.Error(err, "GetOrCreateCluster", "Error when listt profiles")
	}
	clusterDescriptiton, err := r.TransformClusterToClusterDescription(ctx, *cluster, clusterNameSpace, listProfiles.Items)
	loggerCL.Info("Print CLuster Descrition", clusterDescriptiton.Name, clusterDescriptiton)
	// Get Provider COnfig
	OpenStackProviderConfig := GetConfigForOpenStack()
	// Get Credentials of Provider
	configs, err := getCredentialsForOpenStackProvider(CAPIConfigFilePath)
	if err != nil {
		loggerCL.Error(err, "getCredentialsForOpenStackProvider")
	}
	// Add config from CLuster Description to configs variables
	// TODO: Change this function or redesign Profile Specs to easy to transform
	// Current only use mapping variables
	configs = AddInfraConfigsFromClusterDescription(&clusterDescriptiton, configs)

	ownerRefs := map[string]string{
		"CLUSTER_OWNER_API_VERSION": cluster.DeepCopy().APIVersion,
		"CLUSTER_OWNER_KIND":        cluster.DeepCopy().Kind,
		"CLUSTER_OWNER_NAME":        cluster.DeepCopy().Name,
		"CLUSTER_OWNER_UID":         string(cluster.ObjectMeta.DeepCopy().UID),
	}
	configs = AddToConfigs(configs, ownerRefs)

	clusterStr, _ := TranslateFromClusterDescritionToCAPI(&clusterDescriptiton, OpenStackProviderConfig, configs)

	r.ApplyCAPIResourceToKubernertes(clusterDescriptiton.Name, clusterStr)
	return clusterDescriptiton, err
}
func GetOpenStackConfigPath(systemConfigPath string) (string, error) {
	url := "Underconstruction"
	return url, nil
}

// Add config from Infra Part of cluster Description
// TODO: Make it general for most of provider
// Current only support for OpenStack
func AddInfraConfigsFromClusterDescription(clusterDes *intentv1.ClusterDescription, configs map[string]string) map[string]string {
	// Take the first Infra record
	// Fix this later
	infraDes := clusterDes.Spec.Infrastructure[0]
	// Mapping
	configs["WORKER_MACHINE_COUNT"] = infraDes.Spec["workerMachineCount"]
	configs["KUBERNETES_VERSION"] = infraDes.Spec["kubernetesVersion"]
	configs["CONTROL_PLANE_MACHINE_COUNT"] = infraDes.Spec["controlPlaneMachineCount"]
	configs["OPENSTACK_CONTROL_PLANE_MACHINE_FLAVOR"] = infraDes.Spec["controlplaneFlavor"]
	configs["OPENSTACK_NODE_MACHINE_FLAVOR"] = infraDes.Spec["workerFlavor"]
	configs["CAPI_TEMPLATE_URL"] = infraDes.Spec["url"] + infraDes.Spec["filename"]
	// Take the first Network record
	networkDes := clusterDes.Spec.Network[0]
	configs["POD_CIDR"] = networkDes.Spec["podCIDR"]
	configs["CNI_NAME"] = networkDes.Spec["cni"]
	configs["SERVICE_CIDR"] = networkDes.Spec["serviceCIDR"]
	return configs

}

func (r *ClusterReconciler) ApplyCAPIResourceToKubernertes(clusterName, CAPIRes string) error {
	// listCAPIRes, err := SplitYAML([]byte(CAPIRes))
	// if err != nil {
	// 	loggerCL.Error(err, "Error convert CAPI Resources")
	// 	return err
	// }
	filepath, err := ultis.SaveYamlStringToFile(clusterName+".yaml", "/home/ubuntu/l-kaas/L-KaaS/test", &CAPIRes)
	if err != nil {
		loggerCL.Error(err, "Err Save file")
		return err
	}
	defer CleanUpCAPIResource(filepath)
	kubernetesclient.KubectlApplyYamlFile(filepath)
	// for i, capi := range listCAPIRes {
	// 	loggerCL.Info("CAPIRes:", "String: ", string(capi))
	// 	content := string(capi)
	// 	filepath, err := ultis.SaveYamlStringToFile(clusterName+strconv.Itoa(i)+".yaml", "/home/ubuntu/l-kaas/L-KaaS/test", &content)
	// 	if err != nil {
	// 		loggerCL.Error(err, "Error saving file")

	// 	}
	// 	loggerCL.Info("Saving file...", filepath, "Finished")
	// }
	return nil
}
func CleanUpCAPIResource(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		loggerCL.Error(err, "Error clean up file: "+filePath)
		return err
	}
	loggerCL.Info("Cleanup file: " + filePath)
	return nil
}
func SplitYAML(resources []byte) ([][]byte, error) {

	dec := goyaml.NewDecoder(bytes.NewReader(resources))

	var res [][]byte
	for {
		var value interface{}
		err := dec.Decode(&value)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		valueBytes, err := goyaml.Marshal(value)
		if err != nil {
			return nil, err
		}
		res = append(res, valueBytes)
	}
	return res, nil
}
func ChooseKindCAPIResource(resource string) {

}
func (r *ClusterReconciler) getKubeConfigCluster(ctx context.Context, clusterName, nameSpace string) (string, error) {
	secret := &corev1.Secret{}
	secretKey := client.ObjectKey{
		Namespace: nameSpace,
		Name:      Name(clusterName, KubeConfigSecretSuffix),
	}
	if err := r.Get(ctx, secretKey, secret); err != nil {
		return "nil", err
	}
	secretBytes, err := toKubeconfigBytes(secret)
	return string(secretBytes), err
}

// Install Calico CNI
func (r *ClusterReconciler) CalicoInstaller(kubeConfigPath, version, podCIDR string) (*Installer, error) {
	// Calico Version
	// https://raw.githubusercontent.com/projectcalico/calico/v3.26.0/manifests/tigera-operator.yaml
	// Default Calico version v3.26.0
	if len(version) < 1 {
		version = "v3.26.0"
	}
	if len(podCIDR) < 1 {
		podCIDR = "192.168.0.0/16"
	}

	installer := SetUpInstaller(r.Client)
	operatorURL := "https://raw.githubusercontent.com/projectcalico/calico/{VERSION}/manifests/tigera-operator.yaml"
	operatorVersion := version // "v3.26.0"
	operatorComponent := InstallComponent{
		Name:           "calico-operator",
		URL:            operatorURL,
		Version:        operatorVersion,
		KubeConfigPath: kubeConfigPath,
	}
	calicoUrl := "https://raw.githubusercontent.com/projectcalico/calico/{VERSION}/manifests/custom-resources.yaml"
	calicoVersion := operatorVersion

	cniComponent := InstallComponent{
		Name:           "calico-cni",
		URL:            calicoUrl,
		Version:        calicoVersion,
		KubeConfigPath: kubeConfigPath,
	}
	installer.AddInstallComponent(operatorComponent)
	installer.AddInstallComponent(cniComponent)

	return &installer, nil

}

// Install Flannel CNI
func (r *ClusterReconciler) FlannelInstaller(kubeConfigPath, version, podCIDR string) (*Installer, error) {
	// Flannel Version
	// https://raw.githubusercontent.com/flannel-io/flannel/v0.22.0/Documentation/kube-flannel.yml
	// Default in Flannel: POD_CIDR="10.244.0.0/16"
	if len(podCIDR) < 4 {
		podCIDR = "10.244.0.0/16"
	}
	if len(version) < 1 {
		version = "v0.22.0"
	}
	installer := SetUpInstaller(r.Client)
	calicoUrl := "https://raw.githubusercontent.com/flannel-io/flannel/{VERSION}/Documentation/kube-flannel.yml"
	cniComponent := InstallComponent{
		Name:           "flannel-cni",
		URL:            calicoUrl,
		Version:        version,
		KubeConfigPath: kubeConfigPath,
	}
	installer.AddInstallComponent(cniComponent)

	return &installer, nil
}

// Reconcile when cluster provisioned and setting up software components
func (r *ClusterReconciler) ReconcileInstallSoftware(ctx context.Context, request ctrl.Request, cluster *intentv1.Cluster) error {

	// Get CApi CLUSTER
	var CAPICluster capiv1alpha4.Cluster
	err := r.Get(ctx, client.ObjectKey{
		Namespace: request.Namespace,
		Name:      cluster.Name,
	}, &CAPICluster)
	if err != nil {
		loggerCL.Error(err, "Error when get CAPI Cluster, ReconcileInstallSoftware")
		return err
	}

	networkConfig := cluster.Spec.Network
	// Get Blueprint of Network
	var networkProfile intentv1.Profile
	err = r.Get(ctx, client.ObjectKey{
		Namespace: request.Namespace,
		Name:      networkConfig[0].Name,
	}, &networkProfile)
	if err != nil {
		loggerCL.Error(err, "Error when get Network Profile , ReconcileInstallSoftware")
		return err
	}
	podCIDR := CAPICluster.Spec.ClusterNetwork.Pods.CIDRBlocks[0]

	// Install CNI
	// Currently support 2 CNI: Flannel, Calico
	// Cilium is not support install over crds
	CNIType := networkProfile.Spec.Values["cni"]
	CNIVersion := networkProfile.Spec.Values["version"]
	if CNIType == CNIFlannelType {
		// Get kubeconfig
		kubeconfig, err := r.getKubeConfigCluster(ctx, CAPICluster.Name, CAPICluster.Namespace)
		if err != nil {
			loggerCL.Error(err, "Error when get Kubeconfig", "CNICalicoType", CNICalicoType, "Cluster Name: ", CAPICluster.Name)
			return err
		}
		folder := "/tmp/" + CAPICluster.Name + "kube"
		fileName := CAPICluster.Name + "-" + KubeConfigSecretSuffix
		pathKubeconfig, err := ultis.SaveYamlStringToFile(fileName, folder, &kubeconfig)
		if err != nil {
			loggerCL.Error(err, "Error when get Kubeconfig", "CNICalicoType", CNICalicoType, "Cluster Name: ", CAPICluster.Name)
			return err
		}
		defer CleanUpCAPIResource(pathKubeconfig)

		r.FlannelInstaller(pathKubeconfig, CNIVersion, podCIDR)
	} else if CNIType == CNICalicoType {
		// Get kubeconfig
		kubeconfig, err := r.getKubeConfigCluster(ctx, CAPICluster.Name, CAPICluster.Namespace)
		if err != nil {
			loggerCL.Error(err, "Error when get Kubeconfig", "CNICalicoType", CNICalicoType, "Cluster Name: ", CAPICluster.Name)
			return err
		}
		folder := "/tmp/" + CAPICluster.Name + "kube"
		fileName := CAPICluster.Name + "-" + KubeConfigSecretSuffix
		pathKubeconfig, err := ultis.SaveYamlStringToFile(fileName, folder, &kubeconfig)
		if err != nil {
			loggerCL.Error(err, "Error when get Kubeconfig", "CNICalicoType", CNICalicoType, "Cluster Name: ", CAPICluster.Name)
			return err
		}
		defer CleanUpCAPIResource(pathKubeconfig)

		installer, _ := r.CalicoInstaller(pathKubeconfig, CNIVersion, podCIDR)
		installer.Install(CAPICluster.Name)
	}
	return nil
}
