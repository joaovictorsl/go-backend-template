package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/joaovictorsl/go-backend-template/internal/config"
	"github.com/joaovictorsl/go-backend-template/internal/web/handler"
)

func SetupOAuth(r *chi.Mux, cfg *config.Config) {
	providers := handler.GetProviders(cfg)
	r.Get("/oauth/{provider}", handler.HandleOAuth(providers))
	r.Get("/oauth/{provider}/callback", handler.HandleOAuthCallback(providers))
}
