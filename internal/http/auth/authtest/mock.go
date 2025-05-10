package authtest

import (
	"context"
	"net/http"
)

type AuthMock struct {
	UserId *string
}

func (m *AuthMock) RequireAuth(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if m.UserId == nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		rCtx := context.WithValue(r.Context(), "user_id", *m.UserId)
		r = r.WithContext(rCtx)
		handler.ServeHTTP(w, r)
	}
}
