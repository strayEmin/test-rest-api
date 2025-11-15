package repo

import (
	"test-task-rest-api/internal/logger"
	"test-task-rest-api/internal/repo/postgres"
)

type Repo struct {
	Logger   logger.Logger
	UserRepo *UserRepo
}

func NewRepo(logger logger.Logger, pool postgres.AtomicPoolClient) *Repo {
	return &Repo{
		Logger:   logger,
		UserRepo: NewUserRepo(logger, pool),
	}
}
