package source_test

import (
	"errors"
	"testing"

	"github.com/manala/manala/internal/errors/source"
	"github.com/manala/manala/internal/output"
	"github.com/manala/manala/internal/testing/heredoc"

	"github.com/stretchr/testify/suite"
)

type ErrorSuite struct{ suite.Suite }

func TestErrorSuite(t *testing.T) {
	suite.Run(t, new(ErrorSuite))
}

func (s *ErrorSuite) TestRender() {
	file := "/foo/bar"
	src := heredoc.Doc(`
		Lorem ipsum dolor sit amet, consectetur adipiscing elit.
		Etiam in sapien diam.

		Phasellus eu mi gravida, semper nisi non, euismod ex.
		Integer nec ornare mi.
	`)

	tests := []struct {
		test     string
		err      source.Error
		expected string
	}{
		{
			test: "All",
			err: source.Error{
				Origin:   source.Origin{File: file, Source: src},
				Position: testError{errors.New("message")},
				Line:     2,
				Column:   7,
			},
			expected: heredoc.Doc(`

				at /foo/bar:2:7

				  1 │ Lorem ipsum dolor sit amet, consectetur adipiscing elit.
				▶ 2 │ Etiam in sapien diam.
				    ├───────╯ message
				  3 │
				  4 │ Phasellus eu mi gravida, semper nisi non, euismod ex.
				  5 │ Integer nec ornare mi.
			`),
		},
		{
			test: "NoColumn",
			err: source.Error{
				Origin:   source.Origin{File: file, Source: src},
				Position: testError{errors.New("message")},
				Line:     2,
			},
			expected: heredoc.Doc(`

				at /foo/bar:2

				  1 │ Lorem ipsum dolor sit amet, consectetur adipiscing elit.
				▶ 2 │ Etiam in sapien diam.
				    ├ message
				  3 │
				  4 │ Phasellus eu mi gravida, semper nisi non, euismod ex.
				  5 │ Integer nec ornare mi.
			`),
		},
		{
			test: "NoPosition",
			err: source.Error{
				Origin:   source.Origin{File: file, Source: src},
				Position: testError{errors.New("message")},
			},
			expected: heredoc.Doc(`

				at /foo/bar

				  1 │ Lorem ipsum dolor sit amet, consectetur adipiscing elit.
				  2 │ Etiam in sapien diam.

				message
			`),
		},
		{
			test: "NoFile",
			err: source.Error{
				Origin:   source.Origin{Source: src},
				Position: testError{errors.New("message")},
				Line:     2,
				Column:   7,
			},
			expected: heredoc.Doc(`

				  1 │ Lorem ipsum dolor sit amet, consectetur adipiscing elit.
				▶ 2 │ Etiam in sapien diam.
				    ├───────╯ message
				  3 │
				  4 │ Phasellus eu mi gravida, semper nisi non, euismod ex.
				  5 │ Integer nec ornare mi.
			`),
		},
		{
			test: "NoFileNoColumn",
			err: source.Error{
				Origin:   source.Origin{Source: src},
				Position: testError{errors.New("message")},
				Line:     2,
			},
			expected: heredoc.Doc(`

				  1 │ Lorem ipsum dolor sit amet, consectetur adipiscing elit.
				▶ 2 │ Etiam in sapien diam.
				    ├ message
				  3 │
				  4 │ Phasellus eu mi gravida, semper nisi non, euismod ex.
				  5 │ Integer nec ornare mi.
			`),
		},
		{
			test: "NoFileNoPosition",
			err: source.Error{
				Origin:   source.Origin{Source: src},
				Position: testError{errors.New("message")},
			},
			expected: heredoc.Doc(`

				  1 │ Lorem ipsum dolor sit amet, consectetur adipiscing elit.
				  2 │ Etiam in sapien diam.

				message
			`),
		},
		{
			test: "EmptySource",
			err: source.Error{
				Origin:   source.Origin{File: file},
				Position: testError{errors.New("message")},
				Line:     2,
				Column:   7,
			},
			expected: heredoc.Doc(`

				at /foo/bar:2:7

				  1 │

				message
			`),
		},
		{
			test: "EmptySourceNoColumn",
			err: source.Error{
				Origin:   source.Origin{File: file},
				Position: testError{errors.New("message")},
				Line:     2,
			},
			expected: heredoc.Doc(`

				at /foo/bar:2

				  1 │

				message
			`),
		},
		{
			test: "EmptySourceNoPosition",
			err: source.Error{
				Origin:   source.Origin{File: file},
				Position: testError{errors.New("message")},
			},
			expected: heredoc.Doc(`

				at /foo/bar

				  1 │

				message
			`),
		},
		{
			test: "EmptySourceNoFile",
			err: source.Error{
				Position: testError{errors.New("message")},
				Line:     2,
				Column:   7,
			},
			expected: heredoc.Doc(`

				  1 │

				message
			`),
		},
		{
			test: "EmptySourceNoFileNoColumn",
			err: source.Error{
				Position: testError{errors.New("message")},
				Line:     2,
			},
			expected: heredoc.Doc(`

				  1 │

				message
			`),
		},
		{
			test: "EmptySourceNoFileNoPosition",
			err: source.Error{
				Position: testError{errors.New("message")},
			},
			expected: heredoc.Doc(`

				  1 │

				message
			`),
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			s.Equal(test.expected, test.err.Render(output.Plain))
		})
	}
}

type testError struct{ error }

func (e testError) Error() string {
	return e.error.Error()
}

func (e testError) Position() (int, int) {
	return 0, 0
}

func (e testError) Unwrap() error {
	return e.error
}
