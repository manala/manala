package unmarshaler

import (
	"bytes"
	"encoding/json"
)

func Unmarshal(data []byte, value any) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()

	if err := decoder.Decode(value); err != nil {
		return ErrorFrom(err, string(data))
	}

	return nil
}
