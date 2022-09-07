package recipe

import (
	"fmt"
	"github.com/goccy/go-yaml"
	yamlAst "github.com/goccy/go-yaml/ast"
	"github.com/stretchr/testify/suite"
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

func (s *OptionSuite) TestSet() {
	tests := []struct {
		name          string
		node          string
		value         interface{}
		expectedValue interface{}
		err           error
	}{
		{
			name:  "Unsupported",
			node:  `node: ~`,
			value: []string{},
			err:   fmt.Errorf(""),
		},
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

			var value map[string]interface{}
			_ = yaml.NewDecoder(node).Decode(&value)

			s.IsType(test.err, err)
			if test.err == nil {
				s.Equal(test.expectedValue, value["node"])
			}
		})
	}
}
