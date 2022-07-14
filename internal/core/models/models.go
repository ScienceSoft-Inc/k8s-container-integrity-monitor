package models

import (
	"k8s.io/client-go/kubernetes"
)

type HashDataFromDB struct {
	ID             int
	Hash           string
	FileName       string
	FullFilePath   string
	Algorithm      string
	ImageContainer string
	NamePod        string
	NameDeployment string
}

type ConnectionDB struct {
	Dbdriver   string
	DbUser     string
	DbPassword string
	DbPort     string
	DbHost     string
	DbName     string
}

type KuberData struct {
	Clientset  *kubernetes.Clientset
	Namespace  string
	TargetName string
	TargetType string
}

type DeploymentData struct {
	Image                string
	NamePod              string
	Timestamp            string
	NameDeployment       string
	LabelMainProcessName string
}

type ConfigMapData struct {
	ProcName  string
	MountPath string
}

type DataFromK8sAPI struct {
	KuberData      *KuberData
	DeploymentData *DeploymentData
	ConfigMapData  *ConfigMapData
}
