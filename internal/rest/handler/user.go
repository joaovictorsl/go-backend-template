package handler

import (
	"encoding/json"
	"net/http"

	user "github.com/joaovictorsl/go-backend-template/internal/core/user/usecase"
	"github.com/joaovictorsl/go-backend-template/internal/rest/request"
)

func HandleGetUser(u *user.UseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := request.GetUserId(r)

		u, err := u.Get(r.Context(), userId)
		HandleError(err)

		raw, _ := json.Marshal(u)
		w.Write(raw)
	}
}
