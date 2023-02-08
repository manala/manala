package yaml

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	internalReport "manala/internal/report"
	"testing"
)

type ErrorsSuite struct{ suite.Suite }

func TestErrorsSuite(t *testing.T) {
	suite.Run(t, new(ErrorsSuite))
}

func (s *ErrorsSuite) Test() {
	s.Run("Unformatted", func() {
		_err := fmt.Errorf("error")
		err := NewError(_err)

		var _error *Error
		s.ErrorAs(err, &_error)

		report := internalReport.NewErrorReport(err)

		reportAssert := &internalReport.Assert{
			Err: "error",
		}
		reportAssert.Equal(&s.Suite, report)
	})

	s.Run("Formatted", func() {
		_, _err := NewParser().ParseBytes([]byte("&foo"))
		err := NewError(_err)

		var _error *Error
		s.ErrorAs(err, &_error)

		report := internalReport.NewErrorReport(err)

		reportAssert := &internalReport.Assert{
			Err: "unexpected anchor. anchor value is undefined",
			Fields: map[string]interface{}{
				"line":   1,
				"column": 2,
			},
			Trace: ">  1 | &foo\n        ^\n",
		}
		reportAssert.Equal(&s.Suite, report)
	})
}

func (s *ErrorsSuite) TestNode() {
	content := `foo: bar`
	contentNode, _ := NewParser().ParseBytes([]byte(content))

	err := NewNodeError("message", contentNode)

	var _nodeError *NodeError
	s.ErrorAs(err, &_nodeError)

	report := internalReport.NewErrorReport(err)

	reportAssert := &internalReport.Assert{
		Err: "message",
		Fields: map[string]interface{}{
			"line":   1,
			"column": 4,
		},
		Trace: ">  1 | foo: bar\n          ^\n",
	}
	reportAssert.Equal(&s.Suite, report)
}
