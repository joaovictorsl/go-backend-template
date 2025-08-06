package repository

import (
	"context"
	_ "embed"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joaovictorsl/go-backend-template/internal/core"
	"github.com/joaovictorsl/go-backend-template/internal/core/entity"
	"github.com/joaovictorsl/go-backend-template/internal/database"
)

type Repository interface {
	GetUserById(ctx context.Context, id uuid.UUID) (u entity.User, err error)
}

func New(db *pgxpool.Pool) Repository {
	return repositoryImpl{
		db: db,
	}
}

var (
	//go:embed sql/get_user_by_id.sql
	SQLGetUserById string
)

type repositoryImpl struct {
	db *pgxpool.Pool
}

func (r repositoryImpl) GetUserById(ctx context.Context, id uuid.UUID) (u entity.User, err error) {
	row := r.db.QueryRow(ctx, SQLGetUserById, id)
	err = row.Scan(
		&u.Id,
		&u.Email,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if database.IsNotFoundError(err) {
		err = core.ErrNotFound
	}

	return u, err
}
