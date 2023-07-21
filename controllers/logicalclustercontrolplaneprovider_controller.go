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
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/pkg/errors"

	"github.com/go-logr/logr"
	intentv1 "github.com/ntnguyencse/L-KaaS/api/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	capiv1alpha4 "sigs.k8s.io/cluster-api/api/v1alpha4"

	// Required for Watching
	ctrl "sigs.k8s.io/controller-runtime"       // Required for Watching
	"sigs.k8s.io/controller-runtime/pkg/client" // Required for Watching
	"sigs.k8s.io/controller-runtime/pkg/log"

	// Required for Watching
	"sigs.k8s.io/controller-runtime/pkg/reconcile" // Required for Watching
	// Required for Watching
	// intentv1 "github.com/ntnguyencse/L-KaaS/api/v1"
	// emcoctl
	cloudfile "github.com/alexflint/go-cloudfile"
	emcoctl "github.com/ntnguyencse/L-KaaS/pkg/emcoclient"

	// kubernetesclient "github.com/ntnguyencse/L-KaaS/pkg/kubernetes-client"
	grabfile "github.com/cavaliergopher/grab/v3"
	helminstaller "github.com/ntnguyencse/L-KaaS/pkg/helm"
	randomstring "github.com/ntnguyencse/L-KaaS/pkg/randstring"
)

// LogicalClusterControlPlaneProviderReconciler reconciles a LogicalClusterControlPlaneProvider object
type LogicalClusterControlPlaneProviderReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	l      logr.Logger
	s      *json.Serializer
}

const (
	// KubeconfigDataName is the key used to store a Kubeconfig in the secret's data field.
	KubeconfigDataName = "value"

	// TLSKeyDataName is the key used to store a TLS private key in the secret's data field.
	TLSKeyDataName = "tls.key"

	// TLSCrtDataName is the key used to store a TLS certificate in the secret's data field.
	TLSCrtDataName                         = "tls.crt"
	KubeConfigSecretSuffix                 = "kubeconfig"
	EMCOApplyFlag                          = "apply"
	EMCODeleteFlag                         = "delete"
	EMCOConfigPath                         = "/home/ubuntu/l-kaas/L-KaaS/pkg/emcoclient/.emco.yaml"
	EMCO_DEFAULT_CONFIG_FILE_PATH          = "/.l-kaas/config/emco/.emco.yaml"
	EMCO_RESOURCE_DEFAULT_CONFIG_FILE_PATH = "/.l-kaas/config/emco/emco-resource.yml"
)

const timeoutRetryCreateLogicalCluster = 10 * time.Minute

var (
	loggerLKP                             = ctrl.Log.WithName("L-KaaS Control Plane Provider")
	EMCOConfigFile                        = "/home/ubuntu/l-kaas/L-KaaS/pkg/emcoclient/.emco.yaml"
	prerequistiesFilePath                 = "/home/ubuntu/l-kaas/L-KaaS/templates/emco/dcm/prerequisites.yaml"
	instantiateFilePath                   = "/home/ubuntu/l-kaas/L-KaaS/templates/emco/dcm/prerequisites.yaml"
	addClusterToLogicalClusterFilePath    = "/home/ubuntu/l-kaas/L-KaaS/templates/emco/dcm/prerequisites.yaml"
	updateClusterToLogicalClusterFilePath = "/home/ubuntu/l-kaas/L-KaaS/templates/emco/dcm/prerequisites.yaml"
	// For Register Logical Cloud
	// Prerequisite Value => Prerequisite LC => Instatiate LC => Update LC
	EMCOConfigFilePath                    string
	EMCOResourceConfigFilePath            string
	prerequisitesTemplateUrl              = "https://raw.githubusercontent.com/ntnguyencse/L-KaaS/main/templates/emco/dcm/prerequisites.yaml"
	registrationLogicalClusterTemplateUrl = "https://raw.githubusercontent.com/ntnguyencse/L-KaaS/main/templates/emco/dcm/1stcluster.yaml"
	instantiateLogicalClusterTemplateUrl  = "https://raw.githubusercontent.com/ntnguyencse/L-KaaS/main/templates/emco/dcm/instantiate-lc.yaml"
	prerequisitesValuesTemplateUrl        = "https://raw.githubusercontent.com/ntnguyencse/L-KaaS/dev/templates/emco/dcm/values/prerequisites-values.yaml"
	UpdateLogicalClusterTemplateUrl       = "https://raw.githubusercontent.com/ntnguyencse/L-KaaS/main/templates/emco/dcm/update-lc.yaml"
	// For Install Flannel CNI
	EmptyProfileForCompositeAppURL = "https://raw.githubusercontent.com/ntnguyencse/L-KaaS/main/templates/emco/vfw/vfw/profile.tar.gz"
	CompositeAppTemplateURL        = "https://raw.githubusercontent.com/ntnguyencse/L-KaaS/main/templates/emco/vfw/vfw/test-vfw.yaml"
	PrerequisiteCompositeAppURL    = "https://raw.githubusercontent.com/ntnguyencse/L-KaaS/main/templates/emco/vfw/vfw/prerequisites.yaml"
)

