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
	clusterctlclient "github.com/ntnguyencse/cluster-api-sdk/client"
	kubernetesclient "github.com/ntnguyencse/cluster-api-sdk/kubernetes-client"
	intentv1 "github.com/ntnguyencse/intent-kaas/api/v1"
	config "github.com/ntnguyencse/intent-kaas/pkg/config"
)

// ClusterDescriptionReconciler reconciles a ClusterDescription object
type ClusterDescriptionReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	l      logr.Logger
	s      *json.Serializer
}

//+kubebuilder:rbac:groups=intent.automation.dcn.ssu.ac.kr,resources=clusterdescriptions,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=intent.automation.dcn.ssu.ac.kr,resources=clusterdescriptions/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=intent.automation.dcn.ssu.ac.kr,resources=clusterdescriptions/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ClusterDescription object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *ClusterDescriptionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	r.l = log.FromContext(ctx)
	r.l.Info("Reconciling.... CLusterDescription")
	// Load Openstack configuration from file
	openstackConfig := config.LoadOpenStackConfig(config.DEFAULT_OPENSTACKCONFIG_PATH)

	// CLuster Resource object get from Kubernetes API Server
	var deploy intentv1.ClusterDescription
	err := r.Get(context.Background(), req.NamespacedName, &deploy)
	if err != nil {
		if errors.IsNotFound(err) {
			// The Cluster Resources has been deleted, so we need to delete the cluster resource description corresponding
			logger1.V(1).Info("The Cluster Description has been deleted, so we need to delete the physical cluster corresponding")
			/////
			// TO-DO: Delete the physical cluster
			////
			////
			return ctrl.Result{}, nil
		}
		// There was an error getting the Deployment, so we'll retry later
		logger1.V(1).Info("There was an error getting the Cluster Description, so we'll retry later")
		return ctrl.Result{}, err
	}
	// Applying the changes of cluster description to openstack server
	if deploy.Status.Revision != deploy.Generation {
		// Create Kubernetes ctl client
		var kubeconfigFile = "./admin.conf"
		// clientset, _
		_, _ = kubernetesclient.CreateKubernetesClient(&kubeconfigFile)
		// Init the Openstack configuration map string
		var configs = map[string]string{
			"OPENSTACK_IMAGE_NAME":                   openstackConfig.OPENSTACK_IMAGE_NAME,
			"OPENSTACK_EXTERNAL_NETWORK_ID":          openstackConfig.OPENSTACK_EXTERNAL_NETWORK_ID,
			"OPENSTACK_DNS_NAMESERVERS":              openstackConfig.OPENSTACK_DNS_NAMESERVERS,
			"OPENSTACK_SSH_KEY_NAME":                 openstackConfig.OPENSTACK_SSH_KEY_NAME,
			"OPENSTACK_CLOUD_CACERT_B64":             openstackConfig.OPENSTACK_CLOUD_CACERT_B64,
			"OPENSTACK_CLOUD_PROVIDER_CONF_B64":      openstackConfig.OPENSTACK_CLOUD_PROVIDER_CONF_B64,
			"OPENSTACK_CLOUD_YAML_B64":               openstackConfig.OPENSTACK_CLOUD_YAML_B64,
			"OPENSTACK_FAILURE_DOMAIN":               openstackConfig.OPENSTACK_FAILURE_DOMAIN,
			"OPENSTACK_CLOUD":                        openstackConfig.OPENSTACK_CLOUD,
			"OPENSTACK_CONTROL_PLANE_MACHINE_FLAVOR": openstackConfig.OPENSTACK_CONTROL_PLANE_MACHINE_FLAVOR,
			"OPENSTACK_NODE_MACHINE_FLAVOR":          openstackConfig.OPENSTACK_NODE_MACHINE_FLAVOR,
		}
		providerConfigs := clusterctlclient.CreateProviderConfig(clusterctlclient.OPENSTACK, clusterctlclient.OPENSTACK_URL, clusterctlclient.InfrastructureProviderType)
		// Create cluster ctl client
		c, err := clusterctlclient.CreateNewClient(kubeconfigFile, configs, providerConfigs)
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ClusterDescriptionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&intentv1.ClusterDescription{}).
		Complete(r)
}
