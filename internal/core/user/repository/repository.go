package repository

import (
	"context"
	_ "embed"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joaovictorsl/go-backend-template/internal/core/entity"
	"github.com/joaovictorsl/go-backend-template/internal/database"
)

type Repository interface {
	InsertUser(ctx context.Context, u entity.User) (id string, err error)
	SelectUserById(ctx context.Context, id string) (u entity.User, err error)
	SelectUserByProviderId(ctx context.Context, providerId string) (u entity.User, err error)
	DeleteUserById(ctx context.Context, id string) error
}

func New(db *pgxpool.Pool) Repository {
	return repositoryImpl{
		Database: db,
	}
}

var (
	//go:embed sql/insert_user.sql
	insertUserSQL string
	//go:embed sql/select_user_by_id.sql
	selectUserByIdSQL string
	//go:embed sql/select_user_by_provider_id.sql
	selectUserByProviderIdSQL string
	//go:embed sql/delete_user_by_id.sql
	deleteUserByIdSQL string
)

type repositoryImpl struct {
	Database *pgxpool.Pool
}

func (r repositoryImpl) InsertUser(ctx context.Context, u entity.User) (string, error) {
	var id string
	err := r.Database.QueryRow(
		ctx,
		insertUserSQL,
		u.ProviderId,
		u.Email,
	).Scan(&id)
	if err != nil {
		return "", database.MapDatabaseError(err)
	}

	return id, nil
}

func (r repositoryImpl) SelectUserById(ctx context.Context, id string) (entity.User, error) {
	var u entity.User
	err := r.Database.QueryRow(
		ctx,
		selectUserByIdSQL,
		id,
	).Scan(
		&u.Id,
		&u.ProviderId,
		&u.Email,
	)
	if err != nil {
		return entity.User{}, database.MapDatabaseError(err)
	}

	return u, nil
}

func (r repositoryImpl) SelectUserByProviderId(ctx context.Context, providerId string) (entity.User, error) {
	var u entity.User
	err := r.Database.QueryRow(
		ctx,
		selectUserByProviderIdSQL,
		providerId,
	).Scan(
		&u.Id,
		&u.ProviderId,
		&u.Email,
	)
	if err != nil {
		return entity.User{}, database.MapDatabaseError(err)
	}

	return u, nil
}

func (r repositoryImpl) DeleteUserById(ctx context.Context, id string) error {
	_, err := r.Database.Exec(
		ctx,
		deleteUserByIdSQL,
		id,
	)
	if err != nil {
		return database.MapDatabaseError(err)
	}

	return nil
}
