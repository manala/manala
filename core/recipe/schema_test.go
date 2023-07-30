package recipe

import (
	"encoding/json"
	goYamlAst "github.com/goccy/go-yaml/ast"
	"github.com/stretchr/testify/suite"
	"manala/internal/errors/serrors"
	"manala/internal/testing/heredoc"
	"manala/internal/yaml"
	"testing"
)

type SchemaInferrerSuite struct{ suite.Suite }

func TestSchemaInferrerSuite(t *testing.T) {
	suite.Run(t, new(SchemaInferrerSuite))
}

func (s *SchemaInferrerSuite) TestSchemaChainInferrerErrors() {
	tests := []struct {
		test      string
		node      string
		inferrers []schemaInferrerInterface
		expected  *serrors.Assert
	}{
		{
			test: "Error",
			node: `node: string`,
			inferrers: []schemaInferrerInterface{
				NewSchemaCallbackInferrer(func(node goYamlAst.Node, _ map[string]interface{}) error {
					return yaml.NewNodeError("foo", node)
				}),
			},
			expected: &serrors.Assert{
				Type:    &yaml.NodeError{},
				Message: "foo",
				Arguments: []any{
					"line", 1,
					"column", 5,
				},
				Details: heredoc.Doc(`
					>  1 | node: string
					           ^
				`),
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			node, _ := yaml.NewParser(yaml.WithComments()).ParseBytes([]byte(test.node))

			schema := map[string]interface{}{"foo": "bar"}
			err := NewSchemaChainInferrer(test.inferrers...).Infer(node, schema)

			serrors.Equal(s.Assert(), test.expected, err)
		})
	}
}

func (s *SchemaInferrerSuite) TestSchemaChainInferrer() {
	tests := []struct {
		test      string
		node      string
		inferrers []schemaInferrerInterface
		expected  map[string]interface{}
	}{
		{
			test:      "NoInferrers",
			node:      `node: string`,
			inferrers: []schemaInferrerInterface{},
			expected:  map[string]interface{}{"foo": "bar"},
		},
		{
			test: "Inferrers",
			node: `node: string`,
			inferrers: []schemaInferrerInterface{
				NewSchemaTypeInferrer(),
				NewSchemaCallbackInferrer(func(_ goYamlAst.Node, schema map[string]interface{}) error {
					schema["foo"] = "baz"
					return nil
				}),
			},
			expected: map[string]interface{}{"foo": "baz", "type": "string"},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			node, _ := yaml.NewParser(yaml.WithComments()).ParseBytes([]byte(test.node))

			schema := map[string]interface{}{"foo": "bar"}
			err := NewSchemaChainInferrer(test.inferrers...).Infer(node, schema)

			s.NoError(err)

			s.Equal(test.expected, schema)
		})
	}
}

func (s *SchemaInferrerSuite) TestSchemaTypeInferrerErrors() {
	tests := []struct {
		test     string
		node     string
		expected *serrors.Assert
	}{
		{
			test: "UninferableNode",
			node: `string`,
			expected: &serrors.Assert{
				Type:    &yaml.NodeError{},
				Message: "unable to infer schema type",
				Arguments: []any{
					"line", 1,
					"column", 1,
				},
				Details: heredoc.Doc(`
					>  1 | string
					       ^
				`),
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			node, _ := yaml.NewParser(yaml.WithComments()).ParseBytes([]byte(test.node))

			schema := map[string]interface{}{"foo": "bar"}
			err := NewSchemaTypeInferrer().Infer(node, schema)

			serrors.Equal(s.Assert(), test.expected, err)
		})
	}
}

