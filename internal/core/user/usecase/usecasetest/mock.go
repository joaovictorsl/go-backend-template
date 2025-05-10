package usecasetest

import (
	"context"
	"testing"

	"github.com/joaovictorsl/go-backend-template/internal/core/entity"
	"github.com/stretchr/testify/mock"
)

type UseCaseMock struct {
	mock.Mock
}

func NewMock(t *testing.T) *UseCaseMock {
	m := &UseCaseMock{}
	t.Cleanup(func() {
		m.AssertExpectations(t)
	})
	return m
}

func (m *UseCaseMock) CreateUser(ctx context.Context, u entity.User) (string, error) {
	args := m.Called(ctx, u)
	return args.String(0), args.Error(1)
}

func (m *UseCaseMock) GetUserById(ctx context.Context, userId string) (entity.User, error) {
	args := m.Called(ctx, userId)
	return args.Get(0).(entity.User), args.Error(1)
}

func (m *UseCaseMock) GetUserByProviderId(ctx context.Context, providerId string) (entity.User, error) {
	args := m.Called(ctx, providerId)
	return args.Get(0).(entity.User), args.Error(1)
}

func (m *UseCaseMock) DeleteUserById(ctx context.Context, userId string) error {
	args := m.Called(ctx, userId)
	return args.Error(0)
}
