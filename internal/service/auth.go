package service

import (
	"context"
	"errors"
	"test-task-rest-api/internal/apperrors"
	"test-task-rest-api/internal/domain"
	"test-task-rest-api/internal/logger"
	"test-task-rest-api/internal/utils/hasher"
)

type AuthUserRepo interface {
	CreateUser(ctx context.Context, user domain.UserModel) (domain.UserModel, error)
	ByEmail(ctx context.Context, email string) (domain.UserModel, error)
}

type TokenGenerator interface {
	GenerateToken(user domain.User) (string, error)
}

type AuthService struct {
	logger       logger.Logger
	tokenService TokenGenerator
	authUserRepo AuthUserRepo
}

func NewAuthService(logger logger.Logger, authUserRepo AuthUserRepo, tokenService TokenGenerator) *AuthService {
	return &AuthService{
		logger:       logger,
		authUserRepo: authUserRepo,
		tokenService: tokenService,
	}
}

func (s *AuthService) SignUp(ctx context.Context, name, email, password string) (string, error) {
	const op = "AuthService.SignUp"

	hashedPassword, err := hasher.HashPassword(password)
	if err != nil {
		return "", apperrors.Wrap(op, apperrors.ErrHashingPasswordFailed)
	}

	userModel, err := s.authUserRepo.CreateUser(ctx, domain.UserModel{
		Name:    name,
		Email:   email,
		PwdHash: hashedPassword,
	})
	if err != nil {
		if errors.Is(err, apperrors.ErrUserAlreadyExists) {
			return "", apperrors.Wrap(op, apperrors.ErrUserAlreadyExists)
		}
		s.logger.Error(op, err)

		return "", apperrors.Wrap(op, err)
	}

	token, err := s.tokenService.GenerateToken(domain.User{
		ID:      userModel.ID,
		Name:    userModel.Name,
		Email:   userModel.Email,
		PwdHash: hashedPassword,
	})
	if err != nil {
		s.logger.Error(op, err)

		return "", apperrors.Wrap(op, apperrors.ErrUnsuccesffulTokenGeneration)
	}

	return token, nil
}

func (s *AuthService) SignIn(ctx context.Context, email, password string) (string, error) {
	const op = "AuthService.SignIn"

	userModel, err := s.authUserRepo.ByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			return "", apperrors.Wrap(op, apperrors.ErrUserNotFound)
		}
		s.logger.Error(op, err)

		return "", apperrors.Wrap(op, err)
	}

	if !hasher.IsValidPassword(userModel.PwdHash, password) {
		return "", apperrors.Wrap(op, apperrors.ErrWrongPassword)
	}

	token, err := s.tokenService.GenerateToken(domain.User{
		ID:      userModel.ID,
		Name:    userModel.Name,
		Email:   userModel.Email,
		PwdHash: userModel.PwdHash,
	})
	if err != nil {
		s.logger.Error(op, err)

		return "", apperrors.Wrap(op, apperrors.ErrUnsuccesffulTokenGeneration)
	}

	return token, nil
}