func (s *SchemaInferrerSuite) TestSchemaTypeInferrer() {
	tests := []struct {
		test     string
		node     string
		expected map[string]interface{}
	}{
		{
			test:     "String",
			node:     `node: string`,
			expected: map[string]interface{}{"foo": "bar", "type": "string"},
		},
		{
			test:     "Integer",
			node:     `node: 12`,
			expected: map[string]interface{}{"foo": "bar", "type": "integer"},
		},
		{
			test:     "Number",
			node:     `node: 3.4`,
			expected: map[string]interface{}{"foo": "bar", "type": "number"},
		},
		{
			test:     "Boolean",
			node:     `node: true`,
			expected: map[string]interface{}{"foo": "bar", "type": "boolean"},
		},
		{
			test:     "Null",
			node:     `node: ~`,
			expected: map[string]interface{}{"foo": "bar"},
		},
		{
			test: "ArrayEmpty",
			node: `
node: []
`,
			expected: map[string]interface{}{"foo": "bar", "type": "array"},
		},
		{
			test: "ArraySingle",
			node: `
node:
  - single
`,
			expected: map[string]interface{}{"foo": "bar", "type": "array"},
		},
		{
			test: "ArrayMultiple",
			node: `
node:
  - first
  - second
`,
			expected: map[string]interface{}{"foo": "bar", "type": "array"},
		},
		{
			test: "ObjectEmpty",
			node: `
node: {}
`,
			expected: map[string]interface{}{"foo": "bar", "type": "object"},
		},
		{
			test: "ObjectSingle",
			node: `
node:
  single: foo
`,
			expected: map[string]interface{}{"foo": "bar", "type": "object"},
		},
		{
			test: "ObjectMultiple",
			node: `
node:
  first: foo
  second: foo
`,
			expected: map[string]interface{}{"foo": "bar", "type": "object"},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			node, _ := yaml.NewParser(yaml.WithComments()).ParseBytes([]byte(test.node))

			schema := map[string]interface{}{"foo": "bar"}
			err := NewSchemaTypeInferrer().Infer(node, schema)

			s.NoError(err)

			s.Equal(test.expected, schema)
		})
	}
}

func (s *SchemaInferrerSuite) TestSchemaTagsInferrerErrors() {
	tests := []struct {
		test     string
		tags     *yaml.Tags
		expected *serrors.Assert
	}{
		{
			test: "Syntax",
			tags: &yaml.Tags{
				&yaml.Tag{Name: "Tag", Value: `foo`},
			},
			expected: &serrors.Assert{
				Type:    &json.SyntaxError{},
				Message: "invalid character 'o' in literal false (expecting 'a')",
			},
		},
		{
			test: "Type",
			tags: &yaml.Tags{
				&yaml.Tag{Name: "Tag", Value: `[]`},
			},
			expected: &serrors.Assert{
				Type:    &json.UnmarshalTypeError{},
				Message: "json: cannot unmarshal array into Go value of type map[string]interface {}",
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			schema := map[string]interface{}{"foo": "bar"}

			err := NewSchemaTagsInferrer(test.tags).Infer(nil, schema)

			serrors.Equal(s.Assert(), test.expected, err)
		})
	}
}

func (s *SchemaInferrerSuite) TestSchemaTagsInferrer() {
	tests := []struct {
		test     string
		tags     *yaml.Tags
		expected map[string]interface{}
	}{
		{
			test: "Extend",
			tags: &yaml.Tags{
				&yaml.Tag{Name: "Tag", Value: `{"bar": "baz"}`},
			},
			expected: map[string]interface{}{"foo": "bar", "bar": "baz"},
		},
		{
			test: "Override",
			tags: &yaml.Tags{
				&yaml.Tag{Name: "Tag", Value: `{"foo": "baz"}`},
			},
			expected: map[string]interface{}{"foo": "baz"},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			schema := map[string]interface{}{"foo": "bar"}

			err := NewSchemaTagsInferrer(test.tags).Infer(nil, schema)

			s.NoError(err)

			s.Equal(test.expected, schema)
		})
	}
}

func (s *SchemaInferrerSuite) TestSchemaCallbackInferrerErrors() {
	tests := []struct {
		test     string
		callback func(node goYamlAst.Node, schema map[string]interface{}) error
		expected *serrors.Assert
	}{
		{
			test: "Error",
			callback: func(node goYamlAst.Node, schema map[string]interface{}) error {
				return yaml.NewNodeError("foo", node)
			},
			expected: &serrors.Assert{
				Type:    &yaml.NodeError{},
				Message: "foo",
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			schema := map[string]interface{}{"foo": "bar"}

			err := NewSchemaCallbackInferrer(test.callback).Infer(nil, schema)

			serrors.Equal(s.Assert(), test.expected, err)
		})
	}
}

