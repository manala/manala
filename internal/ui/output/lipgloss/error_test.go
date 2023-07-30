package lipgloss

import (
	"bytes"
	"manala/internal/errors/serrors"
	"manala/internal/testing/heredoc"
)

func (s *Suite) TestError() {
	tests := []struct {
		test     string
		err      error
		expected string
	}{
		{
			test:     "Empty",
			err:      serrors.New(""),
			expected: "",
		},
		{
			test: "Error",
			err:  serrors.New("error"),
			expected: heredoc.Doc(`
				  ⨯ error
			`),
		},
		{
			test: "Arguments",
			err: serrors.New("error").
				WithArguments("foo", "bar"),
			expected: heredoc.Doc(`
				  ⨯ error                              foo=bar
			`),
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			out := &bytes.Buffer{}
			err := &bytes.Buffer{}

			output := New(out, err)

			output.Error(test.err)

			s.Empty(out)
			s.Equal(test.expected, err.String())
		})
	}
}
