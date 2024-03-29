# K8s Health checker

[![N|Solid](https://img.icons8.com/color/48/000000/golang.png)]() Golang app

health checker is an application to generate report on health status of all deployed microservices for the K8s environments.


  - UI is in Front-end Branch



### Tech

Healthchecker uses a number of open source projects to work properly:

* Go
* Docker
* Helm


The Docker images are available in Docker Hub public Repository https://hub.docker.com/repository/docker/hemantdua/k8s-service-healthchecker

### Kubernetes Deployment

prerequisites

* Helm
* Access to K8s Cluster

### Follow the below steps

1. Create a values file as infra/helm/values/{ENV}.yaml and update the values as described in infra/helm/values/dev.yaml

2. Run below commands from the terminal from where you can access the K8s cluster
```shell
git clone $this_repo
# From the root project
helm upgrade --install  healthchecker -f infra/helm/values/{ENV}.yaml ./infra/helm/
```

