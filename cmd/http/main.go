package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joaovictorsl/aegis/token"
	"github.com/joaovictorsl/go-backend-template/internal/config"
	userrepository "github.com/joaovictorsl/go-backend-template/internal/core/user/repository"
	userusecase "github.com/joaovictorsl/go-backend-template/internal/core/user/usecase"
	"github.com/joaovictorsl/go-backend-template/internal/database"
	"github.com/joaovictorsl/go-backend-template/internal/http/auth"
	"github.com/joaovictorsl/go-backend-template/internal/http/router"
	userhttp "github.com/joaovictorsl/go-backend-template/internal/http/user"
)

func main() {
	cfg := config.NewConfig()
	db := database.NewDatabase(cfg)

	userRepository := userrepository.New(db)
	userUseCase := userusecase.New(userRepository)

	a := auth.New(cfg, userUseCase, token.NewInMemoryRepository())

	r := router.SetupRouter()

	r.Get("/auth/google", a.GoogleLoginHandler())
	r.Get("/auth/google/callback", a.GoogleCallbackHandler())

	userhttp.SetupRoutes(r, cfg, a.RequireAuth, userUseCase)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.PORT), r))
}
