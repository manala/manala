package serrors

import (
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type OsErrorSuite struct{ suite.Suite }

func TestOsErrorSuite(t *testing.T) {
	suite.Run(t, new(OsErrorSuite))
}

func (s *OsErrorSuite) Test() {
	tests := []struct {
		test     string
		err      error
		expected *Assert
	}{
		{
			test: "Unknown",
			err:  New("message"),
			expected: &Assert{
				Type:    &OsError{},
				Message: "message",
			},
		},
		{
			test: "PathError",
			err: &os.PathError{
				Op:   "operation",
				Path: "path",
				Err:  New("message"),
			},
			expected: &Assert{
				Type:    &OsError{},
				Message: "message",
				Arguments: []any{
					"operation", "operation",
					"path", "path",
				},
			},
		},
		{
			test: "SyscallError",
			err: &os.SyscallError{
				Syscall: "syscall",
				Err:     New("message"),
			},
			expected: &Assert{
				Type:    &OsError{},
				Message: "message",
				Arguments: []any{
					"syscall", "syscall",
				},
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			err := NewOs(test.err)

			Equal(s.Assert(), test.expected, err)
		})
	}
}

func (s *OsErrorSuite) TestWrap() {
	err := WrapOs("message", New("wrap"))

	Equal(s.Assert(), &Assert{
		Type:    &WrapOsError{},
		Message: "message",
		Error: &Assert{
			Type:    &OsError{},
			Message: "wrap",
		},
	}, err)

	err = err.WithArguments(
		"foo", "bar",
	)

	Equal(s.Assert(), &Assert{
		Type:    &WrapOsError{},
		Message: "message",
		Arguments: []any{
			"foo", "bar",
		},
		Error: &Assert{
			Type:    &OsError{},
			Message: "wrap",
		},
	}, err)
}
