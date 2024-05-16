package json

import (
	"bytes"
	"encoding/json"
	"manala/internal/serrors"
)

func Unmarshal(data []byte, value any) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()

	if err := decoder.Decode(value); err != nil {
		return serrors.NewJson(err)
	}

	return nil
}
