package user

import (
	"context"

	"github.com/joaovictorsl/go-backend-template/internal/core/entity"
)

type Service interface {
	CreateUser(ctx context.Context, u entity.User) (id uint, err error)
	GetUserById(ctx context.Context, userId uint) (u entity.User, err error)
	GetUserByProviderId(ctx context.Context, providerId string) (u entity.User, err error)
}

func NewService(userRepository Repository) Service {
	return &ServiceImpl{
		userRepository: userRepository,
	}
}

type ServiceImpl struct {
	userRepository Repository
}

func (s *ServiceImpl) CreateUser(ctx context.Context, u entity.User) (id uint, err error) {
	return s.userRepository.CreateUser(ctx, u)
}

func (s *ServiceImpl) GetUserById(ctx context.Context, userId uint) (u entity.User, err error) {
	return s.GetUserById(ctx, userId)
}

func (s *ServiceImpl) GetUserByProviderId(ctx context.Context, providerId string) (u entity.User, err error) {
	return s.GetUserByProviderId(ctx, providerId)
}
