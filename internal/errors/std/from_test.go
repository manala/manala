package std_test

import (
	"os"
	"testing"

	"github.com/manala/manala/internal/errors/serror"
	"github.com/manala/manala/internal/errors/serror/serrortest"
	"github.com/manala/manala/internal/errors/std"
	"github.com/manala/manala/internal/testing/expectation"

	"github.com/stretchr/testify/suite"
)

type FromSuite struct{ suite.Suite }

func TestFromSuite(t *testing.T) {
	suite.Run(t, new(FromSuite))
}

func (s *FromSuite) Test() {
	tests := []struct {
		test     string
		err      error
		expected expectation.ErrorExpectation
	}{
		{
			test: "Unknown",
			err:  serror.New("unknown"),
			expected: serrortest.Expectation{
				Msg: "unknown",
			},
		},
		{
			test: "PathError",
			err: &os.PathError{
				Op:   "operation",
				Path: "path",
				Err:  serror.New("path"),
			},
			expected: serrortest.Expectation{
				Msg: "path",
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
			expected: serrortest.Expectation{
				Msg: "file does not exist",
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
				Err:     serror.New("syscall"),
			},
			expected: serrortest.Expectation{
				Msg: "syscall",
				Attrs: [][2]any{
					{"syscall", "syscall"},
				},
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			err := std.From(test.err)

			expectation.ExpectError(s.T(), test.expected, err)
		})
	}
}
