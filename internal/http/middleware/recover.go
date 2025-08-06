package middleware

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/joaovictorsl/go-backend-template/internal/http/errs"
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
				log.Fatalf("Expected error, got %v from recover", rawErr)
			}

			err := errs.FromError(castedErr)

			raw, _ := json.Marshal(err)
			w.WriteHeader(err.Status())
			w.Write(raw)

			if err.Status() == http.StatusInternalServerError {
				slog.Error(
					"Unexpected error",
					slog.Any("error", err.Unwrap()),
					slog.String("stack", string(debug.Stack())),
				)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