//+kubebuilder:rbac:groups=intent.automation.dcn.ssu.ac.kr,resources=logicalclustercontrolplaneproviders,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=intent.automation.dcn.ssu.ac.kr,resources=logicalclustercontrolplaneproviders/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=intent.automation.dcn.ssu.ac.kr,resources=logicalclustercontrolplaneproviders/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the LogicalClusterControlPlaneProvider object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *LogicalClusterControlPlaneProviderReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	EMCOConfigFilePath := os.Getenv("OPENSTACKCONFIGFILEPATH")
	if len(EMCOConfigFilePath) < 1 {
		EMCOConfigFilePath = EMCO_DEFAULT_CONFIG_FILE_PATH
	}
	loggerLKP.Info("Print EMCOConfigFilePath", "Value", EMCOConfigFilePath)
	EMCOResourceConfigFilePath := os.Getenv("EMCORESOURCEFILEPATH")
	if len(EMCOResourceConfigFilePath) < 1 {
		EMCOResourceConfigFilePath = EMCO_RESOURCE_DEFAULT_CONFIG_FILE_PATH
	}
	// EMCOResourceConfig := config.LoadEMCOResourceConfiguration(EMCOResourceConfigFilePath)
	// loggerLKP.Info("Print EMCOResourceConfig", "Value", EMCOResourceConfig)
	// prerequisitesTemplateUrl = EMCOResourceConfig.PrerequisiteTemplateUrl
	// registrationLogicalClusterTemplateUrl = EMCOResourceConfig.RegistrationTemplateUrl
	// instantiateLogicalClusterTemplateUrl = EMCOResourceConfig.InstantiateTemplateUrl
	// prerequisitesValuesTemplateUrl = EMCOResourceConfig.PrerequisiteValuesTemplateUrl
	// UpdateLogicalClusterTemplateUrl = EMCOResourceConfig.UpdateTemplateUrl
	_ = log.FromContext(ctx)

	// TODO(user): your logic here
	// Get/ Fetch Cluster Instance Logical cluster
	CAPOClusters := &capiv1alpha4.Cluster{}
	if err := r.Client.Get(ctx, req.NamespacedName, CAPOClusters); err != nil {
		if apierrors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return ctrl.Result{}, nil
		}

		// Error reading the object - requeue the request.
		return ctrl.Result{}, err
	}
	// Print a CAPO CLuster
	// loggerLKP.Info("Print CAPO Cluster", "CAPO:", CAPOClusters)
	// Get status of CAPO Cluster
	CAPOStatus := CAPOClusters.DeepCopy().Status
	// Print CAPO Status
	loggerLKP.Info("Print CAPO STATUS", "CAPO:", CAPOStatus)

	// If cluster does not contain any ref, it's not belong to any logical cluster
	if len(CAPOClusters.ObjectMeta.OwnerReferences) > 0 {
		// Get Logical Cluister form OwnerRef of CAPO Cluster
		ownerRef := CAPOClusters.ObjectMeta.OwnerReferences[0]
		loggerLKP.Info("Print OwnerReferences", "OwnerReferences:", CAPOClusters.ObjectMeta.OwnerReferences)
		ownerLCluster, err := r.GetClusterOwnerObject(ctx, req, &ownerRef)
		if err != nil {
			loggerLKP.Error(err, "Error when get cluster owner's object")
			return ctrl.Result{}, nil
		}
		// Print Owner's CAPO Cluster
		loggerLKP.Info("Print Owner's CAPO cluster", "Owner", ownerLCluster.Name)
		// Get Logical Cluster Owner CAPO Cluster
		RefOfLogicalCluster := ownerLCluster.ObjectMeta.OwnerReferences[0]
		logicalCluster, err := r.GetLogicalClusterOwnerObject(ctx, req, &RefOfLogicalCluster)
		if err != nil {
			loggerLKP.Error(err, "Error when get cluster owner's object")
			return ctrl.Result{}, nil
		}
		// Print Owner's CAPO Cluster
		loggerLKP.Info("Print Owner's CAPO cluster", "Owner", logicalCluster.Name)

		//
		// ------ CHECK STATUS OF CAPO CLUSTER------------------//
		//
		// If Ready status of logical cluster and at least have one status member
		//
		// Base on Capo Cluster Status, Register logical cluster to EMCO
		// Compare Cluster Member status with Each Member Cluster Status
		//
		lClusterMembersStatus := logicalCluster.Status.ClusterMemberStates
		// Find the record of current Cluster inside status:
		idx, _, err := FindMemberStatusCorresspondToClusterName(&lClusterMembersStatus, CAPOClusters.Name)
		// TODO: How to check status of current CLuster
		lenOfLogicalClusterStatus := len(lClusterMembersStatus)
		lenOfLogicalClusterMember := len(logicalCluster.Spec.Clusters)
		capoStatus := CAPOClusters.Status
		if capoStatus.Phase == string(capiv1alpha4.ClusterPhaseProvisioned) {
			ownerLCluster.Status.Ready = true
			ownerLCluster.Status.Phase = intentv1.ConditionType(capoStatus.Phase)

		} else {
			ownerLCluster.Status.Ready = false
			ownerLCluster.Status.Phase = intentv1.ConditionType(capoStatus.Phase)
		}
		if len(capoStatus.Conditions) > 0 {
			lastCondition := capoStatus.Conditions[len(capoStatus.Conditions)-1]
			ownerLCluster.Status.FailureMessage = lastCondition.Message
			ownerLCluster.Status.FailureReason = lastCondition.Reason
			ownerLCluster.Status.Status = string(lastCondition.Status)

		}
		if lenOfLogicalClusterStatus < lenOfLogicalClusterMember {
			if err != nil {
				// Add cluster status to Logical Cluster status
				// Update status of cluster
				loggerLKP.Info("Print ownerLCluster", "ownerLCluster", ownerLCluster)
				loggerLKP.Info("Print capoStatus", "capoStatus", capoStatus)
				ownerLCluster.Status.Phase = intentv1.ConditionType(capoStatus.Phase)
				// if string(ownerLCluster.Status.Phase) != string(capiv1alpha4.ClusterPhaseProvisioned) {
				// 	if capoStatus.FailureMessage != nil && capoStatus.FailureReason != nil {
				// 		ownerLCluster.Status.FailureMessage = *capoStatus.FailureMessage
				// 		ownerLCluster.Status.FailureReason = string(*capoStatus.FailureReason)
				// 	}
				// } else {
				if len(capoStatus.Conditions) > 0 {
					ownerLCluster.Status.FailureMessage = capoStatus.Conditions[0].Message
					ownerLCluster.Status.FailureReason = capoStatus.Conditions[0].Reason
					ownerLCluster.Status.Status = string(capoStatus.Conditions[0].Status)

				}

				// }
				// Update status of L-kaas logical cluster
				memberState := intentv1.ClusterMemberStatus{
					ClusterName:    ownerLCluster.Name,
					Ready:          ownerLCluster.Status.Ready,
					FailureMessage: ownerLCluster.Status.FailureMessage,
					FailureReason:  ownerLCluster.Status.FailureReason,
					// Registration:   ownerLCluster.Status.Registration,
				}

				logicalCluster.Status.ClusterMemberStates = append(logicalCluster.Status.ClusterMemberStates, memberState)
			}

		} else {
			// Update Logical Cluster status
			if err != nil && idx != -1 {
				logicalCluster.Status.ClusterMemberStates[idx].Ready = ownerLCluster.Status.Ready
				logicalCluster.Status.ClusterMemberStates[idx].FailureMessage = ownerLCluster.Status.FailureMessage
				logicalCluster.Status.ClusterMemberStates[idx].FailureReason = ownerLCluster.Status.FailureReason
				logicalCluster.Status.ClusterMemberStates[idx].Registration = ownerLCluster.Status.Registration
			}
		}
		loggerLKP.Info("Print ownerLCluster STATUS:", "ownerLCluster", ownerLCluster.Status)
		if capoStatus.Phase == string(capiv1alpha4.ClusterPhaseProvisioned) {
			// Register Logical Cluster if at least one cluster turn "Ready" and not yet registration
			// if len(logicalCluster.Status.ClusterMemberStates) == 1 {
			// if !logicalCluster.Status.ClusterMemberStates[0].Registration {
			// var err error
			// TODO Create EMCO Cluster Provider
			folderCAPOCluster := "/tmp/" + logicalCluster.Name + "/"
			// TODO Get kubeconfig of Cluster
			// Get KubeCOnfig
			kubeconfig, err := r.getKubeConfigCluster(ctx, CAPOClusters.Name, CAPOClusters.Namespace)
			if err != nil {
				loggerLKP.Error(err, "Error when get Kubeconfig: "+CAPOClusters.Name)
			} else {
				kubePath, err := emcoctl.SaveValueFile(Name(CAPOClusters.Name, KubeConfigSecretSuffix+".kubeconfig"), folderCAPOCluster, &kubeconfig)
				prereString, err := CreateLogicalClusterPrerequisitesValueContent(&logicalCluster, CAPOClusters, kubePath)
				// TODO Add Cluster to EMCO Logical CLuster
				prereValueFilePath, err := emcoctl.SaveValueFile("prerequisitesValues.yaml", folderCAPOCluster, &prereString)
				prereTemplateFileContent, err := GetTemplateFile(prerequisitesTemplateUrl)
				if err != nil {
					loggerLKP.Error(err, "Error when get Prerequisite Template file: "+prerequisitesTemplateUrl)
				}
				prereTemplateFilePath, err := emcoctl.SaveValueFile("prequisitesTemplate.yaml", folderCAPOCluster, &prereTemplateFileContent)
				//
				// TODO Create Logical Cluster in EMCO
				//
				PrerequisitesLogicalCluster(EMCOApplyFlag, EMCOConfigFilePath, prereTemplateFilePath, prereValueFilePath)
				//
				// TODO Add Cluster to Logical CLuster
				//
				AddClusterTemplateFileContent, err := GetTemplateFile(registrationLogicalClusterTemplateUrl)
				if err != nil {
					loggerLKP.Error(err, "Error when get Add Cluster Template file: "+registrationLogicalClusterTemplateUrl)
				}
				AddClusterTemplateFilePath, err := emcoctl.SaveValueFile("AddClusterTemplate.yaml", folderCAPOCluster, &AddClusterTemplateFileContent)
				AddClusterToLogicalCluster(EMCOApplyFlag, EMCOConfigFilePath, AddClusterTemplateFilePath, prereValueFilePath)
				//
				// TODO: Instantiate Logical Cluster
				//
				InstantiateLogicalClusterTemplateFileContent, err := GetTemplateFile(registrationLogicalClusterTemplateUrl)
				if err != nil {
					loggerLKP.Error(err, "Error when get Instantiate Template file: "+registrationLogicalClusterTemplateUrl)
				}
				InstantiateLogicalClusterTemplateFilePath, err := emcoctl.SaveValueFile("InstantiateLogicalClusterTemplate.yaml", folderCAPOCluster, &InstantiateLogicalClusterTemplateFileContent)
				// emcoctl --config emco-cfg.yaml apply -f instantiate-lc.yaml -v values.yaml
				InstansiateLogicalCluster(EMCOApplyFlag, EMCOConfigFilePath, InstantiateLogicalClusterTemplateFilePath, prereValueFilePath)
				// }
				// } else {
				// var err error
				// // TODO Create EMCO Cluster Provider
				// folderCAPOCluster := "/tmp/" + logicalCluster.Name
				// // TODO Get kubeconfig of Cluster
				// // Get KubeCOnfig
				// kubeconfig, err := r.getKubeConfigCluster(ctx, CAPOClusters.Name, CAPOClusters.Namespace)
				// if err != nil {
				// 	loggerLKP.Error(err, "Error when get Kubeconfig (add cluster to logical cluster): "+CAPOClusters.Name)
				// }
				// kubePath, err := emcoctl.SaveValueFile(Name(CAPOClusters.Name, KubeConfigSecretSuffix+".yaml"), folderCAPOCluster, &kubeconfig)
				// // Create values file for Apply to EMCO server
				// valueFileContentString, err := CreateLogicalClusterPrerequisitesValueContent(&logicalCluster, CAPOClusters, kubePath)
				// ValuesFilePath, err := emcoctl.SaveValueFile("prerequisitesValues.yaml", folderCAPOCluster, &valueFileContentString)
				// AddClusterTemplateFileContent, err := GetTemplateFile(registrationLogicalClusterTemplateUrl)
				// if err != nil {
				// 	loggerLKP.Error(err, "Error when get Add Cluster Template file (add cluster to logical cluster)")
				// }
				AddClusterTemplateFilePath, err = emcoctl.SaveValueFile("AddClusterTemplate.yaml", folderCAPOCluster, &AddClusterTemplateFileContent)
				// emcoctl --config emco-cfg.yaml apply -f 2ndcluster.yaml -v values.yaml
				AddClusterToLogicalCluster(EMCOApplyFlag, EMCOConfigFilePath, AddClusterTemplateFilePath, prereValueFilePath)

				// TODO: Update Logical Cluster
				UpdateLogicalClusterTemplateFileContent, err := GetTemplateFile(UpdateLogicalClusterTemplateUrl)
				if err != nil {
					loggerLKP.Error(err, "Error when get Add Cluster Template file (add cluster to logical cluster): "+UpdateLogicalClusterTemplateUrl)
				}
				UpdateLogicalClusterTemplateFilePath, err := emcoctl.SaveValueFile("AddClusterTemplate.yaml", folderCAPOCluster, &UpdateLogicalClusterTemplateFileContent)
				// TODO: Update Logical Cluster
				// emcoctl --config emco-cfg.yaml apply -f update-lc.yaml -v values.yaml
				UpdateLogicalCluster(EMCOApplyFlag, EMCOConfigFilePath, UpdateLogicalClusterTemplateFilePath, prereValueFilePath)
				// Clean UP
				// defer emcoctl.CleanUp(kubePath)
				// defer emcoctl.CleanUp(AddClusterTemplateFilePath)
				// defer emcoctl.CleanUp(InstantiateLogicalClusterTemplateFileContent)

				//************************Install Software*****************//
				// First check status install of cluster
				if ownerLCluster.Status.Ready && string(ownerLCluster.Status.Phase) == string(capiv1alpha4.ClusterPhaseProvisioned) && !ownerLCluster.Status.Registration {
					loggerLKP.Info("Start installing software")
					// Install CNI
					r.ReconcileInstallSoftware(ctx, req, kubePath, &ownerLCluster, CAPOClusters)

				}

				// Install Prometheus
				// Developing
			}

		}

		// }

		// Find status of cluster and check if registration is false
		// Only self add cluster to logical cluster
		// if len(logicalCluster.Status.ClusterMemberStates) > 1 {

		// }

		//------CACULATE THE STATUS OF LOGICAL CLUSTER----------//
		// Separate status object:
		// CAPOPhaseStatus := CAPOStatus.Phase

		// DO Update the changes to API Server
		// DO Update status L-KaaS Cluster
		errUpdate := r.Client.Status().Update(ctx, &ownerLCluster)
		if errUpdate != nil {
			loggerLKP.Error(errUpdate, "Error when update LKaaS cluster status")
			return ctrl.Result{}, errUpdate
		}
		// Do Update status L-KaaS Logical Cluster
		//
		errUpdate = r.Client.Status().Update(ctx, &logicalCluster)
		if errUpdate != nil {
			loggerLKP.Error(errUpdate, "Error when update LKaaS logical cluster status")
			return ctrl.Result{}, errUpdate
		}
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *LogicalClusterControlPlaneProviderReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		// Uncomment the following line adding a pointer to an instance of the controlled resource as an argument
		// For().
		// For(&capiv1beta1.Cluster{}).
		For(&capiv1alpha4.Cluster{}).
		// For(&corev1.
		// Watches(
		// 	&source.Kind{Type: &capiv1beta1.ClusterClass{}},
		// 	handler.EnqueueRequestsFromMapFunc(r.findObjectsForConfigMap),
		// ).
		// Watches(
		// 	&source.Kind{Type: &capiv1beta1.MachineDeployment{}},
		// 	handler.EnqueueRequestsFromMapFunc(r.findObjectsForConfigMap),
		// ).
		Complete(r)
}
func (r *LogicalClusterControlPlaneProviderReconciler) findObjectsForConfigMap(configMap client.Object) []reconcile.Request {
	return []reconcile.Request{}
}

