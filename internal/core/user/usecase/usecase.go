package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/joaovictorsl/go-backend-template/internal/core/entity"
	userrepository "github.com/joaovictorsl/go-backend-template/internal/core/user/repository"
)

type UseCase interface {
	GetUserById(ctx context.Context, id uuid.UUID) (entity.User, error)
}

func New(userRepository userrepository.Repository) UseCase {
	return useCaseImpl{
		userRepository,
	}
}

type useCaseImpl struct {
	userRepository userrepository.Repository
}

func (u useCaseImpl) GetUserById(ctx context.Context, id uuid.UUID) (entity.User, error) {
	return u.userRepository.GetUserById(ctx, id)
}
