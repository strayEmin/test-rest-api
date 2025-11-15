package postgres

import (
	"context"
	"fmt"
	"test-task-rest-api/internal/config"
	"test-task-rest-api/internal/logger"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type AtomicPoolClient interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	BeginTxFunc(ctx context.Context, txOptions pgx.TxOptions, f func(pgx.Tx) error) error
	Ping(ctx context.Context) error
	Close()
}

type AtomicPool struct {
	pool *pgxpool.Pool
}

func NewPostgres(ctx context.Context, maxAttempts int, cfg *config.Config, logger logger.Logger) (*pgxpool.Pool, error) {
	var err error
	var pool *pgxpool.Pool

	// TODO таймзон

	dsn := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Name,
	)

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	config.MaxConns = int32(cfg.DB.MaxOpenConns)
	config.MaxConnLifetime = cfg.DB.ConnMaxLifetime
	config.MaxConnIdleTime = cfg.DB.ConnMaxIdleTime

	err = DoWithTries(func() error {
		ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()

		pool, err = pgxpool.ConnectConfig(ctx, config)
		if err != nil {
			logger.Errorf("DB connection error. %v", err)
			return err
		}

		if err := pool.Ping(ctx); err != nil {
			logger.Errorf("DB ping error. %v\n", err)
			return err
		}

		return nil
	}, maxAttempts, 5*time.Second)

	return pool, err
}
