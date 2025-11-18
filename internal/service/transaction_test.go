package service

import (
	"context"
	"errors"
	"test-task-rest-api/internal/apperrors"
	"test-task-rest-api/internal/domain"
	"testing"
	"time"
)

type transactionRepoMock struct {
	createFn      func(ctx context.Context, model domain.TransactionModel) (domain.TransactionModel, error)
	allByUserIdFn func(ctx context.Context, userId int64) ([]domain.TransactionModel, error)
	byIdFn        func(ctx context.Context, userId int64) (domain.TransactionModel, error)
	updateFn      func(ctx context.Context, transactionModel domain.TransactionModel) (domain.TransactionModel, error)
}

func (r *transactionRepoMock) Create(ctx context.Context, model domain.TransactionModel) (domain.TransactionModel, error) {
	if r.createFn == nil {
		panic("createFn is nil")
	}
	return r.createFn(ctx, model)
}

func (r *transactionRepoMock) AllByUserId(ctx context.Context, userId int64) ([]domain.TransactionModel, error) {
	if r.allByUserIdFn == nil {
		panic("allByUserIdFn is nil")
	}
	return r.allByUserIdFn(ctx, userId)
}

func (r *transactionRepoMock) ById(ctx context.Context, id int64) (domain.TransactionModel, error) {
	if r.byIdFn == nil {
		panic("byIdFn is nil")
	}
	return r.byIdFn(ctx, id)
}

func (r *transactionRepoMock) Update(ctx context.Context, transactionModel domain.TransactionModel) (domain.TransactionModel, error) {
	if r.updateFn == nil {
		panic("updateFn is nil")
	}
	return r.updateFn(ctx, transactionModel)
}

type loggerTransactionMock struct{}

func (m *loggerTransactionMock) InitLogger()                                  {}
func (m *loggerTransactionMock) Debug(args ...interface{})                    {}
func (m *loggerTransactionMock) Debugf(template string, args ...interface{})  {}
func (m *loggerTransactionMock) Info(args ...interface{})                     {}
func (m *loggerTransactionMock) Infof(template string, args ...interface{})   {}
func (m *loggerTransactionMock) Warn(args ...interface{})                     {}
func (m *loggerTransactionMock) Warnf(template string, args ...interface{})   {}
func (m *loggerTransactionMock) Error(args ...interface{})                    {}
func (m *loggerTransactionMock) Errorf(template string, args ...interface{})  {}
func (m *loggerTransactionMock) DPanic(args ...interface{})                   {}
func (m *loggerTransactionMock) DPanicf(template string, args ...interface{}) {}
func (m *loggerTransactionMock) Fatal(args ...interface{})                    {}
func (m *loggerTransactionMock) Fatalf(template string, args ...interface{})  {}

func TestTransactionService_Create_Success(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userId := int64(1)
	name := "transaction 1"
	description := "transaction description"
	now := time.Now()
	repoMock := &transactionRepoMock{
		createFn: func(ctx context.Context, model domain.TransactionModel) (domain.TransactionModel, error) {
			return domain.TransactionModel{
				ID:          123,
				Name:        name,
				Description: description,
				UserID:      userId,
				Status:      domain.TransactionStatusCreated,
				CreatedAt:   now,
				UpdatedAt:   now,
			}, nil
		},
	}
	svc := NewTransactionService(&loggerTransactionMock{}, repoMock)

	transaction, err := svc.Create(ctx, userId, name, description)

	if err != nil {
		t.Fatal("transaction service create error: ", err)
	}
	if transaction.Name != name {
		t.Fatal("transaction.Name should be ", name)
	}
	if transaction.Description != description {
		t.Fatal("transaction.Description should be ", name)
	}
	if transaction.UserID != userId {
		t.Fatal("transaction.UserId should be ", userId)
	}
	if transaction.Status != domain.TransactionStatusCreated {
		t.Fatal("transaction.Status should be ", domain.TransactionStatusCreated)
	}
	if transaction.CreatedAt != now {
		t.Fatal("transaction.CreatedAt should be ", now)
	}
	if transaction.UpdatedAt != now {
		t.Fatal("transaction.UpdatedAt should be ", now)
	}
}

func TestTransactionService_Create_UserNotFound(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userId := int64(1)
	name := "transaction 1"
	description := "transaction description"
	repoMock := &transactionRepoMock{
		createFn: func(ctx context.Context, model domain.TransactionModel) (domain.TransactionModel, error) {
			return domain.TransactionModel{}, apperrors.ErrUserNotFound
		},
	}
	svc := NewTransactionService(&loggerTransactionMock{}, repoMock)

	_, err := svc.Create(ctx, userId, name, description)

	if err == nil {
		t.Fatal("error should be occurred")
	}
	if !errors.Is(err, apperrors.ErrUserNotFound) {
		t.Fatal("expected ErrUserNotFound, got ", err)
	}
}

