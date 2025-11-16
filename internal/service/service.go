package service

import (
	"test-task-rest-api/internal/config"
	"test-task-rest-api/internal/logger"
)

type Deps struct {
	UserRepo        AuthUserRepo
	TransactionRepo TransactionRepo
}

type Service struct {
	config             *config.Config
	AuthService        *AuthService
	TokenService       *TokenService
	TransactionService *TransactionService
}

func NewService(config *config.Config, logger logger.Logger, deps Deps) *Service {
	tokenService := NewTokenService(config)
	return &Service{
		config:             config,
		TokenService:       tokenService,
		AuthService:        NewAuthService(logger, deps.UserRepo, tokenService),
		TransactionService: NewTransactionService(logger, deps.TransactionRepo),
	}
}
