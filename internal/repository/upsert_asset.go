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
const InsertUpdateInHistory = `	
	INSERT INTO portfolio_asset_changes (portfolio_id, symbol,amount,updated_at)
	VALUES (?,?,?,?)
`

func (r *PortfolioRepository) UpsertAsset(ctx context.Context, portfolioID int, symbol string, amount float64) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, UpsertAsset, portfolioID, symbol, amount, time.Now().UTC()); err != nil {
		return fmt.Errorf("error creating/updating asset: %w", err)
	}

	if err := r.clickhouse.Exec(ctx, InsertUpdateInHistory, portfolioID, symbol, amount, time.Now().UTC()); err != nil {
		return fmt.Errorf("clickhouse insert error: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit failed: %w", err)
	}
	return nil
}
