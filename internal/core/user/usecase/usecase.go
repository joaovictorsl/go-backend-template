package usecase

import (
	"context"

	"github.com/joaovictorsl/go-backend-template/internal/core/entity"
	"github.com/joaovictorsl/go-backend-template/internal/core/user/repository"
)

type UseCase interface {
	CreateUser(ctx context.Context, u entity.User) (id string, err error)
	GetUserById(ctx context.Context, userId string) (u entity.User, err error)
	GetUserByProviderId(ctx context.Context, providerId string) (u entity.User, err error)
	DeleteUserById(ctx context.Context, userId string) error
}

func New(userRepository repository.Repository) UseCase {
	return &useCaseImpl{
		userRepository: userRepository,
	}
}

type useCaseImpl struct {
	userRepository repository.Repository
}

func (s *useCaseImpl) CreateUser(ctx context.Context, u entity.User) (id string, err error) {
	return s.userRepository.InsertUser(ctx, u)
}

func (s *useCaseImpl) GetUserById(ctx context.Context, userId string) (u entity.User, err error) {
	return s.userRepository.SelectUserById(ctx, userId)
}

func (s *useCaseImpl) GetUserByProviderId(ctx context.Context, providerId string) (u entity.User, err error) {
	return s.userRepository.SelectUserByProviderId(ctx, providerId)
}

func (s *useCaseImpl) DeleteUserById(ctx context.Context, userId string) error {
	return s.userRepository.DeleteUserById(ctx, userId)
}
