package serrors

import (
	"encoding/json"
)

func NewJSON(err error) Error {
	var arguments []any
	message := err.Error()

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

	return New(message).
		WithArguments(arguments...)
}
