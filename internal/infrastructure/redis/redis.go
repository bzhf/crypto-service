package redis

import (
	"context"
	"fmt"
	"portfolio-service/internal/config"
	"portfolio-service/internal/infrastructure/logger"
	"strconv"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/redis/go-redis/v9"
)

const (
	InitialInterval = 500 * time.Millisecond
	MaxElapsedTime  = 10 * time.Second
	MaxInterval     = 2 * time.Second
)

func NewRedisClient(ctx context.Context) (*redis.Client, error) {
	dbNum, err := strconv.Atoi(config.AppConfig.RedisDB)
	if err != nil {
		return nil, fmt.Errorf("invalid RedisDB number: %w", err)
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.AppConfig.RedisHost, config.AppConfig.RedisPort),
		Password: config.AppConfig.RedisPassword,
		DB:       dbNum,
	})

	ebo := backoff.NewExponentialBackOff()
	ebo.InitialInterval = InitialInterval
	ebo.MaxElapsedTime = MaxElapsedTime
	ebo.MaxInterval = MaxInterval
	bo := backoff.WithContext(ebo, ctx)

	ping := func() error {
		status := rdb.Ping(ctx)
		if err := status.Err(); err != nil {
			return fmt.Errorf("redis ping error: %w", err)
		}
		return nil
	}

	if err := backoff.Retry(ping, bo); err != nil {
		return nil, fmt.Errorf("redis not available: %w", err)
	}

	logger.FromContext(ctx).Infoln("Redis is initialized")
	return rdb, nil
}
