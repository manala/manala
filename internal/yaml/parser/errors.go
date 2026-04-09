package parser

import (
	"errors"
	"fmt"

	"github.com/manala/manala/internal/parsing"
	"github.com/manala/manala/internal/serrors"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/token"
)

// ErrorAt creates a parsing.Error positioned at the given token.
func ErrorAt(err error, token *token.Token) *parsing.Error {
	return &parsing.Error{
		Err:    err,
		Line:   token.Position.Line,
		Column: token.Position.Column,
	}
}

// ErrorFrom converts a go-yaml error into a parsing.Error, extracting position from the error itself.
func ErrorFrom(err error) *parsing.Error {
	// Type error
	if err, ok := errors.AsType[*yaml.TypeError](err); ok {
		return ErrorAt(
			serrors.New(fmt.Sprintf("field must be a %s", err.DstType)),
			err.GetToken(),
		)
	}

	if err, ok := errors.AsType[yaml.Error](err); ok {
		return ErrorAt(serrors.New(err.GetMessage()), err.GetToken())
	}

	// Exceeded max depth error
	if errors.Is(err, yaml.ErrExceededMaxDepth) {
		return &parsing.Error{
			Err: serrors.New("yaml exceeded max depth"),
		}
	}

	return &parsing.Error{Err: err}
}
