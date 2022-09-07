package report

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"testing"
)

type AssertSuite struct{ suite.Suite }

func TestAssertSuite(t *testing.T) {
	suite.Run(t, new(AssertSuite))
}

func (s *AssertSuite) Test() {
	report := &Report{
		message: "message",
		fields: map[string]interface{}{
			"foo": "bar",
		},
		err:   fmt.Errorf("error"),
		trace: "trace",
		reports: []*Report{
			{
				message: "message 0",
				fields: map[string]interface{}{
					"foo": "bar 0",
				},
				err:     fmt.Errorf("error 0"),
				trace:   "trace 0",
				reports: []*Report{},
			},
		},
	}

	assert := &Assert{
		Message: "message",
		Fields: map[string]interface{}{
			"foo": "bar",
		},
		Err:   "error",
		Trace: "trace",
		Reports: []Assert{
			{
				Message: "message 0",
				Fields: map[string]interface{}{
					"foo": "bar 0",
				},
				Err:   "error 0",
				Trace: "trace 0",
			},
		},
	}

	assert.Equal(&s.Suite, report)
}
