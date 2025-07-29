package db

import (
	"context"
	"fmt"
	"time"

	"portfolio-service/internal/config"
	"portfolio-service/internal/infrastructure/logger"

	"github.com/cenkalti/backoff/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	InitialInterval = 500 * time.Millisecond //начальный интервал
	MaxElapsedTime  = 10 * time.Second       //макс время на коннект
	MaxInterval     = 2 * time.Second        // макс интервал
)

func NewPostgresDB(ctx context.Context) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disabled",
		config.AppConfig.DBHost,
		config.AppConfig.DBPort,
		config.AppConfig.DBUser,
		config.AppConfig.DBPassword,
		config.AppConfig.DBName,
	)
	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgresDB - pgxpool.ParseConfig: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("невалидный DSN: %w", err)
	}

	ebo := backoff.NewExponentialBackOff()
	ebo.InitialInterval = InitialInterval
	ebo.MaxElapsedTime = MaxElapsedTime
	ebo.MaxInterval = MaxInterval
	bo := backoff.WithContext(ebo, ctx)

	ping := func() error {
		if err = pool.Ping(ctx); err != nil {
			return fmt.Errorf("БД не отвечает,oшибка:%w", err)
		}
		return nil
	}

	if err := backoff.Retry(ping, bo); err != nil {
		return nil, fmt.Errorf("база не доступна: %w", err)
	}
	logger.FromContext(ctx).Infoln("Postgres is initialized")
	return pool, nil
}
