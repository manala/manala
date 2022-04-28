package yaml

import (
	"fmt"
	"github.com/goccy/go-yaml/parser"
	"github.com/stretchr/testify/suite"
	internalErrors "manala/internal/errors"
	"testing"
)

var internalError *internalErrors.InternalError

type ErrorsSuite struct{ suite.Suite }

func TestErrorsSuite(t *testing.T) {
	suite.Run(t, new(ErrorsSuite))
}

func (s *ErrorsSuite) TestError() {
	s.Run("Unformatted", func() {
		_err := fmt.Errorf("error")
		err := Error("file", _err)

		s.ErrorAs(err, &internalError)
		s.Equal("yaml processing error", internalError.Message)
		s.Equal("file", internalError.Fields["file"])
		s.NotContains(internalError.Fields, "line")
		s.NotContains(internalError.Fields, "column")
		s.NotContains(internalError.Fields, "message")
		s.Equal(_err, internalError.Err)
	})

	s.Run("Formatted", func() {
		_, _err := parser.ParseBytes([]byte("&foo"), 0)
		err := Error("file", _err)

		s.ErrorAs(err, &internalError)
		s.Equal("yaml processing error", internalError.Message)
		s.Equal("file", internalError.Fields["file"])
		s.Equal(1, internalError.Fields["line"])
		s.Equal(2, internalError.Fields["column"])
		s.Equal("unexpected anchor. anchor value is undefined", internalError.Fields["message"])
		s.Equal(">  1 | \x1b[93m&\x1b[0m\x1b[93mfoo\x1b[0m\n        ^\n", internalError.Trace)
	})
}

func (s *ErrorsSuite) TestCommentTagError() {
	_err := fmt.Errorf("error")
	err := CommentTagError("path", _err)

	s.ErrorAs(err, &internalError)
	s.Equal("yaml comment tag error", internalError.Message)
	s.Equal("path", internalError.Fields["path"])
	s.Equal(_err, internalError.Err)
}
