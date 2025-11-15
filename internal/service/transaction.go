package service

import (
	"context"
	"test-task-rest-api/internal/domain"
	"test-task-rest-api/internal/logger"
)

type TransactionRepo interface {
	Create(ctx context.Context, userId int64, name, description string) (domain.TransactionModel, error)
}

type TransactionUserRepo interface {
	ById(ctx context.Context, userId int64) (domain.TransactionModel, error)
}

type TransactionService struct {
	logger          logger.Logger
	transactionRepo TransactionRepo
	userRepo        TransactionUserRepo
}

func NewTransactionService(logger logger.Logger, transactionRepo TransactionRepo, userRepo TransactionUserRepo) *TransactionService {
	return &TransactionService{
		logger:          logger,
		transactionRepo: transactionRepo,
		userRepo:        userRepo,
	}
}

func (s *TransactionService) Create(ctx context.Context, userId int64, name, description string) (domain.Transaction, error) {
	const op = "transactionService.Create"

}
