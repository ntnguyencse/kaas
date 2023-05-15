#!/usr/bin/env bash
cat << EOF > capictl-config.yml
openstackimagename: $OPENSTACK_IMAGE_NAME
openstackexternalnetworkid: $OPENSTACK_IMAGE_NAME
openstackdnsnameservers: $OPENSTACK_IMAGE_NAME
openstacksshkeyname: $OPENSTACK_IMAGE_NAME
openstackcloudcacertb64: $OPENSTACK_IMAGE_NAME
openstackcloudproviderconfb64: $OPENSTACK_IMAGE_NAME
openstackcloudyamlb64: $OPENSTACK_IMAGE_NAME
openstackfailuredomain: $OPENSTACK_IMAGE_NAME
openstackcloud:  $OPENSTACK_IMAGE_NAME
openstackcontrolplanemachineflavor: $OPENSTACK_IMAGE_NAME
openstacknodemachineflavor : $OPENSTACK_IMAGE_NAME
kubernetesversion: $OPENSTACK_IMAGE_NAME
EOF