package cleanup

import (
	"context"
	"time"

	"portfolio-service/internal/infrastructure/logger"

	"github.com/go-co-op/gocron/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

const deleteQuery = `
DELETE FROM portfolio_assets
WHERE portfolio_id IN (
  SELECT id FROM portfolios WHERE created_at < NOW() - INTERVAL '2 months'
);
DELETE FROM portfolios
WHERE created_at < NOW() - INTERVAL '2 months';
`

func StartCleanupCron(ctx context.Context, pool *pgxpool.Pool) (gocron.Scheduler, error) {
	s, err := gocron.NewScheduler(
		gocron.WithLocation(time.UTC),
	)
	if err != nil {
		return nil, err
	}

	_, err = s.NewJob(
		gocron.CronJob("0 0 3 * * *", true),
		gocron.NewTask(func() {
			logger.FromContext(ctx).Info("Starting cleanup task")
			if _, err := pool.Exec(ctx, deleteQuery); err != nil {
				logger.FromContext(ctx).Errorw("Cleanup error", "error", err)
				return
			}
			logger.FromContext(ctx).Info("Cleanup completed successfully")
		}),
	)
	if err != nil {
		return nil, err
	}

	s.Start()

	go func() {
		<-ctx.Done()
		logger.FromContext(ctx).Info("Stopping cron due to context cancellation")
		ctxStop, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.Shutdown(); err != nil {
			logger.FromContext(ctx).Errorw("Error shutting down scheduler", "error", err)
		}
		<-ctxStop.Done()
	}()

	return s, nil
}
