package router

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joaovictorsl/go-backend-template/internal/http/errs"
)

func SetupRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(recoverMiddleware)
	r.Use(middleware.Logger)
	return r
}

func recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				httpErr, ok := err.(errs.HTTPError)
				if !ok {
					log.Printf("ERROR: %v", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				log.Println(httpErr.Err)
				w.WriteHeader(httpErr.Code)
				w.Write([]byte(httpErr.Message))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
