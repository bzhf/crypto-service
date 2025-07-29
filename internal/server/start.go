package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	gen "portfolio-service/gen"
	"portfolio-service/internal/controller"
	"portfolio-service/internal/infrastructure/logger"
	"portfolio-service/internal/usecase"
	"syscall"

	"google.golang.org/grpc"
)

func Start(ctx context.Context, uc usecase.PortfolioInterface, port string) error {
	GrpcServer := grpc.NewServer()
	Controller := controller.NewController(uc)
	gen.RegisterPortfolioServiceServer(GrpcServer, Controller)

	const op = "server.Start"
	l := logger.FromContext(ctx)
	l.Infow("grpc server starting", "port:", port)
	listener, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	l.Info("grpc server is running")
	if err := GrpcServer.Serve(listener); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sig := <-stop
	GrpcServer.GracefulStop()
	logger.FromContext(ctx).Infow("app stopped gracefully", "signal:", sig)
	return nil
}
