package charm_test

import (
	"bytes"
	"testing"

	"manala/internal/serrors"
	"manala/internal/testing/heredoc"
	"manala/internal/ui/adapters/charm"

	"github.com/stretchr/testify/suite"
)

type ErrorSuite struct{ suite.Suite }

func TestErrorSuite(t *testing.T) {
	suite.Run(t, new(ErrorSuite))
}

func (s *ErrorSuite) Test() {
	tests := []struct {
		test     string
		err      error
		expected string
	}{
		{
			test: "Empty",
			err:  serrors.New(""),
			expected: `
			`,
		},
		{
			test: "Error",
			err:  serrors.New("error"),
			expected: `
				 ⨯ error
			`,
		},
		{
			test: "Arguments",
			err: serrors.New("error").
				WithArguments("foo", "bar"),
			expected: `
				 ⨯ error                            foo=bar
			`,
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			out := &bytes.Buffer{}
			err := &bytes.Buffer{}

			adapter := charm.New(nil, out, err)

			adapter.Error(test.err)

			s.Empty(out)
			heredoc.Equal(s.T(), test.expected, err)
		})
	}
}
