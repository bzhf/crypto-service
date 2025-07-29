package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type PortfolioRepository struct {
	pool *pgxpool.Pool
}

func NewPortfolioRepository(pool *pgxpool.Pool) *PortfolioRepository {
	return &PortfolioRepository{pool: pool}
}
