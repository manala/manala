package schema

import "github.com/manala/manala/internal/json"

type Schema map[string]any

func MustParse(source []byte) Schema {
	var schema Schema

	err := json.Unmarshal(source, &schema)
	if err != nil {
		panic(err)
	}

	return schema
}
