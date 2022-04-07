package os

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	internalErrors "manala/internal/errors"
	"os"
	"testing"
)

var internalError *internalErrors.InternalError

type ErrorsSuite struct{ suite.Suite }

func TestErrorsSuite(t *testing.T) {
	suite.Run(t, new(ErrorsSuite))
}

func (s *ErrorsSuite) TestFileSystemError() {
	_err := fmt.Errorf("error")
	err := FileSystemError(_err)

	s.ErrorAs(err, &internalError)
	s.Equal("file system error", internalError.Message)
	s.Equal(_err, internalError.Err)

	s.Run("Path Error", func() {
		_err := &os.PathError{
			Op:   "operation",
			Path: "path",
			Err:  fmt.Errorf("error"),
		}
		err := FileSystemError(_err)

		s.ErrorAs(err, &internalError)
		s.Equal("file system error", internalError.Message)
		s.Equal(_err.Err, internalError.Err)
		s.Equal("operation", internalError.Fields["operation"])
		s.Equal("path", internalError.Fields["path"])
	})
}
