package controller

import (
	"context"
	gen "portfolio-service/gen"
	"portfolio-service/internal/infrastructure/logger"
	"portfolio-service/internal/server/middleware"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

func (s *Controller) GetAllPortfolios(ctx context.Context, e *emptypb.Empty) (*gen.GetAllPortfoliosResponse, error) {

	userID, err := middleware.GetUserID(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "no user_id in context: %v", err)
	}
	portfolios, err := s.uc.GetAllPortfolios(ctx, userID)
	if err != nil {
		logger.FromContext(ctx).Errorw("GetAllPortfolios", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to get all portfolios: %v", err)
	}
	var protoPortfolios []*gen.AllUserPortfolio
	for _, portfolio := range portfolios {
		protoPortfolios = append(protoPortfolios, &gen.AllUserPortfolio{
			Id:       int32(portfolio.PortfolioID),
			Name:     portfolio.Name,
			IsPublic: wrapperspb.Bool(portfolio.IsPublic),
		})
	}
	res := &gen.GetAllPortfoliosResponse{
		Portfolios: protoPortfolios,
	}
	return res, nil
}
