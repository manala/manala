package inferrer_test

import (
	"testing"

	"github.com/manala/manala/internal/schema"
	"github.com/manala/manala/internal/schema/inferrer"
	"github.com/manala/manala/internal/serrors"
	"github.com/manala/manala/internal/testing/errors"

	"github.com/stretchr/testify/suite"
)

type FuncSuite struct{ suite.Suite }

func TestFuncSuite(t *testing.T) {
	suite.Run(t, new(FuncSuite))
}

func (s *FuncSuite) TestErrors() {
	tests := []struct {
		test       string
		schemaFunc func(schema schema.Schema) error
		expected   errors.Assertion
	}{
		{
			test: "Error",
			schemaFunc: func(_ schema.Schema) error {
				return serrors.New("foo")
			},
			expected: &serrors.Assertion{
				Message: "foo",
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			schema := schema.Schema{"foo": "bar"}

			err := inferrer.NewFunc(test.schemaFunc).Infer(schema)

			errors.Equal(s.T(), test.expected, err)
		})
	}
}

func (s *FuncSuite) Test() {
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

			err := inferrer.NewFunc(test.schemaFunc).Infer(schema)

			s.Require().NoError(err)
			s.Equal(test.expected, schema)
		})
	}
}
