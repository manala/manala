package yaml

import (
	"encoding/json"
	"manala/internal/schema"
	"manala/internal/serrors"
)

/*************/
/* Inferrers */
/*************/

func (s *Suite) TestNodeSchemaInferrerErrors() {
	tests := []struct {
		test     string
		node     string
		expected *serrors.Assert
	}{
		{
			test: "NonMap",
			node: `string`,
			expected: &serrors.Assert{
				Type:    serrors.Error{},
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
			test: "MisplacedTag",
			node: `
node: ~  # @schema {"type": "string", "minLength": 1}
`,
			expected: &serrors.Assert{
				Type:    serrors.Error{},
				Message: "misplaced schema tag",
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
			test: "TagError",
			node: `
# @schema foo
node: ~
`,
			expected: &serrors.Assert{
				Type:    serrors.Error{},
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
			node, _ := NewParser(WithComments()).ParseBytes([]byte(test.node))
			schema := schema.Schema{}

			err := NewNodeSchemaInferrer(node).Infer(schema)

			serrors.Equal(s.Assert(), test.expected, err)
		})
	}
}

func (s *Suite) TestNodeSchemaInferrer() {
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
			test: "ScalarsTags",
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
			test: "ArraysTags",
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
			node, _ := NewParser(WithComments()).ParseBytes([]byte(test.node))
			schema := schema.Schema{}

			err := NewNodeSchemaInferrer(node).Infer(schema)

			s.NoError(err)

			s.Equal(test.expected, schema)
		})
	}
}

func (s *Suite) TestNodeTypeSchemaInferrerErrors() {
	tests := []struct {
		test     string
		node     string
		expected *serrors.Assert
	}{
		{
			test: "UninferableNode",
			node: `string`,
			expected: &serrors.Assert{
				Type:    serrors.Error{},
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
			node, _ := NewParser(WithComments()).ParseBytes([]byte(test.node))

			schema := schema.Schema{"foo": "bar"}
			err := NewNodeTypeSchemaInferrer(node).Infer(schema)

			serrors.Equal(s.Assert(), test.expected, err)
		})
	}
}

func (s *Suite) TestNodeTypeSchemaInferrer() {
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
			node, _ := NewParser(WithComments()).ParseBytes([]byte(test.node))

			err := NewNodeTypeSchemaInferrer(node).Infer(test.schema)

			s.NoError(err)

			s.Equal(test.expected, test.schema)
		})
	}
}

func (s *Suite) TestNodeTagsSchemaInferrerErrors() {
	tests := []struct {
		test     string
		tags     *Tags
		expected *serrors.Assert
	}{
		{
			test: "Syntax",
			tags: &Tags{
				&Tag{Name: "Tag", Value: `foo`},
			},
			expected: &serrors.Assert{
				Type:    serrors.Error{},
				Message: "invalid character 'o' in literal false (expecting 'a')",
				Arguments: []any{
					"offset", int64(2),
				},
			},
		},
		{
			test: "Type",
			tags: &Tags{
				&Tag{Name: "Tag", Value: `[]`},
			},
			expected: &serrors.Assert{
				Type:    serrors.Error{},
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

			err := NewNodeTagsSchemaInferrer(nil, test.tags).Infer(schema)

			serrors.Equal(s.Assert(), test.expected, err)
		})
	}
}

func (s *Suite) TestNodeTagsSchemaInferrer() {
	tests := []struct {
		test     string
		tags     *Tags
		expected schema.Schema
	}{
		{
			test: "Extend",
			tags: &Tags{
				&Tag{Name: "Tag", Value: `{"bar": "baz"}`},
			},
			expected: schema.Schema{"foo": "bar", "bar": "baz"},
		},
		{
			test: "Override",
			tags: &Tags{
				&Tag{Name: "Tag", Value: `{"foo": "baz"}`},
			},
			expected: schema.Schema{"foo": "baz"},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			schema := schema.Schema{"foo": "bar"}

			err := NewNodeTagsSchemaInferrer(nil, test.tags).Infer(schema)

			s.NoError(err)

			s.Equal(test.expected, schema)
		})
	}
}
