package annotation

import "github.com/manala/manala/internal/parsing"

// ErrorAt creates a parsing.Error positioned at the given token.
func ErrorAt(err error, token Token) *parsing.Error {
	return &parsing.Error{
		Err:    err,
		Line:   token.Line,
		Column: token.Column,
	}
}
