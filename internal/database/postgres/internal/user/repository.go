package user

import (
	"context"
	_ "embed"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joaovictorsl/go-backend-template/internal/core/entity"
	"github.com/joaovictorsl/go-backend-template/internal/database/postgres/internal"
)

var (
	//go:embed sql/get_user_by_id.sql
	SQLGetUserById string
)

type Repository struct {
	DB *pgxpool.Pool
}

func (r *Repository) Get(ctx context.Context, id uuid.UUID) (u entity.User, err error) {
	row := r.DB.QueryRow(ctx, SQLGetUserById, id)
	err = row.Scan(
		&u.Id,
		&u.Email,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	err, isMapped := internal.MapError(err)
	if err != nil && !isMapped {
		slog.Error(
			"postgres error: failed to read user",
			slog.Any("error", err),
		)
	}

	return u, err
}
