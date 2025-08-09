package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/joaovictorsl/go-backend-template/internal/config"
	userservice "github.com/joaovictorsl/go-backend-template/internal/core/user/service"
	userusecase "github.com/joaovictorsl/go-backend-template/internal/core/user/usecase"
	"github.com/joaovictorsl/go-backend-template/internal/database/postgres"
	"github.com/joaovictorsl/go-backend-template/internal/web/middleware"
	"github.com/joaovictorsl/go-backend-template/internal/web/routes"
)

func main() {
	cfg := config.New()

	logger := slog.New(
		slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{AddSource: true, Level: cfg.LogLevel}),
	)
	slog.SetDefault(logger)

	db, err := postgres.New(cfg)
	if err != nil {
		return
	}
	defer db.Close()

	userRepository := &postgres.UserRepository{
		DB: db,
	}
	userService := &userservice.Service{
		UserStore: userRepository,
	}
	userUseCase := &userusecase.UseCase{
		UserService: userService,
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(chimiddleware.Timeout(cfg.RequestTimeout))
	r.Use(middleware.Recover)

	routes.SetupUser(userUseCase, r, cfg)
	routes.SetupOAuth(r, cfg)

	addr := fmt.Sprintf(":%d", cfg.Port)
	runServer(r, addr, cfg.ShutdownTimeout)
}

func runServer(r *chi.Mux, addr string, timeout time.Duration) {
	srv := http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		slog.Info("web server started", slog.String("addr", addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error(
				"web server failed",
				slog.Any("error", err),
			)
		}
	}()

	shutdown, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()
	<-shutdown.Done()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	slog.Info("shutting down", slog.String("timeout", timeout.String()))
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error(
			"shutdown failed",
			slog.Any("error", err),
		)
	} else {
		slog.Info("successful shutdown")
	}
}
