package v1

type ClusterPhase string
type ConditionType string
type ConditionReason string

const (
	ClusterPhasePending = ClusterPhase("Pending")

	ClusterPhaseProvisioning = ClusterPhase("Provisioning")

	ClusterPhaseProvisioned = ClusterPhase("Provisioned")

	ClusterPhaseDeleting = ClusterPhase("Deleting")

	ClusterPhaseFailed = ClusterPhase("Failed")

	ClusterPhaseUnknown = ClusterPhase("Unknown")
)
const (
	InstanceReadyCondition ConditionReason = "ClusterReady"

	WaitingForClusterInfrastructureReason ConditionReason = "WaitingForInfrastructure"

	InvalidClusterSpecReason ConditionReason = "InvalidClusterSpec"

	ClusterCreateFailedReason ConditionReason = "ClusterCreateFailed"

	ClusterNotFoundReason ConditionReason = "ClusterNotFound"

	ClusterStateErrorReason ConditionReason = "ClusterStateError"

	ClusterDeletedReason ConditionReason = "ClusterDeleted"

	ClusterNotReadyReason ConditionReason = "ClusterNotReady"

	ClusterDeleteFailedReason ConditionReason = "ClusterDeleteFailed"
)
