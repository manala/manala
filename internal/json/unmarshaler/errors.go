package unmarshaler

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/manala/manala/internal/parsing"
)

// ErrorAt creates a parsing.Error positioned at the given offset.
func ErrorAt(err error, src string, offset int64) *parsing.Error {
	e := &parsing.Error{
		Err: err,
	}

	if src == "" {
		return e
	}

	// Compute position
	e.Line, e.Column = 1, 1
	for _, r := range src[:offset-1] {
		if r == '\n' {
			e.Line++
			e.Column = 1
		} else {
			e.Column++
		}
	}

	return e
}

// ErrorFrom converts a json error into a parsing.Error, extracting position from the error offset.
func ErrorFrom(err error, src string) *parsing.Error {
	// Syntax error
	if err, ok := errors.AsType[*json.SyntaxError](err); ok {
		return ErrorAt(
			err,
			src, err.Offset,
		)
	}

	// Unmarshal type error
	if err, ok := errors.AsType[*json.UnmarshalTypeError](err); ok {
		if err.Struct != "" || err.Field != "" {
			return ErrorAt(
				fmt.Errorf("wrong %s type for field \"%s\"", err.Value, err.Field),
				src, err.Offset,
			)
		}

		return ErrorAt(
			fmt.Errorf("wrong %s value type", err.Value),
			src, err.Offset,
		)
	}

	// Unmarshal type error
	if _, ok := errors.AsType[*json.InvalidUnmarshalError](err); ok {
		return &parsing.Error{
			Err: errors.New("invalid unmarshal argument"),
		}
	}

	return &parsing.Error{Err: err}
}
