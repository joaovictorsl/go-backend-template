package database

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joaovictorsl/go-backend-template/internal/config"
)

func NewDatabase(cfg *config.Config) *pgxpool.Pool {
	dbpool, err := pgxpool.New(context.Background(), cfg.DatabaseUrl)
	if err != nil {
		slog.Error("Unable to create connection pool",
			slog.String("error", err.Error()),
		)
		panic(0)
	}
	return dbpool
}

func IsNotFoundError(err error) bool {
	return err == pgx.ErrNoRows
}
