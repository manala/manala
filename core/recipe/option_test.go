package recipe

import (
	"encoding/json"
	goYaml "github.com/goccy/go-yaml"
	goYamlAst "github.com/goccy/go-yaml/ast"
	"github.com/stretchr/testify/suite"
	"manala/internal/errors/serrors"
	"manala/internal/validation"
	"manala/internal/yaml"
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
		test     string
		node     string
		value    interface{}
		expected *serrors.Assert
	}{
		{
			test:  "Unsupported",
			node:  `node: ~`,
			value: []string{},
			expected: &serrors.Assert{
				Type:    &serrors.Error{},
				Message: "unsupported option value type",
				Arguments: []any{
					"value", []string{},
				},
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			node, _ := yaml.NewParser(yaml.WithComments()).ParseBytes([]byte(test.node))

			option := &option{
				node: node.(*goYamlAst.MappingValueNode),
			}

			err := option.Set(test.value)

			var value map[string]interface{}
			_ = goYaml.NewDecoder(node).Decode(&value)

			serrors.Equal(s.Assert(), test.expected, err)
		})
	}
}

func (s *OptionSuite) TestSet() {
	tests := []struct {
		test     string
		node     string
		value    interface{}
		expected interface{}
	}{
		{
			test:     "Nil",
			node:     `node: ""`,
			value:    nil,
			expected: nil,
		},
		{
			test:     "True",
			node:     `node: ~`,
			value:    true,
			expected: true,
		},
		{
			test:     "False",
			node:     `node: ~`,
			value:    false,
			expected: false,
		},
		{
			test:     "String",
			node:     `node: ~`,
			value:    "string",
			expected: "string",
		},
		{
			test:     "StringEmpty",
			node:     `node: ~`,
			value:    "",
			expected: "",
		},
		{
			test:     "StringAsterisk",
			node:     `node: ~`,
			value:    "*",
			expected: "*",
		},
		{
			test:     "StringInt",
			node:     `node: ~`,
			value:    "12",
			expected: "12",
		},
		{
			test:     "StringFloat",
			node:     `node: ~`,
			value:    "2.3",
			expected: "2.3",
		},
		{
			test:     "StringFloatInt",
			node:     `node: ~`,
			value:    "3.0",
			expected: "3.0",
		},
		{
			test:     "IntegerUint64",
			node:     `node: ~`,
			value:    uint64(12),
			expected: uint64(12),
		},
		{
			test:     "IntegerFloat64",
			node:     `node: ~`,
			value:    float64(12),
			expected: uint64(12),
		},
		{
			test:     "Float",
			node:     `node: ~`,
			value:    2.3,
			expected: 2.3,
		},
		{
			test:     "FloatInt",
			node:     `node: ~`,
			value:    3.0,
			expected: uint64(3),
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			node, _ := yaml.NewParser(yaml.WithComments()).ParseBytes([]byte(test.node))

			option := &option{
				node: node.(*goYamlAst.MappingValueNode),
			}

			err := option.Set(test.value)

			s.NoError(err)

			var value map[string]interface{}
			_ = goYaml.NewDecoder(node).Decode(&value)

			s.Equal(test.expected, value["node"])
		})
	}
}

func (s *OptionSuite) TestUnmarshalJSONErrors() {
	tests := []struct {
		test     string
		data     string
		expected *serrors.Assert
	}{
		{
			test: "Syntax",
			data: `foo`,
			expected: &serrors.Assert{
				Type:    &json.SyntaxError{},
				Message: "invalid character 'o' in literal false (expecting 'a')",
			},
		},
		{
			test: "Type",
			data: `[]`,
			expected: &serrors.Assert{
				Type:    &json.UnmarshalTypeError{},
				Message: "json: cannot unmarshal array into Go value of type map[string]interface {}",
			},
		},
		{
			test: "Validation",
			data: `{"foo": "bar"}`,
			expected: &serrors.Assert{
				Type:    &validation.Error{},
				Message: "invalid option",
				Errors: []*serrors.Assert{
					{
						Type:    &validation.ResultError{},
						Message: "missing label field",
						Arguments: []any{
							"property", "label",
						},
					},
					{
						Type:    &validation.ResultError{},
						Message: "don't support additional properties",
						Arguments: []any{
							"property", "foo",
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			var option *option

			err := json.Unmarshal([]byte(test.data), &option)

			serrors.Equal(s.Assert(), test.expected, err)
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
