package serrors_test

import (
	"os"
	"testing"

	"github.com/manala/manala/internal/serrors"
	"github.com/manala/manala/internal/testing/expect"

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
		expected expect.ErrorExpectation
	}{
		{
			test: "Unknown",
			err:  serrors.New("unknown"),
			expected: serrors.Expectation{
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
			expected: serrors.Expectation{
				Message: "path",
				Attrs: [][2]any{
					{"operation", "operation"},
					{"path", "path"},
				},
			},
		},
		{
			test: "PathErrorNotExist",
			err: &os.PathError{
				Op:   "open",
				Path: "/foo/bar",
				Err:  os.ErrNotExist,
			},
			expected: serrors.Expectation{
				Message: "file does not exist",
				Attrs: [][2]any{
					{"operation", "open"},
					{"path", "/foo/bar"},
				},
			},
		},
		{
			test: "SyscallError",
			err: &os.SyscallError{
				Syscall: "syscall",
				Err:     serrors.New("syscall"),
			},
			expected: serrors.Expectation{
				Message: "syscall",
				Attrs: [][2]any{
					{"syscall", "syscall"},
				},
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			err := serrors.FromOs(test.err)

			expect.Error(s.T(), test.expected, err)
		})
	}
}
