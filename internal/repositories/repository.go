package repositories

import (
	"fmt"
	"os"

	"github.com/k8s-container-integrity-monitor/internal/core/ports"
	"github.com/sirupsen/logrus"
)

type AppRepository struct {
	ports.IHashRepository
	logger *logrus.Logger
}

func NewAppRepository(logger *logrus.Logger) *AppRepository {
	return &AppRepository{
		IHashRepository: NewHashRepository(logger),
		logger:          logger,
	}
}

//CheckIsEmptyDB checks if the base is empty
func (ar AppRepository) IsExistDeploymentNameInDB(deploymentName string) (bool, error) {
	db, err := ConnectionToDB(ar.logger)
	if err != nil {
		ar.logger.Error("failed to connection to database %s", err)
		return false, err
	}
	defer db.Close()

	var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE name_deployment=$1 LIMIT 1;", os.Getenv("TABLE_NAME"))
	row := db.QueryRow(query, deploymentName)
	err = row.Scan(&count)
	if err != nil {
		ar.logger.Error("err while scan row in database ", err)
		return false, err
	}

	if count < 1 {
		return true, nil
	}
	return false, nil
}
