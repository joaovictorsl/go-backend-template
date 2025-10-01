package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/joaovictorsl/go-backend-template/internal/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func GetProviders(cfg *config.Config) map[string]Provider {
	providers := make(map[string]Provider)

	providers["google"] = &providerImpl{
		Config: &oauth2.Config{
			ClientID:     cfg.GoogleClientId,
			ClientSecret: cfg.GoogleClientSecret,
			RedirectURL:  cfg.GoogleClientRedirectUrl,
			Endpoint:     google.Endpoint,
			Scopes:       []string{"email"},
		},
		UserInfoURL:       "https://www.googleapis.com/oauth2/v2/userinfo",
		Name:              "google",
		ParseProviderUser: parseGoogleUser,
	}

	return providers
}

type ProviderUser struct {
	ID    string
	Email string
}

type Provider interface {
	AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string
	Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error)
	GetUser(ctx context.Context, tok *oauth2.Token) (*ProviderUser, error)
}

type providerImpl struct {
	*oauth2.Config
	UserInfoURL       string
	Name              string
	ParseProviderUser func(raw []byte) (*ProviderUser, error)
}

func (p *providerImpl) GetUser(ctx context.Context, tok *oauth2.Token) (*ProviderUser, error) {
	client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(tok))

	resp, err := client.Get(p.UserInfoURL)
	if err != nil {
		return nil, fmt.Errorf("getting user from oauth provider: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body from %s with %d status code: %w", p.Name, resp.StatusCode, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected 200 status code, got %d from %s", resp.StatusCode, p.Name)
	}

	u, err := p.ParseProviderUser(raw)
	if err != nil {
		return nil, fmt.Errorf("parsing provider user: %w", err)
	}

	return u, nil
}

func parseGoogleUser(raw []byte) (*ProviderUser, error) {
	var userInfo struct {
		ID    string `json:"id"`
		Email string `json:"email"`
	}
	if err := json.Unmarshal(raw, &userInfo); err != nil {
		return nil, fmt.Errorf("unmarshaling google user: %w %s", err, string(raw))
	}
	return &ProviderUser{
		ID:    userInfo.ID,
		Email: userInfo.Email,
	}, nil
}
