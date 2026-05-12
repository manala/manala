package errors

type Error struct {
	error

	line, column int
}

func (e Error) Error() string {
	return e.error.Error()
}

func (e Error) Position() (int, int) {
	return e.line, e.column
}

func (e Error) Unwrap() error {
	return e.error
}
