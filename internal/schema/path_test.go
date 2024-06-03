package schema_test

import (
	"manala/internal/schema"
	"testing"

	"github.com/stretchr/testify/suite"
)

type PathSuite struct{ suite.Suite }

func TestPathSuite(t *testing.T) {
	suite.Run(t, new(PathSuite))
}

func (s *PathSuite) TestFieldPath() {
	tests := []struct {
		test     string
		field    string
		expected string
	}{
		{
			test:     "Root",
			field:    "(root)",
			expected: "",
		},
		{
			test:     "FirstLevel",
			field:    "foo",
			expected: "foo",
		},
		{
			test:     "Levels",
			field:    "foo.bar",
			expected: "foo.bar",
		},
		{
			test:     "Index",
			field:    "foo.0.bar",
			expected: "foo[0].bar",
		},
		{
			test:     "IndexLast",
			field:    "foo.0",
			expected: "foo[0]",
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			s.Equal(test.expected, schema.FieldPath(test.field).String())
		})
	}
}
