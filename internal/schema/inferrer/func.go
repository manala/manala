package inferrer

import "github.com/manala/manala/internal/schema"

type Func struct {
	schemaFunc func(schema schema.Schema) error
}

func NewFunc(schemaFunc func(schema schema.Schema) error) *Func {
	return &Func{
		schemaFunc: schemaFunc,
	}
}

func (inf *Func) Infer(schema schema.Schema) error {
	return inf.schemaFunc(schema)
}
