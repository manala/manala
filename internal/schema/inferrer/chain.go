package inferrer

import "github.com/manala/manala/internal/schema"

type Chain struct {
	inferrers []Inferrer
}

func NewChain(inferrers ...Inferrer) *Chain {
	return &Chain{
		inferrers: inferrers,
	}
}

func (inf *Chain) Infer(schema schema.Schema) error {
	// Range over inferrers
	for _, inferrer := range inf.inferrers {
		// Infer
		if err := inferrer.Infer(schema); err != nil {
			return err
		}
	}

	return nil
}
