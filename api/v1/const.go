package v1

// Delete Reason
type DeleteReason string

const (
	DeletingReason DeleteReason = "Deleting"

	DeletionFailedReason DeleteReason = "DeletionFailed"

	DeletedReason DeleteReason = "Deleted"
)

type ClusterLabel string

const (
	ClusterNameLabel ClusterLabel = "intent.automation.dcn.ssu.ac.kr/cluster-name"

	ClusterTopologyOwnedLabel ClusterLabel = "topology.intent.automation.dcn.ssu.ac.kr/owned"

	ClusterTopologyClusterMemberNameLabel ClusterLabel = "topology.intent.automation.dcn.ssu.ac.kr/cluster-member-name"

	ProviderNameLabel ClusterLabel = "intent.automation.dcn.ssu.ac.kr/provider"

	ClusterNameAnnotation ClusterLabel = "intent.automation.dcn.ssu.ac.kr/cluster-name"

	ClusterNamespaceAnnotation ClusterLabel = "intent.automation.dcn.ssu.ac.kr/cluster-namespace"

	ClusterAnnotation ClusterLabel = "intent.automation.dcn.ssu.ac.kr/cluster"

	OwnerKindAnnotation ClusterLabel = "intent.automation.dcn.ssu.ac.kr/owner-kind"

	LabelsFromClusterMemberAnnotation ClusterLabel = "intent.automation.dcn.ssu.ac.kr/labels-from-cluster-member"

	OwnerNameAnnotation ClusterLabel = "intent.automation.dcn.ssu.ac.kr/owner-name"

	PausedAnnotation ClusterLabel = "intent.automation.dcn.ssu.ac.kr/paused"

	DisableClusterCreateAnnotation ClusterLabel = "intent.automation.dcn.ssu.ac.kr/disable-cluster-create"

	WatchLabel ClusterLabel = "intent.automation.dcn.ssu.ac.kr/watch-filter"

	DeleteClusterAnnotation ClusterLabel = "intent.automation.dcn.ssu.ac.kr/delete-cluster"

	TemplateClonedFromGroupKindAnnotation ClusterLabel = "intent.automation.dcn.ssu.ac.kr/cloned-from-groupkind"

	ClusterSecretType ClusterLabel = "intent.automation.dcn.ssu.ac.kr/secret"

	ManagedByAnnotation ClusterLabel = "intent.automation.dcn.ssu.ac.kr/managed-by"

	TopologyDryRunAnnotation ClusterLabel = "topology.intent.automation.dcn.ssu.ac.kr/dry-run"

	ReplicasManagedByAnnotation ClusterLabel = "intent.automation.dcn.ssu.ac.kr/replicas-managed-by"

	VariableDefinitionFromInline ClusterLabel = "inline"
)

const (
	ControlPlaneInitializedCondition ConditionType = "ControlPlaneInitialized"

	WaitingForControlPlaneProviderInitializedReason ConditionType = "WaitingForControlPlaneProviderInitialized"

	ControlPlaneReadyCondition ConditionType = "ControlPlaneReady"

	WaitingForControlPlaneFallbackReason ConditionType = "WaitingForControlPlane"

	WaitingForControlPlaneAvailableReason ConditionType = "WaitingForControlPlaneAvailable"
)

const (
	ClusterHealthCheckSucceededCondition ConditionType = "HealthCheckSucceeded"

	ClusterHasFailureReason ConditionType = "ClusterHasFailure"

	ClusterStartupTimeoutReason ConditionType = "ClusterStartupTimeout"

	UnhealthyClusterConditionReason ConditionType = "UnhealthyCluster"
)
