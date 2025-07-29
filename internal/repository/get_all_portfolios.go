package repository

import (
	"context"
	"fmt"
	"portfolio-service/internal/entity"
)

const GetAllPortfolios = `
SELECT * 
FROM portfolios
WHERE user_id=$1`

func (r *PortfolioRepository) GetAllPortfolios(ctx context.Context, user_id int) ([]entity.Portfolio, error) {
	rows, err := r.pool.Query(ctx, GetAllPortfolios, user_id)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	var portfolios []entity.Portfolio
	for rows.Next() {
		var portfolio entity.Portfolio
		if err := rows.Scan(&portfolio.PortfolioID, &portfolio.Name, &portfolio.IsPublic, &portfolio.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		portfolios = append(portfolios, portfolio)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	return portfolios, nil
}
