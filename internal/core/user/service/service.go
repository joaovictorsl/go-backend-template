package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/joaovictorsl/go-backend-template/internal/core/entity"
)

type UserStore interface {
	Get(context.Context, uuid.UUID) (entity.User, error)
}

type Service struct {
	userStore UserStore
}

func New(userStore UserStore) *Service {
	return &Service{
		userStore: userStore,
	}
}

func (s Service) Get(ctx context.Context, id uuid.UUID) (entity.User, error) {
	return s.userStore.Get(ctx, id)
}
