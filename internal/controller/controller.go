package controller

import (
	"context"
	gen "portfolio-service/gen"

	"portfolio-service/internal/infrastructure/logger"
	"portfolio-service/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

type Controller struct {
	gen.UnimplementedPortfolioServiceServer
	uc usecase.PortfolioInterface
}

func NewController(uc usecase.PortfolioInterface) *Controller {
	return &Controller{uc: uc}
}

func (s *Controller) CreateNewPortfolio(ctx context.Context, req *gen.CreateNewPortfolioRequest) (*gen.CreateNewPortfolioResponse, error) {
	if req.UserId == 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "portfolio name is required")
	}
	if req.IsPublic == nil {
		return nil, status.Error(codes.InvalidArgument, "is_public must be set explicitly")
	}
	id, name, isPublic, err := s.uc.CreateNewPortfolio(ctx, int(req.UserId), req.Name, req.IsPublic.Value)
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
func (s *Controller) GetPortfolioContentById(ctx context.Context, req *gen.GetPortfolioContentByIdRequest) (*gen.GetPortfolioContentByIdResponse, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "portfolio_id is required")
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
func (s *Controller) UpsertAsset(ctx context.Context, req *gen.UpsertAssetRequest) (*emptypb.Empty, error) {
	if req.PortfolioId == 0 {
		return nil, status.Error(codes.InvalidArgument, "portfolio_id is required")
	}
	if req.Symbol == "" {
		return nil, status.Error(codes.InvalidArgument, "asset symbol is required")
	}
	if req.Amount == 0 {
		return nil, status.Error(codes.InvalidArgument, "asset amount is required")
	}
	if err := s.uc.UpsertAsset(ctx, int(req.PortfolioId), req.Symbol, req.Amount); err != nil {
		logger.FromContext(ctx).Errorw("UpsertAsset", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to add/update asset: %v", err)
	}
	return nil, nil
}

func (s *Controller) DeleteAsset(ctx context.Context, req *gen.DeleteAssetRequest) (*emptypb.Empty, error) {
	if req.PortfolioId == 0 {
		return nil, status.Error(codes.InvalidArgument, "portfolio_id is required")
	}
	if req.Symbol == "" {
		return nil, status.Error(codes.InvalidArgument, "asset symbol is required")
	}
	if err := s.uc.DeleteAsset(ctx, int(req.PortfolioId), req.Symbol); err != nil {
		logger.FromContext(ctx).Errorw("DeleteAsset", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to delete asset: %v", err)
	}
	return nil, nil
}

func (s *Controller) GetAllPortfolios(ctx context.Context, req *gen.GetAllPortfoliosRequest) (*gen.GetAllPortfoliosResponse, error) {
	if req.UserId == 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	portfolios, err := s.uc.GetAllPortfolios(ctx, int(req.UserId))
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
func (s *Controller) GetPortfolioHistory(context.Context, *gen.GetPortfolioHistoryRequest) (*gen.GetPortfolioHistoryResponse, error) {

}

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
