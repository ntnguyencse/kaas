# Comparision with Cluster API
Advantages compared to Cluster API:
* Current initializing, setting up, and maintenance of a Kubernetes conformant cluster up and running are too 
complex because of highly integrated with low-level infrastructure exposed by each cloud provider or on-premise.

* Users can’t initialize their own cluster up and running without becoming technical experts and performing manual configuring actions. They want to be able to:

    *  Create their own Kubernetes cluster with their own configurations via a simple interface without any requirements of low-level infrastructure expert knowledge.

    * Have their Kubernetes cluster up and running automatically without any manual actions or maintenance.

* Standing on the existing Cluster API ecosystems, utilize the existing functionalities, in some cases, we can extend Cluster API by adding more functions to help the logical cluster lifecycle management process. 