func TestTransactionService_Create_UnexpectedRepoErr(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userId := int64(1)
	name := "transaction 1"
	description := "transaction description"
	repoMock := &transactionRepoMock{
		createFn: func(ctx context.Context, model domain.TransactionModel) (domain.TransactionModel, error) {
			return domain.TransactionModel{}, errors.New("could not connect to database")
		},
	}
	svc := NewTransactionService(&loggerTransactionMock{}, repoMock)

	_, err := svc.Create(ctx, userId, name, description)

	if err == nil {
		t.Fatal("error should be occurred")
	}
}

func TestTransactionService_TransactionsByUserId_Success(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userId := int64(1)
	updatedAt := time.Now()
	createdAt := updatedAt.Add(-5 * time.Second)
	transactionModels := []domain.TransactionModel{
		{
			ID:          123,
			Name:        "transaction 1",
			Description: "transaction description",
			UserID:      userId,
			Status:      domain.TransactionStatusCreated,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
		},
		{
			ID:          456,
			Name:        "transaction 2",
			Description: "transaction description",
			UserID:      userId,
			Status:      domain.TransactionStatusCreated,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
		},
	}
	repoMock := &transactionRepoMock{
		allByUserIdFn: func(ctx context.Context, userId int64) ([]domain.TransactionModel, error) {
			return transactionModels, nil
		},
	}
	svc := NewTransactionService(&loggerTransactionMock{}, repoMock)

	transactions, err := svc.TransactionsByUserId(ctx, userId)

	if err != nil {
		t.Fatal("transaction service transactionsByUserId error: ", err)
	}
	if len(transactions) != len(transactionModels) {
		t.Fatal("invalid transactions count: ", len(transactions))
	}
	for i, transaction := range transactions {
		if transaction.Name != transactionModels[i].Name {
			t.Fatal("transaction.Name should be ", transactionModels[i].Name)
		}
		if transaction.Description != transactionModels[i].Description {
			t.Fatal("transaction.Description should be ", transactionModels[i].Description)
		}
		if transaction.UserID != transactionModels[i].UserID {
			t.Fatal("transaction.UserId should be ", transactionModels[i].UserID)
		}
		if transaction.Status != transactionModels[i].Status {
			t.Fatal("transaction.Status should be ", transactionModels[i].Status)
		}
		if transaction.CreatedAt != transactionModels[i].CreatedAt {
			t.Fatal("transaction.CreatedAt should be ", transactionModels[i].CreatedAt)
		}
	}
}

func TestTransactionService_TransactionsByUserId_UserNotFound(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userId := int64(1)
	repoMock := &transactionRepoMock{
		allByUserIdFn: func(ctx context.Context, userId int64) ([]domain.TransactionModel, error) {
			return []domain.TransactionModel{}, apperrors.ErrUserNotFound
		},
	}
	svc := NewTransactionService(&loggerTransactionMock{}, repoMock)

	_, err := svc.TransactionsByUserId(ctx, userId)

	if err == nil {
		t.Fatal("error should be occurred")
	}
	if !errors.Is(err, apperrors.ErrUserNotFound) {
		t.Fatal("expected ErrUserNotFound, got ", err)
	}
}

func TestTransactionService_TransactionsByUserId_UnexpectedRepoErr(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userId := int64(1)
	repoMock := &transactionRepoMock{
		allByUserIdFn: func(ctx context.Context, userId int64) ([]domain.TransactionModel, error) {
			return []domain.TransactionModel{}, errors.New("could not connect to database")
		},
	}
	svc := NewTransactionService(&loggerTransactionMock{}, repoMock)

	_, err := svc.TransactionsByUserId(ctx, userId)

	if err == nil {
		t.Fatal("error should be occurred")
	}
}

