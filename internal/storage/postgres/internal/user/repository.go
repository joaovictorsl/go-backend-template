package user

import (
	"context"
	_ "embed"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joaovictorsl/go-backend-template/internal/core/entity"
	"github.com/joaovictorsl/go-backend-template/internal/storage/postgres/internal"
)

var (
	//go:embed sql/get_user_by_id.sql
	SQLGetUserById string
	//go:embed sql/get_user_by_provider.sql
	SQLGetUserByProvider string
	//go:embed sql/new_user.sql
	SQLNewUser string
	//go:embed sql/new_linked_account.sql
	SQLNewLinkedAccount string
)

type Repository struct {
	DB *pgxpool.Pool
}

func (r *Repository) Get(ctx context.Context, id uuid.UUID) (u entity.User, err error) {
	err = r.DB.QueryRow(ctx, SQLGetUserById, id).
		Scan(
			&u.Id,
			&u.Email,
			&u.CreatedAt,
			&u.UpdatedAt,
		)

	return u, internal.MapError(err)
}

func (r *Repository) GetByProvider(ctx context.Context, provider, providerID string) (u entity.User, err error) {
	err = r.DB.QueryRow(ctx, SQLGetUserByProvider, provider, providerID).
		Scan(
			&u.Id,
			&u.Email,
			&u.CreatedAt,
			&u.UpdatedAt,
		)

	return u, internal.MapError(err)
}

func (r *Repository) Insert(ctx context.Context, email, provider, providerUserId string) (id uuid.UUID, err error) {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return uuid.Nil, internal.MapError(err)
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx, SQLNewUser, email).Scan(&id)
	if err != nil {
		return uuid.Nil, internal.MapError(err)
	}

	_, err = tx.Exec(ctx, SQLNewLinkedAccount, id, provider, providerUserId)
	if err != nil {
		return uuid.Nil, internal.MapError(err)
	}

	if err = tx.Commit(ctx); err != nil {
		return uuid.Nil, internal.MapError(err)
	}

	return id, nil
}