func (r *LogicalClusterControlPlaneProviderReconciler) GetClusterOwnerObject(ctx context.Context, req ctrl.Request, ownerRef *metav1.OwnerReference) (intentv1.Cluster, error) {
	// lclusters := intentv1.LogicalClusterList{}
	lcluster := intentv1.Cluster{}
	// r.Client.List(ctx, &lclusters)
	err := r.Client.Get(ctx, client.ObjectKey{
		Name:      ownerRef.Name,
		Namespace: req.Namespace,
	}, &lcluster)

	// Check error when get logical cluster corresponding in ownerRef
	if err != nil {
		if apierrors.IsNotFound(err) {
			loggerLKP.Error(err, "Error when get Cluster in OwnerRef not Found: ")
			return lcluster, err
		} else {
			loggerLKP.Error(err, "Error when get Cluster in OwnerRef")
			return lcluster, err
		}
	}
	return lcluster, nil

}
func (r *LogicalClusterControlPlaneProviderReconciler) GetLogicalClusterOwnerObject(ctx context.Context, req ctrl.Request, ownerRef *metav1.OwnerReference) (intentv1.LogicalCluster, error) {
	// lclusters := intentv1.LogicalClusterList{}
	lcluster := intentv1.LogicalCluster{}
	// r.Client.List(ctx, &lclusters)
	err := r.Client.Get(ctx, client.ObjectKey{
		Name:      ownerRef.Name,
		Namespace: req.Namespace,
	}, &lcluster)

	// Check error when get logical cluster corresponding in ownerRef
	if err != nil {
		if apierrors.IsNotFound(err) {
			loggerLKP.Error(err, "Error when get Logical Cluster in OwnerRef not Found: ")
			return lcluster, err
		} else {
			loggerLKP.Error(err, "Error when get Logical Cluster in OwnerRef")
			return lcluster, err
		}
	}
	return lcluster, nil

}
func (r *LogicalClusterControlPlaneProviderReconciler) RegisterLogicalCLusterToEMCO(ctx context.Context, flag string, logicalCluster intentv1.LogicalCluster) error {

	// Init EMCO Client
	// Set Config File
	emcoctl.SetConfigFilePath(EMCOConfigFile)
	// Create values file
	var valueFileString string
	// Save values file
	filePath, err := emcoctl.SaveValueFile("test-values.yaml", "/home/ubuntu/l-kaas/L-KaaS/pkg/emcoclient/test", &valueFileString)
	if err != nil {
		loggerLKP.Error(err, "Error: emcoctl.SaveValueFile test-values.yaml")
	}
	defer emcoctl.CleanUp(filePath)
	// Set Arg
	return nil
}

