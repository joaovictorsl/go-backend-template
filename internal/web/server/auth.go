package server

import (
	"github.com/joaovictorsl/go-backend-template/internal/web/auth"
)

func (app *Server) setupAuth() {
	providers := auth.GetProviders(app.Config)
	app.mux.Get("/oauth/{provider}", auth.HandleOAuth(providers, app.OAuthStore))
	app.mux.Get("/oauth/{provider}/callback", auth.HandleOAuthCallback(
		providers,
		app.Config.RefreshTokenTTL,
		app.OAuthStore,
		app.UserStore,
		app.RefreshTokenStore,
		app.JwtManager,
	))
	app.mux.Get("/auth/refresh", auth.HandleRefresh(
		app.Config.RefreshTokenTTL,
		app.RefreshTokenStore,
		app.JwtManager,
	))
}
