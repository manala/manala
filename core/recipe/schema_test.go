package recipe

import (
	yamlAst "github.com/goccy/go-yaml/ast"
	"github.com/stretchr/testify/suite"
	internalReport "manala/internal/report"
	internalYaml "manala/internal/yaml"
	"testing"
)

type SchemaInferrerSuite struct{ suite.Suite }

func TestSchemaInferrerSuite(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	suite.Run(t, new(SchemaInferrerSuite))
}

func (s *SchemaInferrerSuite) TestSchemaChainInferrerErrors() {
	tests := []struct {
		name      string
		node      string
		inferrers []schemaInferrerInterface
		err       string
	}{
		{
			name: "Error",
			node: `node: string`,
			inferrers: []schemaInferrerInterface{
				NewSchemaCallbackInferrer(func(node yamlAst.Node, _ map[string]interface{}) error {
					return internalYaml.NewNodeError("foo", node)
				}),
			},
			err: "foo",
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			node, _ := internalYaml.NewParser(internalYaml.WithComments()).ParseBytes([]byte(test.node))

			schema := map[string]interface{}{"foo": "bar"}
			err := NewSchemaChainInferrer(test.inferrers...).Infer(node, schema)

			s.EqualError(err, test.err)
		})
	}
}

func (s *SchemaInferrerSuite) TestSchemaChainInferrer() {
	tests := []struct {
		name      string
		node      string
		inferrers []schemaInferrerInterface
		schema    map[string]interface{}
	}{
		{
			name:      "No Inferrers",
			node:      `node: string`,
			inferrers: []schemaInferrerInterface{},
			schema:    map[string]interface{}{"foo": "bar"},
		},
		{
			name: "Inferrers",
			node: `node: string`,
			inferrers: []schemaInferrerInterface{
				NewSchemaTypeInferrer(),
				NewSchemaCallbackInferrer(func(_ yamlAst.Node, schema map[string]interface{}) error {
					schema["foo"] = "baz"

					return nil
				}),
			},
			schema: map[string]interface{}{"foo": "baz", "type": "string"},
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			node, _ := internalYaml.NewParser(internalYaml.WithComments()).ParseBytes([]byte(test.node))

			schema := map[string]interface{}{"foo": "bar"}
			err := NewSchemaChainInferrer(test.inferrers...).Infer(node, schema)

			s.NoError(err)
			s.Equal(test.schema, schema)
		})
	}
}

func (s *SchemaInferrerSuite) TestSchemaTypeInferrerErrors() {
	tests := []struct {
		name   string
		node   string
		report *internalReport.Assert
	}{
		{
			name: "Uninferable Node",
			node: `string`,
			report: &internalReport.Assert{
				Err: "unable to infer schema type",
				Fields: map[string]interface{}{
					"line":   1,
					"column": 1,
				},
				Trace: ">  1 | string\n       ^\n",
			},
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			node, _ := internalYaml.NewParser(internalYaml.WithComments()).ParseBytes([]byte(test.node))

			schema := map[string]interface{}{"foo": "bar"}
			err := NewSchemaTypeInferrer().Infer(node, schema)

			s.Error(err)

			report := internalReport.NewErrorReport(err)

			test.report.Equal(&s.Suite, report)
		})
	}
}

