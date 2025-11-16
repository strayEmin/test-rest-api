package repo

import (
	"context"
	"errors"
	"test-task-rest-api/internal/apperrors"
	"test-task-rest-api/internal/domain"
	"test-task-rest-api/internal/repo/postgres"

	"github.com/jackc/pgx/v4"
)

type UserRepo struct {
	pool postgres.AtomicPoolClient
}

func NewUserRepo(pool postgres.AtomicPoolClient) *UserRepo {
	return &UserRepo{
		pool: pool,
	}
}

func (r *UserRepo) IsExistByEmail(ctx context.Context, email string) (bool, error) {
	const op = "UserRepo.IsExistByEmail"
	q := `SELECT EXISTS(SELECT * FROM users WHERE email = $1)`

	var isExist bool
	if err := r.pool.QueryRow(ctx, q, email).Scan(&isExist); err != nil {
		err = postgres.ParsePgError(err)

		return false, apperrors.Wrap(op, err)
	}

	return isExist, nil
}

func (r *UserRepo) CreateUser(ctx context.Context, userModel domain.UserModel) (domain.UserModel, error) {
	const op = "UserRepo.CreateUser"
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
			return domain.UserModel{}, apperrors.Wrap(op, apperrors.ErrUserAlreadyExists)
		}
		err = postgres.ParsePgError(err)

		return domain.UserModel{}, apperrors.Wrap(op, err)
	}

	return userModel, nil
}

func (r *UserRepo) ByEmail(ctx context.Context, email string) (domain.UserModel, error) {
	const op = "UserRepo.ByEmail"
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
			return domain.UserModel{}, apperrors.Wrap(op, apperrors.ErrUserNotFound)
		}
		err = postgres.ParsePgError(err)

		return domain.UserModel{}, apperrors.Wrap(op, err)
	}

	return userModel, nil
}
