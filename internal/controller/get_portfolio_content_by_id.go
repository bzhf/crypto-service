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

// должна быть проверка айди юзера. если портфель принадлежит юзеру - он может смотреть личные и публичные,
//
//	если нет, только публичные

func (s *Controller) GetPortfolioContentById(ctx context.Context, req *gen.GetPortfolioContentByIdRequest) (*gen.GetPortfolioContentByIdResponse, error) {
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

	assets, err := s.uc.GetPortfolioContentById(ctx, int(req.Id))
	if err != nil {
		logger.FromContext(ctx).Errorw("GetPortfolioContentById", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to get portfolio content: %v", err)
	}
	res := &gen.GetPortfolioContentByIdResponse{
		Assets: assets,
	}
	return res, nil
}
