package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id        uuid.UUID
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
