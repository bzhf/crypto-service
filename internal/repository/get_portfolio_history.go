package repository

import (
	"context"
	"fmt"
	portfolio "portfolio-service/gen"

	"time"
)

func (r *PortfolioRepository) GetPortfolioHistory(ctx context.Context, portfolioID int, page, pageSize int) (map[string]*portfolio.PricePoints, error) {
	limit := pageSize
	offset := (page - 1) * pageSize

	// ClickHouse-запрос: получаем стоимость актива по времени с учетом истории изменения количества
	query := fmt.Sprintf(`
		WITH asset_with_amount AS (
			SELECT
				ac.symbol,
				ac.amount,
				ac.changed_at,
				lead(ac.changed_at, 1, now())
					OVER (
						PARTITION BY ac.symbol 
						ORDER BY ac.changed_at 
					) AS next_change
			FROM portfolio_asset_changes ac
			WHERE ac.portfolio_id = %d
		)
		SELECT	
			p.symbol,
			p.ts,
			p.price * a.amount AS value
		FROM asset_prices p
		JOIN asset_with_amount a
			ON p.symbol = a.symbol
			AND p.ts >= a.changed_at
			AND p.ts < a.next_change
		ORDER BY p.symbol, p.ts DESC
		LIMIT %d OFFSET %d
	`, portfolioID, limit, offset)

	rows, err := r.clickhouse.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("clickhouse query error: %w", err)
	}
	defer rows.Close()

	result := make(map[string]*portfolio.PricePoints)
	for rows.Next() {
		var symbol string
		var value float64
		var ts time.Time
		if err := rows.Scan(&symbol, &ts, &value); err != nil {
			return nil, fmt.Errorf("clickhouse scan error: %w", err)
		}

		if _, ok := result[symbol]; !ok {
			result[symbol] = &portfolio.PricePoints{}
		}
		result[symbol].Points = append(result[symbol].Points, &portfolio.PricePoint{
			Timestamp: ts.Format(time.RFC3339),
			Value:     value,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("clickhouse rows error: %w", err)
	}

	return result, nil
}
