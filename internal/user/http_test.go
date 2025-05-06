package user_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/joaovictorsl/go-backend-template/internal/auth/authtest"
	"github.com/joaovictorsl/go-backend-template/internal/config"
	"github.com/joaovictorsl/go-backend-template/internal/router"
	"github.com/joaovictorsl/go-backend-template/internal/user"
	"github.com/joaovictorsl/go-backend-template/internal/user/usertest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteHandler(t *testing.T) {
	path := "/users"
	userId := uuid.NewString()

	t.Run(
		"should run successfully",
		func(t *testing.T) {
			// Setup
			cfg := &config.Config{TIMEOUT: 5 * time.Second}
			userServiceMock := new(usertest.ServiceMock)
			authMock := authtest.AuthMock{UserId: &userId}
			userDeleteHandler := user.DeleteUserHTTPHandler(cfg, userServiceMock)

			router := router.New()
			router.DELETE(path, authMock.RequireAuth(userDeleteHandler))

			userServiceMock.On(
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
		"should return unauthorized when auth fails",
		func(t *testing.T) {
			// Setup
			cfg := &config.Config{TIMEOUT: 5 * time.Second}
			userServiceMock := new(usertest.ServiceMock)
			authMock := authtest.AuthMock{UserId: nil}
			userDeleteHandler := user.DeleteUserHTTPHandler(cfg, userServiceMock)

			router := router.New()
			router.DELETE(path, authMock.RequireAuth(userDeleteHandler))

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodDelete, path, nil)
			// Action
			router.ServeHTTP(w, r)
			// Assert
			assert.Equal(t, http.StatusUnauthorized, w.Code)
		},
	)

	t.Run(
		"should return internal server error when DeleteUserById fails",
		func(t *testing.T) {
			// Setup
			cfg := &config.Config{TIMEOUT: 5 * time.Second}
			userServiceMock := new(usertest.ServiceMock)
			authMock := authtest.AuthMock{UserId: &userId}
			userDeleteHandler := user.DeleteUserHTTPHandler(cfg, userServiceMock)

			userServiceMock.On(
				"DeleteUserById",
				mock.Anything,
				userId,
			).Return(errors.New("something went wrong"))

			router := router.New()
			router.DELETE(path, authMock.RequireAuth(userDeleteHandler))

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodDelete, path, nil)
			// Action
			router.ServeHTTP(w, r)
			// Assert
			assert.Equal(t, http.StatusInternalServerError, w.Code)
		},
	)
}
