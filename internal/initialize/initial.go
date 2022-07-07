package initialize

import (
	"context"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/k8s-container-integrity-monitor/internal/core/services"
	"github.com/k8s-container-integrity-monitor/internal/repositories"
	"github.com/sirupsen/logrus"
)

func Initialize(ctx context.Context, logger *logrus.Logger, sig chan os.Signal) {
	// Initialize database
	logger.Info("Starting database connection")
	db, err := repositories.InitializeDB(logger)
	if err != nil {
		logger.Fatalf("failed to connection to database %s", err)
	}

	// Initialize repository
	repository := repositories.NewAppRepository(db, logger)

	// Initialize service
	algorithm := os.Getenv("AlGORITHM")

	service, err := services.NewAppService(repository, algorithm, logger)
	if err != nil {
		logger.Fatalf("can't init service: %s", err)
	}

	// Initialize kubernetesAPI
	dataK8sAPI, deploymentDataK8s, configDataK8s, err := service.GetDataFromK8sAPI()
	if err != nil {
		logger.Fatalf("can't get data from K8sAPI: %s", err)
	}

	//Getting pid
	pid, err := service.GetPID(configDataK8s)
	if err != nil {
		logger.Fatalf("err while getting pid %s", err)
	}
	if pid == 0 {
		logger.Fatalf("proc with name %s not exist", configDataK8s.ProcName)
	}

	//Getting current directory path
	dirPath := "../proc/" + strconv.Itoa(pid) + "/root/" + configDataK8s.MountPath

	duration, err := strconv.Atoi(os.Getenv("DURATION_TIME"))
	if err != nil {
		duration = 15
	}
	ticker := time.NewTicker(time.Duration(duration) * time.Second)

	var wg sync.WaitGroup
	wg.Add(1)
	go func(ctx context.Context, ticker *time.Ticker) {
		defer wg.Done()
		for {
			if service.IsExistDeploymentNameInDB(dataK8sAPI.TargetName) {
				logger.Info("Deployment name does not exist in database, save data")
				err := service.Start(ctx, dirPath, sig, deploymentDataK8s)
				if err != nil {
					logger.Fatalf("Error when starting to get and save hash data %s", err)
				}
			} else {
				logger.Info("Deployment name exists in database, checking data")
				for range ticker.C {
					err := service.Check(ctx, dirPath, sig, deploymentDataK8s, dataK8sAPI)
					if err != nil {
						logger.Fatalf("Error when starting to check hash data %s", err)
					}
					logger.Info("Check completed")
				}
			}
		}
	}(ctx, ticker)
	wg.Wait()
	ticker.Stop()
}
