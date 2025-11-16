package service

import (
	"context"
	"errors"
	"test-task-rest-api/internal/apperrors"
	"test-task-rest-api/internal/domain"
	"test-task-rest-api/internal/logger"
)

type TransactionRepo interface {
	Create(ctx context.Context, model domain.TransactionModel) (domain.TransactionModel, error)
	AllByUserId(ctx context.Context, userId int64) ([]domain.TransactionModel, error)
}

type TransactionService struct {
	logger          logger.Logger
	transactionRepo TransactionRepo
}

func NewTransactionService(logger logger.Logger, transactionRepo TransactionRepo) *TransactionService {
	return &TransactionService{
		logger:          logger,
		transactionRepo: transactionRepo,
	}
}

func (s *TransactionService) Create(ctx context.Context, userId int64, name, description string) (domain.Transaction, error) {
	const op = "transactionService.Create"

	statusDefault := domain.TransactionStatusCreated
	transactionModel, err := s.transactionRepo.Create(ctx, domain.TransactionModel{
		Name:        name,
		Description: description,
		UserID:      userId,
		Status:      statusDefault,
	})
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			return domain.Transaction{}, apperrors.Wrap(op, apperrors.ErrUserNotFound)
		}
		s.logger.Error(op, err)

		return domain.Transaction{}, apperrors.Wrap(op, err)
	}

	return domain.Transaction{
		ID:          transactionModel.ID,
		Name:        transactionModel.Name,
		Description: transactionModel.Description,
		UserID:      transactionModel.UserID,
		Status:      transactionModel.Status,
		CreatedAt:   transactionModel.CreatedAt,
		UpdatedAt:   transactionModel.UpdatedAt,
	}, nil
}

func (s *TransactionService) TransactionsByUserId(ctx context.Context, userId int64) ([]domain.Transaction, error) {
	const op = "transactionService.TransactionsByUserId"

	transactionModels, err := s.transactionRepo.AllByUserId(ctx, userId)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			return nil, apperrors.Wrap(op, apperrors.ErrUserNotFound)
		}
		s.logger.Error(op, err)

		return nil, apperrors.Wrap(op, err)
	}

	transactions := make([]domain.Transaction, len(transactionModels))
	for i, transactionModel := range transactionModels {
		transactions[i] = domain.Transaction{
			ID:          transactionModel.ID,
			Name:        transactionModel.Name,
			Description: transactionModel.Description,
			UserID:      transactionModel.UserID,
			Status:      transactionModel.Status,
			CreatedAt:   transactionModel.CreatedAt,
			UpdatedAt:   transactionModel.UpdatedAt,
		}
	}

	return transactions, nil
}
