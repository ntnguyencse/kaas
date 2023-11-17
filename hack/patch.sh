

kubectl patch machines/edge-1-control-plane-vpqsc --type json --patch='[ { "op": "remove", "path": "/metadata/finalizers" } ]'
kubectl patch openstackcluster/edge-2 --type json --patch='[ { "op": "remove", "path": "/metadata/finalizers" } ]'
kubectl patch openstackcluster/edge-1 --type json --patch='[ { "op": "remove", "path": "/metadata/finalizers" } ]'
kubectl patch kubeadmcontrolplanes/mec-1-control-plane --type json --patch='[ { "op": "remove", "path": "/metadata/finalizers" } ]'
kubectl patch kubeadmcontrolplanes/mec-2-control-plane --type json --patch='[ { "op": "remove", "path": "/metadata/finalizers" } ]'
kubectl patch openstackclusters.infrastructure.cluster.x-k8s.io/edge-1 --type json --patch='[ { "op": "remove", "path": "/metadata/finalizers" } ]'
kubectl patch openstackclusters.infrastructure.cluster.x-k8s.io/edge-2 --type json --patch='[ { "op": "remove", "path": "/metadata/finalizers" } ]'


kubectl patch logicalcluster/intra-edges --type json --patch='[ { "op": "remove", "path": "/metadata/finalizers" } ]'

kubectl patch clusters.cluster.x-k8s.io/edge-1 --type json --patch='[ { "op": "remove", "path": "/metadata/finalizers" } ]'

kubectl  logs capo-controller-manager-7d454d8df9-2vfb6 -n capo-system



clusters.cluster.x-k8s.io    
machines.cluster.x-k8s.io   
machinedeployments.cluster.x-k8s.io  
kubeadmconfigs.bootstrap.cluster.x-k8s.io
kubeadmcontrolplanes.controlplane.cluster.x-k8s.io
kubeadmcontrolplanetemplates.controlplane.cluster.x-k8s.io
openstackclusters.infrastructure.cluster.x-k8s.io
watch kubectl get pod -A --kubeconfig /tmp/intra-edges/edge-1-kubeconfig.kubeconfig
watch kubectl get pod -A --kubeconfig /tmp/intra-edges/edge-2-kubeconfig.kubeconfig
watch kubectl get node --kubeconfig /tmp/intra-edges/edge-2-kubeconfig.kubeconfig