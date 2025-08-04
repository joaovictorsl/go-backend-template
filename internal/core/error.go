package core

type AppError struct {
	message string
}

func (ae *AppError) Error() string {
	return ae.message
}

var (
	ErrNotFound = &AppError{message: "not found"}
)
