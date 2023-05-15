#!/usr/bin/env bash
NAMESPACE = "namespace"
CLUSTER_NAME = "cluster_name"
kubectl patch crd/clusters.cluster.x-k8s.io -p ${CLUSTER_NAME} '{"metadata":{"finalizers":[]}}' --type=merge -n $NAMESPACE
kubectl patch crd/clusterclasses.cluster.x-k8s.io -p   '{"metadata":{"finalizers":[]}}' --type=merge -n $NAMESPACE
kubectl patch crd/machinedeployments.cluster.x-k8s.io -p  '{"metadata":{"finalizers":[]}}' --type=merge -n $NAMESPACE
kubectl patch crd/openstackmachines.infrastructure.cluster.x-k8s.io  -p  '{"metadata":{"finalizers":[]}}' --type=merge -n $NAMESPACE 
kubectl patch crd/openstackmachinetemplates.infrastructure.cluster.x-k8s.io -p ${CLUSTER_NAME}-md-0 '{"metadata":{"finalizers":[]}}' --type=merge -n $NAMESPACE