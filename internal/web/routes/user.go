package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/joaovictorsl/go-backend-template/internal/config"
	user "github.com/joaovictorsl/go-backend-template/internal/core/user/usecase"
	"github.com/joaovictorsl/go-backend-template/internal/web/handler"
	"github.com/joaovictorsl/go-backend-template/internal/web/middleware"
)

func SetupUser(u *user.UseCase, r *chi.Mux, cfg *config.Config) {
	authMiddleware := middleware.RequiresAuthentication("temp")

	r.Group(func(r chi.Router) {
		r.Use(authMiddleware)

		r.Get("/users/me", handler.HandleGetUser(u))
	})
}
