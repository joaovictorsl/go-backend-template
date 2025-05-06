package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joaovictorsl/aegis/token"
	"github.com/joaovictorsl/go-backend-template/internal/auth"
	"github.com/joaovictorsl/go-backend-template/internal/config"
	"github.com/joaovictorsl/go-backend-template/internal/database"
	"github.com/joaovictorsl/go-backend-template/internal/router"
	"github.com/joaovictorsl/go-backend-template/internal/user"
)

func main() {
	cfg := config.NewConfig()
	db := database.NewDatabase(cfg)

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	deleteUserHandler := user.DeleteUserHTTPHandler(cfg, userService)

	a := auth.New(cfg, userService, token.NewInMemoryRepository())

	router := router.New()

	router.GET("/auth/google", a.GoogleLoginHandler())
	router.GET("/auth/google/callback", a.GoogleCallbackHandler())

	router.DELETE("/users/me", a.RequireAuth(deleteUserHandler))

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.PORT), router))
}
