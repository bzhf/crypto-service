package main

import (
	"context"
	"os"
	"os/signal"
	"portfolio-service/internal/config"
	"portfolio-service/internal/infrastructure/clickhouse"
	"portfolio-service/internal/infrastructure/logger"
	db "portfolio-service/internal/infrastructure/postgres"
	"portfolio-service/internal/infrastructure/postgres/cleanup"
	"portfolio-service/internal/infrastructure/redis"
	"portfolio-service/internal/migrations"
	"portfolio-service/internal/repository"
	"portfolio-service/internal/server"
	"portfolio-service/internal/usecase"
	"syscall"
)

func main() {
	log := logger.NewStdOut(nil)
	logger.SetLogger(log)
	defer log.Sync()

	if err := config.LoadConfig(); err != nil {
		log.Fatalw("Ошибка загрузки конфигурации", "error", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx = logger.WithLogger(ctx, log)
	pool, err := db.NewPostgresDB(ctx)
	if err != nil {
		log.Fatalw("Ошибка подключения к БД", "error", err)
	}

	defer pool.Close()
	if err = migrations.RunPostgresMigrations(ctx); err != nil {
		log.Fatalw("postgres migrations failed", "error", err)
	}
	_, err = cleanup.StartCleanupCron(ctx, pool)
	if err != nil {
		log.Fatalw("failed to start cleanup cron: %v", err)
	}
	ch, err := clickhouse.NewClickhouse(ctx)
	if err != nil {
		log.Fatalw("failed to start clickhouse: %v", err)
	}
	if err := migrations.RunClickhouseMigrations(ctx, ch); err != nil {
		log.Fatalw("clickhouse migrations failed", "error", err)
	}

	redis, err := redis.NewRedisClient(ctx)
	if err != nil {
		log.Fatalw("redis start failed", "error", err)
	}
	repo := repository.NewPortfolioRepository(pool, ch, redis)
	uc := usecase.NewPortfolioUsecase(repo)

	grpcServer, err := server.StartGRPC(ctx, uc, config.AppConfig.GrpcPort)
	if err != nil {
		log.Fatalw("failed to start gRPC server", "error", err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	sig := <-stop
	log.Infow("shutdown signal received", "signal", sig)
	grpcServer.GracefulStop()
	defer cancel()
	log.Info("server stopped gracefully")
}
