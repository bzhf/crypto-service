package entity

import (
	"time"
)

type Asset struct {
	PortfolioID int
	Symbol      string
	Amount      float64
	UpdatedAt   time.Time
}
