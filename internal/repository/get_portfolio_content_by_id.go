package repository

import (
	"context"
	"fmt"
	"portfolio-service/internal/entity"
)

const GetPortfolioContentById = `
SELECT symbol,amount,updated_at
FROM portfolio_assets
WHERE portfolio_id=$1
`

func (r *PortfolioRepository) GetPortfolioContentById(ctx context.Context, portfolioID int) (map[string]float64, error) {
	rows, err := r.pool.Query(ctx, GetPortfolioContentById, portfolioID)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()
	var assets []entity.Asset
	for rows.Next() {
		var asset entity.Asset
		if err := rows.Scan(&asset.Symbol, &asset.Amount, asset.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		assets = append(assets, asset)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	result := make(map[string]float64)
	for _, asset := range assets {
		result[asset.Symbol] = asset.Amount
	}
	return result, nil
}
