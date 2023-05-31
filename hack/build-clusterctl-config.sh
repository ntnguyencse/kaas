#!/usr/bin/env bash
cat << EOF > clusterctl-config.yml
OPENSTACK_CLOUD: $OPENSTACK_CLOUD
OPENSTACK_CLOUD_CACERT_B64: $OPENSTACK_CLOUD_CACERT_B64
OPENSTACK_CLOUD_PROVIDER_CONF_B64: $OPENSTACK_CLOUD_PROVIDER_CONF_B64
OPENSTACK_CLOUD_YAML_B64: $OPENSTACK_CLOUD_YAML_B64
# The list of nameservers for OpenStack Subnet being created.
# Set this value when you need create a new network/subnet while the access through DNS is required.
OPENSTACK_DNS_NAMESERVERS: 8.8.8.8
# FailureDomain is the failure domain the machine will be created in.
OPENSTACK_FAILURE_DOMAIN: compute
# OPENSTACK_FAILURE_DOMAIN: $(OPENSTACK_FAILURE_DOMAIN)
# The flavor reference for the flavor for your server instance.
OPENSTACK_CONTROL_PLANE_MACHINE_FLAVOR: m1.medium
# The flavor reference for the flavor for your server instance.
OPENSTACK_NODE_MACHINE_FLAVOR: m1.medium
# The name of the image to use for your server instance. If the RootVolume is specified, this will be ignored and use rootVolume directly.
OPENSTACK_IMAGE_NAME: ubuntu-k8-1.22
# The SSH key pair name
OPENSTACK_SSH_KEY_NAME: $OPENSTACK_SSH_KEY_NAME
# The external network
OPENSTACK_EXTERNAL_NETWORK_ID: $OPENSTACK_EXTERNAL_NETWORK_ID
# Enabling Feature Gates
CLUSTER_TOPOLOGY: true
# Kubernetes Version
KUBERNETES_VERSION: $KUBERNETES_VERSION
# CLUSTER_TOPOLOGY=true
EOF