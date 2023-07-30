package yaml

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type ValidationSuite struct{ suite.Suite }

func TestValidationSuite(t *testing.T) {
	suite.Run(t, new(ValidationSuite))
}

func (s *ValidationSuite) TestNormalizePath() {
	decorator := NewNodeValidationResultPathErrorDecorator(nil)

	tests := []struct {
		test     string
		path     string
		expected string
	}{
		{
			test:     "Root",
			path:     "(root)",
			expected: "$",
		},
		{
			test:     "FirstLevel",
			path:     "foo",
			expected: "$.foo",
		},
		{
			test:     "Levels",
			path:     "foo.bar",
			expected: "$.foo.bar",
		},
		{
			test:     "Index",
			path:     "foo.0.bar",
			expected: "$.foo[0].bar",
		},
		{
			test:     "IndexLast",
			path:     "foo.0",
			expected: "$.foo[0]",
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			s.Equal(test.expected, decorator.normalizePath(test.path))
		})
	}
}
