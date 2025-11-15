package postgres

import (
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgconn"
)

func DoWithTries(fn func() error, attempts int, delay time.Duration) (err error) {
	for attempts > 0 {
		if err = fn(); err != nil {
			time.Sleep(delay)
			attempts--

			continue
		}

		return nil
	}

	return
}

func ParsePgError(err error) error {
	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		return fmt.Errorf(
			"database error: message=%s, details=%s, where=%s, sqlstate=%s: %w",
			pgErr.Message,
			pgErr.Detail,
			pgErr.Where,
			pgErr.SQLState(),
			err,
		)
	}

	return err
}