// Prerequisites Logical Cluster
// emcoctl --config emco-cfg.yaml apply -f prerequisites.yaml -v values.yaml
// https://github.com/ntnguyencse/L-KaaS/tree/main/templates/emco/dcm
func PrerequisitesLogicalCluster(flag string, emcoConfigPath string, prerequisitePath string, valuePath string) error {
	emcoctl.SetConfigFilePath(emcoConfigPath)
	args := []string{
		"--config",
		emcoConfigPath,
		flag,
		"-f",
		prerequisitePath,
		"-v",
		valuePath,
	}
	emcoctl.SetArgs(args)
	// emcoctl.SetDebugFlags()
	return emcoctl.ExecWithError()
}

// Add cluster to a prerequisite cluster
// emcoctl --config emco-cfg.yaml apply -f 1stcluster.yaml -v values.yaml
// https://github.com/ntnguyencse/L-KaaS/tree/main/templates/emco/dcm
func AddClusterToLogicalCluster(flag string, emcoConfigPath string, AddTemplateclusterPath string, valuePath string) error {
	emcoctl.SetConfigFilePath(emcoConfigPath)
	args := []string{
		"--config",
		emcoConfigPath,
		flag,
		"-f",
		AddTemplateclusterPath,
		"-v",
		valuePath,
	}
	emcoctl.SetArgs(args)
	// emcoctl.SetDebugFlags()
	return emcoctl.ExecWithError()
}

