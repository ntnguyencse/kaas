clusterctl generate cluster test \
  --kubernetes-version v1.28.0 \
  --control-plane-machine-count=3 \
  --worker-machine-count=3 \
  --configure /home/ubuntu/l-kaas/L-KaaS/config/capi/clusterctl-config.yaml \
  > capi-quickstart.yaml
