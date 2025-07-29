package repository

import (
	"context"
	"fmt"
	"time"
)

const UpsertAsset = `
	INSERT INTO portfolio_assets (portfolio_id, symbol,amount,updated_at)
	VALUES ($1,$2,$3,$4)
	ON CONFLICT (portfolio_id,symbol)
	DO UPDATE SET 
		amount=EXCLUDED.amount,
		updated_at=EXCLUDED.updated_at
`

func (r *PortfolioRepository) UpsertAsset(ctx context.Context, portfolioID int, symbol string, amount float64) error {
	if _, err := r.pool.Exec(ctx, UpsertAsset, portfolioID, symbol, amount, time.Now().UTC()); err != nil {
		return fmt.Errorf("error creating/updating asset: %w", err)
	}
	return nil
}