// Instansiate Logical Cluster
// emcoctl --config emco-cfg.yaml apply -f instantiate-lc.yaml -v values.yaml
// https://github.com/ntnguyencse/L-KaaS/tree/main/templates/emco/dcm
func InstansiateLogicalCluster(flag string, emcoConfigPath string, InstansiteLogicalClusterPath string, valuePath string) error {
	emcoctl.SetConfigFilePath(emcoConfigPath)
	args := []string{
		"--config",
		emcoConfigPath,
		flag,
		"-f",
		InstansiteLogicalClusterPath,
		"-v",
		valuePath,
	}
	emcoctl.SetArgs(args)
	// emcoctl.SetDebugFlags()
	return emcoctl.ExecWithError()
}

// Update cluster after add a new cluster to EMCO
// // emcoctl --config emco-cfg.yaml apply -f update-lc.yaml -v values.yaml
// https://github.com/ntnguyencse/L-KaaS/tree/main/templates/emco/dcm
func UpdateLogicalCluster(flag string, emcoConfigPath string, UpdateLogicalClusterPath string, valuePath string) error {
	emcoctl.SetConfigFilePath(emcoConfigPath)
	args := []string{
		"--config",
		emcoConfigPath,
		flag,
		"-f",
		UpdateLogicalClusterPath,
		"-v",
		valuePath,
	}
	emcoctl.SetArgs(args)
	// emcoctl.SetDebugFlags()
	return emcoctl.ExecWithError()
}

