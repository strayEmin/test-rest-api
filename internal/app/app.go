package app

import (
	"context"
	"test-task-rest-api/internal/config"
	"test-task-rest-api/internal/logger"
	"test-task-rest-api/internal/repo"
	"test-task-rest-api/internal/repo/postgres"
	"test-task-rest-api/internal/service"
	"test-task-rest-api/internal/transport/rest"
	"test-task-rest-api/internal/transport/rest/handler"
)

type App struct {
	cfg *config.Config
	log logger.Logger
}

func NewApp(cfg *config.Config, log logger.Logger) *App {
	return &App{
		cfg: cfg,
		log: log,
	}
}

func (a *App) Run(ctx context.Context) {
	a.log.Info("Starting app")

	a.log.Info("DB initializing")
	pool, err := postgres.NewPostgres(ctx, 3, a.cfg, a.log)
	if err != nil {
		a.log.Fatal("Error initializing DB pool: %w", err)
	}
	defer pool.Close()

	repos := repo.NewRepo(a.log, pool)

	services := service.NewService(a.cfg, a.log, service.Deps{
		UserRepo:        repos.UserRepo,
		TransactionRepo: repos.TransactionRepo,
	})

	restHandlers := handler.NewHandler(a.cfg, handler.Deps{
		AuthService:        services.AuthService,
		TransactionService: services.TransactionService,
	})
	restApp := restHandlers.Init(ctx)

	a.log.Info("Starting rest app")
	restSrv := rest.NewServer(restApp, a.cfg.HTTPServer.Address, a.log)

	a.log.Info("Starting http server")
	restSrv.StartWithGracefulShutdown()
}
