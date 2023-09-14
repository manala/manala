package inferrer

import (
	"manala/internal/schema"
	"manala/internal/serrors"
)

func (s *Suite) TestChainErrors() {
	tests := []struct {
		test      string
		inferrers []Inferrer
		expected  *serrors.Assert
	}{
		{
			test: "Error",
			inferrers: []Inferrer{
				NewFunc(func(_ schema.Schema) error {
					return serrors.New("foo")
				}),
			},
			expected: &serrors.Assert{
				Type:    serrors.Error{},
				Message: "foo",
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			schema := schema.Schema{"foo": "bar"}

			err := NewChain(test.inferrers...).Infer(schema)

			serrors.Equal(s.Assert(), test.expected, err)
		})
	}
}

func (s *Suite) TestChain() {
	tests := []struct {
		test      string
		inferrers []Inferrer
		expected  schema.Schema
	}{
		{
			test:      "Empty",
			inferrers: []Inferrer{},
			expected:  schema.Schema{"foo": "bar"},
		},
		{
			test: "Inferrers",
			inferrers: []Inferrer{
				NewFunc(func(schema schema.Schema) error {
					schema["foo"] = "baz"
					return nil
				}),
				NewFunc(func(schema schema.Schema) error {
					schema["bar"] = "baz"
					return nil
				}),
			},
			expected: schema.Schema{"foo": "baz", "bar": "baz"},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			schema := schema.Schema{"foo": "bar"}

			err := NewChain(test.inferrers...).Infer(schema)

			s.NoError(err)
			s.Equal(test.expected, schema)
		})
	}
}
