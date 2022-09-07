package os

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	internalReport "manala/internal/report"
	"os"
	"testing"
)

type ErrorsSuite struct{ suite.Suite }

func TestErrorsSuite(t *testing.T) {
	suite.Run(t, new(ErrorsSuite))
}

func (s *ErrorsSuite) TestFileSystemError() {
	s.Run("Error", func() {
		_err := fmt.Errorf("error")
		err := NewError(_err)

		var _error *Error
		s.ErrorAs(err, &_error)

		s.EqualError(err, "error")
	})

	s.Run("Path Error", func() {
		_err := &os.PathError{
			Op:   "operation",
			Path: "path",
			Err:  fmt.Errorf("error"),
		}
		err := NewError(_err)

		var _error *Error
		s.ErrorAs(err, &_error)

		report := internalReport.NewErrorReport(err)

		reportAssert := &internalReport.Assert{
			Err: "error",
			Fields: map[string]interface{}{
				"operation": "operation",
			},
		}
		reportAssert.Equal(&s.Suite, report)
	})

	s.Run("SyscallError Error", func() {
		_err := &os.SyscallError{
			Syscall: "syscall",
			Err:     fmt.Errorf("error"),
		}
		err := NewError(_err)

		var _error *Error
		s.ErrorAs(err, &_error)

		report := internalReport.NewErrorReport(err)

		reportAssert := &internalReport.Assert{
			Err: "error",
			Fields: map[string]interface{}{
				"syscall": "syscall",
			},
		}
		reportAssert.Equal(&s.Suite, report)
	})
}
