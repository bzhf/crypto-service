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

func (s *Controller) GetPortfolioHistory(ctx context.Context, req *gen.GetPortfolioHistoryRequest) (*gen.GetPortfolioHistoryResponse, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "portfolio_id is required")
	}
	if req.Page == 0 {
		return nil, status.Error(codes.InvalidArgument, "page number is required")
	}
	if req.PageSize == 0 {
		return nil, status.Error(codes.InvalidArgument, "page size is required")
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
	result, err := s.uc.GetPortfolioHistory(ctx, int(req.Id), int(req.Page), int(req.PageSize))
	if err != nil {
		logger.FromContext(ctx).Errorw("GetPortfolioHistory", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to get portfolio history: %v", err)
	}
	res := &gen.GetPortfolioHistoryResponse{
		History: result,
	}
	return res, nil

}
