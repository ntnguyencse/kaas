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

	"github.com/go-logr/logr"
	intentv1 "github.com/ntnguyencse/L-KaaS/api/v1"
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
)

// LogicalClusterControlPlaneProviderReconciler reconciles a LogicalClusterControlPlaneProvider object
type LogicalClusterControlPlaneProviderReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	l      logr.Logger
	s      *json.Serializer
}

const timeoutRetryCreateLogicalCluster = 10 * time.Minute

var (
	loggerLKP = ctrl.Log.WithName("L-KaaS Control Plane Provider")
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

	// Separate status object:
	// CAPOPhaseStatus := CAPOStatus.Phase

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

func (r *LogicalClusterControlPlaneProviderReconciler) GetOwnerObject(ctx context.Context, req ctrl.Request, ownerRef *metav1.OwnerReference) intentv1.LogicalCluster, error {
	// lclusters := intentv1.LogicalClusterList{}
	lcluster := intentv1.LogicalCluster{}
	// r.Client.List(ctx, &lclusters)
	err := r.Client.Get(ctx, client.ObjectKey{
		Name:      ownerRef.Name,
		Namespace: req.Namespace,
	}, &lcluster)

	// Check error when get logical cluster corresponding in ownerRef
	if err != nil {
		loggerLKP.Error(err, "Error when get Logical Cluster in OwnerRef")
		return  lcluster, err
	}
	if lcluster != nil {
		return lcluster, error{"Empty Logical Cluster"}
	} else {
		return lcluster, nil
	}
}
func (r *LogicalClusterControlPlaneProviderReconciler) RegisterLogicalCLusterToEMCO(ctx context.Context, logicalCluster intentv1.LogicalCluster ) error {

	// Init EMCO Client
	
	return nil
}