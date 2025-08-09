package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	gen "portfolio-service/gen"
	"portfolio-service/internal/controller"
	"portfolio-service/internal/infrastructure/logger"
	"portfolio-service/internal/server/middleware"
	"portfolio-service/internal/usecase"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
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

func StartRESTProxy(ctx context.Context, grpcPort string, httpPort string) (*http.Server, error) {
	l := logger.FromContext(ctx)
	l.Infow("starting REST proxy", "port", ":"+httpPort)

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := gen.RegisterPortfolioServiceHandlerFromEndpoint(ctx, mux, grpcPort, opts)
	if err != nil {
		return nil, fmt.Errorf("register grpc-gateway handler failed: %w", err)
	}

	srv := &http.Server{
		Addr:    ":" + httpPort,
		Handler: mux,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			l.Errorw("REST proxy serve failed", "error", err)
		}
	}()

	l.Infow("REST proxy started")
	return srv, nil
}
