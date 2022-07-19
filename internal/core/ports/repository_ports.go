package ports

import (
	"github.com/k8s-container-integrity-monitor/internal/core/models"
	"github.com/k8s-container-integrity-monitor/pkg/api"
)

//go:generate mockgen -source=repository_ports.go -destination=mocks/mock_repository.go

type IAppRepository interface {
	IsExistDeploymentNameInDB(deploymentName string) (bool, error)
}

type IHashRepository interface {
	SaveHashData(allHashData []*api.HashData, deploymentData *models.DeploymentData) error
	GetHashData(dirFiles string, algorithm string, deploymentData *models.DeploymentData) ([]*models.HashDataFromDB, error)
	DeleteFromTable(nameDeployment string) error
}
