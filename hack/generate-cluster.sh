clusterctl generate cluster test \
  --infrastructure openstack \
  --kubernetes-version v1.28.0 \
  --control-plane-machine-count=3 \
  --worker-machine-count=3 \
  --configure /home/ubuntu/l-kaas/L-KaaS/config/capi/clusterctl-config.yaml \
  > capi-quickstart.yaml


clusterctl generate cluster my-cluster --kubernetes-version v1.28.0 \
    --infrastructure aws --config /home/ubuntu/aws-capi/capi-config/config.yaml > my-cluster.yaml

### Openstack
clusterctl generate cluster my-cluster --kubernetes-version v1.28.0     --infrastructure openstack:v0.7.0  --config /home/ubuntu/aws-capi/capi-config/config.yaml > my-cluster-openstack.yaml


### AWS
clusterctl generate cluster my-cluster --kubernetes-version v1.28.0     --infrastructure aws:v2.2.2  --config /home/ubuntu/aws-capi/capi-config/config.yaml > my-cluster-aws.yaml