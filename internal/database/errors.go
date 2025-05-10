package database

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/joaovictorsl/go-backend-template/internal/core/errs"
)

const (
	UniqueViolationCode = "23505"
)

func isUniqueViolationError(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == UniqueViolationCode
}

func isNoRowError(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}

func MapDatabaseError(err error) error {
	mappedErr := errors.New("unexpected error")

	if isUniqueViolationError(err) {
		mappedErr = errs.ErrDuplicated
	} else if isNoRowError(err) {
		mappedErr = errs.ErrNotFound
	}

	return mappedErr
}