func TestTransactionService_UpdateStatus_Success(t *testing.T) {
	ctx := context.Background()
	userId := int64(1)
	name := "test transaction"
	description := "test transaction description"
	updatedAt := time.Now()
	createdAt := updatedAt.Add(-5 * time.Second)
	tests := []struct {
		testDescription string
		transactionId   int64
		currentStatus   domain.TransactionStatus
		newStatus       domain.TransactionStatus
	}{
		{
			testDescription: "from created to inProcess",
			transactionId:   1,
			currentStatus:   domain.TransactionStatusCreated,
			newStatus:       domain.TransactionStatusInProcess,
		},
		{
			testDescription: "from created to canceled",
			transactionId:   2,
			currentStatus:   domain.TransactionStatusCreated,
			newStatus:       domain.TransactionStatusCancelled,
		},
		{
			testDescription: "from inProcess to completed",
			transactionId:   3,
			currentStatus:   domain.TransactionStatusInProcess,
			newStatus:       domain.TransactionStatusCompleted,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run("", func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(ctx, time.Second)
			t.Cleanup(cancel)
			repoMock := &transactionRepoMock{
				byIdFn: func(ctx context.Context, transactionId int64) (domain.TransactionModel, error) {
					return domain.TransactionModel{
						ID:          transactionId,
						Name:        name,
						Description: description,
						UserID:      userId,
						Status:      tt.currentStatus,
						CreatedAt:   createdAt,
						UpdatedAt:   createdAt,
					}, nil
				},
				updateFn: func(ctx context.Context, transactionModel domain.TransactionModel) (domain.TransactionModel, error) {
					return transactionModel, nil
				},
			}
			svc := NewTransactionService(&loggerTransactionMock{}, repoMock)

			transaction, err := svc.UpdateStatus(ctx, userId, tt.transactionId, tt.newStatus)
			if err != nil {
				t.Fatal(tt.testDescription, " transaction service updateStatus error: ", err)
			}
			if transaction.ID != tt.transactionId {
				t.Fatal(tt.testDescription, " transaction id should be ", tt.transactionId)
			}
			if transaction.Name != name {
				t.Fatal(tt.testDescription, " transaction name should be ", name)
			}
			if transaction.Description != description {
				t.Fatal(tt.testDescription, " transaction.Description should be ", description)
			}
			if transaction.UserID != userId {
				t.Fatal(tt.testDescription, " transaction.UserId should be ", userId)
			}
			if transaction.Status != tt.newStatus {
				t.Fatal(tt.testDescription, " transaction.Status should be ", tt.newStatus)
			}
		})
	}
}

func TestTransactionService_UpdateStatus_InvalidTransactionStatusTransition(t *testing.T) {
	ctx := context.Background()

	userId := int64(1)
	name := "test transaction"
	description := "test transaction description"
	updatedAt := time.Now()
	createdAt := updatedAt.Add(-5 * time.Second)
	tests := []struct {
		testDescription string
		transactionId   int64
		currentStatus   domain.TransactionStatus
		newStatus       domain.TransactionStatus
	}{
		{
			testDescription: "from created to created",
			transactionId:   1,
			currentStatus:   domain.TransactionStatusCreated,
			newStatus:       domain.TransactionStatusCreated,
		},
		{
			testDescription: "from created to completed",
			transactionId:   2,
			currentStatus:   domain.TransactionStatusCreated,
			newStatus:       domain.TransactionStatusCompleted,
		},
		{
			testDescription: "from inProcess to cancelled",
			transactionId:   3,
			currentStatus:   domain.TransactionStatusInProcess,
			newStatus:       domain.TransactionStatusCancelled,
		},
		{
			testDescription: "from inProcess to inProcess",
			transactionId:   4,
			currentStatus:   domain.TransactionStatusInProcess,
			newStatus:       domain.TransactionStatusInProcess,
		},
		{
			testDescription: "from cancelled to cancelled",
			transactionId:   5,
			currentStatus:   domain.TransactionStatusCancelled,
			newStatus:       domain.TransactionStatusCancelled,
		},
		{
			testDescription: "from cancelled to created",
			transactionId:   6,
			currentStatus:   domain.TransactionStatusCancelled,
			newStatus:       domain.TransactionStatusCreated,
		},
		{
			testDescription: "from cancelled to inProcess",
			transactionId:   7,
			currentStatus:   domain.TransactionStatusCancelled,
			newStatus:       domain.TransactionStatusInProcess,
		},
		{
			testDescription: "from cancelled to completed",
			transactionId:   8,
			currentStatus:   domain.TransactionStatusCancelled,
			newStatus:       domain.TransactionStatusCompleted,
		},
		{
			testDescription: "from completed to cancelled",
			transactionId:   9,
			currentStatus:   domain.TransactionStatusCompleted,
			newStatus:       domain.TransactionStatusCancelled,
		},
		{
			testDescription: "from completed to inProcess",
			transactionId:   10,
			currentStatus:   domain.TransactionStatusCompleted,
			newStatus:       domain.TransactionStatusInProcess,
		},
		{
			testDescription: "from inProcess to created",
			transactionId:   11,
			currentStatus:   domain.TransactionStatusInProcess,
			newStatus:       domain.TransactionStatusCreated,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run("", func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(ctx, time.Second)
			t.Cleanup(cancel)

			repoMock := &transactionRepoMock{
				byIdFn: func(ctx context.Context, transactionId int64) (domain.TransactionModel, error) {
					return domain.TransactionModel{
						ID:          transactionId,
						Name:        name,
						Description: description,
						UserID:      userId,
						Status:      tt.currentStatus,
						CreatedAt:   createdAt,
						UpdatedAt:   createdAt,
					}, nil
				},
				updateFn: func(ctx context.Context, transactionModel domain.TransactionModel) (domain.TransactionModel, error) {
					return transactionModel, nil
				},
			}
			svc := NewTransactionService(&loggerTransactionMock{}, repoMock)

			_, err := svc.UpdateStatus(ctx, userId, tt.transactionId, tt.newStatus)
			if err == nil {
				t.Fatal(tt.testDescription, " transaction service error should be occurred")
			}
			if !errors.Is(err, apperrors.ErrInvalidTransactionStatusTransition) {
				t.Fatal(tt.testDescription, " expected ErrInvalidTransactionStatusTransition, got ", err)
			}
		})
	}
}

