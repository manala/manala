package errors

import (
	"github.com/goccy/go-yaml/token"
)

type Error struct {
	error

	token *token.Token
}

func New(err error, token *token.Token) Error {
	return Error{
		error: err,
		token: token,
	}
}

func (e Error) Error() string {
	return e.error.Error()
}

func (e Error) Position() (int, int) {
	if e.token == nil {
		return 0, 0
	}
	return e.token.Position.Line, e.token.Position.Column
}

func (e Error) Unwrap() error {
	return e.error
}
