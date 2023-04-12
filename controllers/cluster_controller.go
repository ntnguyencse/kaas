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
	// jsonclassic "encoding/json"

	// config "github.com/ntnguyencse/L-KaaS/pkg/config"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/go-logr/logr"
	intentv1 "github.com/ntnguyencse/L-KaaS/api/v1"
	// "github.com/ntnguyencse/L-KaaS/pkg/git"
)

// ClusterReconciler reconciles a Cluster object
type ClusterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	l      logr.Logger
	s      *json.Serializer
}

var (
	loggerCL = ctrl.Log.WithName("Cluster Controller")
)

const NAMESPACE_DEFAULT string = "default"
const STATUS_GENERATED string = "Generated"

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
	r.GetOrCreateCluster(ctx, cluster)
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

func (r *ClusterReconciler) GetOrCreateCluster(ctx context.Context, cluster *intentv1.Cluster) (intentv1.ClusterDescription, error) {
	// Get list Profiles
	listProfiles := intentv1.ProfileList{}
	err := r.Client.List(ctx, &listProfiles)
	if err != nil {
		loggerCL.Error(err, "GetOrCreateCluster", "Error when listt profiles")
	}
	clusterDescriptiton, err := r.TransformClusterToClusterDescription(ctx, *cluster, listProfiles.Items)
	loggerCL.Info("Print CLuster Descrition", clusterDescriptiton.Name, clusterDescriptiton)
	return clusterDescriptiton, err
}
