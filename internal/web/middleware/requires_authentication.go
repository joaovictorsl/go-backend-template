package middleware

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/joaovictorsl/go-backend-template/internal/web/jwt"
	"github.com/joaovictorsl/go-backend-template/internal/web/request"
)

type JwtValidator interface {
	Validate(tokenString string) (*jwt.Claims, error)
}

func RequiresAuthentication(jwtSecret string, jwtValidator JwtValidator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c, err := r.Cookie("atok")
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			claims, err := jwtValidator.Validate(c.Value)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			sub, err := claims.GetSubject()
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
