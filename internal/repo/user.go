package repo

import (
	"context"
	"errors"
	"test-task-rest-api/internal/apperrors"
	"test-task-rest-api/internal/domain"
	"test-task-rest-api/internal/logger"
	"test-task-rest-api/internal/repo/postgres"

	"github.com/jackc/pgx/v4"
)

type UserRepo struct {
	logger logger.Logger
	pool   postgres.AtomicPoolClient
}

func NewUserRepo(logger logger.Logger, pool postgres.AtomicPoolClient) *UserRepo {
	return &UserRepo{
		logger: logger,
		pool:   pool,
	}
}

func (r *UserRepo) IsExistByEmail(ctx context.Context, email string) (bool, error) {
	q := `SELECT EXISTS(SELECT * FROM users WHERE email = $1)`

	var isExist bool
	if err := r.pool.QueryRow(ctx, q, email).Scan(&isExist); err != nil {
		err = postgres.ParsePgError(err)
		r.logger.Errorf("isExist query error: %v", err)
		return false, err
	}

	return isExist, nil
}

func (r *UserRepo) CreateUser(ctx context.Context, userModel domain.UserModel) (domain.UserModel, error) {
	q := `INSERT INTO users (name, email, pwd_hash) VALUES ($1, $2, $3)
		RETURNING id, name, email, pwd_hash`

	if err := r.pool.QueryRow(
		ctx,
		q,
		userModel.Name,
		userModel.Email,
		userModel.PwdHash,
	).Scan(
		&userModel.ID,
		&userModel.Name,
		&userModel.Email,
		&userModel.PwdHash,
	); err != nil {
		if postgres.IsUniqueViolation(err) {
			return domain.UserModel{}, apperrors.ErrUserAlreadyExists
		}

		err = postgres.ParsePgError(err)
		r.logger.Errorf("CreateUser query error: %v", err)
		return domain.UserModel{}, err
	}

	return userModel, nil
}

func (r *UserRepo) ByEmail(ctx context.Context, email string) (domain.UserModel, error) {
	q := `SELECT id, name, email, pwd_hash FROM users WHERE email = $1`

	var userModel domain.UserModel
	if err := r.pool.QueryRow(
		ctx,
		q,
		email,
	).Scan(&userModel.ID,
		&userModel.Name,
		&userModel.Email,
		&userModel.PwdHash,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.UserModel{}, apperrors.ErrUserNotFound
		}

		err = postgres.ParsePgError(err)
		r.logger.Errorf("CreateUser query error: %v", err)
		return domain.UserModel{}, err
	}

	return userModel, nil
}
