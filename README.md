# Kubernetes API Server Discovery

[![GoReportCard Widget]][GoReportCard] [![Travis Widget]][Travis]

[Travis]: https://travis-ci.org/egergo/kube-discovery
[Travis Widget]: https://travis-ci.org/egergo/kube-discovery.svg?branch=master
[GoReportCard]: https://goreportcard.com/report/github.com/egergo/kube-discovery
[GoReportCard Widget]: https://goreportcard.com/badge/github.com/egergo/kube-discovery


Utility to discover a working Kubernetes API server address and hand it over to other services.

In case of running a [high-availability (HA) Kubernetes cluster](http://kubernetes.io/docs/admin/high-availability/), the kube-apiserver endpoints must be load balanced, since kubelet and kube-proxy accept only one API endpoint [(kubernetes #18154)](https://github.com/kubernetes/kubernetes/issues/18174). In certain environments an internal load balancer is not desirable, and Kubernetes' load balancer is not available before a running a properly configured kube-proxy.

kube-discovery queries a set of seed API for a service definition and prints a randomly selected API server IP and port and returns with a zero exit status. If none of the seed servers responsed appropriately, or no functional API server was found in the service definition, the process returns with a non-zero exit status.

### License

[BSD 3-Clause License](LICENSE)
