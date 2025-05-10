package errs

type HTTPError struct {
	Code    int
	Message string
	Err     error
}

func HandleIfError(err error, code int, message string) {
	if err != nil {
		panic(HTTPError{
			Code:    code,
			Message: message,
			Err:     err,
		})
	}
}
