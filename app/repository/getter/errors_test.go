package getter_test

import (
	"errors"
	"manala/app/repository/getter"
	"manala/internal/serrors"
	"testing"

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
type awsErrorTest struct{}

func (awsErrorTest) Error() string   { return "error" }
func (awsErrorTest) Code() string    { return "code" }
func (awsErrorTest) Message() string { return "message" }
func (awsErrorTest) OrigErr() error  { return nil }

func (s *ErrorsSuite) TestError() {
	tests := []struct {
		test     string
		err      error
		expected *serrors.Assertion
	}{
		{
			test: "Any",
			err:  errors.New("foo"),
			expected: &serrors.Assertion{
				Message: "unable to handle repository",
				Arguments: []any{
					"error", "foo",
				},
			},
		},
		{
			test: "SubdirOutOfRepository",
			err:  errors.New("subdirectory component contain path traversal out of the repository"),
			expected: &serrors.Assertion{
				Message: "subdir out of repository",
			},
		},
		{
			test: "Aws",
			err:  awsErrorTest{},
			expected: &serrors.Assertion{
				Message: "aws error",
				Details: "error",
				Arguments: []any{
					"code", "code",
					"message", "message",
				},
			},
		},
		{
			test: "CommandErrorCode",
			err:  errors.New("foo exited with 123: bar"),
			expected: &serrors.Assertion{
				Message: "command error",
				Details: "bar",
				Arguments: []any{
					"command", "foo",
					"code", 123,
				},
			},
		},
		{
			test: "CommandError",
			err:  errors.New("error running foo: bar"),
			expected: &serrors.Assertion{
				Message: "command error",
				Details: "bar",
				Arguments: []any{
					"command", "foo",
				},
			},
		},
		{
			test: "MultiError",
			//revive:disable:error-strings
			err: errors.New("error downloading 'foo': 123 errors occurred:\nbar\nbaz\n\n"),
			expected: &serrors.Assertion{
				Message: "unable to handle repository",
				Details: "bar\nbaz",
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			serrors.Equal(s.T(), test.expected, getter.NewError(test.err))
		})
	}
}