func (s *SchemaInferrerSuite) TestSchemaCallbackInferrer() {
	tests := []struct {
		test     string
		callback func(node goYamlAst.Node, schema map[string]interface{}) error
		expected map[string]interface{}
	}{
		{
			test: "Extend",
			callback: func(_ goYamlAst.Node, schema map[string]interface{}) error {
				schema["bar"] = "baz"

				return nil
			},
			expected: map[string]interface{}{"foo": "bar", "bar": "baz"},
		},
		{
			test: "Override",
			callback: func(_ goYamlAst.Node, schema map[string]interface{}) error {
				schema["foo"] = "baz"

				return nil
			},
			expected: map[string]interface{}{"foo": "baz"},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			schema := map[string]interface{}{"foo": "bar"}

			err := NewSchemaCallbackInferrer(test.callback).Infer(nil, schema)

			s.NoError(err)

			s.Equal(test.expected, schema)
		})
	}
}

func (s *SchemaInferrerSuite) TestSchemaInferrerErrors() {
	tests := []struct {
		test     string
		node     string
		expected *serrors.Assert
	}{
		{
			test: "NonMap",
			node: `string`,
			expected: &serrors.Assert{
				Type:    &yaml.NodeError{},
				Message: "unable to infer schema type",
				Arguments: []any{
					"line", 1,
					"column", 1,
				},
				Details: heredoc.Doc(`
					>  1 | string
					       ^
				`),
			},
		},
		{
			test: "MisplacedTag",
			node: `
node: ~  # @schema {"type": "string", "minLength": 1}
`,
			expected: &serrors.Assert{
				Type:    &yaml.NodeError{},
				Message: "misplaced schema tag",
				Arguments: []any{
					"line", 2,
					"column", 10,
				},
				Details: heredoc.Doc(`
					>  2 | node: ~  # @schema {"type": "string", "minLength": 1}
					                ^
				`),
			},
		},
		{
			test: "TagError",
			node: `
# @schema foo
node: ~
`,
			expected: &serrors.Assert{
				Type:    &yaml.NodeError{},
				Message: "invalid character 'o' in literal false (expecting 'a')",
				Arguments: []any{
					"line", 2,
					"column", 1,
				},
				Details: heredoc.Doc(`
					>  2 | # @schema foo
					       ^
					   3 | node: ~
				`),
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			node, _ := yaml.NewParser(yaml.WithComments()).ParseBytes([]byte(test.node))
			schema := map[string]interface{}{}

			err := NewSchemaInferrer().Infer(node, schema)

			serrors.Equal(s.Assert(), test.expected, err)
		})
	}
}

