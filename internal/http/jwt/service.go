package jwt

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joaovictorsl/go-backend-template/internal/config"
)

type Claims struct {
	jwt.RegisteredClaims
	UserId uint
}
type Token struct {
	UserId    uint
	JWT       string
	CreatedAt time.Time
	ExpiresAt time.Time
}

type Service interface {
	NewToken(userId uint, iss string, secret string, exp time.Duration) (Token, error)
	StoreRefreshToken(ctx context.Context, tok Token) error
	ValidateAccessToken(tok string) (Claims, error)
	ValidateRefreshToken(ctx context.Context, tok string) (Claims, error)
}

func NewService(cfg *config.Config, tokenRepository Repository) Service {
	return &ServiceImpl{
		cfg:             cfg,
		tokenRepository: tokenRepository,
	}
}

type ServiceImpl struct {
	cfg             *config.Config
	tokenRepository Repository
}

func (s *ServiceImpl) NewToken(userId uint, iss string, secret string, exp time.Duration) (Token, error) {
	createdAt := time.Now().UTC()
	expiresAt := createdAt.Add(exp)
	claims := Claims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(createdAt),
			NotBefore: jwt.NewNumericDate(createdAt),
			Issuer:    iss,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return Token{}, fmt.Errorf("failed to generate token: %w", err)
	}

	return Token{
		UserId:    userId,
		JWT:       signedToken,
		CreatedAt: createdAt,
		ExpiresAt: expiresAt,
	}, nil
}

func (s *ServiceImpl) StoreRefreshToken(ctx context.Context, tok Token) error {
	return s.tokenRepository.StoreToken(ctx, tok)
}

func (s *ServiceImpl) ValidateAccessToken(tok string) (Claims, error) {
	return parseToken(tok, s.cfg.JWT_SECRET)
}

func (s *ServiceImpl) ValidateRefreshToken(ctx context.Context, tok string) (Claims, error) {
	claims, err := parseToken(tok, s.cfg.JWT_SECRET)
	if err != nil {
		return Claims{}, err
	}

	token, err := s.tokenRepository.GetToken(ctx, claims.UserId)
	if err != nil {
		return Claims{}, err
	}

	if token.JWT != tok {
		return Claims{}, ErrInvalidToken{Reason: "refresh token was replaced"}
	} else if token.ExpiresAt.Before(time.Now()) {
		return Claims{}, ErrInvalidToken{Reason: "refresh token is expired"}
	}

	return claims, nil
}

func parseToken(tokString string, secret string) (Claims, error) {
	token, err := jwt.ParseWithClaims(tokString, nil, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnexpectedSigningMethod{Algorithm: token.Header["alg"]}
		}
		return []byte(secret), nil
	})
	if errors.Is(err, jwt.ErrTokenMalformed) {
		return Claims{}, ErrFailedToParseToken{Reason: err.Error()}
	}

	claims, ok := token.Claims.(Claims)
	if !ok || !token.Valid {
		return claims, ErrInvalidToken{Reason: err.Error()}
	}

	return claims, nil
}
