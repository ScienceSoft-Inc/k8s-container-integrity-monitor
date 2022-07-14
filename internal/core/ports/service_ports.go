package ports

import (
	"context"
	"github.com/k8s-container-integrity-monitor/internal/core/models"
	"github.com/k8s-container-integrity-monitor/pkg/api"
	"os"
	"sync"
)

//go:generate mockgen -source=service_ports.go -destination=mocks/mock_service.go

type IAppService interface {
	GetPID(configData *models.ConfigMapData) (int, error)
	IsExistDeploymentNameInDB(deploymentName string) bool
	LaunchHasher(ctx context.Context, dirPath string, sig chan os.Signal) []api.HashData
	Start(ctx context.Context, dirPath string, sig chan os.Signal, deploymentData *models.DeploymentData) error
	Check(ctx context.Context, dirPath string, sig chan os.Signal, deploymentData *models.DeploymentData, kuberData *models.KuberData) error
}

type IHashService interface {
	SaveHashData(allHashData []api.HashData, deploymentData *models.DeploymentData) error
	GetHashData(dirPath string, deploymentData *models.DeploymentData) ([]models.HashDataFromDB, error)
	DeleteFromTable(nameDeployment string) error
	IsDataChanged(currentHashData []api.HashData, hashSumFromDB []models.HashDataFromDB, deploymentData *models.DeploymentData) (bool, error)
	CreateHash(path string) (api.HashData, error)
	WorkerPool(ctx context.Context, jobs chan string, results chan api.HashData)
	Worker(ctx context.Context, wg *sync.WaitGroup, jobs <-chan string, results chan<- api.HashData)
}

type IKuberService interface {
	GetDataFromK8sAPI() (*models.DataFromK8sAPI, error)
	ConnectionToK8sAPI() (*models.KuberData, error)
	GetDataFromDeployment(kuberData *models.KuberData) (*models.DeploymentData, error)
	GetDataFromConfigMap(kuberData *models.KuberData, label string) (*models.ConfigMapData, error)
	RolloutDeployment(kuberData *models.KuberData) error
}
