package web

import (
	"fmt"
	"net/http"

	"github.com/joaovictorsl/go-backend-template/internal/core"
)

type HttpError struct {
	status  int
	message string
	error   error
}

func HttpErrorFrom(err error) HttpError {
	var (
		status  int
		message string
	)

	switch err {
	case core.ErrNotFound:
		status = http.StatusNotFound
		message = "We couldn't find what you were looking for"
	default:
		status = http.StatusInternalServerError
		message = "Something went wrong on our side"
	}

	return HttpError{
		status,
		message,
		err,
	}
}

func (err HttpError) Error() string {
	return err.message
}

func (err HttpError) Unwrap() error {
	return err.error
}

func (err HttpError) Status() int {
	return err.status
}

func (err HttpError) MarshalJSON() ([]byte, error) {
	jsonErr := fmt.Sprintf(`{"message": "%s"}`, err.message)
	return []byte(jsonErr), nil
}
