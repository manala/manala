package serrors_test

import (
	"os"
	"testing"

	"github.com/manala/manala/internal/serrors"
	"github.com/manala/manala/internal/testing/errors"

	"github.com/stretchr/testify/suite"
)

type OsSuite struct{ suite.Suite }

func TestOsSuite(t *testing.T) {
	suite.Run(t, new(OsSuite))
}

func (s *OsSuite) Test() {
	tests := []struct {
		test     string
		err      error
		expected errors.Assertion
	}{
		{
			test: "Unknown",
			err:  serrors.New("unknown"),
			expected: &serrors.Assertion{
				Type:    serrors.Error{},
				Message: "unknown",
			},
		},
		{
			test: "PathError",
			err: &os.PathError{
				Op:   "operation",
				Path: "path",
				Err:  serrors.New("path"),
			},
			expected: &serrors.Assertion{
				Type:    serrors.Error{},
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
				Err:     serrors.New("syscall"),
			},
			expected: &serrors.Assertion{
				Type:    serrors.Error{},
				Message: "syscall",
				Arguments: []any{
					"syscall", "syscall",
				},
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			err := serrors.NewOs(test.err)

			errors.Equal(s.T(), test.expected, err)
		})
	}
}
