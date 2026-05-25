package log_test

import (
	"bytes"
	"testing"

	"github.com/manala/manala/internal/errors/serror"
	"github.com/manala/manala/internal/log"
	"github.com/manala/manala/internal/output"
	"github.com/manala/manala/internal/testing/heredoc"

	"github.com/stretchr/testify/suite"
)

type ErrorSuite struct{ suite.Suite }

func TestErrorSuite(t *testing.T) {
	suite.Run(t, new(ErrorSuite))
}

func (s *ErrorSuite) TestDump() {
	tests := []struct {
		test     string
		dump     string
		expected string
	}{
		{
			test: "Empty",
			dump: "",
			expected: heredoc.Doc(`
				 ✖ error
			`),
		},
		{
			test: "SingleLine",
			dump: "dump",
			expected: heredoc.Doc(`
				 ✖ error

				    │ dump
			`),
		},
		{
			test: "MultiLine",
			dump: "dump\ndump",
			expected: heredoc.Doc(`
				 ✖ error

				    │ dump
				    │ dump
			`),
		},
		{
			test: "TrailingLine",
			dump: "dump\ndump\n",
			expected: heredoc.Doc(`
				 ✖ error

				    │ dump
				    │ dump
			`),
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			out := &bytes.Buffer{}
			logger := log.New(output.NewDetached(out))

			err := serror.New("error").WithDump(test.dump)
			logger.Error(err)

			s.Equal(test.expected, out.String())
		})
	}
}