// Get Template File
func GetTemplateFile(url string) (string, error) {
	// url = "https://raw.githubusercontent.com/ntnguyencse/L-KaaS/dev/templates/emco/dcm/values/prerequisites-values.yaml"
	r, err := cloudfile.Open(url)
	if err != nil {
		loggerLKP.Error(err, "Error read file from remote url: "+url)
		return "", err
	}

	defer r.Close()
	strBinary, err := io.ReadAll(r)
	if err != nil {
		loggerLKP.Error(err, "Error read file GetTemplateFile io.ReadAll: "+url)
		return "", err
	}
	result := string(strBinary)
	return result, nil

}

// Create Value File for Prerequisites
func CreateLogicalClusterPrerequisitesValueContent(lCluster *intentv1.LogicalCluster, capoCluster *capiv1alpha4.Cluster, kubePath string) (string, error) {
	hostAPIEndpoint := capoCluster.Spec.ControlPlaneEndpoint.Host
	valuesMap := map[string]string{
		"PROJECTNAME":       lCluster.ObjectMeta.Labels["automation.dcn.ssu.ac.kr/project"],
		"CLUSTERPROVIDER":   lCluster.ObjectMeta.Labels["tenant"],
		"CLUSTERNAME":       capoCluster.Name,
		"CLUSTERREF":        capoCluster.Name + "-ref",
		"LOGICALCLOUD":      lCluster.Name,
		"STANDARDNAMESPACE": "default",
		"HOSTIP":            hostAPIEndpoint,
		// "KUBE_PATH":         "/home/ubuntu/l-kaas/L-KaaS/templates/emco/dcm/prerequisites.yaml",
		"KUBE_PATH": kubePath,
	}
	valuestemplateString, err := GetTemplateFile(prerequisitesValuesTemplateUrl)
	if err != nil {
		loggerLKP.Error(err, "Error wwhen get remote file github"+prerequisitesValuesTemplateUrl)
		return "", err
	}
	outputStr, err := emcoctl.InterpolateValueFile(&valuestemplateString, valuesMap)
	if err != nil {
		loggerLKP.Error(err, "Error wwhen interpolate prerequisite value file")
		return "", err
	}
	loggerLKP.Info("Print value file: ", "outString", outputStr)
	return outputStr, nil
}

// Flag:
// "apply"
// "delete"
func (r *LogicalClusterControlPlaneProviderReconciler) CreateLogicalCluster(ctx context.Context, flag string, logicalCluster intentv1.LogicalCluster) error {
	// Insert config to EMCOctl
	emcoctl.SetConfigFilePath(EMCOConfigFile)
	// Get template File
	var templateString string
	var valuesMap map[string]string
	// Create Value file for Logical Cluster
	// values.yaml
	valueString, err := emcoctl.InterpolateValueFile(&templateString, valuesMap)
	if err != nil {
		loggerLKP.Error(err, "Error when interpolate Value File")
	}
	valueFilePath, err := emcoctl.SaveValueFile("values.yaml", "/tmp/"+logicalCluster.Name+"/", &valueString)
	// defer emcoctl.CleanUp(valueFilePath)
	// Apply to EMCO
	var emptyOptions []string
	err = ApplyCommand(ctx, flag, prerequistiesFilePath, valueFilePath, emptyOptions)

	return err
}

func ApplyCommand(ctx context.Context, flag string, fileApplyPath string, valueFilePath string, options []string) error {
	args := []string{
		flag,
		"-f",
		fileApplyPath,
		"-v",
		valueFilePath,
	}
	args = append(args, options...)
	emcoctl.SetArgs(args)
	emcoctl.SetDebugFlags()

	return emcoctl.ExecWithError()
}