func TestTransactionService_UpdateStatus_TransactionNotFound(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userId := int64(1)
	transactionId := int64(2)
	newStatus := domain.TransactionStatusCancelled
	repoMock := &transactionRepoMock{
		byIdFn: func(ctx context.Context, transactionId int64) (domain.TransactionModel, error) {
			return domain.TransactionModel{}, apperrors.ErrTransactionNotFound
		},
		updateFn: func(ctx context.Context, transactionModel domain.TransactionModel) (domain.TransactionModel, error) {
			return domain.TransactionModel{}, apperrors.ErrTransactionNotFound
		},
	}
	svc := NewTransactionService(&loggerTransactionMock{}, repoMock)

	_, err := svc.UpdateStatus(ctx, userId, transactionId, newStatus)
	if err == nil {
		t.Fatal("transaction service error should be occurred")
	}
	if !errors.Is(err, apperrors.ErrTransactionNotFound) {
		t.Fatal("expected ErrTransactionNotFound, got ", err)
	}
}

func TestTransactionService_UpdateStatus_UnexpectedRepoErr(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userId := int64(1)
	transactionId := int64(2)
	newStatus := domain.TransactionStatusCompleted
	repoMock := &transactionRepoMock{
		byIdFn: func(ctx context.Context, transactionId int64) (domain.TransactionModel, error) {
			return domain.TransactionModel{}, errors.New("could not connect to database")
		},
		updateFn: func(ctx context.Context, transactionModel domain.TransactionModel) (domain.TransactionModel, error) {
			return domain.TransactionModel{}, errors.New("could not connect to database")
		},
	}
	svc := NewTransactionService(&loggerTransactionMock{}, repoMock)

	_, err := svc.UpdateStatus(ctx, userId, transactionId, newStatus)

	if err == nil {
		t.Fatal("transaction service error should be occurred")
	}
}

func TestNewTransactionService_UpdateStatus_AnotherUsersTransaction(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userId := int64(1)
	transactionId := int64(2)
	name := "test transaction"
	description := "test transaction description"
	currentStatus := domain.TransactionStatusInProcess
	newStatus := domain.TransactionStatusCompleted
	updatedAt := time.Now()
	createdAt := updatedAt.Add(-5 * time.Second)
	repoMock := &transactionRepoMock{
		byIdFn: func(ctx context.Context, transactionId int64) (domain.TransactionModel, error) {
			return domain.TransactionModel{
				ID:          transactionId,
				Name:        name,
				Description: description,
				UserID:      userId + 1,
				Status:      currentStatus,
				CreatedAt:   createdAt,
				UpdatedAt:   updatedAt,
			}, nil
		},
		updateFn: func(ctx context.Context, transactionModel domain.TransactionModel) (domain.TransactionModel, error) {
			return transactionModel, nil
		},
	}
	svc := NewTransactionService(&loggerTransactionMock{}, repoMock)

	_, err := svc.UpdateStatus(ctx, userId, transactionId, newStatus)
	if err == nil {
		t.Fatal("transaction service error should be occurred")
	}
	if !errors.Is(err, apperrors.ErrAnotherUsersTransaction) {
		t.Fatal("expected ErrAnotherUsersTransaction, got ", err)
	}
}
