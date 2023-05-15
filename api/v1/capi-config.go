package v1

import (
	clusterctlv1 "sigs.k8s.io/cluster-api/cmd/clusterctl/api/v1alpha3"
)

type CAPIConfig map[string]string

// Config for Provider
type ProviderConfig struct {
	Name         string
	URL          string
	ProviderType clusterctlv1.ProviderType
}

// COnfig for CAPI Cluster
type ClusterConfig map[string]string
