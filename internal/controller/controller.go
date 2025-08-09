package controller

import (
	gen "portfolio-service/gen"
	"portfolio-service/internal/usecase"
)

type Controller struct {
	gen.UnimplementedPortfolioServiceServer
	uc usecase.PortfolioInterface
}

func NewController(uc usecase.PortfolioInterface) *Controller {
	return &Controller{uc: uc}
}
