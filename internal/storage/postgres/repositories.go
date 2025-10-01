package postgres

import (
	"github.com/jackc/pgx/v5/pgxpool"
	refreshtoken "github.com/joaovictorsl/go-backend-template/internal/storage/postgres/internal/refresh_token"
	"github.com/joaovictorsl/go-backend-template/internal/storage/postgres/internal/user"
)

type UserRepository = user.Repository

func NewUserRepository(db *pgxpool.Pool) *user.Repository {
	return &user.Repository{
		DB: db,
	}
}

type RefreshTokenRepository = refreshtoken.Repository

func NewRefreshTokenRepository(db *pgxpool.Pool) *refreshtoken.Repository {
	return &refreshtoken.Repository{
		DB: db,
	}
}
