package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/joaovictorsl/go-backend-template/internal/core/entity"
)

type Service interface {
	Get(ctx context.Context, id uuid.UUID) (entity.User, error)
}

type UseCase struct {
	userService Service
}

func New(userService Service) *UseCase {
	return &UseCase{
		userService: userService,
	}
}

func (u *UseCase) Get(ctx context.Context, id uuid.UUID) (entity.User, error) {
	return u.userService.Get(ctx, id)
}
