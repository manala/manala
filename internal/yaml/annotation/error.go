package annotation

type Error struct {
	error

	token Token
}

func NewError(err error, token Token) Error {
	return Error{
		error: err,
		token: token,
	}
}

func (e Error) Error() string {
	return e.error.Error()
}

func (e Error) Position() (int, int) {
	return e.token.Line, e.token.Column
}

func (e Error) Unwrap() error {
	return e.error
}
