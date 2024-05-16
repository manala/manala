package serrors

import (
	"os"
)

func (s *Suite) TestOs() {
	tests := []struct {
		test     string
		err      error
		expected *Assertion
	}{
		{
			test: "Unknown",
			err:  New("unknown"),
			expected: &Assertion{
				Type:    Error{},
				Message: "unknown",
			},
		},
		{
			test: "PathError",
			err: &os.PathError{
				Op:   "operation",
				Path: "path",
				Err:  New("path"),
			},
			expected: &Assertion{
				Type:    Error{},
				Message: "path",
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
				Err:     New("syscall"),
			},
			expected: &Assertion{
				Type:    Error{},
				Message: "syscall",
				Arguments: []any{
					"syscall", "syscall",
				},
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			err := NewOs(test.err)

			Equal(s.T(), test.expected, err)
		})
	}
}
