package inferrer

import "manala/internal/schema"

func NewFunc(schemaFunc func(schema schema.Schema) error) *Func {
	return &Func{
		schemaFunc: schemaFunc,
	}
}

type Func struct {
	schemaFunc func(schema schema.Schema) error
}

func (inf *Func) Infer(schema schema.Schema) error {
	return inf.schemaFunc(schema)
}