func (s *SchemaInferrerSuite) TestSchemaInferrer() {
	tests := []struct {
		test     string
		node     string
		expected map[string]interface{}
	}{
		{
			test: "Scalars",
			node: `
string: string
integer: 12
number: 2.3
boolean: true
`,
			expected: map[string]interface{}{
				"type":                 "object",
				"additionalProperties": false,
				"properties": map[string]interface{}{
					"string":  map[string]interface{}{"type": "string"},
					"integer": map[string]interface{}{"type": "integer"},
					"number":  map[string]interface{}{"type": "number"},
					"boolean": map[string]interface{}{"type": "boolean"},
				},
			},
		},
		{
			test: "Arrays",
			node: `
array_empty: []
array_single:
  - alone
array_multiple:
  - first
  - second
`,
			expected: map[string]interface{}{
				"type":                 "object",
				"additionalProperties": false,
				"properties": map[string]interface{}{
					"array_empty":    map[string]interface{}{"type": "array"},
					"array_single":   map[string]interface{}{"type": "array"},
					"array_multiple": map[string]interface{}{"type": "array"},
				},
			},
		},
		{
			test: "Objects",
			node: `
object_empty: {}
object_single:
  alone: foo
object_multiple:
  first: foo
  second: bar
`,
			expected: map[string]interface{}{
				"type":                 "object",
				"additionalProperties": false,
				"properties": map[string]interface{}{
					"object_empty": map[string]interface{}{
						"type": "object",
					},
					"object_single": map[string]interface{}{
						"type":                 "object",
						"additionalProperties": false,
						"properties": map[string]interface{}{
							"alone": map[string]interface{}{"type": "string"},
						},
					},
					"object_multiple": map[string]interface{}{
						"type":                 "object",
						"additionalProperties": false,
						"properties": map[string]interface{}{
							"first":  map[string]interface{}{"type": "string"},
							"second": map[string]interface{}{"type": "string"},
						},
					},
				},
			},
		},
		{
			test: "ScalarsTags",
			node: `
# @schema {"type": "string", "minLength": 1}
string: ~
# @schema {"minimum": 10}
integer: 12
`,
			expected: map[string]interface{}{
				"type":                 "object",
				"additionalProperties": false,
				"properties": map[string]interface{}{
					"string":  map[string]interface{}{"type": "string", "minLength": float64(1)},
					"integer": map[string]interface{}{"type": "integer", "minimum": float64(10)},
				},
			},
		},
		{
			test: "ArraysTags",
			node: `
# @schema {"items": {"type": "string"}}
array_empty: []
`,
			expected: map[string]interface{}{
				"type":                 "object",
				"additionalProperties": false,
				"properties": map[string]interface{}{
					"array_empty": map[string]interface{}{"type": "array", "items": map[string]interface{}{
						"type": "string",
					}},
				},
			},
		},
		{
			test: "ObjectsTags",
			node: `
# @schema {"additionalProperties": false}
object_empty: {}

object_single:
  # @schema {"type": "string"}
  alone: ~

object_multiple:
  # @schema {"type": "string", "minLength": 1}
  first: ~
  # @schema {"type": "string", "minLength": 2}
  second: ~

# @schema {"additionalProperties": true}
object_single_with_comment:
  # @schema {"type": "string"}
  alone: ~

# @schema {"additionalProperties": true}
object_multiple_with_comment:
  # @schema {"type": "string", "minLength": 1}
  first: ~
  # @schema {"type": "string", "minLength": 2}
  second: ~
`,
			expected: map[string]interface{}{
				"type":                 "object",
				"additionalProperties": false,
				"properties": map[string]interface{}{
					"object_empty": map[string]interface{}{
						"type":                 "object",
						"additionalProperties": false,
					},
					"object_single": map[string]interface{}{
						"type":                 "object",
						"additionalProperties": false,
						"properties": map[string]interface{}{
							"alone": map[string]interface{}{"type": "string"},
						},
					},
					"object_multiple": map[string]interface{}{
						"type":                 "object",
						"additionalProperties": false,
						"properties": map[string]interface{}{
							"first":  map[string]interface{}{"type": "string", "minLength": float64(1)},
							"second": map[string]interface{}{"type": "string", "minLength": float64(2)},
						},
					},
					"object_single_with_comment": map[string]interface{}{
						"type":                 "object",
						"additionalProperties": true,
						"properties": map[string]interface{}{
							"alone": map[string]interface{}{"type": "string"},
						},
					},
					"object_multiple_with_comment": map[string]interface{}{
						"type":                 "object",
						"additionalProperties": true,
						"properties": map[string]interface{}{
							"first":  map[string]interface{}{"type": "string", "minLength": float64(1)},
							"second": map[string]interface{}{"type": "string", "minLength": float64(2)},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			node, _ := yaml.NewParser(yaml.WithComments()).ParseBytes([]byte(test.node))
			schema := map[string]interface{}{}

			err := NewSchemaInferrer().Infer(node, schema)

			s.NoError(err)

			s.Equal(test.expected, schema)
		})
	}
}
