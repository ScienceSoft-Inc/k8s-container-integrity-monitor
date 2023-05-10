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

This program provides integrity monitoring that checks the container's file system to determine if they have been maliciously modified. If the program detects that files have been modified, updated, added, or compromised, it rolls back the deployment to the previous version.
This application consists of two repositories: the [integrity-sum](https://github.com/ScienceSoft-Inc/integrity-sum) and the [integrity-mutator](https://github.com/ScienceSoft-Inc/k8s-container-integrity-mutator) .

Repository [integrity-sum](https://github.com/ScienceSoft-Inc/integrity-sum) injects hasher-sidecar into your modules as a sidecar container. sidecar integrity is a golang implementation of a hasher that calculates the checksum of files using various algorithms in kubernetes:
* MD5
* SHA256
* SHA1
* SHA224
* SHA384
* SHA512
* BEE2 (optional)

Repository [integrity-mutator](https://github.com/ScienceSoft-Inc/k8s-container-integrity-mutator) implements sidecar container for monitoring.

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
$ cd path/to/install
$ git clone https://github.com/ScienceSoft-Inc/k8s-container-integrity-monitor.git
```
Initialize and update submodules
```
$ git submodule init
$ git submodule update --remote
```

## :hammer: Installing components
### Install minikube
The code only works running inside a pod in Kubernetes.
You need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster.
If you do not already have a cluster, you can create one by using `minikube`.  
Example https://minikube.sigs.k8s.io/docs/start/

### Install Docker
To deploy app you need to install docker.  
Example https://docs.docker.com/engine/install/

### Install Helm
Before using helm charts you need to install helm on your local machine.  
You can find the necessary installation information at this link https://helm.sh/docs/intro/install/

### Configuration
To work properly, you first need to set the configuration files:
+ values in the file `helm-charts/database-to-integrity-sum/values.yaml`
+ values in the file `helm-charts/demo-apps-to-monitor/values.yaml`
+ values in the file `helm-charts/mutator/values.yaml`

## Manual start

1) You should go to the [README.md](https://github.com/ScienceSoft-Inc/k8s-container-integrity-mutator) in the `./integrity-mutator` project, set all the configurations and deploy.

2) You should go to the [README.md](https://github.com/ScienceSoft-Inc/integrity-sum) in the `./integrity-sum`project  project, set all the configurations and deploy.
   However, you need to go to the `"Run application"` section and configure the dependencies that are indicated by `"Need to install dependencies"`.

Install helm chart from the project root, for example:
```
helm install app helm-charts/demo-apps-to-monitor
```

## Quick start
### Using Makefile
You can use make function.  
```
make all
```

## Troubleshooting

Sometimes you may find that pod is injected with sidecar container as expected, check the following items:

1) The pod is in running state with `integrity` sidecar container injected and no error logs.
2) Check if the application pod has the correct annotations as described above.
___________________________

## License
This project uses the MIT software license. See [full license file](https://github.com/ScienceSoft-Inc/k8s-container-integrity-monitor/blob/main/LICENSE)
