package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

// From try to convert an encoding/json error into an Error, extracting position from the error offset.
func From(err error, src string) error {
	// Syntax error
	if err, ok := errors.AsType[*json.SyntaxError](err); ok {
		return At(
			errors.New(err.Error()),
			src, err.Offset,
		)
	}

	// Unmarshal type error
	if err, ok := errors.AsType[*json.UnmarshalTypeError](err); ok {
		if err.Struct != "" || err.Field != "" {
			return At(
				fmt.Errorf("wrong %s type for field \"%s\"", err.Value, err.Field),
				src, err.Offset,
			)
		}

		return At(
			fmt.Errorf("wrong %s value type", err.Value),
			src, err.Offset,
		)
	}

	// EOF errors - point at the end of the source
	if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
		return At(err, src, int64(len(src)))
	}

	return err
}
