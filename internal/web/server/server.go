package server

import (
	"context"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/joaovictorsl/go-backend-template/internal/config"
	userusecase "github.com/joaovictorsl/go-backend-template/internal/core/user/usecase"
	"github.com/joaovictorsl/go-backend-template/internal/storage/inmemory"
	"github.com/joaovictorsl/go-backend-template/internal/storage/postgres"
	"github.com/joaovictorsl/go-backend-template/internal/web/jwt"
	"github.com/joaovictorsl/go-backend-template/internal/web/middleware"
)

type Server struct {
	mux            *chi.Mux
	authMiddleware func(http.Handler) http.Handler

	Config            *config.Config
	OAuthStore        *inmemory.KVCache
	RefreshTokenStore *postgres.RefreshTokenRepository
	UserStore         *postgres.UserRepository
	JwtManager        *jwt.TokenManager
	UserUseCase       *userusecase.UseCase
}

func (app *Server) Run(addr string) {
	timeout := app.Config.ShutdownTimeout
	srv := http.Server{
		Addr:    addr,
		Handler: app.mux,
	}

	go func() {
		slog.Info("web server started", slog.String("addr", addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error(
				"starting web server",
				slog.Any("error", err),
				slog.String("addr", srv.Addr),
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
			"shutting down",
			slog.Any("error", err),
		)
	} else {
		slog.Info("successful shutdown")
	}
}

func (app *Server) SetupRoutes() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(chimiddleware.Timeout(app.Config.RequestTimeout))
	r.Use(middleware.Recover)
	r.Use(middleware.PreventCSRF)

	app.mux = r
	app.authMiddleware = middleware.RequiresAuthentication(app.Config.JwtSecret, app.JwtManager)

	app.setupUser()
	app.setupAuth()
}
