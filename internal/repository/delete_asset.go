package repository

import (
	"context"
	"fmt"
	"time"
)

const DeleteAsset = `
DELETE FROM portfolio_assets
WHERE portfolio_id=$1 AND symbol=$2
`
const InsertDeleteInHistory = `
INSERT INTO portfolio_asset_changes (portfolio_id, symbol, amount, changed_at)
VALUES (?, ?, ?, ?)
`

func (r *PortfolioRepository) DeleteAsset(ctx context.Context, portfolioID int, symbol string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}
	cmdTag, err := tx.Exec(ctx, DeleteAsset, portfolioID, symbol)
	if err != nil {
		return fmt.Errorf("error deleting asset: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("asset for deletion not found: %w", err)
	}
	if err := r.clickhouse.Exec(ctx, InsertDeleteInHistory, portfolioID, symbol, 0.0, time.Now().UTC()); err != nil {
		return fmt.Errorf("clickhouse history insert failed: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit failed: %w", err)
	}

	return nil
}
