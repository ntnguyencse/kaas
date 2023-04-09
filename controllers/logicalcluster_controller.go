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

	"github.com/go-logr/logr"
	intentv1 "github.com/ntnguyencse/L-KaaS/api/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/apimachinery/pkg/types"
	capiulti "sigs.k8s.io/cluster-api/util"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// LogicalClusterReconciler reconciles a LogicalCluster object
type LogicalClusterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	l      logr.Logger
	s      *json.Serializer
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

func (r *LogicalClusterReconciler) GetOrCreateCluster(ctx context.Context, lcluster *intentv1.LogicalCluster, clusterName string) (intentv1.Cluster, error) {
	cluster := intentv1.Cluster{}

	// Step 1: Get the Clusters in CAPI
	CAPIClusterList := intentv1.ClusterList{}
	err := r.Client.List(ctx, &CAPIClusterList, client.InNamespace(lcluster.Namespace))
	if err != nil {
		logger.V(1).Error(err, "Error when get List CAPI Cluster")
	}
	// Step 2: If cluster is existed, add OwnerRef to existing cluster
	existingCAPICluster := intentv1.Cluster{}
	isAlreadyExist := false
	for _, clster := range CAPIClusterList.Items {
		if clster.Name == clusterName {
			isAlreadyExist = true
			existingCAPICluster = clster
		}
	}
	if isAlreadyExist {
		// Add OwnerRef to existing Cluster
		// CAPI Cluster is owned by Logical Cluster
		existingCAPICluster.SetOwnerReferences(capiulti.EnsureOwnerRef(existingCAPICluster.GetOwnerReferences(), metav1.OwnerReference{
			APIVersion: lcluster.APIVersion,
			Kind:       lcluster.Kind,
			Name:       lcluster.Name,
			UID:        lcluster.UID,
		}))
		// Update to API Server
		err := r.Client.Update(ctx, &existingCAPICluster)
		if err != nil {
			logger.Error(err, "Error when update OwnerRef of Cluster")
		}
	}
	// Step 3: If cluster is not existed, create a new one
	// Create new cluster
	// If Cluster contain both of Catalog and Profile,
	// We prioritize create Cluster with Profile

	// Check Cluster Member contains Catalog or not
	// TODO: Check Catalog and Profile
	return cluster, nil
}

// Create Cluster from Cluster Catalog
func (r *LogicalClusterReconciler) CreateClusterFromClusterCatalog(ctx context.Context, lcluster *intentv1.LogicalCluster, clusterSpec *intentv1.ClusterMember) (intentv1.Cluster, error) {
	newCluster := intentv1.Cluster{}
	// Get Catalog
	clusterCatalog := intentv1.ClusterCatalog{}
	key := types.NamespacedName{Namespace: lcluster.Namespace, Name: clusterSpec.ClusterCatalog}
	err := r.Client.Get(ctx, key, &clusterCatalog)
	if err != nil {
		logger.V(1).Error(err, "Error when get Cluster Catalog")
		// Check which cause error
		// Catalog does not exist...
		return newCluster, err
	}
	// Create Cluster from Catalog

	// TODO: Create Cluster from Catalog
	// Build a Cluster with profile pieces in Cluster Catalog
	newCluster, err = r.BuildClusterObjectFromCatalog(ctx, lcluster, clusterSpec, &clusterCatalog)
	if err != nil {
		logger.Error(err, "Error when BuildClusterObjectFromCatalog")
		return newCluster, err
	}
	return newCluster, nil
}

// Create Cluster from Cluster Profile
func (r *LogicalClusterReconciler) CreateClusterFromClusterProfile(ctx context.Context, lcluster *intentv1.LogicalCluster, clusterSpec *intentv1.ClusterMember) (intentv1.Cluster, error) {
	newCluster := intentv1.Cluster{}
	// Get Catalog
	clusterCatalog := intentv1.ClusterCatalog{}
	key := types.NamespacedName{Namespace: lcluster.Namespace, Name: clusterSpec.ClusterCatalog}
	err := r.Client.Get(ctx, key, &clusterCatalog)
	if err != nil {
		logger.V(1).Error(err, "Error when get Cluster Catalog")
		// Check which cause error
		// Catalog does not exist...
		return newCluster, err
	}
	// Create Cluster from Catalog

	// TODO: Create Cluster from Catalog
	// Build a Cluster with profile pieces in Cluster Catalog
	newCluster, err = r.BuildClusterObjectFromProfile(ctx, lcluster, clusterSpec, &clusterCatalog)
	if err != nil {
		logger.Error(err, "Error when BuildClusterObjectFromProfile")
		return newCluster, err
	}
	return newCluster, nil
}

func (r *LogicalClusterReconciler) BuildClusterObjectFromCatalog(ctx context.Context, lcluster *intentv1.LogicalCluster, clusterMember *intentv1.ClusterMember, clusterCatalog *intentv1.ClusterCatalog) (intentv1.Cluster, error) {

	// Get Profile from Catalog
	clusterSpec := intentv1.ClusterSpec{
		Infrastructure: clusterCatalog.Spec.Infrastructure,
		Software:       clusterCatalog.Spec.Software,
	}
	// Construct a Object
	clusterObject := intentv1.Cluster{
		TypeMeta: metav1.TypeMeta{
			APIVersion: intentv1.GroupVersion.String(),
			Kind:       "Cluster",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      clusterMember.DeepCopy().Name,
			Namespace: lcluster.DeepCopy().Namespace,
			// OwnerRef from Logical cluster
			OwnerReferences: []metav1.OwnerReference{*&metav1.OwnerReference{
				APIVersion: lcluster.APIVersion,
				Kind:       lcluster.DeepCopy().Kind,
				Name:       lcluster.DeepCopy().Name,
				UID:        lcluster.DeepCopy().UID,
			}},
		},
		Spec: clusterSpec,
	}
	return clusterObject, nil
}
func (r *LogicalClusterReconciler) BuildClusterObjectFromProfile(ctx context.Context, lcluster *intentv1.LogicalCluster, clusterMember *intentv1.ClusterMember, clusterCatalog *intentv1.ClusterCatalog) (intentv1.Cluster, error) {
	// Get Profile from Catalog
	clusterSpec := intentv1.ClusterSpec{
		Infrastructure: clusterMember.ClusterMemberSpec.Infrastructure,
		Software:       clusterMember.ClusterMemberSpec.Software,
	}
	// Construct a Object
	clusterObject := intentv1.Cluster{
		TypeMeta: metav1.TypeMeta{
			APIVersion: intentv1.GroupVersion.String(),
			Kind:       "Cluster",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      clusterMember.DeepCopy().Name,
			Namespace: lcluster.DeepCopy().Namespace,
			// OwnerRef from Logical cluster
			OwnerReferences: []metav1.OwnerReference{*&metav1.OwnerReference{
				APIVersion: lcluster.APIVersion,
				Kind:       lcluster.DeepCopy().Kind,
				Name:       lcluster.DeepCopy().Name,
				UID:        lcluster.DeepCopy().UID,
			}},
		},
		Spec: clusterSpec,
	}
	return clusterObject, nil

}
