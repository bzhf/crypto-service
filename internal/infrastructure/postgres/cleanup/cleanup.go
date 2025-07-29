package cleanup

import (
	"context"
	"portfolio-service/internal/infrastructure/logger"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robfig/cron/v3"
)

const deleteQuery = `
DELETE FROM portfolio_assets
WHERE portfolio_id IN (
  SELECT id FROM portfolios WHERE created_at < NOW() - INTERVAL '2 months'
);
DELETE FROM portfolios
WHERE created_at < NOW() - INTERVAL '2 months';
`

func StartCleanupCron(ctx context.Context, pool *pgxpool.Pool) (*cron.Cron, error) {
	c := cron.New(cron.WithSeconds())
	_, err := c.AddFunc("0 0 3 * * *", func() {
		logger.FromContext(ctx).Info("Starting cleanup task")
		if _, err := pool.Exec(ctx, deleteQuery); err != nil {
			logger.FromContext(ctx).Errorw("Cleanup error", "error", err)
			return
		}
		logger.FromContext(ctx).Info("Cleanup completed successfully")
	})
	if err != nil {
		return nil, err
	}

	c.Start()

	// В фоне остановка по контексту
	go func() {
		<-ctx.Done()
		logger.FromContext(ctx).Info("Stopping cron due to context cancellation")
		ctxStop, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		c.Stop()
		<-ctxStop.Done()
	}()

	return c, nil
}
