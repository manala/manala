package yaml

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type JsonPathNormalizerSuite struct{ suite.Suite }

func TestJsonPathNormalizerSuite(t *testing.T) {
	suite.Run(t, new(JsonPathNormalizerSuite))
}

func (s *JsonPathNormalizerSuite) TestNormalize() {
	tests := []struct {
		name     string
		actual   string
		expected string
	}{
		{
			name:     "Root",
			actual:   "(root)",
			expected: "$",
		},
		{
			name:     "First Level",
			actual:   "foo",
			expected: "$.foo",
		},
		{
			name:     "Levels",
			actual:   "foo.bar",
			expected: "$.foo.bar",
		},
		{
			name:     "Index",
			actual:   "foo.0.bar",
			expected: "$.foo[0].bar",
		},
		{
			name:     "Index Last",
			actual:   "foo.0",
			expected: "$.foo[0]",
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			normalizer := NewJsonPathNormalizer(test.actual)
			s.Equal(test.expected, normalizer.Normalize())
		})
	}
}
