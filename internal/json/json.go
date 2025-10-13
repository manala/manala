package json

import (
	"bytes"
	"encoding/json"

	"manala/internal/serrors"
)

// Unmarshal decodes JSON-encoded data.
func Unmarshal(data []byte, value any) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()

	if err := decoder.Decode(value); err != nil {
		return serrors.NewJSON(err)
	}

	return nil
}
