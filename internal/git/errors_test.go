package git

import (
	"fmt"
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

	s.Run("CloneRepositoryUrlError", func() {
		_err := fmt.Errorf("error")
		err := CloneRepositoryUrlError("dir", "url", _err)

		s.ErrorAs(err, &internalError)
		s.Equal("clone git repository", internalError.Message)
		s.Equal(_err, internalError.Err)
		s.Equal("dir", internalError.Fields["dir"])
		s.Equal("url", internalError.Fields["url"])
	})

	s.Run("InvalidRepositoryError", func() {
		_err := fmt.Errorf("error")
		err := InvalidRepositoryError("dir", _err)

		s.ErrorAs(err, &internalError)
		s.Equal("invalid git repository", internalError.Message)
		s.Equal(_err, internalError.Err)
		s.Equal("dir", internalError.Fields["dir"])
	})

	s.Run("PullRepositoryError", func() {
		_err := fmt.Errorf("error")
		err := PullRepositoryError("dir", _err)

		s.ErrorAs(err, &internalError)
		s.Equal("pull git repository", internalError.Message)
		s.Equal(_err, internalError.Err)
		s.Equal("dir", internalError.Fields["dir"])
	})

	s.Run("DeleteRepositoryError", func() {
		_err := fmt.Errorf("error")
		err := DeleteRepositoryError("dir", _err)

		s.ErrorAs(err, &internalError)
		s.Equal("delete git repository", internalError.Message)
		s.Equal(_err, internalError.Err)
		s.Equal("dir", internalError.Fields["dir"])
	})

	s.Run("OpenRepositoryError", func() {
		_err := fmt.Errorf("error")
		err := OpenRepositoryError("dir", _err)

		s.ErrorAs(err, &internalError)
		s.Equal("open git repository", internalError.Message)
		s.Equal(_err, internalError.Err)
		s.Equal("dir", internalError.Fields["dir"])
	})

}
