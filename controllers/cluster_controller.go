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

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/go-logr/logr"
	intentv1 "github.com/ntnguyencse/intent-kaas/api/v1"
	"github.com/ntnguyencse/intent-kaas/pkg/git"
)

// ClusterReconciler reconciles a Cluster object
type ClusterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	l      logr.Logger
	s      *json.Serializer
}

var (
	logger1 = ctrl.Log.WithName("Cluster Controller")
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
	r.l = log.FromContext(ctx)
	gitclient1, err := git.NewClient("clusters", "ntnguyen-dcn", "", ctx)
	if err != nil {
		r.l.Error(err, "Error while create new Github client")
		_ = gitclient1
		return ctrl.Result{}, nil
	}
	r.l.Info("Reconciling.... CLuster")
	logger1.Info("Start reconciling Cluster Resource")

	var deploy intentv1.Cluster
	err = r.Get(context.Background(), req.NamespacedName, &deploy)
	if err != nil {
		if errors.IsNotFound(err) {
			// The Deployment has been deleted, so we don't need to do anything
			logger1.V(1).Info("The Deployment has been deleted, so we don't need to do anything")
			return ctrl.Result{}, nil
		}
		// There was an error getting the Deployment, so we'll retry later
		logger1.V(1).Info("There was an error getting the Deployment, so we'll retry later")
		return ctrl.Result{}, err
	}
	// Transform Cluster Resource to Cluster Description

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&intentv1.Cluster{}).
		Complete(r)
}

// Draft
// func test() {
// 	var deploy intentv1.Cluster
// 	gitclient1, err := git.NewClient("a", "a", "", context.TODO())
// 	contentLastAppliedConfigurationBytes := deploy.Annotations["kubectl.kubernetes.io/last-applied-configuration"]
// 	logger1.V(1).Info("Print Deploy:", "Cluster", deploy)
// 	var content intentv1.Cluster
// 	err = jsonclassic.Unmarshal([]byte(contentLastAppliedConfigurationBytes), &content)

// 	if err != nil {
// 		logger1.V(0).Info("Error when unmarshal cluster resource")
// 	} else {
// 		if len(deploy.ManagedFields) > 0 {
// 			logger1.V(0).Info("ManagedField", deploy.ManagedFields[0].Manager, deploy.ManagedFields[0].Operation, deploy.ManagedFields[0].Time.String(), deploy.ManagedFields[0].APIVersion)
// 		}
// 	}
// 	// Check status and sha, version, revision
// 	logger1.Info("Print status of object:", "Status", deploy.Status.Status, deploy.Status.SHA, deploy.Status.Sync)

// 	///////////////////////////================
// 	logger1.Info("Print Before status of object:", "Status", deploy.Status.Status, deploy.Status.SHA, deploy.Status.Sync)
// 	if deploy.Status.Status != "Reconciling" {
// 		deploy.Status.Status = "Reconciling"
// 		deploy.Status.Sync = "Not Synced"
// 		deploy.Status.SHA = "SHA Tests"
// 		// Update status of Kubernetes objects
// 		// err = r.Client.Status().Update(ctx, &deploy)
// 		if err != nil {
// 			logger1.V(0).Error(err, "Error while update status object Not Synced")
// 			// return ctrl.Result{}, err
// 		}
// 		var deploy1 intentv1.Cluster
// 		// err = r.Get(context.Background(), req.NamespacedName, &deploy1)
// 		if err != nil {
// 			if errors.IsNotFound(err) {
// 				// The Deployment has been deleted, so we don't need to do anything
// 				logger1.V(1).Info("The Deployment has been deleted, so we don't need to do anything")
// 				// return ctrl.Result{}, nil
// 			}
// 			// There was an error getting the Deployment, so we'll retry later
// 			logger1.V(1).Info("There was an error getting the Deployment, so we'll retry later")
// 			// return ctrl.Result{}, err
// 		}
// 		logger1.Info("Print status of object:", "Status", deploy1.Status.Status, deploy1.Status.SHA, deploy1.Status.Sync)

// 	}

// 	if false {
// 		// if deploy.CreationTimestamp == *deploy.ManagedFields[0].Time {
// 		if len(deploy.ObjectMeta.Labels["SHA"]) < 2 {
// 			isFileNotExist, err := gitclient1.IsFileNotExist(deploy.Name+".yaml", "test/")
// 			if err != nil {
// 				logger1.Error(err, "Error while check file existing..")
// 			} else {
// 				content1, _ := jsonclassic.MarshalIndent(content, " ", "    ")
// 				var sha string
// 				if !isFileNotExist {
// 					resp, err := gitclient1.UpdateFile(deploy.Name+".yaml", "test/", content1)
// 					if err != nil {
// 						logger1.V(0).Error(err, "Error whilde reconciling cluster resource. Update file")
// 					} else {
// 						sha = *resp.SHA
// 						logger1.V(0).Info("Updated SHA", "SHA", sha)
// 					}

// 				} else {
// 					resp, err := gitclient1.CommitNewFile(deploy.Name+".yaml", "main", "test/", content1)
// 					if err != nil {
// 						logger1.V(0).Error(err, "Error whilde reconciling cluster resource. Commit file")
// 					} else {
// 						sha = *resp.SHA
// 						logger1.V(0).Info("Updated SHA", "SHA", sha)
// 					}
// 				}

// 				deploy.ObjectMeta.Labels["SHA"] = sha
// 				deploy.ObjectMeta.Labels["Sync"] = "Synced"
// 				// deploy.Status.Status
// 				// err := r.Update(ctx, &deploy)
// 				if err != nil {
// 					logger1.V(0).Error(err, "Error while update status object Synced")
// 					// return ctrl.Result{}, err
// 				}
// 			}
// 		} else {
// 			// deploy.ObjectMeta.Labels["Sync"] = "UnSynced"
// 			deploy.Status.Sync = "Not Synced"
// 			// err := r.Update(ctx, &deploy)
// 			if err != nil {
// 				logger1.V(0).Error(err, "Error while update status object Not Synced")
// 				// return ctrl.Result{}, err
// 			}
// 			var deploy1 intentv1.Cluster
// 			// err = r.Get(context.Background(), req.NamespacedName, &deploy1)
// 			if err != nil {
// 				if errors.IsNotFound(err) {
// 					// The Deployment has been deleted, so we don't need to do anything
// 					logger1.V(1).Info("The Deployment has been deleted, so we don't need to do anything")
// 					// return ctrl.Result{}, nil
// 				}
// 				// There was an error getting the Deployment, so we'll retry later
// 				logger1.V(1).Info("There was an error getting the Deployment, so we'll retry later")
// 				// return ctrl.Result{}, err
// 			}
// 			logger1.Info("Print status of object:", "Status", deploy1.Status.Status, deploy1.Status.SHA, deploy1.Status.Sync)
// 		}
// 	}
// 	// r.Status().Update()
// }
