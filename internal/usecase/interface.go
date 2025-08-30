package usecase

import (
	"context"
	portfolio "portfolio-service/gen"
	"portfolio-service/internal/entity"
)

type PortfolioUsecase struct {
	repo PortfolioInterface
}

func NewPortfolioUsecase(repo PortfolioInterface) PortfolioInterface {
	return &PortfolioUsecase{repo: repo}
}

type PortfolioInterface interface {
	CreateNewPortfolio(ctx context.Context, userID int, name string, isPublic bool) (int, string, bool, error)
	GetPortfolioContentById(ctx context.Context, portfolioID int) (map[string]float64, error)
	UpsertAsset(ctx context.Context, portfolioID int, symbol string, amount float64) error
	DeleteAsset(ctx context.Context, portfolioID int, symbol string) error
	GetAllPortfolios(ctx context.Context, userID int) ([]entity.Portfolio, error)
	GetPortfolioHistory(ctx context.Context, portfolioID int, page, pageSize int) (map[string]*portfolio.PricePoints, error)
	GetOtherUserPublicPortfolios(ctx context.Context, userID int) ([]*portfolio.PublicPortfolio, error)
	PortfolioBelongsToUser(ctx context.Context, portfolioID, userID int) (bool, error)
	GetPortfolioByID(ctx context.Context, id int) (*entity.Portfolio, error)
	GetPortfolioProfit(ctx context.Context, portfolioID int) ([]*portfolio.AssetProfit, error)
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

func (uc *PortfolioUsecase) GetPortfolioHistory(ctx context.Context, portfolioID int, page, pageSize int) (map[string]*portfolio.PricePoints, error) {
	return uc.repo.GetPortfolioHistory(ctx, portfolioID, page, pageSize)
}

func (uc *PortfolioUsecase) GetOtherUserPublicPortfolios(ctx context.Context, userID int) ([]*portfolio.PublicPortfolio, error) {
	return uc.repo.GetOtherUserPublicPortfolios(ctx, userID)
}

func (uc *PortfolioUsecase) PortfolioBelongsToUser(ctx context.Context, portfolioID, userID int) (bool, error) {
	return uc.repo.PortfolioBelongsToUser(ctx, portfolioID, userID)
}

func (uc *PortfolioUsecase) GetPortfolioByID(ctx context.Context, id int) (*entity.Portfolio, error) {
	return uc.repo.GetPortfolioByID(ctx, id)
}

func (uc *PortfolioUsecase) GetPortfolioProfit(ctx context.Context, portfolioID int) ([]*portfolio.AssetProfit, error) {
	return uc.repo.GetPortfolioProfit(ctx, portfolioID)
}
