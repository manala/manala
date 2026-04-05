package schema

import "github.com/manala/manala/internal/json/unmarshaler"

type Schema map[string]any

func MustParse(source []byte) Schema {
	var schema Schema

	err := unmarshaler.Unmarshal(source, &schema)
	if err != nil {
		panic(err)
	}

	return schema
}
