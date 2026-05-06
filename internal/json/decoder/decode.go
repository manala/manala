package decoder

import (
	"bytes"
	"encoding/json"

	jsonerrors "github.com/manala/manala/internal/json/errors"
)

// Decode decodes JSON bytes into the provided value using json.Number for numeric values,
// and returns an enhanced error with position information if decoding fails.
func Decode(data []byte, value any) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()

	if err := decoder.Decode(value); err != nil {
		return jsonerrors.From(err, string(data))
	}

	return nil
}
