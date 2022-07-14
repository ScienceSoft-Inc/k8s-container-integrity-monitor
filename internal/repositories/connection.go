package repositories

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/k8s-container-integrity-monitor/internal/core/models"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func InitializeDB(logger *logrus.Logger) (*sql.DB, error) {
	connectionDB := models.ConnectionDB{
		Dbdriver:   os.Getenv("DB_DRIVER"),
		DbUser:     os.Getenv("DB_USER"),
		DbPassword: os.Getenv("DB_PASSWORD"),
		DbPort:     os.Getenv("DB_PORT"),
		DbHost:     os.Getenv("DB_HOST"),
		DbName:     os.Getenv("DB_NAME"),
	}
	if connectionDB.Dbdriver == "" {
		logger.Error("DB_DRIVER was not set")
	}
	if connectionDB.DbHost == "" {
		logger.Error("DB_HOST was not set")
	}
	DBURL := fmt.Sprintf("host=%v port=%s user=%s dbname=%s sslmode=disable password=%s", connectionDB.DbHost, connectionDB.DbPort, connectionDB.DbUser, connectionDB.DbName, connectionDB.DbPassword)

	db, err := sql.Open(connectionDB.Dbdriver, DBURL)
	if err != nil {
		logger.Info("Cannot connect to database ", connectionDB.Dbdriver)
		return db, err
	} else {
		logger.Info("Connected to the database ", connectionDB.Dbdriver)
	}

	return db, nil
}
