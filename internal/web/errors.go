package web

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/joaovictorsl/go-backend-template/internal/core"
)

func HttpErrResponse(w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"message": "%s"}`, msg)
}

func HandleError(err error) {
	panic(err)
}

type HttpError struct {
	status  int
	message string
}

func HttpErrorFrom(err error) HttpError {
	var (
		status  int    = http.StatusInternalServerError
		message string = "Something went wrong on our side"
	)

	var coreErr core.Error
	if !errors.As(err, &coreErr) {
		return HttpError{
			status,
			message,
		}
	}

	if errors.Is(coreErr, core.ErrNotFound) {
		status = http.StatusNotFound
		message = "We couldn't find what you were looking for"
	} else {
		slog.Error(
			"matching core error",
			slog.Any("error", coreErr),
		)
	}

	return HttpError{
		status,
		message,
	}
}

func (err HttpError) Error() string {
	return err.message
}

func (err HttpError) Status() int {
	return err.status
}

func (err HttpError) MarshalJSON() ([]byte, error) {
	jsonErr := fmt.Sprintf(`{"message": "%s"}`, err.message)
	return []byte(jsonErr), nil
}
