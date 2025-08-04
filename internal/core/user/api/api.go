package api

import (
	"context"

	"github.com/google/uuid"
	"github.com/joaovictorsl/go-backend-template/internal/core/entity"
	"github.com/joaovictorsl/go-backend-template/internal/core/user/usecase"
)

type Api interface {
	GetUserById(ctx context.Context, id uuid.UUID) (entity.User, error)
}

func New(userUseCase usecase.UseCase) Api {
	return &apiImpl{
		userUseCase: userUseCase,
	}
}

type apiImpl struct {
	userUseCase usecase.UseCase
}

func (api *apiImpl) GetUserById(ctx context.Context, id uuid.UUID) (entity.User, error) {
	return api.userUseCase.GetUserById(ctx, id)
}
