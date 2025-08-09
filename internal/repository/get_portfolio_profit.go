package repository

import (
	"context"
	"fmt"
	portfolio "portfolio-service/gen"
	"strconv"
)

func (r *PortfolioRepository) GetPortfolioProfit(ctx context.Context, portfolioID int) ([]*portfolio.AssetProfit, error) {
	query := fmt.Sprintf(`
		WITH changes AS (
			SELECT
				symbol,
				amount,
				changed_at,
				lead(changed_at, 1, now()) OVER (PARTITION BY symbol ORDER BY changed_at) AS next_change
			FROM portfolio_asset_changes
			WHERE portfolio_id = %d
		),
		prices_at_change AS (
			SELECT
				c.symbol,
				c.amount,
				c.changed_at,
				argMin(p.price, p.ts) AS price
			FROM changes c
			LEFT JOIN asset_prices p
				ON p.symbol = c.symbol
				AND p.ts >= c.changed_at
				AND p.ts < c.next_change
			GROUP BY
        		c.symbol, c.amount, c.changed_at
		),
		invested AS (
			SELECT
				symbol,
				sum(price * amount) AS invested_value
			FROM prices_at_change
			GROUP BY symbol
		),
		latest_amounts AS (
			SELECT
				symbol,
				anyLast(amount) AS current_amount
			FROM portfolio_asset_changes
			WHERE portfolio_id = %d
			GROUP BY symbol
		)
		SELECT
			i.symbol,
			i.invested_value,
			la.current_amount
		FROM invested i
		JOIN latest_amounts la ON i.symbol = la.symbol
	`, portfolioID, portfolioID)

	rows, err := r.clickhouse.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("clickhouse profit query failed: %w", err)
	}
	defer rows.Close()

	var result []*portfolio.AssetProfit
	for rows.Next() {
		var symbol string
		var invested float64
		var amount float64

		if err := rows.Scan(&symbol, &invested, &amount); err != nil {
			return nil, fmt.Errorf("clickhouse scan error: %w", err)
		}

		key := fmt.Sprintf("crypto:%s:price", symbol)
		priceStr, err := r.redis.Get(ctx, key).Result()
		if err != nil {
			return nil, fmt.Errorf("redis get price failed for %s: %w", symbol, err)
		}

		currentPrice, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid price in redis for %s: %w", symbol, err)
		}

		profit := (currentPrice * amount) - invested

		result = append(result, &portfolio.AssetProfit{
			Symbol:       symbol,
			Amount:       amount,
			Invested:     invested,
			CurrentPrice: currentPrice,
			Profit:       profit,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("clickhouse rows error: %w", err)
	}

	return result, nil
}
