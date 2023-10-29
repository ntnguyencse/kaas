

kubectl patch machines/edge-2-control-plane-5tlj9 --type json --patch='[ { "op": "remove", "path": "/metadata/finalizers" } ]'
kubectl patch openstackcluster/edge-2 --type json --patch='[ { "op": "remove", "path": "/metadata/finalizers" } ]'
kubectl patch openstackcluster/edge-1 --type json --patch='[ { "op": "remove", "path": "/metadata/finalizers" } ]'
kubectl patch kubeadmcontrolplanes/mec-1-control-plane --type json --patch='[ { "op": "remove", "path": "/metadata/finalizers" } ]'
kubectl patch kubeadmcontrolplanes/mec-2-control-plane --type json --patch='[ { "op": "remove", "path": "/metadata/finalizers" } ]'
kubectl patch openstackclusters.infrastructure.cluster.x-k8s.io/edge-1 --type json --patch='[ { "op": "remove", "path": "/metadata/finalizers" } ]'
kubectl patch openstackclusters.infrastructure.cluster.x-k8s.io/edge-2 --type json --patch='[ { "op": "remove", "path": "/metadata/finalizers" } ]'


kubectl patch logicalcluster/intra-edges --type json --patch='[ { "op": "remove", "path": "/metadata/finalizers" } ]'

kubectl patch clusters.cluster.x-k8s.io/mec-1 --type json --patch='[ { "op": "remove", "path": "/metadata/finalizers" } ]'




clusters.cluster.x-k8s.io    
machines.cluster.x-k8s.io   
machinedeployments.cluster.x-k8s.io  
kubeadmconfigs.bootstrap.cluster.x-k8s.io
kubeadmcontrolplanes.controlplane.cluster.x-k8s.io
kubeadmcontrolplanetemplates.controlplane.cluster.x-k8s.io
openstackclusters.infrastructure.cluster.x-k8s.io
