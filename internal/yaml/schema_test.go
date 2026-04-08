package yaml_test

import (
	"encoding/json"
	"testing"

	"github.com/manala/manala/internal/schema"
	"github.com/manala/manala/internal/serrors"
	"github.com/manala/manala/internal/testing/errors"
	"github.com/manala/manala/internal/yaml"
	"github.com/manala/manala/internal/yaml/annotation"

	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
	"github.com/stretchr/testify/suite"
)

type SchemaSuite struct{ suite.Suite }

func TestSchemaSuite(t *testing.T) {
	suite.Run(t, new(SchemaSuite))
}

/*************/
/* Inferrers */
/*************/

func (s *SchemaSuite) TestNodeInferrerErrors() {
	tests := []struct {
		test     string
		node     string
		expected errors.Assertion
	}{
		{
			test: "NotMap",
			node: `string`,
			expected: &serrors.Assertion{
				Message: "unable to infer schema type",
				Arguments: []any{
					"line", 1,
					"column", 1,
				},
				Details: `
					>  1 | string
					       ^
				`,
			},
		},
		{
			test: "MisplacedAnnotation",
			node: `
node: ~  # @schema {"type": "string", "minLength": 1}
`,
			expected: &serrors.Assertion{
				Message: "misplaced schema annotation",
				Arguments: []any{
					"line", 2,
					"column", 10,
				},
				Details: `
					>  2 | node: ~  # @schema {"type": "string", "minLength": 1}
					                ^
				`,
			},
		},
		{
			test: "AnnotationError",
			node: `
# @schema foo
node: ~
`,
			expected: &serrors.Assertion{
				Message: "invalid character 'o' in literal false (expecting 'a')",
				Arguments: []any{
					"line", 2,
					"column", 1,
				},
				Details: `
					>  2 | # @schema foo
					       ^
					   3 | node: ~
				`,
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			file, _ := parser.ParseBytes([]byte(test.node), parser.ParseComments)
			node := file.Docs[0].Body

			schema := schema.Schema{}
			err := yaml.NewNodeSchemaInferrer(node).Infer(schema)

			errors.Equal(s.T(), test.expected, err)
		})
	}
}

