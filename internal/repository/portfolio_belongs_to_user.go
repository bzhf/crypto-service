package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

const PortfolioBelongsToUser = `
SELECT 1 FROM portfolios 
WHERE id = $1 AND user_id = $2 
LIMIT 1;
`

func (r *PortfolioRepository) PortfolioBelongsToUser(ctx context.Context, portfolioID, userID int) (bool, error) {
	var exists int
	if err := r.pool.QueryRow(ctx, PortfolioBelongsToUser, portfolioID, userID).Scan(exists); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, fmt.Errorf("portfolio does not belong to user: %w", err)
		}
		return false, fmt.Errorf("error scanning: %w", err)
	}
	return true, nil

}
