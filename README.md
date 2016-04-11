# Kubernetes API Server Discovery

Utility to discover a working Kubernetes API server address and hand it over to other services.

In case of running a [high-availability (HA) Kubernetes cluster](http://kubernetes.io/docs/admin/high-availability/), the kube-apiserver endpoints must be load balanced, since kubelet and kube-proxy accept only one API endpoint [https://github.com/kubernetes/kubernetes/issues/18174]((kubernetes #18154)). In certain environments an internal load balancer is not desirable, and Kubernetes' load balancer is not available before a running a properly configured kube-proxy.

kube-discovery queries a set of seed API for a service definition and prints a randomly selected API server IP and port and returns with a zero exit status. If none of the seed servers responsed appropriately, or no functional API server was found in the service definition, the process returns with a non-zero exit status.

### License

[BSD 3-Clause License](LICENSE)
