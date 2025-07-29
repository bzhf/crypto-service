package dto

import (
	"time"
)

type PricePoint struct {
	Timestamp time.Time `json:"timestamp"`
	Price     float64   `json:"price"`
}

type HistoryResponse map[string][]PricePoint

type HistoryRequest struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

type PortfolioWithAssets struct {
	PortfolioID int                `json:"portfolio_id"`
	Name        string             `json:"name"`
	Assets      map[string]float64 `json:"assets"`
}
