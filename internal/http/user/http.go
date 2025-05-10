package http

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/joaovictorsl/aegis"
	"github.com/joaovictorsl/go-backend-template/internal/config"
	"github.com/joaovictorsl/go-backend-template/internal/core/user/usecase"
	"github.com/joaovictorsl/go-backend-template/internal/http/errs"
)

func SetupRoutes(
	r *chi.Mux,
	cfg *config.Config,
	requireAuth func(http.HandlerFunc) http.HandlerFunc,
	userUseCase usecase.UseCase,
) {
	r.Delete("/users/me", requireAuth(deleteUserHandler(cfg, userUseCase)))
}

func deleteUserHandler(
	cfg *config.Config,
	userUseCase usecase.UseCase,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := aegis.GetUserIdFromContext(r.Context())

		ctx, cancel := context.WithTimeout(context.Background(), cfg.TIMEOUT)
		defer cancel()

		err := userUseCase.DeleteUserById(ctx, userId)
		errs.HandleIfError(err, http.StatusInternalServerError, "failed to delete user")

		w.WriteHeader(http.StatusNoContent)
	}
}
