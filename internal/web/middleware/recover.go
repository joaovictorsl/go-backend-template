package middleware

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/joaovictorsl/go-backend-template/internal/web"
)

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			rawErr := recover()
			if rawErr == nil {
				return
			}

			castedErr, ok := rawErr.(error)
			if !ok {
				slog.Error(
					"casting recover value to error",
					slog.Any("value", rawErr),
				)
				return
			}

			err := web.HttpErrorFrom(castedErr)

			raw, _ := json.Marshal(err)
			w.WriteHeader(err.Status())
			w.Write(raw)
		}()

		next.ServeHTTP(w, r)
	})
}