func (s *SchemaInferrerSuite) TestSchemaTypeInferrer() {
	tests := []struct {
		name   string
		node   string
		schema map[string]interface{}
		report *internalReport.Assert
	}{
		{
			name:   "String",
			node:   `node: string`,
			schema: map[string]interface{}{"foo": "bar", "type": "string"},
		},
		{
			name:   "Integer",
			node:   `node: 12`,
			schema: map[string]interface{}{"foo": "bar", "type": "integer"},
		},
		{
			name:   "Number",
			node:   `node: 3.4`,
			schema: map[string]interface{}{"foo": "bar", "type": "number"},
		},
		{
			name:   "Boolean",
			node:   `node: true`,
			schema: map[string]interface{}{"foo": "bar", "type": "boolean"},
		},
		{
			name:   "Null",
			node:   `node: ~`,
			schema: map[string]interface{}{"foo": "bar"},
		},
		{
			name: "Array Empty",
			node: `
node: []
`,
			schema: map[string]interface{}{"foo": "bar", "type": "array"},
		},
		{
			name: "Array Single",
			node: `
node:
  - single
`,
			schema: map[string]interface{}{"foo": "bar", "type": "array"},
		},
		{
			name: "Array Multiple",
			node: `
node:
  - first
  - second
`,
			schema: map[string]interface{}{"foo": "bar", "type": "array"},
		},
		{
			name: "Object Empty",
			node: `
node: {}
`,
			schema: map[string]interface{}{"foo": "bar", "type": "object"},
		},
		{
			name: "Object Single",
			node: `
node:
  single: foo
`,
			schema: map[string]interface{}{"foo": "bar", "type": "object"},
		},
		{
			name: "Object Multiple",
			node: `
node:
  first: foo
  second: foo
`,
			schema: map[string]interface{}{"foo": "bar", "type": "object"},
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			node, _ := internalYaml.NewParser(internalYaml.WithComments()).ParseBytes([]byte(test.node))

			schema := map[string]interface{}{"foo": "bar"}
			err := NewSchemaTypeInferrer().Infer(node, schema)

			s.NoError(err)
			s.Equal(test.schema, schema)
		})
	}
}

func (s *SchemaInferrerSuite) TestSchemaTagsInferrerErrors() {
	tests := []struct {
		name   string
		tags   *internalYaml.Tags
		report *internalReport.Assert
	}{
		{
			name: "Syntax",
			tags: &internalYaml.Tags{
				&internalYaml.Tag{Name: "Tag", Value: `foo`},
			},
			report: &internalReport.Assert{
				Err: "invalid character 'o' in literal false (expecting 'a')",
			},
		},
		{
			name: "Type",
			tags: &internalYaml.Tags{
				&internalYaml.Tag{Name: "Tag", Value: `[]`},
			},
			report: &internalReport.Assert{
				Err: "json: cannot unmarshal array into Go value of type map[string]interface {}",
			},
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			schema := map[string]interface{}{"foo": "bar"}
			err := NewSchemaTagsInferrer(test.tags).Infer(nil, schema)

			s.Error(err)

			report := internalReport.NewErrorReport(err)

			test.report.Equal(&s.Suite, report)
		})
	}
}

func (s *SchemaInferrerSuite) TestSchemaTagsInferrer() {
	tests := []struct {
		name   string
		tags   *internalYaml.Tags
		schema map[string]interface{}
	}{
		{
			name: "Extend",
			tags: &internalYaml.Tags{
				&internalYaml.Tag{Name: "Tag", Value: `{"bar": "baz"}`},
			},
			schema: map[string]interface{}{"foo": "bar", "bar": "baz"},
		},
		{
			name: "Override",
			tags: &internalYaml.Tags{
				&internalYaml.Tag{Name: "Tag", Value: `{"foo": "baz"}`},
			},
			schema: map[string]interface{}{"foo": "baz"},
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			schema := map[string]interface{}{"foo": "bar"}
			err := NewSchemaTagsInferrer(test.tags).Infer(nil, schema)

			s.NoError(err)
			s.Equal(test.schema, schema)
		})
	}
}

func (s *SchemaInferrerSuite) TestSchemaCallbackInferrerErrors() {
	tests := []struct {
		name     string
		callback func(node yamlAst.Node, schema map[string]interface{}) error
		err      string
	}{
		{
			name: "Error",
			callback: func(node yamlAst.Node, schema map[string]interface{}) error {
				return internalYaml.NewNodeError("foo", node)
			},
			err: "foo",
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			schema := map[string]interface{}{"foo": "bar"}
			err := NewSchemaCallbackInferrer(test.callback).Infer(nil, schema)

			s.EqualError(err, test.err)
		})
	}
}

