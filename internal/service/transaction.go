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
	ById(ctx context.Context, userId int64) (domain.TransactionModel, error)
	Update(ctx context.Context, transactionModel domain.TransactionModel) (domain.TransactionModel, error)
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

func (s *TransactionService) UpdateStatus(ctx context.Context, userId, transactionId int64, newStatus domain.TransactionStatus) (domain.Transaction, error) {
	const op = "transactionService.UpdateStatus"

	if newStatus == domain.TransactionStatusCreated {
		return domain.Transaction{}, apperrors.Wrap(op, apperrors.ErrInvalidTransactionStatusTransition)
	}

	transactionModel, err := s.transactionRepo.ById(ctx, transactionId)
	if err != nil {
		if errors.Is(err, apperrors.ErrTransactionNotFound) {
			return domain.Transaction{}, apperrors.Wrap(op, apperrors.ErrTransactionNotFound)
		}
		s.logger.Error(op, err)

		return domain.Transaction{}, apperrors.Wrap(op, err)
	}

	if transactionModel.UserID != userId {
		return domain.Transaction{}, apperrors.Wrap(op, apperrors.ErrAnotherUsersTransaction)
	}

	switch transactionModel.Status {
	case domain.TransactionStatusCancelled, domain.TransactionStatusCompleted, newStatus:
		return domain.Transaction{}, apperrors.Wrap(op, apperrors.ErrInvalidTransactionStatusTransition)
	case domain.TransactionStatusCreated:
		if newStatus != domain.TransactionStatusInProcess && newStatus != domain.TransactionStatusCancelled {
			return domain.Transaction{}, apperrors.Wrap(op, apperrors.ErrInvalidTransactionStatusTransition)
		}
	case domain.TransactionStatusInProcess:
		if newStatus != domain.TransactionStatusCompleted && newStatus != domain.TransactionStatusCancelled {
			return domain.Transaction{}, apperrors.Wrap(op, apperrors.ErrInvalidTransactionStatusTransition)
		}
	}

	transactionModel.Status = newStatus
	transactionModel, err = s.transactionRepo.Update(ctx, transactionModel)
	if err != nil {
		s.logger.Error(op, err)

		return domain.Transaction{}, apperrors.Wrap(op, err)
	}

	transaction := domain.Transaction{
		ID:          transactionModel.ID,
		Name:        transactionModel.Name,
		Description: transactionModel.Description,
		UserID:      transactionModel.UserID,
		Status:      transactionModel.Status,
		CreatedAt:   transactionModel.CreatedAt,
		UpdatedAt:   transactionModel.UpdatedAt,
	}

	return transaction, nil
}
