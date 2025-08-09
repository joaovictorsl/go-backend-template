package request

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

func WithUserId(r *http.Request, userId uuid.UUID) {
	*r = *r.WithContext(context.WithValue(r.Context(), "user_id", userId))
}

func GetUserId(r *http.Request) uuid.UUID {
	id := r.Context().Value("user_id").(uuid.UUID)
	return id
}
