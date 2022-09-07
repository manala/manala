package yaml

import (
	yamlAst "github.com/goccy/go-yaml/ast"
	"github.com/stretchr/testify/suite"
	internalReport "manala/internal/report"
	internalTesting "manala/internal/testing"
	"testing"
)

type ParserSuite struct{ suite.Suite }

func TestParserSuite(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	suite.Run(t, new(ParserSuite))
}

func (s *ParserSuite) TestEmpty() {
	parser := NewParser()

	node, err := parser.ParseBytes(nil)

	s.Nil(node)
	s.EqualError(err, "empty yaml file")

	report := internalReport.NewErrorReport(err)

	reportAssert := &internalReport.Assert{
		Err: "empty yaml file",
	}
	reportAssert.Equal(&s.Suite, report)
}

func (s *ParserSuite) TestMultipleDocuments() {
	parser := NewParser()

	node, err := parser.ParseFile(internalTesting.DataPath(s, "node.yaml"))

	s.Nil(node)
	s.EqualError(err, "multiple documents yaml file")

	report := internalReport.NewErrorReport(err)

	reportAssert := &internalReport.Assert{
		Err: "multiple documents yaml file",
		Fields: map[string]interface{}{
			"line":   4,
			"column": 1,
		},
	}
	reportAssert.Equal(&s.Suite, report)
}

func (s *ParserSuite) TestMappingComments() {
	parser := NewParser(WithComments())

	node, err := parser.ParseFile(internalTesting.DataPath(s, "node.yaml"))

	s.NoError(err)

	emptyNode := node.(*yamlAst.MappingNode).Values[0]
	s.Equal("# Empty", emptyNode.GetComment().String())

	mappingValueNode := node.(*yamlAst.MappingNode).Values[1]
	s.Equal("# Mapping Value", mappingValueNode.GetComment().String())
	s.Equal("# Mapping Value Foo", mappingValueNode.Value.GetComment().String())

	mappingNode := node.(*yamlAst.MappingNode).Values[2]
	s.Equal("# Mapping", mappingNode.GetComment().String())
	s.Equal("# Mapping Foo", mappingNode.Value.(*yamlAst.MappingNode).Values[0].GetComment().String())
	s.Equal("# Mapping Bar", mappingNode.Value.(*yamlAst.MappingNode).Values[1].GetComment().String())
}

func (s *ParserSuite) TestIrregularMapKeys() {
	tests := []struct {
		name   string
		err    string
		report *internalReport.Assert
	}{
		{
			name: "Integer",
			err:  "irregular map key",
			report: &internalReport.Assert{
				Err: "irregular map key",
				Fields: map[string]interface{}{
					"line":   1,
					"column": 2,
				},
			},
		},
		{
			name: "Integer Anchor",
			err:  "irregular map key",
			report: &internalReport.Assert{
				Err: "irregular map key",
				Fields: map[string]interface{}{
					"line":   2,
					"column": 4,
				},
			},
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			parser := NewParser()

			node, err := parser.ParseFile(internalTesting.DataPath(s, "node.yaml"))

			s.Nil(node)
			s.EqualError(err, test.err)

			report := internalReport.NewErrorReport(err)

			test.report.Equal(&s.Suite, report)
		})
	}
}

func (s *ParserSuite) TestIrregularTypes() {
	tests := []struct {
		name   string
		err    string
		report *internalReport.Assert
	}{
		{
			name: "Inf",
			err:  "irregular type",
			report: &internalReport.Assert{
				Err: "irregular type",
				Fields: map[string]interface{}{
					"line":   1,
					"column": 6,
				},
			},
		},
		{
			name: "Nan",
			err:  "irregular type",
			report: &internalReport.Assert{
				Err: "irregular type",
				Fields: map[string]interface{}{
					"line":   1,
					"column": 6,
				},
			},
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			parser := NewParser()

			node, err := parser.ParseFile(internalTesting.DataPath(s, "node.yaml"))

			s.Nil(node)
			s.EqualError(err, test.err)

			report := internalReport.NewErrorReport(err)

			test.report.Equal(&s.Suite, report)
		})
	}
}

func (s *ParserSuite) TestMappingKey() {
	parser := NewParser()

	node, err := parser.ParseFile(internalTesting.DataPath(s, "node.yaml"))

	s.NoError(err)
	s.IsType((*yamlAst.MappingValueNode)(nil), node)

	keyNode := node.(*yamlAst.MappingValueNode).Key
	s.IsType((*yamlAst.StringNode)(nil), keyNode)
	s.Equal("foo", keyNode.(*yamlAst.StringNode).Value)

	valueNode := node.(*yamlAst.MappingValueNode).Value
	s.IsType((*yamlAst.StringNode)(nil), valueNode)
	s.Equal("bar", valueNode.(*yamlAst.StringNode).Value)
}

func (s *ParserSuite) TestIrregularMappingKey() {
	parser := NewParser()

	node, err := parser.ParseFile(internalTesting.DataPath(s, "node.yaml"))

	s.Nil(node)
	s.EqualError(err, "irregular map key")

	report := internalReport.NewErrorReport(err)

	reportAssert := &internalReport.Assert{
		Err: "irregular map key",
		Fields: map[string]interface{}{
			"line":   1,
			"column": 6,
		},
	}
	reportAssert.Equal(&s.Suite, report)
}

func (s *ParserSuite) TestTags() {
	parser := NewParser()

	node, err := parser.ParseFile(internalTesting.DataPath(s, "node.yaml"))

	s.NoError(err)
	s.IsType((*yamlAst.StringNode)(nil), node)
	s.Equal("foo", node.(*yamlAst.StringNode).Value)
}

