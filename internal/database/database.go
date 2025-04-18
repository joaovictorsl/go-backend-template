package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joaovictorsl/go-backend-template/internal/config"
)

func NewDatabase(cfg *config.Config) *pgxpool.Pool {
	dbpool, err := pgxpool.New(context.Background(), cfg.DATABASE_URL)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}
	return dbpool
}
