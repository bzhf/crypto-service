package entity

import (
	"time"
)

type Portfolio struct {
	PortfolioID int       `json:"portfolio_id"`
	UserID      int       `json:"user_id"`
	Name        string    `json:"name"`
	IsPublic    bool      `json:"is_public"`
	CreatedAt   time.Time `json:"created_at"`
}
