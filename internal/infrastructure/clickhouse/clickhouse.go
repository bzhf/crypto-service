package clickhouse

import (
	"context"
	"fmt"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/cenkalti/backoff/v4"

	"portfolio-service/internal/config"
	"portfolio-service/internal/infrastructure/logger"
)

const (
	InitialInterval = 500 * time.Millisecond
	MaxElapsedTime  = 10 * time.Second
	MaxInterval     = 2 * time.Second
)

func NewClickhouse(ctx context.Context) (clickhouse.Conn, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%s", config.AppConfig.ClickhouseHost, config.AppConfig.ClickhousePort)},
		Auth: clickhouse.Auth{
			Database: config.AppConfig.ClickhouseDatabase,
			Username: config.AppConfig.ClickhouseUser,
			Password: config.AppConfig.ClickhousePassword,
		},
		DialTimeout: 5 * time.Second,
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("clickhouse - NewClickhouse - open: %w", err)
	}

	ebo := backoff.NewExponentialBackOff()
	ebo.InitialInterval = InitialInterval
	ebo.MaxElapsedTime = MaxElapsedTime
	ebo.MaxInterval = MaxInterval
	bo := backoff.WithContext(ebo, ctx)

	ping := func() error {
		if err := conn.Ping(ctx); err != nil {
			return fmt.Errorf("ClickHouse не отвечает: %w", err)
		}
		return nil
	}

	if err := backoff.Retry(ping, bo); err != nil {
		return nil, fmt.Errorf("ClickHouse недоступен: %w", err)
	}

	logger.FromContext(ctx).Infoln("ClickHouse is initialized")
	return conn, nil
}
