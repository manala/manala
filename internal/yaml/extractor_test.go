package yaml

import (
	yamlAst "github.com/goccy/go-yaml/ast"
	"github.com/stretchr/testify/suite"
	internalReport "manala/internal/report"
	internalTesting "manala/internal/testing"
	"testing"
)

type ExtractorSuite struct{ suite.Suite }

func TestExtractorSuite(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	suite.Run(t, new(ExtractorSuite))
}

func (s *ExtractorSuite) TestExtractRootMapErrors() {
	tests := []struct {
		name   string
		err    string
		report *internalReport.Assert
	}{
		{
			name: "Empty",
			err:  "root must be a map",
			report: &internalReport.Assert{
				Err: "root must be a map",
			},
		},
		{
			name: "Non Map",
			err:  "root must be a map",
			report: &internalReport.Assert{
				Err: "root must be a map",
				Fields: map[string]interface{}{
					"line":   1,
					"column": 1,
				},
			},
		},
		{
			name: "Subject Not Found Single",
			err:  "unable to find \"subject\" map",
			report: &internalReport.Assert{
				Err: "unable to find \"subject\" map",
			},
		},
		{
			name: "Subject Not Found Multiple",
			err:  "unable to find \"subject\" map",
			report: &internalReport.Assert{
				Err: "unable to find \"subject\" map",
			},
		},
		{
			name: "Subject Non Map Single",
			err:  "\"subject\" is not a map",
			report: &internalReport.Assert{
				Err: "\"subject\" is not a map",
				Fields: map[string]interface{}{
					"line":   1,
					"column": 10,
				},
			},
		},
		{
			name: "Subject Non Map Multiple",
			err:  "\"subject\" is not a map",
			report: &internalReport.Assert{
				Err: "\"subject\" is not a map",
				Fields: map[string]interface{}{
					"line":   1,
					"column": 10,
				},
			},
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			parser := NewParser()

			node, _ := parser.ParseFile(internalTesting.DataPath(s, "node.yaml"))

			extractor := NewExtractor(&node)
			subjectNode, err := extractor.ExtractRootMap("subject")

			s.Nil(subjectNode)
			s.EqualError(err, test.err)

			report := internalReport.NewErrorReport(err)

			test.report.Equal(&s.Suite, report)
		})
	}
}

func (s *ExtractorSuite) TestExtractRootMap() {
	tests := []struct {
		name    string
		subject interface{}
		node    interface{}
	}{
		{
			name:    "Single Empty",
			subject: (*yamlAst.MappingNode)(nil),
			node:    (*yamlAst.MappingNode)(nil),
		},
		{
			name:    "Single Single",
			subject: (*yamlAst.MappingValueNode)(nil),
			node:    (*yamlAst.MappingNode)(nil),
		},
		{
			name:    "Single Multiple",
			subject: (*yamlAst.MappingNode)(nil),
			node:    (*yamlAst.MappingNode)(nil),
		},
		{
			name:    "Couple Empty",
			subject: (*yamlAst.MappingNode)(nil),
			node:    (*yamlAst.MappingValueNode)(nil),
		},
		{
			name:    "Couple Single",
			subject: (*yamlAst.MappingValueNode)(nil),
			node:    (*yamlAst.MappingValueNode)(nil),
		},
		{
			name:    "Couple Multiple",
			subject: (*yamlAst.MappingNode)(nil),
			node:    (*yamlAst.MappingValueNode)(nil),
		},
		{
			name:    "Multiple Empty",
			subject: (*yamlAst.MappingNode)(nil),
			node:    (*yamlAst.MappingNode)(nil),
		},
		{
			name:    "Multiple Single",
			subject: (*yamlAst.MappingValueNode)(nil),
			node:    (*yamlAst.MappingNode)(nil),
		},
		{
			name:    "Multiple Multiple",
			subject: (*yamlAst.MappingNode)(nil),
			node:    (*yamlAst.MappingNode)(nil),
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			parser := NewParser()

			node, _ := parser.ParseFile(internalTesting.DataPath(s, "node.yaml"))

			extractor := NewExtractor(&node)
			subjectNode, err := extractor.ExtractRootMap("subject")

			s.NoError(err)
			s.IsType(test.subject, subjectNode)
			s.IsType(test.node, node)
		})
	}
}
