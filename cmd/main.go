package main

import (
	"context"

	"portfolio-service/internal/config"
	"portfolio-service/internal/infrastructure/logger"
	db "portfolio-service/internal/infrastructure/postgres"
	"portfolio-service/internal/infrastructure/postgres/cleanup"
	"portfolio-service/internal/repository"
	"portfolio-service/internal/server"
	"portfolio-service/internal/usecase"
	"portfolio-service/migrations"
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
	if err = migrations.RunMigrations(ctx); err != nil {
		log.Fatalw("Ошибка миграции", "error", err)
	}
	_, err = cleanup.StartCleanupCron(ctx, pool)
	if err != nil {
		log.Fatalw("failed to start cleanup cron: %v", err)
	}

	repo := repository.NewPortfolioRepository(pool)
	uc := usecase.NewPortfolioUsecase(repo)

	go func() {
		if err := server.Start(ctx, uc, config.AppConfig.ServerPort); err != nil {
			log.Fatal("Error starting grpc server", "error", err)
		}
	}()

}
