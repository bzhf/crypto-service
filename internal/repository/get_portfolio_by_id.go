package repository

import (
	"context"
	"errors"
	"fmt"
	"portfolio-service/internal/entity"

	"github.com/jackc/pgx/v5"
)

const GetPortfolioByID = `
SELECT user_id,is_public
FROM portfolios
WHERE id=$1 
`

func (r *PortfolioRepository) GetPortfolioByID(ctx context.Context, id int) (*entity.Portfolio, error) {
	portfolio := &entity.Portfolio{}
	if err := r.pool.QueryRow(ctx, GetPortfolioByID, id).Scan(&portfolio.UserID, &portfolio.IsPublic); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("portfolio not found: %w", err)
		}
		return nil, fmt.Errorf("error scanning: %w", err)
	}
	return portfolio, nil
}
