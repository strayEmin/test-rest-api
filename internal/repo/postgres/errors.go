package postgres

import (
	"errors"

	"github.com/jackc/pgconn"
)

var ErrUniqueViolation = errors.New("unique violation")

func IsUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}

	return false
}

func IsFKViolation(err error) bool {
	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		return pgErr.Code == "23503"
	}

	return false
}