func (s *SchemaSuite) TestNodeInferrer() {
	tests := []struct {
		test     string
		node     string
		expected schema.Schema
	}{
		{
			test: "Scalars",
			node: `
string: string
integer: 12
number: 2.3
boolean: true
`,
			expected: schema.Schema{
				"type":                 "object",
				"additionalProperties": false,
				"properties": map[string]any{
					"string":  map[string]any{"type": "string"},
					"integer": map[string]any{"type": "integer"},
					"number":  map[string]any{"type": "number"},
					"boolean": map[string]any{"type": "boolean"},
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
			expected: schema.Schema{
				"type":                 "object",
				"additionalProperties": false,
				"properties": map[string]any{
					"array_empty":    map[string]any{"type": "array"},
					"array_single":   map[string]any{"type": "array"},
					"array_multiple": map[string]any{"type": "array"},
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
			expected: schema.Schema{
				"type":                 "object",
				"additionalProperties": false,
				"properties": map[string]any{
					"object_empty": map[string]any{
						"type": "object",
					},
					"object_single": map[string]any{
						"type":                 "object",
						"additionalProperties": false,
						"properties": map[string]any{
							"alone": map[string]any{"type": "string"},
						},
					},
					"object_multiple": map[string]any{
						"type":                 "object",
						"additionalProperties": false,
						"properties": map[string]any{
							"first":  map[string]any{"type": "string"},
							"second": map[string]any{"type": "string"},
						},
					},
				},
			},
		},
		{
			test: "ScalarsAnnotations",
			node: `
# @schema {"type": "string", "minLength": 1}
string: ~
# @schema {"minimum": 10}
integer: 12
`,
			expected: schema.Schema{
				"type":                 "object",
				"additionalProperties": false,
				"properties": map[string]any{
					"string":  map[string]any{"type": "string", "minLength": json.Number("1")},
					"integer": map[string]any{"type": "integer", "minimum": json.Number("10")},
				},
			},
		},
		{
			test: "ArraysAnnotations",
			node: `
# @schema {"items": {"type": "string"}}
array_empty: []
`,
			expected: schema.Schema{
				"type":                 "object",
				"additionalProperties": false,
				"properties": map[string]any{
					"array_empty": map[string]any{"type": "array", "items": map[string]any{
						"type": "string",
					}},
				},
			},
		},
		{
			test: "ObjectsAnnotations",
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
			expected: schema.Schema{
				"type":                 "object",
				"additionalProperties": false,
				"properties": map[string]any{
					"object_empty": map[string]any{
						"type":                 "object",
						"additionalProperties": false,
					},
					"object_single": map[string]any{
						"type":                 "object",
						"additionalProperties": false,
						"properties": map[string]any{
							"alone": map[string]any{"type": "string"},
						},
					},
					"object_multiple": map[string]any{
						"type":                 "object",
						"additionalProperties": false,
						"properties": map[string]any{
							"first":  map[string]any{"type": "string", "minLength": json.Number("1")},
							"second": map[string]any{"type": "string", "minLength": json.Number("2")},
						},
					},
					"object_single_with_comment": map[string]any{
						"type":                 "object",
						"additionalProperties": true,
						"properties": map[string]any{
							"alone": map[string]any{"type": "string"},
						},
					},
					"object_multiple_with_comment": map[string]any{
						"type":                 "object",
						"additionalProperties": true,
						"properties": map[string]any{
							"first":  map[string]any{"type": "string", "minLength": json.Number("1")},
							"second": map[string]any{"type": "string", "minLength": json.Number("2")},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			file, _ := parser.ParseBytes([]byte(test.node), parser.ParseComments)
			node := file.Docs[0].Body

			schema := schema.Schema{}
			err := yaml.NewNodeSchemaInferrer(node).Infer(schema)

			s.Require().NoError(err)
			s.Equal(test.expected, schema)
		})
	}
}

func (s *SchemaSuite) TestNodeTypeInferrerErrors() {
	tests := []struct {
		test     string
		node     string
		expected errors.Assertion
	}{
		{
			test: "UninferableNode",
			node: `string`,
			expected: &serrors.Assertion{
				Message: "unable to infer schema type",
				Arguments: []any{
					"line", 1,
					"column", 1,
				},
				Details: `
					>  1 | string
					       ^
				`,
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			file, _ := parser.ParseBytes([]byte(test.node), 0)
			node := file.Docs[0].Body

			schema := schema.Schema{"foo": "bar"}
			err := yaml.NewNodeTypeSchemaInferrer(node).Infer(schema)

			errors.Equal(s.T(), test.expected, err)
		})
	}
}

func (s *SchemaSuite) TestNodeTypeInferrer() {
	tests := []struct {
		test     string
		node     string
		schema   schema.Schema
		expected schema.Schema
	}{
		{
			test:     "String",
			node:     `node: string`,
			schema:   schema.Schema{"foo": "bar"},
			expected: schema.Schema{"foo": "bar", "type": "string"},
		},
		{
			test:     "Integer",
			node:     `node: 12`,
			schema:   schema.Schema{"foo": "bar"},
			expected: schema.Schema{"foo": "bar", "type": "integer"},
		},
		{
			test:     "Number",
			node:     `node: 3.4`,
			schema:   schema.Schema{"foo": "bar"},
			expected: schema.Schema{"foo": "bar", "type": "number"},
		},
		{
			test:     "Boolean",
			node:     `node: true`,
			schema:   schema.Schema{"foo": "bar"},
			expected: schema.Schema{"foo": "bar", "type": "boolean"},
		},
		{
			test: "ArrayEmpty",
			node: `
node: []
`,
			schema:   schema.Schema{"foo": "bar"},
			expected: schema.Schema{"foo": "bar", "type": "array"},
		},
		{
			test: "ArraySingle",
			node: `
node:
  - single
`,
			schema:   schema.Schema{"foo": "bar"},
			expected: schema.Schema{"foo": "bar", "type": "array"},
		},
		{
			test: "ArrayMultiple",
			node: `
node:
  - first
  - second
`,
			schema:   schema.Schema{"foo": "bar"},
			expected: schema.Schema{"foo": "bar", "type": "array"},
		},
		{
			test: "ObjectEmpty",
			node: `
node: {}
`,
			schema:   schema.Schema{"foo": "bar"},
			expected: schema.Schema{"foo": "bar", "type": "object"},
		},
		{
			test: "ObjectSingle",
			node: `
node:
  single: foo
`,
			schema:   schema.Schema{"foo": "bar"},
			expected: schema.Schema{"foo": "bar", "type": "object"},
		},
		{
			test: "ObjectMultiple",
			node: `
node:
  first: foo
  second: foo
`,
			schema:   schema.Schema{"foo": "bar"},
			expected: schema.Schema{"foo": "bar", "type": "object"},
		},
		{
			test:     "Null",
			node:     `node: ~`,
			schema:   schema.Schema{"foo": "bar"},
			expected: schema.Schema{"foo": "bar"},
		},
		{
			test:     "AlreadySet",
			node:     `node: string`,
			schema:   schema.Schema{"foo": "bar", "type": "foo"},
			expected: schema.Schema{"foo": "bar", "type": "foo"},
		},
		{
			test:     "Enum",
			node:     `node: string`,
			schema:   schema.Schema{"foo": "bar", "enum": []any{}},
			expected: schema.Schema{"foo": "bar", "enum": []any{}},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			file, _ := parser.ParseBytes([]byte(test.node), 0)
			node := file.Docs[0].Body

			s.Require().IsType((*ast.MappingNode)(nil), node)
			s.Require().Len(node.(*ast.MappingNode).Values, 1)

			err := yaml.NewNodeTypeSchemaInferrer(node.(*ast.MappingNode).Values[0]).Infer(test.schema)

			s.Require().NoError(err)
			s.Equal(test.expected, test.schema)
		})
	}
}

func (s *SchemaSuite) TestNodeAnnotationsInferrerErrors() {
	tests := []struct {
		test     string
		src      string
		expected errors.Assertion
	}{
		{
			test: "Syntax",
			src:  `# @annotation foo`,
			expected: &serrors.Assertion{
				Message: "invalid character 'o' in literal false (expecting 'a')",
				Arguments: []any{
					"offset", int64(2),
				},
			},
		},
		{
			test: "Type",
			src:  `# @annotation []`,
			expected: &serrors.Assertion{
				Message: "cannot unmarshal into value",
				Arguments: []any{
					"offset", int64(1),
					"value", "array",
					"type", "schema.Schema",
				},
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			schema := schema.Schema{"foo": "bar"}

			annots, _ := annotation.Parse(test.src)
			annot, _ := annots.Lookup("annotation")

			err := yaml.NewNodeAnnotationSchemaInferrer(nil, annot).Infer(schema)

			errors.Equal(s.T(), test.expected, err)
		})
	}
}

func (s *SchemaSuite) TestNodeTagsInferrer() {
	tests := []struct {
		test     string
		src      string
		expected schema.Schema
	}{
		{
			test:     "Extend",
			src:      `# @annotation {"bar": "baz"}`,
			expected: schema.Schema{"foo": "bar", "bar": "baz"},
		},
		{
			test:     "Override",
			src:      `# @annotation {"foo": "baz"}`,
			expected: schema.Schema{"foo": "baz"},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			schema := schema.Schema{"foo": "bar"}

			annots, _ := annotation.Parse(test.src)
			annot, _ := annots.Lookup("annotation")

			err := yaml.NewNodeAnnotationSchemaInferrer(nil, annot).Infer(schema)

			s.Require().NoError(err)
			s.Equal(test.expected, schema)
		})
	}
}
