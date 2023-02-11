kubebuilder init --domain automation.dcn.ssu.ac.kr --owner "Nguyen Thanh Nguyen" --project-name "kubernetes-as-a-service" --repo "github.com/ntnguyencse/kaas"
kubebuilder create api --controller true --group intent --version v1 --kind Cluster  --resource true
make manifests