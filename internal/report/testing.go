package report

import (
	"fmt"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/suite"
	internalTesting "manala/internal/testing"
	"os"
)

type Assert struct {
	Message     string
	Fields      map[string]interface{}
	Err         string
	Trace       string
	TraceGolden string
	Reports     []Assert
}

func (assert *Assert) Equal(s *suite.Suite, report *Report) {
	// Message
	s.Equal(assert.Message, report.Message(), "Report message not equal")

	// Fields
	if assert.Fields != nil {
		s.Equal(assert.Fields, report.Fields(), "Report fields not equal")
	} else {
		s.Equal(map[string]interface{}{}, report.Fields(), "Report fields not nil")
	}

	// Error
	if assert.Err != "" {
		s.EqualError(report.Err(), assert.Err, "Report error not equal")
	} else {
		s.NoError(report.Err(), "Report error not nil")
	}

	// Trace golden
	if assert.TraceGolden == "" {
		assert.TraceGolden = "trace"
	}

	// Trace
	if assert.Trace == "" {
		// Trace Golden
		if _, statErr := os.Stat(internalTesting.DataPath(s, assert.TraceGolden+".golden")); statErr == nil {
			g := goldie.New(s.T())
			g.Assert(s.T(), internalTesting.Path(s, assert.TraceGolden), []byte(report.Trace()))
		} else {
			s.Equal("", report.Trace(), "Report trace not equal")
		}
	} else {
		s.Equal(assert.Trace, report.Trace(), "Report trace not equal")
	}

	// Reports
	if assert.Reports != nil {
		s.Len(report.Reports(), len(assert.Reports))
		for i, rep := range assert.Reports {
			rep.TraceGolden = fmt.Sprintf("%s.%d", assert.TraceGolden, i)
			rep.Equal(s, report.Reports()[i])
		}
	} else {
		s.Equal([]*Report{}, report.Reports(), "Report reports not empty")
	}
}