func (r *LogicalClusterControlPlaneProviderReconciler) InstantiateLogicalCluster(ctx context.Context, flag string, logicalCluster intentv1.LogicalCluster) error {
	// Insert config to EMCOctl
	emcoctl.SetConfigFilePath(EMCOConfigFile)
	// Get template File
	var templateString string
	var valuesMap map[string]string
	// Create Value file for Logical Cluster
	// values.yaml
	valueString, err := emcoctl.InterpolateValueFile(&templateString, valuesMap)
	if err != nil {
		loggerLKP.Error(err, "Error when interpolate Value File")
	}
	valueFilePath, err := emcoctl.SaveValueFile("values.yaml", "/tmp/"+logicalCluster.Name+"/", &valueString)
	defer emcoctl.CleanUp(valueFilePath)
	// Apply to EMCO
	var emptyOptions []string
	err = ApplyCommand(ctx, flag, prerequistiesFilePath, valueFilePath, emptyOptions)

	return err
}

func (r *LogicalClusterControlPlaneProviderReconciler) AddClusterToLogicalCluster(ctx context.Context, flag string, cluster intentv1.Cluster) error {
	// Insert config to EMCOctl
	emcoctl.SetConfigFilePath(EMCOConfigFile)
	// Get template File
	var templateString string
	var valuesMap map[string]string
	// Create Value file for Logical Cluster
	// values.yaml
	valueString, err := emcoctl.InterpolateValueFile(&templateString, valuesMap)
	if err != nil {
		loggerLKP.Error(err, "Error when interpolate Value File")
	}
	valueFilePath, err := emcoctl.SaveValueFile("values.yaml", "/tmp/"+cluster.Name+"/", &valueString)
	defer emcoctl.CleanUp(valueFilePath)
	// Apply to EMCO
	var emptyOptions []string
	err = ApplyCommand(ctx, flag, prerequistiesFilePath, valueFilePath, emptyOptions)

	return err
}

func (r *LogicalClusterControlPlaneProviderReconciler) UpdateClusterToLogicalCluster(ctx context.Context, flag string, cluster intentv1.Cluster) error {
	// Insert config to EMCOctl
	emcoctl.SetConfigFilePath(EMCOConfigFile)
	// Get template File
	var templateString string
	var valuesMap map[string]string
	// Create Value file for Logical Cluster
	// values.yaml
	valueString, err := emcoctl.InterpolateValueFile(&templateString, valuesMap)
	if err != nil {
		loggerLKP.Error(err, "Error when interpolate Value File")
	}
	valueFilePath, err := emcoctl.SaveValueFile("values.yaml", "/tmp/"+cluster.Name+"-updatecluster/", &valueString)
	defer emcoctl.CleanUp(valueFilePath)
	// Apply to EMCO
	var emptyOptions []string
	err = ApplyCommand(ctx, flag, prerequistiesFilePath, valueFilePath, emptyOptions)

	return err
}

func (r *LogicalClusterControlPlaneProviderReconciler) InstantiateCompositeApplication(ctx context.Context, flag string, cluster intentv1.Cluster) error {
	// Insert config to EMCOctl
	emcoctl.SetConfigFilePath(EMCOConfigFile)
	// Get template File
	var templateString string
	var valuesMap map[string]string
	// Create Value file for Logical Cluster
	// values.yaml
	valueString, err := emcoctl.InterpolateValueFile(&templateString, valuesMap)
	if err != nil {
		loggerLKP.Error(err, "Error when interpolate Value File")
	}
	valueFilePath, err := emcoctl.SaveValueFile("values.yaml", "/tmp/"+cluster.Name+"-instantiate-compositeapp/", &valueString)
	defer emcoctl.CleanUp(valueFilePath)
	// Apply to EMCO
	var emptyOptions []string
	err = ApplyCommand(ctx, flag, prerequistiesFilePath, valueFilePath, emptyOptions)

	return err
}
func FindMemberStatusCorresspondToClusterName(memberStatus *[]intentv1.ClusterMemberStatus, clusterName string) (int, intentv1.ClusterMemberStatus, error) {
	for index, item := range *memberStatus {
		if item.ClusterName == clusterName {
			return index, item, nil
		}
	}

	return -1, intentv1.ClusterMemberStatus{}, errors.New("Could not find member status in array")
}
func (r *LogicalClusterControlPlaneProviderReconciler) getKubeConfigCluster(ctx context.Context, clusterName, nameSpace string) (string, error) {
	secret := &corev1.Secret{}
	secretKey := client.ObjectKey{
		Namespace: nameSpace,
		Name:      Name(clusterName, KubeConfigSecretSuffix),
	}
	if err := r.Client.Get(ctx, secretKey, secret); err != nil {
		return "nil", err
	}
	secretBytes, err := toKubeconfigBytes(secret)
	return string(secretBytes), err
}

func toKubeconfigBytes(out *corev1.Secret) ([]byte, error) {
	data, ok := out.Data[KubeconfigDataName]
	if !ok {
		return nil, errors.Errorf("missing key %q in secret data", KubeconfigDataName)
	}
	return data, nil
}

// Name returns the name of the secret for a cluster.
func Name(cluster string, suffix string) string {
	return fmt.Sprintf("%s-%s", cluster, suffix)
}

