package handler

import (
	"encoding/json"
	"net/http"

	userapi "github.com/joaovictorsl/go-backend-template/internal/core/user/api"
	"github.com/joaovictorsl/go-backend-template/internal/http/request"
)

func HandleGetUserById(api userapi.Api) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := request.GetUserId(r)

		u, err := api.GetUserById(r.Context(), userId)
		HandleError(err)

		raw, err := json.Marshal(u)
		HandleError(err)

		w.Write(raw)
	}
}
