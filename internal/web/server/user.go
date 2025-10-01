package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/joaovictorsl/go-backend-template/internal/web/handler"
)

func (app *Server) setupUser() {
	app.mux.Group(func(r chi.Router) {
		r.Use(app.authMiddleware)

		r.Get("/users/me", handler.HandleGetUser(app.UserUseCase))
	})
}
