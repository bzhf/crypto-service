package repository

import (
	"context"
	"fmt"
	"time"
	//"portfolio-service/internal/entity"
	//
)

const CreateNewPortfolio = `
INSERT INTO portfolios (user_id, name, is_public,created_at)
VALUES ($1,$2,$3,$4)
RETURNING id
`

func (r *PortfolioRepository) CreateNewPortfolio(ctx context.Context, userID int, name string, isPublic bool) (int, string, bool, error) {
	var id int
	err := r.pool.QueryRow(ctx, CreateNewPortfolio, userID, name, isPublic, time.Now().UTC()).Scan(&id)
	if err != nil {
		return 0, "", false, fmt.Errorf("ошибка создания портфеля: %w", err)
	}

	return id, name, isPublic, nil
}
