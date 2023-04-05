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

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	intentv1 "github.com/ntnguyencse/L-KaaS/api/v1"
)

// LogicalClusterReconciler reconciles a LogicalCluster object
type LogicalClusterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=intent.automation.dcn.ssu.ac.kr,resources=logicalclusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=intent.automation.dcn.ssu.ac.kr,resources=logicalclusters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=intent.automation.dcn.ssu.ac.kr,resources=logicalclusters/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the LogicalCluster object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *LogicalClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// Get/ Fetch Cluster Instance Logical cluster
	cluster := &intentv1.LogicalCluster{}
	if err := r.Client.Get(ctx, req.NamespacedName, cluster); err != nil {
		if apierrors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return ctrl.Result{}, nil
		}

		// Error reading the object - requeue the request.
		return ctrl.Result{}, err
	}

	// Do something
	// defer func() {}
	defer func() {
		// Always reconcile the Status.Phase field.
		// Reconcile phase of logical cluster
		r.ReconcileClusterPhase(ctx, cluster)
	}()

	// Handle deletion reconciliation loop.
	if !cluster.ObjectMeta.DeletionTimestamp.IsZero() {
		return r.ReconcileDelete(ctx, cluster)
	}

	// Handle normal reconciliation loop.
	// return ctrl.Result{}, nil
	return r.ReconcileNormal(ctx, req)
}

// SetupWithManager sets up the controller with the Manager.
func (r *LogicalClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&intentv1.LogicalCluster{}).
		Complete(r)
}

// Update Phase of Logical Cluster
func (r *LogicalClusterReconciler) ReconcileClusterPhase(ctx context.Context, cluster *intentv1.LogicalCluster) error {
	return nil
}

func (r *LogicalClusterReconciler) ReconcileNormal(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	// TODO(user): your logic here
	return ctrl.Result{}, nil
}

func (r *LogicalClusterReconciler) ReconcileDelete(ctx context.Context, cluster *intentv1.LogicalCluster) (ctrl.Result, error) {

	// TODO(user): your logic here
	return ctrl.Result{}, nil
}

func (r *LogicalClusterReconciler) ReconcileCreate(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	// TODO(user): your logic here
	return ctrl.Result{}, nil
}
