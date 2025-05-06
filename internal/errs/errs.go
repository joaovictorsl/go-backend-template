package errs

type Error struct {
	Code    int
	Message string
	Err     error
}

func HandleIfError(err error, code int, message string) {
	if err != nil {
		panic(Error{
			Code:    code,
			Message: message,
			Err:     err,
		})
	}
}
