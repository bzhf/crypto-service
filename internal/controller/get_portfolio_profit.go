package controller

import (
	"context"
	"errors"
	gen "portfolio-service/gen"
	"portfolio-service/internal/infrastructure/logger"
	"portfolio-service/internal/server/middleware"

	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Controller) GetPortfolioProfit(ctx context.Context, req *gen.GetPortfolioProfitRequest) (*gen.GetPortfolioProfitResponse, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "portfolio_id is required")
	}
	userID, err := middleware.GetUserID(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "no user_id in context: %v", err)
	}
	portfolio, err := s.uc.GetPortfolioByID(ctx, int(req.Id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "portfolio not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	if portfolio.UserID != userID && !portfolio.IsPublic {
		return nil, status.Error(codes.PermissionDenied, "access denied")
	}
	profit, err := s.uc.GetPortfolioProfit(ctx, int(req.Id))
	if err != nil {
		logger.FromContext(ctx).Errorw("GetPortfolioProfit", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to get portfolio profit: %v", err)
	}
	result := &gen.GetPortfolioProfitResponse{
		Assets: profit,
	}
	return result, nil

}
