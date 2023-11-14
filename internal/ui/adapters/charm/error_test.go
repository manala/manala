package charm

import (
	"bytes"
	"manala/internal/serrors"
	"manala/internal/testing/heredoc"
)

func (s *Suite) TestError() {
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
				 ⨯ error
				   foo=bar
			`,
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			out := &bytes.Buffer{}
			err := &bytes.Buffer{}

			adapter := New(nil, out, err)

			adapter.Error(test.err)

			s.Empty(out)
			heredoc.Equal(s.Assert(), test.expected, err.String())
		})
	}
}
