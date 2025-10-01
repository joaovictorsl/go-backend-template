package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/joaovictorsl/go-backend-template/internal/core"
	"github.com/joaovictorsl/go-backend-template/internal/core/entity"
	"github.com/joaovictorsl/go-backend-template/internal/web"
	"github.com/justinas/nosurf"
	"golang.org/x/oauth2"
)

type OAuthStore interface {
	Get(state string) (verifier string, err error)
	Insert(state string, verifier string) error
	Remove(state string)
}

type UserStore interface {
	GetByProvider(ctx context.Context, provider, providerID string) (u entity.User, err error)
	Insert(ctx context.Context, email, provider, providerId string) (id uuid.UUID, err error)
}

type RefreshTokenStore interface {
	Insert(ctx context.Context, userId uuid.UUID, rTok uuid.UUID, expiresAt time.Time) error
	Get(ctx context.Context, rTok uuid.UUID) (RefreshToken, error)
	Update(ctx context.Context, oldRTok uuid.UUID, newRTok uuid.UUID, newRTokExpiresAt time.Time) error
	Delete(ctx context.Context, rTok uuid.UUID) error
}

type JwtGenerator interface {
	Generate(userID uuid.UUID) (string, time.Time, error)
}

func HandleOAuth(providers map[string]Provider, oauthStore OAuthStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		providerKey := r.PathValue("provider")
		p, ok := providers[providerKey]
		if !ok {
			web.HttpErrResponse(w, http.StatusBadRequest, fmt.Sprintf("%s is not a valid provider", providerKey))
			return
		}

		state := oauth2.GenerateVerifier()
		verifier := oauth2.GenerateVerifier()
		err := oauthStore.Insert(state, verifier)
		if err != nil {
			slog.Error(
				"inserting state and verifier in oauthStore",
				slog.Any("error", err),
			)
			web.HandleError(err)
		}

		pUrl := p.AuthCodeURL(state, oauth2.AccessTypeOnline, oauth2.S256ChallengeOption(verifier))
		http.Redirect(w, r, pUrl, http.StatusTemporaryRedirect)
	}
}

func HandleOAuthCallback(
	providers map[string]Provider,
	rTokTtl time.Duration,
	oauthStore OAuthStore,
	userStore UserStore,
	refreshTokenStore RefreshTokenStore,
	jwtGenerator JwtGenerator,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		providerKey := r.PathValue("provider")
		p, ok := providers[providerKey]
		if !ok {
			web.HttpErrResponse(w, http.StatusBadRequest, fmt.Sprintf("%s is not a valid provider", providerKey))
			return
		}

		query := r.URL.Query()
		code := query.Get("code")
		state := query.Get("state")
		verifier, err := oauthStore.Get(state)
		if errors.Is(err, core.ErrNotFound) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else if err != nil {
			web.HandleError(err)
		}
		oauthStore.Remove(state)

		tok, err := p.Exchange(r.Context(), code, oauth2.VerifierOption(verifier))
		if err != nil {
			slog.Error(
				"exchanging code",
				slog.Any("error", err),
				slog.String("provider", providerKey),
			)
			web.HandleError(err)
		}

		pu, err := p.GetUser(r.Context(), tok)
		if err != nil {
			slog.Error(
				"getting user from oauth provider",
				slog.Any("error", err),
				slog.String("provider", providerKey),
			)
			web.HandleError(err)
		}

		u, err := userStore.GetByProvider(r.Context(), providerKey, pu.ID)
		if errors.Is(err, core.ErrNotFound) {
			id, err := userStore.Insert(r.Context(), pu.Email, providerKey, pu.ID)
			if err != nil {
				slog.Error(
					"inserting user on oauth",
					slog.Any("error", err),
					slog.String("email", pu.Email),
					slog.String("provider", providerKey),
					slog.String("provider_id", pu.ID),
				)
				web.HandleError(err)
			}
			u.Id = id
		} else if err != nil {
			slog.Error(
				"getting user by provider on oauth",
				slog.Any("error", err),
				slog.String("provider", providerKey),
				slog.String("provider_user_id", pu.ID),
			)
			web.HandleError(err)
		}

		setCookies(w, r, u.Id, rTokTtl, jwtGenerator, refreshTokenStore)

		http.Redirect(w, r, "/home", http.StatusFound)
	}
}

