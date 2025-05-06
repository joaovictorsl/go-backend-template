package user

import (
	"context"
	"net/http"

	"github.com/joaovictorsl/aegis"
	"github.com/joaovictorsl/go-backend-template/internal/config"
	"github.com/joaovictorsl/go-backend-template/internal/errs"
)

func DeleteUserHTTPHandler(
	cfg *config.Config,
	userService Service,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := aegis.GetUserIdFromContext(r.Context())

		ctx, cancel := context.WithTimeout(context.Background(), cfg.TIMEOUT)
		defer cancel()

		err := userService.DeleteUserById(ctx, userId)
		errs.HandleIfError(err, http.StatusInternalServerError, "failed to delete user")

		w.WriteHeader(http.StatusNoContent)
	})
}