func (r *LogicalClusterControlPlaneProviderReconciler) CalicoInstaller(kubeConfigPath, version, podCIDR string) (*Installer, error) {
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
func (r *LogicalClusterControlPlaneProviderReconciler) FlannelInstaller(kubeConfigPath, version, podCIDR string) (*Installer, error) {
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
	flannelUrl := "https://raw.githubusercontent.com/flannel-io/flannel/{VERSION}/Documentation/kube-flannel.yml"
	cniComponent := InstallComponent{
		Name:           "flannel-cni",
		URL:            flannelUrl,
		Version:        version,
		KubeConfigPath: kubeConfigPath,
	}
	installer.AddInstallComponent(cniComponent)

	return &installer, nil
}
func (r *LogicalClusterControlPlaneProviderReconciler) FlannelEMCOInstaller(kubeConfigPath, version, podCIDR string) error {
	// Flannel Version
	// https://raw.githubusercontent.com/flannel-io/flannel/v0.22.0/Documentation/kube-flannel.yml
	// Default in Flannel: POD_CIDR="10.244.0.0/16"
	if len(podCIDR) < 4 {
		podCIDR = "10.244.0.0/16"
	}
	flannelUrl := "https://raw.githubusercontent.com/ntnguyencse/L-KaaS/main/templates/cni/flannel/flannel.tar.gz"
	// Download kubeflannel file
	resp, err := grabfile.Get("/tmp/emco-tmp/", flannelUrl)
	if err != nil {
		loggerLKP.Error(err, "Error get Flannel helm chart: FlannelEMCOInstallerFunc")
	}
	_ = resp
	return nil
}
func (r *LogicalClusterControlPlaneProviderReconciler) ReconcileInstallSoftware(ctx context.Context, req ctrl.Request, kubePath string, cluster *intentv1.Cluster, CAPICluster *capiv1alpha4.Cluster) error {
	// clusterStatus.Ready && string(clusterStatus.Phase) == string(capiv1alpha4.ClusterPhaseProvisioned) && !clusterStatus.RegistrationkubePath
	// Install CNI
	loggerLKP.Info("Begin Install CNI")
	// Get Profiles related to Clusters
	listProfiles, err := r.GetListProfile(ctx, req)
	if err != nil {
		loggerLKP.Error(err, "Error get profiles")
	}
	for _, item := range cluster.Spec.Network {
		// 1. CNI Profiles
		CNIProfileName := item.Name
		CNIProfile, err := FindProfileWithName(listProfiles, CNIProfileName)
		if err == nil {
			// Install CNI
			chartPath := CNIProfile.Spec.Values["url"]
			chartName := CNIProfileName + "-" + randomstring.String(5)
			CNINamespace := CNIProfile.Spec.Values["namespace"]
			valueFilePath := CNIProfile.Spec.Values["value"]
			// if strings.Contains(CNIProfileName, "flannel"){
			// 	chartPath = flannelTemplate
			// }

			err = helminstaller.Install(kubePath, chartName, chartPath, valueFilePath, CNINamespace)
			if err != nil {
				loggerLKP.Error(err, "Error install Network components: "+CNIProfileName)
				return err
			}
		}

	}
	loggerLKP.Info("Finish install CNI")
	loggerLKP.Info("Begin install Software")
	// 2. Software Profiles
	for _, item := range cluster.Spec.Software {
		// 1. CNI Profiles
		SoftwareProfileName := item.Name
		SoftwareProfile, err := FindProfileWithName(listProfiles, SoftwareProfileName)
		if err == nil {
			// Install CNI
			chartPath := SoftwareProfile.Spec.Values["url"]
			chartName := SoftwareProfileName + "-" + randomstring.String(5)
			CNINamespace := SoftwareProfile.Spec.Values["namespace"]
			valueFilePath := SoftwareProfile.Spec.Values["value"]
			// if strings.Contains(CNIProfileName, "flannel"){
			// 	chartPath = flannelTemplate
			// }

			err = helminstaller.Install(kubePath, chartName, chartPath, valueFilePath, CNINamespace)
			if err != nil {
				loggerLKP.Error(err, "Error install Network components: "+SoftwareProfileName)
				return err
			}
		}
	}
	loggerLKP.Info("Finish install Software")

	// Update cluster status
	cluster.Status.Registration = true
	// Clean up
	// defer os.RemoveAll("$HOME/.cache/helm/repository")
	return nil
}
func (r *LogicalClusterControlPlaneProviderReconciler) GetListProfile(ctx context.Context, req ctrl.Request) (*intentv1.ProfileList, error) {
	listProfiles := intentv1.ProfileList{}
	err := r.Client.List(ctx, &listProfiles)
	if err != nil {
		loggerCL.Error(err, "ReconcileInstallSoftware", "Error when list profiles")
		return nil, err
	}
	return &listProfiles, err
}

func FindProfileWithName(listProfiles *intentv1.ProfileList, name string) (intentv1.Profile, error) {
	if len(listProfiles.Items) < 1 {
		return intentv1.Profile{}, errors.New("List Profiles is empty")
	}
	for _, item := range listProfiles.Items {
		if item.Name == name {
			return item, nil
		}
	}
	return intentv1.Profile{}, errors.New("Can not find any profile match with name: " + name)
}
