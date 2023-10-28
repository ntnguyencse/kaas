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
	"time"

	intentv1 "github.com/ntnguyencse/L-KaaS/api/v1"
	emcoctl "github.com/ntnguyencse/L-KaaS/pkg/emcoclient"
	helminstaller "github.com/ntnguyencse/L-KaaS/pkg/helm"
	kubernetesclient "github.com/ntnguyencse/L-KaaS/pkg/kubernetes-client"
	randomstring "github.com/ntnguyencse/L-KaaS/pkg/randstring"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	capiv1alpha4 "sigs.k8s.io/cluster-api/api/v1alpha4"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// SoftwareInstallationReconciler reconciles a SoftwareInstallation object
type SoftwareInstallationReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

var (
	loggerSIC = ctrl.Log.WithName("Software Installation")
)

//+kubebuilder:rbac:groups=intent.automation.dcn.ssu.ac.kr,resources=softwareinstallations,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=intent.automation.dcn.ssu.ac.kr,resources=softwareinstallations/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=intent.automation.dcn.ssu.ac.kr,resources=softwareinstallations/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the SoftwareInstallation object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *SoftwareInstallationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// TODO(user): your logic here
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

	if !CAPOClusters.ObjectMeta.DeletionTimestamp.IsZero() {
		loggerSIC.Info("Cluster is Deleted")
		return ctrl.Result{}, nil
	}
	// CAPOStatus := CAPOClusters.DeepCopy().Status
	if len(CAPOClusters.ObjectMeta.OwnerReferences) > 0 {
		ownerRef := CAPOClusters.ObjectMeta.OwnerReferences[0]

		ownerLCluster, err := r.GetClusterOwnerObject(ctx, req, &ownerRef)
		if err != nil {
			loggerSIC.Error(err, "Error when get cluster owner's object")
			return ctrl.Result{}, nil
		}
		capoStatus := CAPOClusters.Status
		if capoStatus.Phase == string(capiv1alpha4.ClusterPhaseProvisioned) {
			// Get kubeconfig
			kubeconfig, err := r.getKubeConfigCluster(ctx, CAPOClusters.Name, CAPOClusters.Namespace)
			if err != nil {
				loggerSIC.Error(err, "Error when get Kubeconfig: "+CAPOClusters.Name)
				return ctrl.Result{RequeueAfter: 3 * time.Minute}, nil
			}
			if capoStatus.Phase == string(capiv1alpha4.ClusterPhaseProvisioned) {
				if ownerLCluster.Status.Ready && string(ownerLCluster.Status.Phase) == string(capiv1alpha4.ClusterPhaseProvisioned) && !ownerLCluster.Status.Registration {
					loggerSIC.Info("Start installing software " + CAPOClusters.Name)
					folderCAPOCluster := "/tmp/" + ownerLCluster.Name + randomstring.String(5) + "/"
					kubePath, _ := emcoctl.SaveValueFile(Name(CAPOClusters.Name, KubeConfigSecretSuffix+".kubeconfig"), folderCAPOCluster, &kubeconfig)
					// Check healthy of target cluster
					serverVer, errGetVer := kubernetesclient.GetKubernetesServerVersion(kubePath)

					if errGetVer == nil && len(serverVer.String()) > 3 {
						loggerSIC.Info("Cluster Version: ", serverVer.String(), "EOF")
						// Install CNI
						r.ReconcileInstallSoftware(ctx, req, kubePath, &ownerLCluster, CAPOClusters)
						// Update status
						errUpdate := r.Client.Status().Update(ctx, &ownerLCluster)
						if errUpdate != nil {
							loggerSIC.Error(errUpdate, "Error when update LKaaS cluster status")
							return ctrl.Result{RequeueAfter: 3 * time.Minute}, errUpdate
						}
						return ctrl.Result{}, nil
					} else {
						loggerSIC.Error(errGetVer, "Get Cluster version of "+ownerLCluster.Name)
						return ctrl.Result{RequeueAfter: 3 * time.Minute}, nil
					}

				}
			}
		}

	}
	return ctrl.Result{RequeueAfter: 3 * time.Minute}, nil
}
func (r *SoftwareInstallationReconciler) getKubeConfigCluster(ctx context.Context, clusterName, nameSpace string) (string, error) {
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
func (r *SoftwareInstallationReconciler) GetListProfile(ctx context.Context, req ctrl.Request) (*intentv1.ProfileList, error) {
	listProfiles := intentv1.ProfileList{}
	err := r.Client.List(ctx, &listProfiles)
	if err != nil {
		loggerCL.Error(err, "ReconcileInstallSoftware", "Error when list profiles")
		return nil, err
	}
	return &listProfiles, err
}
func (r *SoftwareInstallationReconciler) ReconcileInstallSoftware(ctx context.Context, req ctrl.Request, kubePath string, cluster *intentv1.Cluster, CAPICluster *capiv1alpha4.Cluster) error {
	// clusterStatus.Ready && string(clusterStatus.Phase) == string(capiv1alpha4.ClusterPhaseProvisioned) && !clusterStatus.RegistrationkubePath
	// Install CNI
	loggerSIC.Info("ReconcileInstallSoftware" + CAPICluster.Name)
	loggerSIC.Info("Begin Install CNI to  "+ CAPICluster.Name)
	// Get Profiles related to Clusters
	listProfiles, err := r.GetListProfile(ctx, req)
	if err != nil {
		loggerSIC.Error(err, "Error get profiles")
	}
	for _, item := range cluster.Spec.Network {
		// 1. CNI Profiles
		CNIProfileName := item.Name
		CNIProfile, err := FindProfileWithName(listProfiles, CNIProfileName)
		if err == nil {
			// Install CNI
			chartPath := CNIProfile.Spec.Values["url"]
			chartName := CNIProfileName + "-" + cluster.Name
			CNINamespace := CNIProfile.Spec.Values["namespace"]
			valueFilePath := CNIProfile.Spec.Values["value"]
			// if strings.Contains(CNIProfileName, "flannel"){
			// 	chartPath = flannelTemplate
			// }
			loggerSIC.Info(CAPICluster.Name + " Helm Installer:", kubePath, "chartName:", chartName, chartPath)
			err = helminstaller.Install(kubePath, chartName, chartPath, valueFilePath, CNINamespace)
			if err != nil {
				loggerSIC.Error(err, "Error install Network components: "+CNIProfileName)
				return err
			}
		}

	}
	loggerSIC.Info("Finish install CNI")
	loggerSIC.Info("Begin install Software")
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
			loggerSIC.Info("Helm Installer:", kubePath, "chartName:", chartName, chartPath)
			err = helminstaller.Install(kubePath, chartName, chartPath, valueFilePath, CNINamespace)
			if err != nil {
				loggerSIC.Error(err, "Error install Network components: "+SoftwareProfileName)
				return err
			}
		}
	}
	loggerSIC.Info("Finish install Software")

	// Update cluster status
	cluster.Status.Registration = true
	// Clean up
	// defer os.RemoveAll("$HOME/.cache/helm/repository")
	return nil
}
func (r *SoftwareInstallationReconciler) GetClusterOwnerObject(ctx context.Context, req ctrl.Request, ownerRef *metav1.OwnerReference) (intentv1.Cluster, error) {
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
			loggerSIC.Error(err, "Error when get Cluster in OwnerRef not Found: ")
			return lcluster, err
		} else {
			loggerSIC.Error(err, "Error when get Cluster in OwnerRef")
			return lcluster, err
		}
	}
	return lcluster, nil

}

// SetupWithManager sets up the controller with the Manager.
func (r *SoftwareInstallationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		// Uncomment the following line adding a pointer to an instance of the controlled resource as an argument
		// For().
		For(&capiv1alpha4.Cluster{}).
		Complete(r)
}
