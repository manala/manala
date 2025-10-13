package inferrer_test

import (
	"testing"

	"manala/internal/schema"
	"manala/internal/schema/inferrer"
	"manala/internal/serrors"

	"github.com/stretchr/testify/suite"
)

type ChainSuite struct{ suite.Suite }

func TestChainSuite(t *testing.T) {
	suite.Run(t, new(ChainSuite))
}

func (s *ChainSuite) TestErrors() {
	tests := []struct {
		test      string
		inferrers []inferrer.Inferrer
		expected  *serrors.Assertion
	}{
		{
			test: "Error",
			inferrers: []inferrer.Inferrer{
				inferrer.NewFunc(func(_ schema.Schema) error {
					return serrors.New("foo")
				}),
			},
			expected: &serrors.Assertion{
				Message: "foo",
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			schema := schema.Schema{"foo": "bar"}

			err := inferrer.NewChain(test.inferrers...).Infer(schema)

			serrors.Equal(s.T(), test.expected, err)
		})
	}
}

func (s *ChainSuite) Test() {
	tests := []struct {
		test      string
		inferrers []inferrer.Inferrer
		expected  schema.Schema
	}{
		{
			test:      "Empty",
			inferrers: []inferrer.Inferrer{},
			expected:  schema.Schema{"foo": "bar"},
		},
		{
			test: "Inferrers",
			inferrers: []inferrer.Inferrer{
				inferrer.NewFunc(func(schema schema.Schema) error {
					schema["foo"] = "baz"

					return nil
				}),
				inferrer.NewFunc(func(schema schema.Schema) error {
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

			err := inferrer.NewChain(test.inferrers...).Infer(schema)

			s.Require().NoError(err)
			s.Equal(test.expected, schema)
		})
	}
}