func HandleRefresh(rTokTtl time.Duration, refreshTokenStore RefreshTokenStore, jwtGenerator JwtGenerator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rTokCookie, err := r.Cookie("rtok")
		if err != nil {
			web.HttpErrResponse(w, http.StatusBadRequest, "missing refresh token")
			return
		}

		rTokValue, err := uuid.Parse(rTokCookie.Value)
		if err != nil {
			web.HttpErrResponse(w, http.StatusBadRequest, "invalid refresh token")
			return
		}

		rTok, err := refreshTokenStore.Get(r.Context(), rTokValue)
		if err != nil {
			if errors.Is(err, core.ErrNotFound) {
				web.HttpErrResponse(w, http.StatusUnauthorized, "invalid refresh token")
				return
			}
			slog.Error(
				"retrieving refresh token on refresh",
				slog.Any("error", err),
				slog.String("rTok", rTokValue.String()),
			)
			web.HandleError(err)
		}

		if time.Now().After(rTok.ExpiresAt) {
			web.HttpErrResponse(w, http.StatusUnauthorized, "expired refresh token")
			return
		}

		setCookies(w, r, rTok.UserId, rTokTtl, jwtGenerator, refreshTokenStore)

		w.WriteHeader(http.StatusOK)
	}
}

func HandleSignOut(refreshTokenStore RefreshTokenStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rTokCookie, err := r.Cookie("rtok")
		if err == nil {
			rTokValue, err := uuid.Parse(rTokCookie.Value)
			if err == nil {
				if err := refreshTokenStore.Delete(r.Context(), rTokValue); err != nil {
					slog.Error(
						"deleting refresh token on signout",
						slog.Any("error", err),
						slog.String("rTok", rTokValue.String()),
					)
					web.HandleError(err)
				}
			}
		}

		deleteCookies(w)

		w.WriteHeader(http.StatusOK)
	}
}

func setCookies(
	w http.ResponseWriter,
	r *http.Request,
	userId uuid.UUID,
	rTokTtl time.Duration,
	jwtGenerator JwtGenerator,
	refreshTokenStore RefreshTokenStore,
) {
	rTok, _ := uuid.NewV7()
	aTok, aTokExpiresAt, err := jwtGenerator.Generate(userId)

	rTokExpiresAt := time.Now().Add(rTokTtl)
	err = refreshTokenStore.Insert(r.Context(), userId, rTok, rTokExpiresAt)
	if err != nil {
		slog.Error(
			"inserting refresh token",
			slog.Any("error", err),
			slog.String("user_id", userId.String()),
		)
		web.HandleError(err)
	}

	http.SetCookie(w, configCookie(
		"rtok",
		rTok.String(),
		rTokExpiresAt,
		true,
	))

	http.SetCookie(w, configCookie(
		"atok",
		aTok,
		aTokExpiresAt,
		true,
	))

	http.SetCookie(w, configCookie(
		nosurf.CookieName+"_client",
		nosurf.Token(r),
		aTokExpiresAt,
		false,
	))
}

func configCookie(
	name string,
	value string,
	expiresAt time.Time,
	httpOnly bool,
) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Expires:  expiresAt,
		Secure:   false,
		HttpOnly: httpOnly,
		SameSite: http.SameSiteLaxMode,
	}
}

func deleteCookies(
	w http.ResponseWriter,
) {
	http.SetCookie(w, &http.Cookie{
		Name:   "rtok",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	http.SetCookie(w, &http.Cookie{
		Name:   "atok",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	http.SetCookie(w, &http.Cookie{
		Name:   nosurf.CookieName + "_client",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
}
