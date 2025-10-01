package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	MinSecretSize uint = 32
)

var (
	ErrInvalidToken   = errors.New("token is invalid or expired")
	ErrTokenExpired   = errors.New("token has expired")
	ErrSecretTooShort = fmt.Errorf("invalid secret size: must be at least %d characters", MinSecretSize)
)

type Claims struct {
	jwt.RegisteredClaims
	UserID uuid.UUID `json:"user_id"`
}

type TokenManager struct {
	hmacSecret []byte
	ttl        time.Duration
}

func NewTokenManager(secret string, ttl time.Duration) (*TokenManager, error) {
	if uint(len(secret)) < MinSecretSize {
		return nil, ErrSecretTooShort
	}
	return &TokenManager{hmacSecret: []byte(secret), ttl: ttl}, nil
}

func (tm *TokenManager) Generate(userID uuid.UUID) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(tm.ttl)
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			Subject:   userID.String(),
		},
		UserID: userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokStr, err := token.SignedString(tm.hmacSecret)

	return tokStr, expiresAt, err
}

func (tm *TokenManager) Validate(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return tm.hmacSecret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}

		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
