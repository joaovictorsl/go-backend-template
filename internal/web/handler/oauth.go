package handler

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/joaovictorsl/go-backend-template/internal/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func GetProviders(cfg *config.Config) map[string]*oauth2.Config {
	return map[string]*oauth2.Config{
		"google": {
			ClientID:     cfg.GoogleClientId,
			ClientSecret: cfg.GoogleClientSecret,
			RedirectURL:  "http://localhost:8080/oauth/google/callback",
			Endpoint:     google.Endpoint,
			Scopes:       []string{"email"},
		},
	}
}

func HandleOAuth(providers map[string]*oauth2.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		providerKey := r.PathValue("provider")
		p, ok := providers[providerKey]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// TODO: generate a random state and a random verifier
		// TODO: after that store in a redis K: state V: verifier
		// TODO: when the callback comes we check for the verifier
		// TODO: in redis and then use the verifier associated with the state
		// TODO: to exchange the code for a token
		pUrl := p.AuthCodeURL("backend", oauth2.AccessTypeOffline, oauth2.S256ChallengeOption("verifier"))
		http.Redirect(w, r, pUrl, http.StatusTemporaryRedirect)
	}
}

func HandleOAuthCallback(providers map[string]*oauth2.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		providerKey := r.PathValue("provider")
		p, ok := providers[providerKey]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		query := r.URL.Query()
		code := query.Get("code")
		state := query.Get("state")
		if state != "backend" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		tok, err := p.Exchange(r.Context(), code, oauth2.VerifierOption("verifier"))
		if err != nil {
			slog.Error(
				fmt.Sprintf("oauth callback failed for %s", providerKey),
				slog.Any("error", err),
			)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     fmt.Sprintf("%s_atok", providerKey),
			Value:    tok.AccessToken,
			Path:     "/",
			Domain:   "localhost",
			Expires:  tok.Expiry,
			Secure:   false,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		})

		if tok.RefreshToken != "" {
			http.SetCookie(w, &http.Cookie{
				Name:     fmt.Sprintf("%s_rtok", providerKey),
				Value:    tok.RefreshToken,
				Path:     "/",
				Domain:   "localhost",
				Secure:   false,
				HttpOnly: true,
				SameSite: http.SameSiteLaxMode,
			})
		}
	}
}
