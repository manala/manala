package inferrer

import "manala/internal/schema"

func NewChain(inferrers ...Inferrer) *Chain {
	return &Chain{
		inferrers: inferrers,
	}
}

type Chain struct {
	inferrers []Inferrer
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
