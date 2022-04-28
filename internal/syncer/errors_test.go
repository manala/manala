package syncer

import (
	"github.com/stretchr/testify/suite"
	internalErrors "manala/internal/errors"
	"testing"
)

var internalError *internalErrors.InternalError

type ErrorsSuite struct{ suite.Suite }

func TestErrorsSuite(t *testing.T) {
	suite.Run(t, new(ErrorsSuite))
}

func (s *ErrorsSuite) Test() {

	s.Run("SourceNotExistError", func() {
		err := SourceNotExistError("path")

		s.ErrorAs(err, &internalError)
		s.Equal("no source file or directory", internalError.Message)
		s.Equal("path", internalError.Fields["path"])
	})

}
