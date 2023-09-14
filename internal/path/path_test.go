package path

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type Suite struct{ suite.Suite }

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) TestJoin() {
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
			path := Path(test.path)

			path = path.Join(test.seg)

			s.Equal(test.expected, path.String())
		})
	}
}
