#!/usr/bin/env bash

export KUBECONFIG=~/config
# See log of openstack controller
kubectl logs -f -l=cluster.x-k8s.io/provider=infrastructure-openstack -n capo-system
# See log of cluster api controller
kubectl logs -f -l=cluster.x-k8s.io/provider=infrastructure-openstack -n capi-system