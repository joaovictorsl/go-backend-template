package jwt

import (
	"context"
	_ "embed"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	StoreToken(ctx context.Context, tok Token) error
	GetToken(ctx context.Context, userId uint) (Token, error)
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &RepositoryImpl{
		Database: db,
	}
}

var (
	//go:embed sql/store_token.sql
	StoreTokenSQL string
	//go:embed sql/get_token_by_user_id.sql
	GetTokenByUserIdSQL string
)

type RepositoryImpl struct {
	Database *pgxpool.Pool
}

func (r *RepositoryImpl) StoreToken(ctx context.Context, tok Token) error {
	_, err := r.Database.Exec(ctx, StoreTokenSQL, tok.UserId, tok.JWT, tok.CreatedAt, tok.ExpiresAt)
	return err
}

func (r *RepositoryImpl) GetToken(ctx context.Context, userId uint) (Token, error) {
	var tok Token
	row := r.Database.QueryRow(ctx, GetTokenByUserIdSQL, userId)
	err := row.Scan(
		&tok.UserId,
		&tok.JWT,
		&tok.CreatedAt,
		&tok.ExpiresAt,
	)

	return tok, err
}
