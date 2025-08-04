package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/joaovictorsl/go-backend-template/internal/config"
	userapi "github.com/joaovictorsl/go-backend-template/internal/core/user/api"
	"github.com/joaovictorsl/go-backend-template/internal/http/handler"
	"github.com/joaovictorsl/go-backend-template/internal/http/middleware"
)

func SetupUserRoutes(api userapi.Api, r *chi.Mux, cfg *config.Config) {
	authMiddleware := middleware.RequiresAuthentication(cfg.JwtSecret)

	r.Group(func(r chi.Router) {
		r.Use(authMiddleware)

		r.Get("/users/me", handler.HandleGetUserById(api))
	})
}
