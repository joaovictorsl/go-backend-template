package handler

import (
	"encoding/json"
	"net/http"

	user "github.com/joaovictorsl/go-backend-template/internal/core/user/usecase"
	"github.com/joaovictorsl/go-backend-template/internal/web"
	"github.com/joaovictorsl/go-backend-template/internal/web/request"
)

func HandleGetUser(u *user.UseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := request.GetUserId(r)

		u, err := u.Get(r.Context(), userId)
		if err != nil {
			web.HandleError(err)
		}

		raw, _ := json.Marshal(u)
		w.Write(raw)
	}
}
