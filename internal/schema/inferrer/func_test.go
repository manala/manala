package inferrer

import (
	"manala/internal/schema"
	"manala/internal/serrors"
)

func (s *Suite) TestFuncErrors() {
	tests := []struct {
		test       string
		schemaFunc func(schema schema.Schema) error
		expected   *serrors.Assert
	}{
		{
			test: "Error",
			schemaFunc: func(_ schema.Schema) error {
				return serrors.New("foo")
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

			err := NewFunc(test.schemaFunc).Infer(schema)

			serrors.Equal(s.Assert(), test.expected, err)
		})
	}
}

func (s *Suite) TestFunc() {
	tests := []struct {
		test       string
		schemaFunc func(schema schema.Schema) error
		expected   schema.Schema
	}{
		{
			test: "Extend",
			schemaFunc: func(schema schema.Schema) error {
				schema["bar"] = "baz"
				return nil
			},
			expected: schema.Schema{"foo": "bar", "bar": "baz"},
		},
		{
			test: "Override",
			schemaFunc: func(schema schema.Schema) error {
				schema["foo"] = "baz"
				return nil
			},
			expected: schema.Schema{"foo": "baz"},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			schema := schema.Schema{"foo": "bar"}

			err := NewFunc(test.schemaFunc).Infer(schema)

			s.NoError(err)
			s.Equal(test.expected, schema)
		})
	}
}
