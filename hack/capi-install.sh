curl -L https://github.com/kubernetes-sigs/cluster-api/releases/download/v1.4.3/clusterctl-linux-amd64 -o clusterctl
sudo install -o root -g root -m 0755 clusterctl /usr/local/bin/clusterctl
clusterctl version
export CLUSTER_TOPOLOGY=true
# Initialize the management cluster
clusterctl init --infrastructure=openstack:v0.6.4
