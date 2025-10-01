package auth

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	UserId    uuid.UUID
	Value     uuid.UUID
	ExpiresAt time.Time
}
