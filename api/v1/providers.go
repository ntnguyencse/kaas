package v1

type CloudProvider string

const (
	DEFAULT   CloudProvider = "default"
	AWS       CloudProvider = "aws"
	OPENSTACK CloudProvider = "openstack"
)
