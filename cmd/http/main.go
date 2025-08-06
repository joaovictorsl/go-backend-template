package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/joaovictorsl/go-backend-template/internal/config"
	userapi "github.com/joaovictorsl/go-backend-template/internal/core/user/api"
	userrepository "github.com/joaovictorsl/go-backend-template/internal/core/user/repository"
	userusecase "github.com/joaovictorsl/go-backend-template/internal/core/user/usecase"
	"github.com/joaovictorsl/go-backend-template/internal/database"
	"github.com/joaovictorsl/go-backend-template/internal/http/middleware"
	"github.com/joaovictorsl/go-backend-template/internal/http/router"
	"github.com/joho/godotenv"
)

func main() {
	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, Level: slog.LevelWarn}),
	)
	slog.SetDefault(logger)

	godotenv.Load()

	cfg := config.New()
	db := database.NewDatabase(cfg)

	userRepository := userrepository.New(db)
	userService := userusecase.New(userRepository)
	userApi := userapi.New(userService)

	r := chi.NewRouter()
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Timeout(cfg.Timeout))
	r.Use(middleware.Recover)

	router.SetupUserRoutes(userApi, r, cfg)

	addr := fmt.Sprintf(":%d", cfg.Port)
	if err := http.ListenAndServe(addr, r); err != nil {
		panic(err)
	}
}
