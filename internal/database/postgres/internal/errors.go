package internal

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/joaovictorsl/go-backend-template/internal/core"
)

func MapError(err error) (error, bool) {
	isMapped := true

	if errors.Is(err, pgx.ErrNoRows) {
		err = core.ErrNotFound
	} else if errors.Is(err, pgx.ErrTooManyRows) {
		err = core.ErrTooManyResults
	} else {
		isMapped = false
	}

	return err, isMapped
}
