package auth

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/joaovictorsl/aegis"
	"github.com/joaovictorsl/aegis/oauth"
	"github.com/joaovictorsl/aegis/token"
	"github.com/joaovictorsl/go-backend-template/internal/config"
	"github.com/joaovictorsl/go-backend-template/internal/entity"
	"github.com/joaovictorsl/go-backend-template/internal/user"
)

type Auth interface {
	GoogleLoginHandler() http.Handler
	GoogleCallbackHandler() http.Handler
	RequireAuth(http.Handler) http.Handler
}

func New(
	cfg *config.Config,
	userService user.Service,
	tokenRepository token.Repository,
) Auth {
	jwtManager := token.NewJWTManager(cfg.JWT_ISS, cfg.JWT_SECRET, cfg.ACCESS_TOKEN_EXP, cfg.REFRESH_TOKEN_EXP)

	a := aegis.New(
		jwtManager,
		func(ctx context.Context, pu oauth.ProviderUser) (string, error) {
			u, err := userService.GetUserByProviderId(ctx, pu.Id)
			if err != nil {
				if !errors.Is(err, sql.ErrNoRows) {
					return "", err
				}

				newUserId, err := userService.CreateUser(ctx, entity.User{
					ProviderId: pu.Provider + "|" + pu.Id,
					Email:      pu.Email,
				})
				if err != nil {
					return "", err
				}

				u.Id = newUserId
			}

			return u.Id, nil
		},
		tokenRepository,
	)

	gh, err := a.NewGoogleHandlers(
		cfg.GOOGLE_CLIENT_ID,
		cfg.GOOGLE_CLIENT_SECRET,
		cfg.GOOGLE_CLIENT_REDIRECT_URI,
	)
	if err != nil {
		panic(err)
	}

	requireAuth := aegis.RequireAuthMiddleware(jwtManager)

	return &AegisAuth{
		googleHandlers: gh,
		requireAuth:    requireAuth,
	}
}

type AegisAuth struct {
	googleHandlers aegis.Handlers
	requireAuth    func(http.Handler) http.Handler
}

func (a *AegisAuth) GoogleLoginHandler() http.Handler {
	return a.googleHandlers.Login
}

func (a *AegisAuth) GoogleCallbackHandler() http.Handler {
	return a.googleHandlers.Callback
}

func (a *AegisAuth) RequireAuth(handler http.Handler) http.Handler {
	return a.requireAuth(handler)
}
