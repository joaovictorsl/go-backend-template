package internal

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/joaovictorsl/go-backend-template/internal/core"
)

func MapError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		err = fmt.Errorf("%w: %w", core.ErrNotFound, err)
	}

	return err
}
