package core

import "errors"

var (
	ErrNotFound       = errors.New("not found")
	ErrTooManyResults = errors.New("too many results")
)
