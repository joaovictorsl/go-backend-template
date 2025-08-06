package middleware

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/joaovictorsl/go-backend-template/internal/http/request"
)

func RequiresAuthentication(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c, err := r.Cookie("access_token")
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			tok, err := jwt.Parse(c.Value, func(t *jwt.Token) (any, error) {
				return []byte(jwtSecret), nil
			})
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			sub, err := tok.Claims.GetSubject()
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			userId, err := uuid.Parse(sub)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			request.WithUserId(r, userId)
			next.ServeHTTP(w, r)
		})
	}
}