func (s *ParserSuite) TestUnknownAnchors() {
	parser := NewParser()

	node, err := parser.ParseFile(internalTesting.DataPath(s, "node.yaml"))

	s.Nil(node)
	s.EqualError(err, "cannot find anchor \"anchor\"")

	report := internalReport.NewErrorReport(err)

	reportAssert := &internalReport.Assert{
		Err: "cannot find anchor \"anchor\"",
		Fields: map[string]interface{}{
			"line":   1,
			"column": 2,
		},
	}
	reportAssert.Equal(&s.Suite, report)
}

func (s *ParserSuite) TestAnchors() {
	s.Run("Anchors", func() {
		parser := NewParser()

		node, err := parser.ParseFile(internalTesting.DataPath(s, "node.yaml"))

		s.NoError(err)

		anchorNode := node.(*yamlAst.MappingNode).Values[0]
		s.IsType((*yamlAst.StringNode)(nil), anchorNode.Value)
		s.Equal("foo", anchorNode.Value.(*yamlAst.StringNode).Value)

		aliasNode := node.(*yamlAst.MappingNode).Values[1]
		s.IsType((*yamlAst.StringNode)(nil), aliasNode.Value)
		s.Equal("foo", aliasNode.Value.(*yamlAst.StringNode).Value)
	})
	s.Run("Merge Keys", func() {
		parser := NewParser()

		node, err := parser.ParseFile(internalTesting.DataPath(s, "node.yaml"))

		s.NoError(err)

		emptyAnchorNode := node.(*yamlAst.MappingNode).Values[0]
		s.IsType((*yamlAst.MappingNode)(nil), emptyAnchorNode.Value)
		s.Len(emptyAnchorNode.Value.(*yamlAst.MappingNode).Values, 0)

		mappingValueAnchorNode := node.(*yamlAst.MappingNode).Values[1]
		s.IsType((*yamlAst.MappingValueNode)(nil), mappingValueAnchorNode.Value)

		mappingAnchorNode := node.(*yamlAst.MappingNode).Values[2]
		s.IsType((*yamlAst.MappingNode)(nil), mappingAnchorNode.Value)
		s.Len(mappingAnchorNode.Value.(*yamlAst.MappingNode).Values, 2)

		mappingValueAliasEmptyAnchorNode := node.(*yamlAst.MappingNode).Values[3]
		s.IsType((*yamlAst.MappingNode)(nil), mappingValueAliasEmptyAnchorNode.Value)
		s.Len(mappingValueAliasEmptyAnchorNode.Value.(*yamlAst.MappingNode).Values, 0)

		mappingValueAliasMappingValueAnchorNode := node.(*yamlAst.MappingNode).Values[4]
		s.IsType((*yamlAst.MappingValueNode)(nil), mappingValueAliasMappingValueAnchorNode.Value)

		mappingValueAliasMappingAnchorNode := node.(*yamlAst.MappingNode).Values[5]
		s.IsType((*yamlAst.MappingNode)(nil), mappingValueAliasMappingAnchorNode.Value)
		s.Len(mappingValueAliasMappingAnchorNode.Value.(*yamlAst.MappingNode).Values, 2)

		mappingAliasEmptyAnchorNode := node.(*yamlAst.MappingNode).Values[6]
		s.IsType((*yamlAst.MappingValueNode)(nil), mappingAliasEmptyAnchorNode.Value)

		mappingAliasMappingValueAnchorNode := node.(*yamlAst.MappingNode).Values[7]
		s.IsType((*yamlAst.MappingNode)(nil), mappingAliasMappingValueAnchorNode.Value)
		s.Len(mappingAliasMappingValueAnchorNode.Value.(*yamlAst.MappingNode).Values, 2)

		mappingValueAliasMappingNode := node.(*yamlAst.MappingNode).Values[8]
		s.IsType((*yamlAst.MappingNode)(nil), mappingValueAliasMappingNode.Value)
		s.Len(mappingValueAliasMappingNode.Value.(*yamlAst.MappingNode).Values, 3)
	})
	s.Run("Merge Keys Duplicated", func() {
		parser := NewParser()

		node, err := parser.ParseFile(internalTesting.DataPath(s, "node.yaml"))

		s.NoError(err)

		singleMappingAliasMappingValueAnchorNode := node.(*yamlAst.MappingNode).Values[2]
		s.IsType((*yamlAst.MappingValueNode)(nil), singleMappingAliasMappingValueAnchorNode.Value)
		s.Equal("bar", singleMappingAliasMappingValueAnchorNode.Value.(*yamlAst.MappingValueNode).Value.(*yamlAst.StringNode).Value)

		multipleMappingAliasMappingValueAnchorNode := node.(*yamlAst.MappingNode).Values[3]
		s.IsType((*yamlAst.MappingNode)(nil), multipleMappingAliasMappingValueAnchorNode.Value)
		s.Len(multipleMappingAliasMappingValueAnchorNode.Value.(*yamlAst.MappingNode).Values, 2)
		s.Equal("bar", multipleMappingAliasMappingValueAnchorNode.Value.(*yamlAst.MappingNode).Values[0].Value.(*yamlAst.StringNode).Value)

		mappingAliasMappingAnchorNode := node.(*yamlAst.MappingNode).Values[4]
		s.IsType((*yamlAst.MappingNode)(nil), mappingAliasMappingAnchorNode.Value)
		s.Len(mappingAliasMappingAnchorNode.Value.(*yamlAst.MappingNode).Values, 3)
		s.Equal("bar", multipleMappingAliasMappingValueAnchorNode.Value.(*yamlAst.MappingNode).Values[1].Value.(*yamlAst.StringNode).Value)
	})
}
