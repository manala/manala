package inferrer

import "github.com/manala/manala/internal/schema"

type Inferrer interface {
	Infer(schema schema.Schema) error
}
