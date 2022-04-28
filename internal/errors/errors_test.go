package errors

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"testing"
)

var internalError *InternalError

type ErrorsSuite struct{ suite.Suite }

func TestErrorsSuite(t *testing.T) {
	suite.Run(t, new(ErrorsSuite))
}

func (s *ErrorsSuite) Test() {
	err := New("message")

	s.ErrorAs(err, &internalError)
	s.Equal("message", internalError.Message)
	s.Empty(internalError.Fields)
	s.Nil(internalError.Err)
	s.Empty(internalError.Errs)
	s.Empty(internalError.Trace)

	s.Run("Error", func() {
		s.Equal("message", internalError.Error())
	})

	s.Run("With", func() {
		_ = err.With("foo")
		s.Equal("foo", internalError.Message)
	})

	s.Run("With Field", func() {
		_ = err.WithField("key", "value")
		s.Equal("value", internalError.Fields["key"])
	})

	s.Run("With Error", func() {
		_err := fmt.Errorf("error")
		_ = err.WithError(_err)
		s.Equal(_err, internalError.Err)
	})

	s.Run("With Errors", func() {
		_err := New("error")
		_ = err.WithErrors([]*InternalError{_err})
		s.Equal(_err, internalError.Errs[0])
	})

	s.Run("With Trace", func() {
		_ = err.WithTrace("trace")
		s.Equal("trace", internalError.Trace)
	})
}
