package v1

type CloudProvider string

const (
	DEFAULT   CloudProvider = "Default"
	AWS       CloudProvider = "AWS"
	OPENSTACK CloudProvider = "OpenStack"
)
