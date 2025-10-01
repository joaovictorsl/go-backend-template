package refreshtoken

import (
	"context"
	_ "embed"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joaovictorsl/go-backend-template/internal/storage/postgres/internal"
	"github.com/joaovictorsl/go-backend-template/internal/web/auth"
)

var (
	//go:embed sql/new_refresh_token.sql
	SQLNewRefreshToken string
	//go:embed sql/get_refresh_token_by_value.sql
	SQLGetRefreshToken string
	//go:embed sql/update_refresh_token.sql
	SQLUpdateRefreshToken string
	//go:embed sql/delete_refresh_token.sql
	SQLDeleteRefreshToken string
)

type Repository struct {
	DB *pgxpool.Pool
}

func (r *Repository) Insert(ctx context.Context, userId uuid.UUID, rTok uuid.UUID, expiresAt time.Time) error {
	_, err := r.DB.Exec(ctx, SQLNewRefreshToken, userId, rTok, expiresAt)
	return internal.MapError(err)
}

func (r *Repository) Get(ctx context.Context, rTokValue uuid.UUID) (rTok auth.RefreshToken, err error) {
	err = r.DB.QueryRow(ctx, SQLGetRefreshToken, rTokValue).
		Scan(
			&rTok.UserId,
			&rTok.Value,
			&rTok.ExpiresAt,
		)
	return rTok, internal.MapError(err)
}

func (r *Repository) Update(ctx context.Context, oldRTok uuid.UUID, newRTok uuid.UUID, newRTokExpiresAt time.Time) error {
	_, err := r.DB.Exec(ctx, SQLUpdateRefreshToken, oldRTok, newRTok, newRTokExpiresAt)
	return internal.MapError(err)
}

func (r *Repository) Delete(ctx context.Context, rTok uuid.UUID) error {
	_, err := r.DB.Exec(ctx, SQLDeleteRefreshToken, rTok)
	return internal.MapError(err)
}
