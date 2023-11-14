package json

import (
	"encoding/json"
	"manala/internal/serrors"
)

func NewError(err error) serrors.Error {
	message := err.Error()
	arguments := []any{}

	switch err := err.(type) {
	case *json.SyntaxError:
		arguments = append(arguments, "offset", err.Offset)
	case *json.UnmarshalTypeError:
		arguments = append(arguments, "offset", err.Offset)
		if err.Struct != "" || err.Field != "" {
			message = "cannot unmarshal into struct field"
			arguments = append(arguments,
				"value", err.Value,
				"struct", err.Struct,
				"field", err.Field,
				"type", err.Type.String(),
			)

		} else {
			message = "cannot unmarshal into value"
			arguments = append(arguments,
				"value", err.Value,
				"type", err.Type.String(),
			)
		}
	case *json.InvalidUnmarshalError:
		message = "invalid unmarshal argument"
		arguments = append(arguments, "type", err.Type.String())
	}

	return serrors.New(message).
		WithArguments(arguments...)
}
