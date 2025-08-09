package controller

import (
	"context"
	gen "portfolio-service/gen"
	"portfolio-service/internal/infrastructure/logger"
	"portfolio-service/internal/server/middleware"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Controller) DeleteAsset(ctx context.Context, req *gen.DeleteAssetRequest) (*emptypb.Empty, error) {
	userID, err := middleware.GetUserID(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "no user_id in context: %v", err)
	}
	if req.PortfolioId == 0 {
		return nil, status.Error(codes.InvalidArgument, "portfolio_id is required")
	}
	if req.Symbol == "" {
		return nil, status.Error(codes.InvalidArgument, "asset symbol is required")
	}
	ok, err := s.uc.PortfolioBelongsToUser(ctx, int(req.PortfolioId), userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check portfolio ownership: %v", err)
	}
	if !ok {
		return nil, status.Error(codes.PermissionDenied, "access to portfolio denied")
	}
	if err := s.uc.DeleteAsset(ctx, int(req.PortfolioId), req.Symbol); err != nil {
		logger.FromContext(ctx).Errorw("DeleteAsset", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to delete asset: %v", err)
	}
	return nil, nil
}
