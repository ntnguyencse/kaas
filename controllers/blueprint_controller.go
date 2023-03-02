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
	jsonclassic "encoding/json"

	"github.com/go-logr/logr"
	intentv1 "github.com/ntnguyencse/intent-kaas/api/v1"
	config "github.com/ntnguyencse/intent-kaas/pkg/config"
	git "github.com/ntnguyencse/intent-kaas/pkg/git"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// BlueprintReconciler reconciles a Blueprint object
type BlueprintReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	l      logr.Logger
	s      *json.Serializer
}

const repo = "blueprints"
const owner = "ntnguyen-dcn"

// const config_path = "./config.yml"

var (
	logger = ctrl.Log.WithName("Blueprint Controller")
)

//+kubebuilder:rbac:groups=intent.automation.dcn.ssu.ac.kr,resources=blueprints,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=intent.automation.dcn.ssu.ac.kr,resources=blueprints/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=intent.automation.dcn.ssu.ac.kr,resources=blueprints/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Blueprint object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *BlueprintReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// Load Configurations
	configuration := config.LoadConfig(config.DEFAULT_CONFIG_PATH)

	r.l = log.FromContext(ctx)
	// Get all blueprint
	githubclient, _ := git.NewClient(configuration.BlueprintRepo, configuration.Owner, configuration.GitHubToken, ctx)
	logger.V(0).Info("Reconciling.... Blueprint\n")
	var blueprintList intentv1.BlueprintList
	err := r.Client.List(ctx, &blueprintList)
	if err != nil {
		logger.V(3).Error(err, "Error when list Blueprints", "error")

	}
	logger.V(1).Info("Print List Blueorint:\n")
	for _, item := range blueprintList.Items {
		logger.V(3).Info("Item: ", item.Name, ".\n")
	}
	var bp intentv1.Blueprint
	err = r.Get(ctx, req.NamespacedName, &bp)

	if err != nil {
		logger.V(3).Error(err, "unable to fetch PackageDeployment")
	} else {
		logger.V(3).Info("Print blue print: %s\n", "blueprint", bp, "Blueprint Name: %s\n", bp.Name)
		// filecontent, err := jsonclassic.Marshal(bp)
		// // content, err := r.s.Encode()
		// if err != nil {
		// 	logger.Error(err, "unable to Decode Json file")
		// } else {
		// 	gitClient.UpdateFile(bp.Name+".yaml", "blueprint/", filecontent)
		// }
		var bp1 intentv1.Blueprint
		err := jsonclassic.Unmarshal([]byte(bp.ObjectMeta.Annotations["kubectl.kubernetes.io/last-applied-configuration"]), &bp1)

		if err != nil {
			logger.V(3).Error(err, "Error when convert object")
		} else {
			logger.V(3).Info("Blueprint...", "blueprint", bp1)

		}
		content, err := jsonclassic.MarshalIndent(bp1, " ", "    ")
		if err != nil {
			logger.V(3).Error(err, "Error when marshal blueprint...")

		} else {
			logger.V(3).Info("Marshed blueprint:", "content", string(content))
			isYamlFileExist, err := githubclient.IsFileNotExist(bp1.Name+".yaml", "blueprints/")
			if err == nil {
				if !isYamlFileExist {
					githubclient.UpdateFile(bp1.Name+".yaml", "blueprints/", content)
				} else {
					githubclient.CommitNewFile(bp1.Name+".yaml", "main", "blueprints/", content)
				}
			}

		}

	}
	// // TODO(user): your logic here

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *BlueprintReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&intentv1.Blueprint{}).
		Complete(r)
}
