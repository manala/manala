package watcher

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
	_err := fmt.Errorf("error")
	err := Error(_err)

	s.ErrorAs(err, &internalError)
	s.Equal("watcher error", internalError.Message)
	s.Equal(_err, internalError.Err)
}
