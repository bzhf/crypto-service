package repository

import (
	"context"
	"fmt"
	gen "portfolio-service/gen"
)

const GetPublicPortfolios = `
SELECT portfolios.id,portfolios.name,portfolio_assets.symbol,portfolio_assets.amount
FROM portfolios
JOIN portfolio_assets
ON portfolio_assets.portfolio_id=portfolios.id
WHERE portfolios.user_id=$1 AND portfolios.is_public=TRUE
ORDER BY portfolios.id
`

func (r *PortfolioRepository) GetOtherUserPublicPortfolios(ctx context.Context, userID int) ([]*gen.PublicPortfolio, error) {
	rows, err := r.pool.Query(ctx, GetPublicPortfolios, userID)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	portfoliosMap := make(map[int32]*gen.PublicPortfolio)

	for rows.Next() {
		var id int32
		var name string
		var symbol string
		var amount float64

		if err := rows.Scan(&id, &name, &symbol, &amount); err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}

		if _, ok := portfoliosMap[id]; !ok {
			portfoliosMap[id] = &gen.PublicPortfolio{
				PortfolioId: id,
				Name:        name,
				Assets:      make(map[string]float64),
			}
		}
		portfoliosMap[id].Assets[symbol] = amount
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	result := make([]*gen.PublicPortfolio, 0, len(portfoliosMap))
	for _, portfolio := range portfoliosMap {
		result = append(result, portfolio)
	}

	return result, nil
}
