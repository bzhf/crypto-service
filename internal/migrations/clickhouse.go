package migrations

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"sort"
	"strings"

	"github.com/ClickHouse/clickhouse-go/v2"

	"portfolio-service/internal/infrastructure/logger"

	_ "github.com/ClickHouse/clickhouse-go/v2"
)

//go:embed sql/clickhouse/*.sql
var clickhouseMigrations embed.FS

func RunClickhouseMigrations(ctx context.Context, conn clickhouse.Conn) error {
	l := logger.FromContext(ctx)
	files, err := fs.ReadDir(clickhouseMigrations, "sql/clickhouse")
	if err != nil {
		return fmt.Errorf("failed to read migrations dir: %w", err)
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".sql") {
			continue
		}

		path := "sql/clickhouse/" + f.Name()
		content, err := clickhouseMigrations.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", f.Name(), err)
		}

		query := string(content)
		if strings.TrimSpace(query) == "" {
			continue
		}

		if err := conn.Exec(ctx, query); err != nil {
			return fmt.Errorf("error executing migration %s: %w", f.Name(), err)
		}
	}
	l.Infow("Migrations applied successfully")
	return nil
}
