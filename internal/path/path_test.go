package path_test

import (
	"testing"

	"github.com/manala/manala/internal/path"

	"github.com/stretchr/testify/suite"
)

type PathSuite struct{ suite.Suite }

func TestPathSuite(t *testing.T) {
	suite.Run(t, new(PathSuite))
}

func (s *PathSuite) TestJoin() {
	tests := []struct {
		test     string
		path     string
		seg      string
		expected string
	}{
		{
			test:     "Root",
			path:     "",
			seg:      "foo",
			expected: "foo",
		},
		{
			test:     "Leaf",
			path:     "foo",
			seg:      "bar",
			expected: "foo.bar",
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			path := path.Path(test.path)

			path = path.Join(test.seg)

			s.Equal(test.expected, path.String())
		})
	}
}
