# Problems

* Up til now to spin up a new Kubernetes cluster: using Kubernetes as a service of public cloud providers (e.g EKS, GKE), using command-line, using Terraform/Ansible. But *all of them require to have a deep knowledge of Kubernetes and a ton of Kubernetes-related configuration*.
* Administrator must maintain the cluster after spinning up, caring about workload, and keep an eye on Kubernetes cluster at all times to respond if any incident happens. Therefore, *if system consists of multiple clusters, or tens of clusters, administrators must put in a lot of effort to maintain, upgrade,  and monitor*.
* Cluster API is a project that provides declarative APIs and tooling to simplify provisioning, upgrading, and operating Kubernetes clusters but *there is no project that provides provisioning, managing, or operating a system that has multiple clusters operating and running in a distributed environment*.
