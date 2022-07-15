package repositories

import (
	"fmt"
	"github.com/k8s-container-integrity-monitor/internal/core/models"
	"github.com/k8s-container-integrity-monitor/pkg/api"
	"github.com/sirupsen/logrus"
	"os"
)

type HashRepository struct {
	logger *logrus.Logger
}

func NewHashRepository(logger *logrus.Logger) *HashRepository {
	return &HashRepository{
		logger: logger,
	}
}

// SaveHashData iterates through all elements of the slice and triggers the save to database function
func (hr HashRepository) SaveHashData(allHashData []*api.HashData, deploymentData *models.DeploymentData) error {
	db, err := ConnectionToDB(hr.logger)
	if err != nil {
		hr.logger.Error("failed to connection to database %s", err)
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		hr.logger.Error("err while saving data in database ", err)
		return err
	}
	query := fmt.Sprintf(`
		INSERT INTO %s (file_name,full_file_path,hash_sum,algorithm,name_pod,image_tag,time_of_creation, name_deployment) 
		VALUES($1,$2,$3,$4,$5,$6,$7,$8) ON CONFLICT (full_file_path,algorithm) 
		DO UPDATE SET hash_sum=EXCLUDED.hash_sum`, os.Getenv("TABLE_NAME"))

	for _, hash := range allHashData {
		_, err = tx.Exec(query, hash.FileName, hash.FullFilePath, hash.Hash, hash.Algorithm, deploymentData.NamePod, deploymentData.Image, deploymentData.Timestamp, deploymentData.NameDeployment)
		if err != nil {
			err := tx.Rollback()
			if err != nil {
				hr.logger.Error("err in Rollback", err)
				return err
			}
			hr.logger.Error("err while save data in database ", err)
			return err
		}
	}

	return tx.Commit()
}

// GetHashData retrieves data from the database using the path and algorithm
func (hr HashRepository) GetHashData(dirFiles, algorithm string, deploymentData *models.DeploymentData) ([]*models.HashDataFromDB, error) {
	db, err := ConnectionToDB(hr.logger)
	if err != nil {
		hr.logger.Error("failed to connection to database %s", err)
		return nil, err
	}
	defer db.Close()

	var allHashDataFromDB []*models.HashDataFromDB

	query := fmt.Sprintf("SELECT id,file_name,full_file_path,hash_sum,algorithm,image_tag,name_pod,name_deployment FROM %s WHERE full_file_path LIKE $1 and algorithm=$2 and name_pod=$3", os.Getenv("TABLE_NAME"))

	rows, err := db.Query(query, "%"+dirFiles+"%", algorithm, deploymentData.NamePod)
	if err != nil {
		hr.logger.Error(err)
		return nil, err
	}
	for rows.Next() {
		var hashDataFromDB models.HashDataFromDB
		err := rows.Scan(&hashDataFromDB.ID, &hashDataFromDB.FileName, &hashDataFromDB.FullFilePath, &hashDataFromDB.Hash, &hashDataFromDB.Algorithm, &hashDataFromDB.ImageContainer, &hashDataFromDB.NamePod, &hashDataFromDB.NameDeployment)
		if err != nil {
			hr.logger.Error(err)
			return nil, err
		}
		allHashDataFromDB = append(allHashDataFromDB, &hashDataFromDB)
	}

	return allHashDataFromDB, nil
}

func (hr HashRepository) DeleteFromTable(nameDeployment string) error {
	db, err := ConnectionToDB(hr.logger)
	if err != nil {
		hr.logger.Error("failed to connection to database %s", err)
		return err
	}
	defer db.Close()

	query := fmt.Sprintf("DELETE FROM %s WHERE name_deployment=$1;", os.Getenv("TABLE_NAME"))
	_, err = db.Exec(query, nameDeployment)
	if err != nil {
		hr.logger.Error("err while deleting rows in database", err)
		return err
	}
	return nil
}
