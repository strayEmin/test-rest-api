package main

import (
	"context"
	"log"
	"test-task-rest-api/internal/app"
	"test-task-rest-api/internal/config"
	"test-task-rest-api/internal/logger"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Println("Starting api")

	cfg := config.MustLoadConfig()

	appLogger := logger.NewLogger(*cfg)
	appLogger.InitLogger()
	appLogger.Info("Starting http server")

	a := app.NewApp(cfg, appLogger)

	a.Run(ctx)
}
