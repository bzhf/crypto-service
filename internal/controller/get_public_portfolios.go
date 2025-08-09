package controller

import (
	"context"
	gen "portfolio-service/gen"
	"portfolio-service/internal/infrastructure/logger"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Controller) GetPublicPortfolios(ctx context.Context, req *gen.GetPublicPortfoliosRequest) (*gen.GetPublicPortfoliosResponse, error) {
	if req.UserId == 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	portfolios, err := s.uc.GetOtherUserPublicPortfolios(ctx, int(req.UserId))
	if err != nil {
		logger.FromContext(ctx).Errorw("GetPublicPortfolios", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to get all portfolios: %v", err)
	}
	res := &gen.GetPublicPortfoliosResponse{
		Portfolios: portfolios,
	}
	return res, nil

}
