package http_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/joaovictorsl/go-backend-template/internal/config"
	"github.com/joaovictorsl/go-backend-template/internal/core/user/usecase/usecasetest"
	"github.com/joaovictorsl/go-backend-template/internal/http/auth/authtest"
	"github.com/joaovictorsl/go-backend-template/internal/http/router"
	userhttp "github.com/joaovictorsl/go-backend-template/internal/http/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupRouter(
	cfg *config.Config,
	requireAuth func(http.HandlerFunc) http.HandlerFunc,
	userServiceMock *usecasetest.UseCaseMock,
) *chi.Mux {
	r := router.SetupRouter()
	userhttp.SetupRoutes(r, cfg, requireAuth, userServiceMock)
	return r
}

func TestDeleteHandler(t *testing.T) {
	path := "/users/me"
	userId := uuid.NewString()

	t.Run(
		"should return no content status when successfully deletes user",
		func(t *testing.T) {
			// Setup
			cfg := &config.Config{TIMEOUT: 5 * time.Second}
			userUseCaseMock := usecasetest.NewMock(t)
			authMock := authtest.AuthMock{UserId: &userId}

			router := setupRouter(cfg, authMock.RequireAuth, userUseCaseMock)

			userUseCaseMock.On(
				"DeleteUserById",
				mock.Anything,
				userId,
			).Return(nil)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodDelete, path, nil)
			// Action
			router.ServeHTTP(w, r)
			// Assert
			assert.Equal(t, http.StatusNoContent, w.Code)
		},
	)

	t.Run(
		"should return unauthorized status when auth fails",
		func(t *testing.T) {
			// Setup
			cfg := &config.Config{TIMEOUT: 5 * time.Second}
			userUseCaseMock := usecasetest.NewMock(t)
			authMock := authtest.AuthMock{UserId: nil}

			router := setupRouter(cfg, authMock.RequireAuth, userUseCaseMock)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodDelete, path, nil)
			// Action
			router.ServeHTTP(w, r)
			// Assert
			assert.Equal(t, http.StatusUnauthorized, w.Code)
		},
	)

	t.Run(
		"should return internal server error status when usecase fails",
		func(t *testing.T) {
			// Setup
			cfg := &config.Config{TIMEOUT: 5 * time.Second}
			userUseCaseMock := usecasetest.NewMock(t)
			authMock := authtest.AuthMock{UserId: &userId}

			userUseCaseMock.On(
				"DeleteUserById",
				mock.Anything,
				userId,
			).Return(errors.New("failed"))

			router := setupRouter(cfg, authMock.RequireAuth, userUseCaseMock)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodDelete, path, nil)
			// Action
			router.ServeHTTP(w, r)
			// Assert
			assert.Equal(t, http.StatusInternalServerError, w.Code)
		},
	)
}
