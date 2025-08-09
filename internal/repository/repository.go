package repository

import (
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type PortfolioRepository struct {
	pool       *pgxpool.Pool
	clickhouse clickhouse.Conn
	redis      *redis.Client
}

func NewPortfolioRepository(pool *pgxpool.Pool, clickhouse clickhouse.Conn, redis *redis.Client) *PortfolioRepository {
	return &PortfolioRepository{pool: pool, clickhouse: clickhouse, redis: redis}
}
