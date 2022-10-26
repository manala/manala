package recipe

import (
	"encoding/json"
	"github.com/goccy/go-yaml"
	yamlAst "github.com/goccy/go-yaml/ast"
	"github.com/stretchr/testify/suite"
	internalReport "manala/internal/report"
	internalYaml "manala/internal/yaml"
	"testing"
)

type OptionSuite struct{ suite.Suite }

func TestOptionSuite(t *testing.T) {
	suite.Run(t, new(OptionSuite))
}

func (s *OptionSuite) Test() {
	label := "label"
	schema := map[string]interface{}{"foo": "bar"}

	option := &option{
		label:  label,
		schema: schema,
	}

	s.Equal(label, option.Label())
	s.Equal(schema, option.Schema())
}

func (s *OptionSuite) TestSetErrors() {
	tests := []struct {
		name  string
		node  string
		value interface{}
		err   string
	}{
		{
			name:  "Unsupported",
			node:  `node: ~`,
			value: []string{},
			err:   "unsupported option value type: []",
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			node, _ := internalYaml.NewParser(internalYaml.WithComments()).ParseBytes([]byte(test.node))

			option := &option{
				node: node.(*yamlAst.MappingValueNode),
			}

			err := option.Set(test.value)

			var value map[string]interface{}
			_ = yaml.NewDecoder(node).Decode(&value)

			s.EqualError(err, test.err)
		})
	}
}

func (s *OptionSuite) TestSet() {
	tests := []struct {
		name          string
		node          string
		value         interface{}
		expectedValue interface{}
	}{
		{
			name:          "Nil",
			node:          `node: ""`,
			value:         nil,
			expectedValue: nil,
		},
		{
			name:          "True",
			node:          `node: ~`,
			value:         true,
			expectedValue: true,
		},
		{
			name:          "False",
			node:          `node: ~`,
			value:         false,
			expectedValue: false,
		},
		{
			name:          "String",
			node:          `node: ~`,
			value:         "string",
			expectedValue: "string",
		},
		{
			name:          "String Empty",
			node:          `node: ~`,
			value:         "",
			expectedValue: "",
		},
		{
			name:          "String Asterisk",
			node:          `node: ~`,
			value:         "*",
			expectedValue: "*",
		},
		{
			name:          "String Int",
			node:          `node: ~`,
			value:         "12",
			expectedValue: "12",
		},
		{
			name:          "String Float",
			node:          `node: ~`,
			value:         "2.3",
			expectedValue: "2.3",
		},
		{
			name:          "String Float Int",
			node:          `node: ~`,
			value:         "3.0",
			expectedValue: "3.0",
		},
		{
			name:          "Integer Uint64",
			node:          `node: ~`,
			value:         uint64(12),
			expectedValue: uint64(12),
		},
		{
			name:          "Integer Float64",
			node:          `node: ~`,
			value:         float64(12),
			expectedValue: uint64(12),
		},
		{
			name:          "Float",
			node:          `node: ~`,
			value:         2.3,
			expectedValue: 2.3,
		},
		{
			name:          "Float Int",
			node:          `node: ~`,
			value:         3.0,
			expectedValue: uint64(3),
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			node, _ := internalYaml.NewParser(internalYaml.WithComments()).ParseBytes([]byte(test.node))

			option := &option{
				node: node.(*yamlAst.MappingValueNode),
			}

			err := option.Set(test.value)

			s.NoError(err)

			var value map[string]interface{}
			_ = yaml.NewDecoder(node).Decode(&value)

			s.Equal(test.expectedValue, value["node"])
		})
	}
}

func (s *OptionSuite) TestUnmarshalJSONErrors() {
	tests := []struct {
		name   string
		data   string
		report *internalReport.Assert
	}{
		{
			name: "Syntax",
			data: `foo`,
			report: &internalReport.Assert{
				Err: "invalid character 'o' in literal false (expecting 'a')",
			},
		},
		{
			name: "Type",
			data: `[]`,
			report: &internalReport.Assert{
				Err: "json: cannot unmarshal array into Go value of type map[string]interface {}",
			},
		},
		{
			name: "Validation",
			data: `{"foo": "bar"}`,
			report: &internalReport.Assert{
				Err: "invalid option",
				Reports: []internalReport.Assert{
					{
						Message: "missing label field",
						Fields: map[string]interface{}{
							"property": "label",
						},
					},
					{
						Message: "don't support additional properties",
						Fields: map[string]interface{}{
							"property": "foo",
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			var option *option

			err := json.Unmarshal([]byte(test.data), &option)

			s.Error(err)

			report := internalReport.NewErrorReport(err)

			test.report.Equal(&s.Suite, report)
		})
	}
}

func (s *OptionSuite) TestUnmarshalJSON() {
	var option *option

	data := `{"label": "foo"}`

	err := json.Unmarshal([]byte(data), &option)

	s.NoError(err)
	s.Equal("foo", option.Label())
}
