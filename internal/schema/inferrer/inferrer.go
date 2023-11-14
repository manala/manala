package inferrer

import "manala/internal/schema"

type Inferrer interface {
	Infer(schema schema.Schema) error
}
