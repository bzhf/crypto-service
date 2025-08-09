package controller

import (
	"context"
	gen "portfolio-service/gen"
	"portfolio-service/internal/infrastructure/logger"
	"portfolio-service/internal/server/middleware"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

func (s *Controller) CreateNewPortfolio(ctx context.Context, req *gen.CreateNewPortfolioRequest) (*gen.CreateNewPortfolioResponse, error) {
	userID, err := middleware.GetUserID(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "no user_id in context: %v", err)
	}
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "portfolio name is required")
	}
	if req.IsPublic == nil {
		return nil, status.Error(codes.InvalidArgument, "is_public must be set explicitly")
	}
	id, name, isPublic, err := s.uc.CreateNewPortfolio(ctx, userID, req.Name, req.IsPublic.Value)
	if err != nil {
		logger.FromContext(ctx).Errorw("CreateNewPortfolio", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to create new portfolio: %v", err)
	}
	res := &gen.CreateNewPortfolioResponse{
		Id:       int32(id),
		Name:     name,
		IsPublic: wrapperspb.Bool(isPublic),
	}
	return res, nil
}
