![GitHub contributors](https://img.shields.io/github/contributors/ScienceSoft-Inc/k8s-container-integrity-monitor)
![GitHub last commit](https://img.shields.io/github/last-commit/ScienceSoft-Inc/k8s-container-integrity-monitor)
![GitHub](https://img.shields.io/github/license/ScienceSoft-Inc/k8s-container-integrity-monitor)
![GitHub issues](https://img.shields.io/github/issues/ScienceSoft-Inc/k8s-container-integrity-monitor)
![GitHub forks](https://img.shields.io/github/forks/ScienceSoft-Inc/k8s-container-integrity-monitor)

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
### Running locally
The code only works running inside a pod in Kubernetes
You need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster.
If you do not already have a cluster, you can create one by using `minikube`.  
Example https://minikube.sigs.k8s.io/docs/start/

### Configuration

## :hammer: Installing components
```
docker build -t mutator ./k8s-container-integrity-mutator
helm intall mutator

make integrity sum

helm install demo-apps-to-monitor
```
### Install Helm
Before using helm charts you need to install helm on your local machine.  
You can find the necessary installation information at this link https://helm.sh/docs/intro/install/

Then update the on-disk dependencies to mirror Chart.yaml.
```
helm dependency update helm-charts/database-to-integrity-sum
```
INSERT TEXT HERE
```
helm install helm-charts/database-to-integrity-sum
```
INSERT TEXT HERE
```
helm install db helm-charts/mutator
```
INSERT TEXT HERE
```
helm install app helm-charts/demo-apps-to-monitor
```

## Quick start

## Troubleshooting
___________________________
### :notebook_with_decorative_cover: Godoc extracts and generates documentation for Go programs
#### Presents the documentation as a web page.
```go
godoc -http=:6060/sha256sum
go doc packge.function_name
```
for example
```go
go doc pkg/api.Result
```

### :mag: Running tests

You need to go to the folder where the file is located *_test.go and run the following command:
```go
go test -v
```

for example
```go
cd ../pkg/api
go test -v
```

##License
This project uses the MIT software license. See [full license file](https://github.com/ScienceSoft-Inc/k8s-container-integrity-monitor/blob/main/LICENSE)
