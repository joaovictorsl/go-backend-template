package usertest

import (
	"context"

	"github.com/joaovictorsl/go-backend-template/internal/entity"
	"github.com/stretchr/testify/mock"
)

type ServiceMock struct {
	mock.Mock
}

func (m *ServiceMock) CreateUser(ctx context.Context, u entity.User) (id string, err error) {
	args := m.Called(ctx, u)
	return args.String(0), args.Error(1)
}

func (m *ServiceMock) GetUserById(ctx context.Context, userId string) (u entity.User, err error) {
	args := m.Called(ctx, userId)
	return args.Get(0).(entity.User), args.Error(1)
}

func (m *ServiceMock) GetUserByProviderId(ctx context.Context, providerId string) (u entity.User, err error) {
	args := m.Called(ctx, providerId)
	return args.Get(0).(entity.User), args.Error(1)
}

func (m *ServiceMock) DeleteUserById(ctx context.Context, userId string) error {
	args := m.Called(ctx, userId)
	return args.Error(0)
}

