package server

import (
	"context"
	"fmt"
	"net"
	gen "portfolio-service/gen"
	"portfolio-service/internal/controller"
	"portfolio-service/internal/infrastructure/logger"
	"portfolio-service/internal/server/middleware"
	"portfolio-service/internal/usecase"

	"google.golang.org/grpc"
)

func StartGRPC(ctx context.Context, uc usecase.PortfolioInterface, port string) (*grpc.Server, error) {
	l := logger.FromContext(ctx)
	l.Infow("starting gRPC server", "port", port)

	server := grpc.NewServer(grpc.UnaryInterceptor(middleware.UnaryInterceptor()))
	handler := controller.NewController(uc)
	gen.RegisterPortfolioServiceServer(server, handler)

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return nil, fmt.Errorf("grpc listen failed: %w", err)
	}

	go func() {
		if err := server.Serve(listener); err != nil {
			l.Errorw("grpc serve failed", "error", err)
		}
	}()

	l.Infow("gRPC server started")
	return server, nil
}
