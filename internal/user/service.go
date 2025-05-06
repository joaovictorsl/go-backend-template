package user

import (
	"context"

	"github.com/joaovictorsl/go-backend-template/internal/entity"
)

type Service interface {
	CreateUser(ctx context.Context, u entity.User) (id string, err error)
	GetUserById(ctx context.Context, userId string) (u entity.User, err error)
	GetUserByProviderId(ctx context.Context, providerId string) (u entity.User, err error)
	DeleteUserById(ctx context.Context, userId string) error
}

func NewService(userRepository Repository) Service {
	return &ServiceImpl{
		userRepository: userRepository,
	}
}

type ServiceImpl struct {
	userRepository Repository
}

func (s *ServiceImpl) CreateUser(ctx context.Context, u entity.User) (id string, err error) {
	return s.userRepository.CreateUser(ctx, u)
}

func (s *ServiceImpl) GetUserById(ctx context.Context, userId string) (u entity.User, err error) {
	return s.userRepository.GetUserById(ctx, userId)
}

func (s *ServiceImpl) GetUserByProviderId(ctx context.Context, providerId string) (u entity.User, err error) {
	return s.userRepository.GetUserByProviderId(ctx, providerId)
}

func (s *ServiceImpl) DeleteUserById(ctx context.Context, userId string) error {
	return s.userRepository.DeleteUserById(ctx, userId)
}
