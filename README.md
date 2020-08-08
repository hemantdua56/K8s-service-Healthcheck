# K8s Health checker

[![N|Solid](https://img.icons8.com/color/48/000000/golang.png)]() Golang app

GLP health checker is an application to generate report on health status of all deployed microservices for the GLP environment.


  - UI is in Front-end Branch



### Tech

Healthchecker uses a number of open source projects to work properly:

* Go
* Docker




### Deployment

##### Docker
#
```shell
git clone $this_repo
# From the root project
docker build -t healthcheck:[latest|custom-version ] .
# Run
docker run -d -p 8080:8080 healthcheck:[latest|custom-version]
```

##### Kubernetes - Image available in ECR
#
```yaml
---
apiVersion: apps/v1beta1 # for versions before 1.8.0 use apps/v1beta1
kind: Deployment
metadata:
  name: healthcheck
spec:
  selector:
    matchLabels:
      app: healthcheck
      tier: backend
  replicas: 1
  template:
    metadata:
      labels:
        app: healthcheck
        tier: backend
    namespace: int 
    #Change the namespace accordingly
    spec:
      containers:
      - name: healthcheck
        image: healthcheck:latest
        resources:
          requests:
            cpu: 1000m
            memory: 2048Mi
        env:
        - secret: smtp-user
          value: smtp/smtp-user
        - secret: smtp-pass
          value: smtp/smtp-pass
        #For email notification add a secret smtp in environment with user smtp-user and password smtp-pass
        - name: alert
          value: false
        - name: namespace
          valueFrom:
            fieldRef:
              apiVersion: v1
              fieldPath: metadata.namespace
          # Passing Namespace in env is required by
          # the application. It uses namespace to get
          # pod details and use Incluster kubeconfig
        ports:
        - containerPort: 80
---

apiVersion: v1
kind: Service
metadata:
  name: healthcheck
  labels:
    app: healthcheck
    tier: backend
spec:
  # comment or delete the following line if you want to use a LoadBalancer
  # if your cluster supports it, uncomment the following to automatically create
  # an external load-balanced IP for the frontend service.
  # type: LoadBalancer
  ports:
  - port: 8080
  selector:
    app: healthcheck
    tier: backend

---

apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: healthcheck
  annotations:
    ingress.kubernetes.io/rewrite-target: /
spec:
  rules:
  - host: sample-app.poc.com
    http:
      paths:
      - path: /
        backend:
          serviceName: healthchecker
          servicePort: 80
```




