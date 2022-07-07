package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
	config "github.com/k8s-container-integrity-monitor/internal/configs"
	"github.com/k8s-container-integrity-monitor/internal/initialize"
)

func main() {
	//Load values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	//Initialize config
	_, logger, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Error during loading from config file", err)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		signal.Stop(sig)
		cancel()
	}()

	initialize.Initialize(ctx, logger, sig)
}
