# Sample-Broker
This broker follows [OSB](https://github.com/openservicebrokerapi/servicebroker/blob/v2.16/spec.md) specifications to a simpler manner. 
#### Note: This broker creates in-memory storage resources to handle the broker operations and not actual services.

## Contents
- [Local Set-up](#local-set-up)
- [Deploy on K8s](#deploy-on-k8s)
- [Register on svcat](#register-on-svcat)
  - [Install broker](#install-broker)
  - [Services](#services)
  - [Plans](#plans)
  - [Provision](#provision)
  - [Binding](#binding)
 
## Local Set-Up
- To run the broker in your system, make sure you have `go`(>=1.16.6) installed.
- Clone this repository and run the following commands:
```
go mod tidy
go run broker.go
```
This will install dependencies as per `go.mod` file and start the web-app running on port `8080`

## Deploy on K8s 
#### Pre-requisites
- Install Kubernetes locally using minikube or Docker-desktop. One can also use K8s by any cloud-provider.
- Install Kubectl and have your cluster ready to use.

### Steps to deploy on a k8s cluster
- Make sure you are in `sample-broker/k8s` directory and run the following commands:
```
kubectl apply -f deployment.yml
```
##### deployment.yml
```
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: go-broker-app
spec:
  replicas: 1
  selector:
    matchLabels:
      name: go-broker-app
  template:
    metadata:
      labels:
        name: go-broker-app
    spec:
      containers:
      - name: broker-app-container
        image: visargsoneji/go-simple-broker
        imagePullPolicy: IfNotPresent
        ports:
          - containerPort: 8080 # App listeining on this port inside container
        resources:
          limits:
            memory: 512Mi
            cpu: "1"
```
- This will create the `deployment` & `replicaset` resources. It will also start the pods.
- Now, to expose the pods add `service` by running:
```
kubectl apply -f service.yml
```
##### service.yml
```
---
apiVersion: v1
kind: Service
metadata:
  name: go-broker-service
spec:
  type: LoadBalancer
  ports:
  - name: http
    port: 80  # Actually use the app via k8s
    targetPort: 8080 # App listeining on this port inside container 
  selector:
    name: go-broker-app
```
- Get the external-ip and port of the app using:
``` 
kubectl get service/go-broker-service
```
Voila! Start hitting the OSB-API endpoints and monitor the logs inside pod.

## Register on svcat
#### Pre-requisites
- Follow this [guide](https://svc-cat.io/docs/install/) and make sure you have svcat installed along with its CLI.

### Install broker
- Make sure you are in `sample-broker/k8s/svcat` directory.
- Change the url in `broker.yml` where your broker is exposed and run the following command.<br>
PS: If you have deployed broker in the same cluster, its the ClusterIP.
```
kubectl apply -f broker.yml
```
##### broker.yml
```
apiVersion: servicecatalog.k8s.io/v1beta1
kind: ClusterServiceBroker
metadata:
  name: visargbroker4
spec:
        url: http://10.104.243.160:80
```
This will register the broker with svcat.
### Services
- To view services offered by all the registered brokers, use
```
svcat get classes
kubectl get clusterserviceclasses
```
### Plans
- To view various plans, use
```
svcat get plans
kubectl get clusterserviceplans
```
This will internally hit the `GET` `/v2/catalog` endpoint of the broker.
### Provision
- To Provision a service-instance, first we need to create a namespace and use it create instance in that namespace:
```
kubectl create namespace test-ns
```
```
kubectl apply -f instance.yml
```
##### instance.yml
```
apiVersion: servicecatalog.k8s.io/v1beta1
kind: ServiceInstance
metadata:
  name: instance55
  namespace: test-ns
spec:
  clusterServiceClassExternalName: redis
  clusterServicePlanExternalName: basic
```
This will internally hit the `PUT` `/v2/service_instances/:instance_id` & provision our instance with service and plan mentioned as per `instance.yml` 
- Check the status of the instance using:
```
svcat describe instance -n test-ns instance55
```
This will internally hit the `GET` `/v2/service_instances/:instance_id` 
### Binding
- To create a binding for your service instance using:
```
kubectl apply -f binding.yml
```
##### binding.yml
```
apiVersion: servicecatalog.k8s.io/v1beta1
kind: ServiceBinding
metadata:
  name: binding555
  namespace: test-ns
spec:
  instanceRef:
    name: instance55
```
This will internally hit `PUT` `/v2/service_instances/:instance_id/service_bindings/:binding_id` & create a binding to `instance55`.
- To update the plan of the instance, modify the new-plan in `instance.yml` and run `kubectl apply -f instance.yml` to get the changes.
- For more operations with svcat, use `svcat --help`.

#### Arigato!
