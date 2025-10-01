package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/joaovictorsl/go-backend-template/internal/config"
	userservice "github.com/joaovictorsl/go-backend-template/internal/core/user/service"
	userusecase "github.com/joaovictorsl/go-backend-template/internal/core/user/usecase"
	"github.com/joaovictorsl/go-backend-template/internal/storage/inmemory"
	"github.com/joaovictorsl/go-backend-template/internal/storage/postgres"
	"github.com/joaovictorsl/go-backend-template/internal/web/jwt"
	"github.com/joaovictorsl/go-backend-template/internal/web/server"
)

func main() {
	cfg := config.New()

	logger := slog.New(
		slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{AddSource: true, Level: cfg.LogLevel}),
	)
	slog.SetDefault(logger)

	db, err := postgres.New(cfg)
	if err != nil {
		slog.Error(
			"creating new postgres instance",
			slog.Any("error", err),
		)
		return
	}
	defer db.Close()

	oauthStore := inmemory.New()
	jwtManager, err := jwt.NewTokenManager(cfg.JwtSecret, cfg.AccessTokenTTL)
	if err != nil {
		slog.Error(
			"creating new jwt TokenManager instance",
			slog.Any("error", err),
		)
		return
	}

	refreshTokenRepository := postgres.NewRefreshTokenRepository(db)

	userRepository := postgres.NewUserRepository(db)
	userService := userservice.New(userRepository)
	userUseCase := userusecase.New(userService)

	app := &server.Server{
		Config:            cfg,
		OAuthStore:        oauthStore,
		RefreshTokenStore: refreshTokenRepository,
		UserStore:         userRepository,
		JwtManager:        jwtManager,
		UserUseCase:       userUseCase,
	}
	app.SetupRoutes()

	addr := fmt.Sprintf(":%d", cfg.Port)
	app.Run(addr)
}
