package ports

import (
	"context"

	"github.com/k8s-container-integrity-monitor/internal/core/models"
	"github.com/k8s-container-integrity-monitor/pkg/api"
)

//go:generate mockgen -source=repository_ports.go -destination=mocks/mock_repository.go

type IAppRepository interface {
	IsExistDeploymentNameInDB(deploymentName string) (bool, error)
}

type IHashRepository interface {
	SaveHashData(ctx context.Context, allHashData []api.HashData, deploymentData models.DeploymentData) error
	GetHashData(ctx context.Context, dirFiles string, algorithm string, deploymentData models.DeploymentData) ([]models.HashDataFromDB, error)
	DeleteFromTable(nameDeployment string) error
}
