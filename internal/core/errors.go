package core

var (
	ErrNotFound = Error{"not found"}
)

type Error struct {
	msg string
}

func (err Error) Error() string {
	return err.msg
}
