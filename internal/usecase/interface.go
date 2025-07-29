package usecase

import (
	"context"
	gen "portfolio-service/gen"
	"portfolio-service/internal/entity"
	"portfolio-service/internal/repository"
)

type PortfolioUsecase struct {
	repo *repository.PortfolioRepository
}

func NewPortfolioUsecase(repo *repository.PortfolioRepository) *PortfolioUsecase {
	return &PortfolioUsecase{repo: repo}
}

type PortfolioInterface interface {
	CreateNewPortfolio(ctx context.Context, userID int, name string, isPublic bool) (int, string, bool, error)
	GetPortfolioContentById(ctx context.Context, portfolioID int) (map[string]float64, error)
	UpsertAsset(ctx context.Context, portfolioID int, symbol string, amount float64) error
	DeleteAsset(ctx context.Context, portfolioID int, symbol string) error
	GetAllPortfolios(ctx context.Context, userID int) ([]entity.Portfolio, error)
	GetPortfolioHistory()
	GetOtherUserPublicPortfolios(ctx context.Context, userID int) ([]*gen.PublicPortfolio, error)
}

func (uc *PortfolioUsecase) CreateNewPortfolio(ctx context.Context, userID int, name string, isPublic bool) (int, string, bool, error) {
	return uc.repo.CreateNewPortfolio(ctx, userID, name, isPublic)
}

func (uc *PortfolioUsecase) GetPortfolioContentById(ctx context.Context, portfolioID int) (map[string]float64, error) {
	return uc.repo.GetPortfolioContentById(ctx, portfolioID)
}

func (uc *PortfolioUsecase) UpsertAsset(ctx context.Context, portfolioID int, symbol string, amount float64) error {
	return uc.repo.UpsertAsset(ctx, portfolioID, symbol, amount)
}

func (uc *PortfolioUsecase) DeleteAsset(ctx context.Context, portfolioID int, symbol string) error {
	return uc.repo.DeleteAsset(ctx, portfolioID, symbol)
}

func (uc *PortfolioUsecase) GetAllPortfolios(ctx context.Context, userID int) ([]entity.Portfolio, error) {
	return uc.repo.GetAllPortfolios(ctx, userID)
}

func (uc *PortfolioUsecase) GetPortfolioHistory() {
}

func (uc *PortfolioUsecase) GetOtherUserPublicPortfolios(ctx context.Context, userID int) ([]*gen.PublicPortfolio, error) {
	return uc.repo.GetOtherUserPublicPortfolios(ctx, userID)
}
