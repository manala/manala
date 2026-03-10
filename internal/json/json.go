package json

import (
	"bytes"
	"encoding/json"

	"github.com/manala/manala/internal/serrors"
)

func Unmarshal(data []byte, value any) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()

	if err := decoder.Decode(value); err != nil {
		return serrors.NewJSON(err)
	}

	return nil
}