func (s *SchemaInferrerSuite) TestSchemaCallbackInferrer() {
	tests := []struct {
		name     string
		callback func(node yamlAst.Node, schema map[string]interface{}) error
		schema   map[string]interface{}
	}{
		{
			name: "Extend",
			callback: func(_ yamlAst.Node, schema map[string]interface{}) error {
				schema["bar"] = "baz"

				return nil
			},
			schema: map[string]interface{}{"foo": "bar", "bar": "baz"},
		},
		{
			name: "Override",
			callback: func(_ yamlAst.Node, schema map[string]interface{}) error {
				schema["foo"] = "baz"

				return nil
			},
			schema: map[string]interface{}{"foo": "baz"},
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			schema := map[string]interface{}{"foo": "bar"}
			err := NewSchemaCallbackInferrer(test.callback).Infer(nil, schema)

			s.NoError(err)
			s.Equal(test.schema, schema)
		})
	}
}

func (s *SchemaInferrerSuite) TestSchemaInferrerErrors() {
	tests := []struct {
		name   string
		node   string
		report *internalReport.Assert
	}{
		{
			name: "Non Map",
			node: `string`,
			report: &internalReport.Assert{
				Err: "unable to infer schema type",
				Fields: map[string]interface{}{
					"line":   1,
					"column": 1,
				},
				Trace: ">  1 | string\n       ^\n",
			},
		},
		{
			name: "Misplaced Tag",
			node: `
node: ~  # @schema {"type": "string", "minLength": 1}
`,
			report: &internalReport.Assert{
				Err: "misplaced schema tag",
				Fields: map[string]interface{}{
					"line":   2,
					"column": 10,
				},
				Trace: ">  2 | node: ~  # @schema {\"type\": \"string\", \"minLength\": 1}\n                ^\n",
			},
		},
		{
			name: "Tag Error",
			node: `
# @schema foo
node: ~
`,
			report: &internalReport.Assert{
				Err: "invalid character 'o' in literal false (expecting 'a')",
				Fields: map[string]interface{}{
					"line":   2,
					"column": 1,
				},
				Trace: ">  2 | # @schema foo\n       ^\n   3 | node: ~",
			},
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			node, _ := internalYaml.NewParser(internalYaml.WithComments()).ParseBytes([]byte(test.node))

			schema := map[string]interface{}{}
			err := NewSchemaInferrer().Infer(node, schema)

			s.Error(err)

			report := internalReport.NewErrorReport(err)

			test.report.Equal(&s.Suite, report)
		})
	}
}

func (s *SchemaInferrerSuite) TestSchemaInferrer() {
	tests := []struct {
		name   string
		node   string
		schema map[string]interface{}
	}{
		{
			name: "Scalars",
			node: `
string: string
integer: 12
number: 2.3
boolean: true
`,
			schema: map[string]interface{}{
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
			name: "Arrays",
			node: `
array_empty: []
array_single:
  - alone
array_multiple:
  - first
  - second
`,
			schema: map[string]interface{}{
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
			name: "Objects",
			node: `
object_empty: {}
object_single:
  alone: foo
object_multiple:
  first: foo
  second: bar
`,
			schema: map[string]interface{}{
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
			name: "Scalars Tags",
			node: `
# @schema {"type": "string", "minLength": 1}
string: ~
# @schema {"minimum": 10}
integer: 12
`,
			schema: map[string]interface{}{
				"type":                 "object",
				"additionalProperties": false,
				"properties": map[string]interface{}{
					"string":  map[string]interface{}{"type": "string", "minLength": float64(1)},
					"integer": map[string]interface{}{"type": "integer", "minimum": float64(10)},
				},
			},
		},
		{
			name: "Arrays Tags",
			node: `
# @schema {"items": {"type": "string"}}
array_empty: []
`,
			schema: map[string]interface{}{
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
			name: "Objects Tags",
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
			schema: map[string]interface{}{
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
		s.Run(test.name, func() {
			node, _ := internalYaml.NewParser(internalYaml.WithComments()).ParseBytes([]byte(test.node))

			schema := map[string]interface{}{}
			err := NewSchemaInferrer().Infer(node, schema)

			s.NoError(err)
			s.Equal(test.schema, schema)
		})
	}
}
