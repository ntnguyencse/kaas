helm install cilium https://github.com/ntnguyencse/helm-mod/raw/main/test/charts/cilium-1.13.0.tar.gz --namespace cilium-system --kubeconfig /tmp/



helm install prometheus https://github.com/prometheus-community/helm-charts/releases/download/kube-prometheus-stack-46.7.0/kube-prometheus-stack-46.7.0.tgz --namespace prometheus-system --kubeconfig /tmp/