package jwt_test

import (
	"crypto/rand"
	"testing"
	"time"

	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/joaovictorsl/go-backend-template/internal/web/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func randomString(length uint) string {
	s := make([]byte, length)
	rand.Read(s)
	return string(s)
}

func TestNewTokenManager(t *testing.T) {
	tests := []struct {
		name   string
		secret string
		ttl    time.Duration
		err    error
	}{
		{
			"should return error when secret size is too short",
			randomString(jwt.MinSecretSize - 1),
			time.Second,
			jwt.ErrSecretTooShort,
		},
		{
			"should return jwt manager when secret size is minimum",
			randomString(jwt.MinSecretSize),
			time.Second,
			nil,
		},
		{
			"should return jwt manager when secret size is greater than minimum",
			randomString(jwt.MinSecretSize + 1),
			time.Second,
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := jwt.NewTokenManager(tt.secret, tt.ttl)
			assert.ErrorIs(t, tt.err, gotErr)
			if tt.err != nil {
				assert.Nil(t, got)
			} else {
				assert.NotNil(t, got)
			}
		})
	}
}

func TestGenerateAndValidateTokenManager(t *testing.T) {
	defaultSecret := randomString(jwt.MinSecretSize)
	defaultUserId, _ := uuid.NewV7()
	gojwt.TimePrecision = time.Nanosecond

	tests := []struct {
		name    string
		ttl     time.Duration
		wantErr bool
	}{
		{
			"should return error when token expired",
			time.Nanosecond,
			true,
		},
		{
			"should return claims when token valid",
			time.Second,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm, err := jwt.NewTokenManager(defaultSecret, tt.ttl)
			require.NoError(t, err)

			tokStr, expiresAt, err := tm.Generate(defaultUserId)
			require.NoError(t, err)

			time.Sleep(time.Millisecond)

			claims, err := tm.Validate(tokStr)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, defaultUserId, claims.UserID)
				assert.Equal(t, defaultUserId.String(), claims.Subject)
				assert.WithinDuration(t, expiresAt, claims.ExpiresAt.Time, time.Second)
				assert.WithinDuration(t, expiresAt, claims.IssuedAt.Add(tt.ttl), time.Second)
			}
		})
	}
}
