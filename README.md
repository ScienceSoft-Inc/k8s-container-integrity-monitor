[![GitHub contributors](https://img.shields.io/github/contributors/ScienceSoft-Inc/k8s-container-integrity-monitor)](https://github.com/ScienceSoft-Inc/k8s-container-integrity-monitor)
[![GitHub last commit](https://img.shields.io/github/last-commit/ScienceSoft-Inc/k8s-container-integrity-monitor)](https://github.com/ScienceSoft-Inc/k8s-container-integrity-monitor)
[![GitHub](https://img.shields.io/github/license/ScienceSoft-Inc/k8s-container-integrity-monitor)](https://github.com/ScienceSoft-Inc/k8s-container-integrity-monitor/blob/main/LICENSE)
[![GitHub issues](https://img.shields.io/github/issues/ScienceSoft-Inc/k8s-container-integrity-monitor)](https://github.com/ScienceSoft-Inc/k8s-container-integrity-monitor/issues)
[![GitHub forks](https://img.shields.io/github/forks/ScienceSoft-Inc/k8s-container-integrity-monitor)](https://github.com/ScienceSoft-Inc/k8s-container-integrity-monitor/network/members)

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![Kubernetes](https://img.shields.io/badge/kubernetes-%23326ce5.svg?style=for-the-badge&logo=kubernetes&logoColor=white)
![Postgres](https://img.shields.io/badge/postgres-%23316192.svg?style=for-the-badge&logo=postgresql&logoColor=white)
![GitHub](https://img.shields.io/badge/github-%23121011.svg?style=for-the-badge&logo=github&logoColor=white)

# k8s-container-integrity-monitor

This program provides integrity monitoring that checks file or   directory of container to determine whether or not they have been tampered with or corrupted.  
k8s-container-integrity-monitor, which is a type of change auditing, verifies and validates these files by comparing them to the stored data in the database.

If program detects that files have been altered, updated, added or compromised, it rolls back deployments to a previous version.

k8s-container-integrity-monitor injects a `hasher-webhook-injector-sidecar` to your pods with the "hasher-webhook-injector-sidecar" label.  
`hasher-webhook-injector-sidecar` the implementation of a hasher in golang, which calculates the checksum of files using different algorithms in kubernetes:
* MD5
* SHA256
* SHA1
* SHA224
* SHA384
* SHA512

## Architecture
### Component diagram
![File location: docs/diagrams/projectStructure.png](/docs/diagrams/projectStructure.png?raw=true "Component diagram")
### Activity diagram
![File location: docs/diagrams/deployDiagram.png](/docs/diagrams/deployDiagram.png?raw=true "Activity diagram")
### Statechart diagram
![File location: docs/diagrams/appStatechartDiagram.png](/docs/diagrams/appStatechartDiagram.png?raw=true "Statechart diagram")
### Sequence diagram
![File location: docs/diagrams/appSequenceDiagram.png](/docs/diagrams/appSequenceDiagram.png?raw=true "Sequence diagram")
## Getting Started

### Clone repository and install dependencies
```
$ git clone https://github.com/ScienceSoft-Inc/k8s-container-integrity-monitor.git
$ cd path/to/install
```
Download the named modules into the module cache
```
go mod download
```

## :hammer: Installing components
### Running locally
The code only works running inside a pod in Kubernetes.
You need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster.
If you do not already have a cluster, you can create one by using `minikube`.  
Example https://minikube.sigs.k8s.io/docs/start/

### Install Helm
Before using helm charts you need to install helm on your local machine.  
You can find the necessary installation information at this link https://helm.sh/docs/intro/install/

### Configuration
To work properly, you first need to set the configuration files:
+ values in the file `helm-charts/database-to-integrity-sum/values.yaml`
+ values in the file `helm-charts/demo-apps-to-monitor/values.yaml`
+ values in the file `helm-charts/mutator/values.yaml`

## Manual start
+ Minikube start
```
minikube start
```
1) You should go to the `README.md` in the `./k8s-container-integrity-mutator` project and set all the settings and certificates.

Build docker images mutator:
```
eval $(minikube docker-env)
cd k8s-container-integrity-mutator
docker build -t mutator
```
or
```
eval $(minikube docker-env)
docker build -t mutator -f k8s-container-integrity-mutator/Dockerfile .
```
Install helm chart:
for example
```
helm install mutator helm-charts/mutator
```
2) You need to install the database using helm charts.  
   Update the on-disk dependencies to mirror Chart.yaml.
```
helm dependency update helm-charts/database-to-integrity-sum
```
Install helm chart:
for example
```
helm install db helm-charts/database-to-integrity-sum
```

3) You should go to the `README.md` in the `./integrity-sum` project and set all the settings.
   Build docker images mutator:
```
eval $(minikube docker-env)
cd integrity-sum
docker build -t hasher
```
or
```
eval $(minikube docker-env)
docker build -t hasher -f integrity-sum/Dockerfile .
```
Install helm chart:
for example
```
helm install app helm-charts/demo-apps-to-monitor
```

## Quick start
### Using Makefile

## Troubleshooting
Sometimes you may find that pod is injected with sidecar container as expected, check the following items:

1) The pod is in running state with `hasher-sidecar` sidecar container injected and no error logs.
2) Check if the application demo-pod has he correct labels `hasher-certificates-injector-sidecar: "true"` and installed `main-process-name`.
___________________________
### :mag: Running tests

You need to go to the folder where the file is located *_test.go and run the following command:
```go
go test -v ./..
```

## License
This project uses the MIT software license. See [full license file](https://github.com/ScienceSoft-Inc/k8s-container-integrity-monitor/blob/main/LICENSE)
