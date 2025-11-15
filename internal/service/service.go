package service

import (
	"test-task-rest-api/internal/config"
	"test-task-rest-api/internal/logger"
)

type Deps struct {
	AuthUserRepo AuthUserRepo
}

type Service struct {
	config       *config.Config
	AuthService  *AuthService
	TokenService *TokenService
}

func NewService(config *config.Config, logger logger.Logger, deps Deps) *Service {
	tokenService := NewTokenService(config)
	return &Service{
		config:       config,
		TokenService: tokenService,
		AuthService:  NewAuthService(logger, deps.AuthUserRepo, tokenService),
	}
}
