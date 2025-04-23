package user

import (
	"context"
	_ "embed"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joaovictorsl/go-backend-template/internal/core/entity"
)

type Repository interface {
	CreateUser(ctx context.Context, u entity.User) (id uint, err error)
	GetUserById(ctx context.Context, id uint) (u entity.User, err error)
	GetUserByProviderId(ctx context.Context, googleId string) (u entity.User, err error)
}

func NewRepository(db *pgxpool.Pool) Repository {
	return RepositoryImpl{
		Database: db,
	}
}

var (
	//go:embed sql/create_user.sql
	CreateUserSQL string
	//go:embed sql/get_user_by_id.sql
	GetUserByIdSQL string
	//go:embed sql/get_user_by_provider_id.sql
	GetUserByProviderIdSQL string
)

type RepositoryImpl struct {
	Database *pgxpool.Pool
}

func (r RepositoryImpl) CreateUser(ctx context.Context, u entity.User) (id uint, err error) {
	row := r.Database.QueryRow(ctx, CreateUserSQL, u.ProviderId, u.Email)
	err = row.Scan(&id)
	return id, err
}

func (r RepositoryImpl) GetUserById(ctx context.Context, id uint) (u entity.User, err error) {
	row := r.Database.QueryRow(ctx, GetUserByIdSQL, id)
	err = row.Scan(&u.Id, &u.ProviderId, &u.Email)
	return u, err
}

func (r RepositoryImpl) GetUserByProviderId(ctx context.Context, googleId string) (u entity.User, err error) {
	row := r.Database.QueryRow(ctx, GetUserByProviderIdSQL, googleId)
	err = row.Scan(&u.Id, &u.ProviderId, &u.Email)
	return u, err
}
