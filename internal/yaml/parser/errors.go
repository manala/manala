package parser

import (
	"errors"

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
	// Syntax error
	if err, ok := errors.AsType[*yaml.SyntaxError](err); ok {
		return ErrorAt(serrors.New(err.Message), err.Token)
	}

	// Type error
	if err, ok := errors.AsType[*yaml.TypeError](err); ok {
		return ErrorAt(serrors.New(err.Error()), err.Token)
	}

	// Overflow error
	if err, ok := errors.AsType[*yaml.OverflowError](err); ok {
		return ErrorAt(serrors.New(err.Error()), err.Token)
	}

	// Duplicate key error
	if err, ok := errors.AsType[*yaml.DuplicateKeyError](err); ok {
		return ErrorAt(serrors.New(err.Message), err.Token)
	}

	// Unknown field error
	if err, ok := errors.AsType[*yaml.UnknownFieldError](err); ok {
		return ErrorAt(serrors.New(err.Message), err.Token)
	}

	// Unexpected node type error
	if err, ok := errors.AsType[*yaml.UnexpectedNodeTypeError](err); ok {
		return ErrorAt(serrors.New(err.Error()), err.Token)
	}

	// Exceeded max depth error
	if errors.Is(err, yaml.ErrExceededMaxDepth) {
		return &parsing.Error{
			Err: serrors.New(err.Error()),
		}
	}

	return &parsing.Error{Err: err}
}
