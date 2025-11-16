package repo

import (
	"context"
	"errors"
	"test-task-rest-api/internal/apperrors"
	"test-task-rest-api/internal/domain"
	"test-task-rest-api/internal/repo/postgres"

	"github.com/jackc/pgx/v4"
)

type TransactionRepo struct {
	pool postgres.AtomicPoolClient
}

func NewTransactionRepo(pool postgres.AtomicPoolClient) *TransactionRepo {
	return &TransactionRepo{
		pool: pool,
	}
}

func (r *TransactionRepo) Create(ctx context.Context, transactionModel domain.TransactionModel) (domain.TransactionModel, error) {
	const op = "TransactionRepo.Create"
	q := `INSERT INTO transactions (name, description, user_id, status) VALUES ($1, $2, $3, $4)
		RETURNING id, name, description, user_id, status, created_at, updated_at`

	if err := r.pool.QueryRow(
		ctx,
		q,
		transactionModel.Name,
		transactionModel.Description,
		transactionModel.UserID,
		transactionModel.Status,
	).Scan(
		&transactionModel.ID,
		&transactionModel.Name,
		&transactionModel.Description,
		&transactionModel.UserID,
		&transactionModel.Status,
		&transactionModel.CreatedAt,
		&transactionModel.UpdatedAt,
	); err != nil {
		if postgres.IsFKViolation(err) {
			return domain.TransactionModel{}, apperrors.Wrap(op, apperrors.ErrUserNotFound)
		}
		err = postgres.ParsePgError(err)

		return domain.TransactionModel{}, apperrors.Wrap(op, err)
	}

	return transactionModel, nil
}

func (r *TransactionRepo) AllByUserId(ctx context.Context, userId int64) ([]domain.TransactionModel, error) {
	const op = "TransactionRepo.AllByUserId"
	q := `SELECT id, name, description, user_id, status, created_at, updated_at FROM transactions WHERE user_id = $1`

	transactionModels := make([]domain.TransactionModel, 0)

	rows, err := r.pool.Query(ctx, q, userId)
	if err != nil {
		err = apperrors.Wrap(op, err)
	}

	defer rows.Close()

	for rows.Next() {
		transactionModel := domain.TransactionModel{}

		if err := rows.Scan(
			&transactionModel.ID,
			&transactionModel.Name,
			&transactionModel.Description,
			&transactionModel.UserID,
			&transactionModel.Status,
			&transactionModel.CreatedAt,
			&transactionModel.UpdatedAt,
		); err != nil {
			err = postgres.ParsePgError(err)

			return transactionModels, apperrors.Wrap(op, err)
		}

		transactionModels = append(transactionModels, transactionModel)
	}

	return transactionModels, nil
}

func (r *TransactionRepo) ById(ctx context.Context, userId int64) (domain.TransactionModel, error) {
	const op = "TransactionRepo.ById"
	q := `SELECT id, name, description, user_id, status, created_at, updated_at FROM transactions WHERE id = $1`

	transactionModel := domain.TransactionModel{}
	if err := r.pool.QueryRow(
		ctx,
		q,
		userId,
	).Scan(
		&transactionModel.ID,
		&transactionModel.Name,
		&transactionModel.Description,
		&transactionModel.UserID,
		&transactionModel.Status,
		&transactionModel.CreatedAt,
		&transactionModel.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.TransactionModel{}, apperrors.Wrap(op, apperrors.ErrTransactionNotFound)
		}
		err = postgres.ParsePgError(err)

		return domain.TransactionModel{}, apperrors.Wrap(op, err)
	}

	return transactionModel, nil
}

func (r *TransactionRepo) Update(ctx context.Context, transactionModel domain.TransactionModel) (domain.TransactionModel, error) {
	const op = "TransactionRepo.Update"
	q := `
	UPDATE transactions
	SET
		name = $1,
		description = $2,
		user_id = $3,
		status = $4,
		updated_at = NOW()
	WHERE id = $5
	RETURNING id, name, description, user_id, status, created_at, updated_at
	`

	if err := r.pool.QueryRow(
		ctx,
		q,
		transactionModel.Name,
		transactionModel.Description,
		transactionModel.UserID,
		transactionModel.Status,
		transactionModel.ID,
	).Scan(
		&transactionModel.ID,
		&transactionModel.Name,
		&transactionModel.Description,
		&transactionModel.UserID,
		&transactionModel.Status,
		&transactionModel.CreatedAt,
		&transactionModel.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.TransactionModel{}, apperrors.Wrap(op, apperrors.ErrTransactionNotFound)
		}
		err = postgres.ParsePgError(err)

		return domain.TransactionModel{}, apperrors.Wrap(op, err)
	}

	return transactionModel, nil
}
