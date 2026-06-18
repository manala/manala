package getter_test

import (
	"errors"
	"testing"

	"github.com/manala/manala/app/repository/getter"
	"github.com/manala/manala/internal/errors/serror/serrortest"
	"github.com/manala/manala/internal/testing/expectation"

	"github.com/stretchr/testify/suite"
)

type ErrorsSuite struct{ suite.Suite }

func TestErrorsSuite(t *testing.T) {
	suite.Run(t, new(ErrorsSuite))
}

func (s *ErrorsSuite) TestIsNotDetected() {
	tests := []struct {
		test     string
		err      error
		expected bool
	}{
		{
			test:     "Any",
			err:      errors.New("foo"),
			expected: false,
		},
		{
			test:     "Exact",
			err:      errors.New("error downloading 'foo'"),
			expected: true,
		},
		{
			test:     "CarriageReturn",
			err:      errors.New("error downloading 'foo\nbar'"),
			expected: true,
		},
		{
			test:     "Almost",
			err:      errors.New("error downloading 'foo': bar"),
			expected: false,
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			if test.expected {
				s.True(getter.IsNotDetected(test.err))
			} else {
				s.False(getter.IsNotDetected(test.err))
			}
		})
	}
}

// Mimic an aws sdk error to avoid direct dependency on it.
type awsError struct{}

func (awsError) Error() string   { return "error" }
func (awsError) Code() string    { return "code" }
func (awsError) Message() string { return "message" }
func (awsError) OrigErr() error  { return nil }

func (s *ErrorsSuite) TestError() {
	tests := []struct {
		test     string
		err      error
		expected expectation.ErrorExpectation
	}{
		{
			test: "Any",
			err:  errors.New("foo"),
			expected: serrortest.Expectation{
				Msg: "unable to handle repository",
				Attrs: [][2]any{
					{"error", "foo"},
				},
			},
		},
		{
			test: "SubdirOutOfRepository",
			err:  errors.New("subdirectory component contain path traversal out of the repository"),
			expected: serrortest.Expectation{
				Msg: "subdir out of repository",
			},
		},
		{
			test: "Aws",
			err:  awsError{},
			expected: serrortest.Expectation{
				Msg:  "aws error",
				Dump: "error",
				Attrs: [][2]any{
					{"code", "code"},
					{"message", "message"},
				},
			},
		},
		{
			test: "CommandErrorCode",
			err:  errors.New("foo exited with 123: bar"),
			expected: serrortest.Expectation{
				Msg:  "command error",
				Dump: "bar",
				Attrs: [][2]any{
					{"command", "foo"},
					{"code", 123},
				},
			},
		},
		{
			test: "CommandError",
			err:  errors.New("error running foo: bar"),
			expected: serrortest.Expectation{
				Msg:  "command error",
				Dump: "bar",
				Attrs: [][2]any{
					{"command", "foo"},
				},
			},
		},
		{
			test: "MultiError",
			//revive:disable:error-strings
			err: errors.New("error downloading 'foo': 123 errors occurred:\nbar\nbaz\n\n"),
			expected: serrortest.Expectation{
				Msg:  "unable to handle repository",
				Dump: "bar\nbaz",
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			expectation.ExpectError(s.T(), test.expected, getter.ErrorFrom(test.err))
		})
	}
}
