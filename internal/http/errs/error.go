package errs

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

func FromError(err error) HttpError {
	var (
		status  int
		message string
	)

	castedErr, ok := err.(*core.AppError)
	if !ok {
		status = http.StatusInternalServerError
		message = "Something went wrong on our side"
	} else if castedErr == core.ErrNotFound {
		status = http.StatusNotFound
		message = "We didn't find what you were looking for"
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
	return []byte(fmt.Sprintf(`{"message": "%s"}`, err.message)), nil
}
