package postgres

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joaovictorsl/go-backend-template/internal/config"
)

func New(cfg *config.Config) (*pgxpool.Pool, error) {
	dbpool, err := pgxpool.New(context.Background(), cfg.DatabaseUrl)
	if err != nil {
		slog.Error("unable to create connection pool",
			slog.Any("error", err),
		)
		return nil, err
	}

	if err := dbpool.Ping(context.Background()); err != nil {
		slog.Error("unable to ping database",
			slog.Any("error", err),
		)
		return nil, err
	}

	return dbpool, nil
}
