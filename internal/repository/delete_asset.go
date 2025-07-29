package repository

import (
	"context"
	"fmt"
)

const DeleteAsset = `
DELETE FROM portfolio_assets
WHERE portfolio_id=$1 AND symbol=$2
`

func (r *PortfolioRepository) DeleteAsset(ctx context.Context, portfolioID int, symbol string) error {
	cmdTag, err := r.pool.Exec(ctx, DeleteAsset, portfolioID, symbol)
	if err != nil {
		return fmt.Errorf("error deleting asset: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("asset for deletion not found: %w", err)
	}
	return nil
}
