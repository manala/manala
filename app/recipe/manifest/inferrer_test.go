package manifest_test

import (
	"encoding/json"
	"testing"

	"github.com/manala/manala/app"
	"github.com/manala/manala/app/recipe/manifest"
	"github.com/manala/manala/app/recipe/option"
	"github.com/manala/manala/internal/parsing"
	"github.com/manala/manala/internal/schema"
	"github.com/manala/manala/internal/serrors"
	"github.com/manala/manala/internal/testing/errors"
	"github.com/manala/manala/internal/yaml/parser"

	"github.com/stretchr/testify/suite"
)

type InferrerSuite struct{ suite.Suite }

func TestInferrerSuite(t *testing.T) {
	suite.Run(t, new(InferrerSuite))
}

func (s *InferrerSuite) TestSchemaErrors() {
	tests := []struct {
		test     string
		src      string
		expected errors.Assertion
	}{
		{
			test: "Annotation",
			src: `
# @schema foo
node: ~
`,
			expected: &parsing.FlattenErrorAssertion{
				Line:   2,
				Column: 12,
				Err: &serrors.Assertion{
					Message: "invalid character 'o' in literal false (expecting 'a')",
				},
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			node, err := parser.Parse([]byte(test.src))
			s.Require().NoError(err)

			inf := manifest.Inferrer{}
			err = inf.Infer(node)

			errors.Equal(s.T(), test.expected, err)
		})
	}
}

func (s *InferrerSuite) TestSchema() {
	tests := []struct {
		test     string
		src      string
		expected schema.Schema
	}{
		{
			test: "Scalars",
			src: `
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
			test: "ScalarsAnnotations",
			src: `
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
			test: "Arrays",
			src: `
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
			test: "ArraysAnnotations",
			src: `
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
			test: "Objects",
			src: `
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
						"type":                 "object",
						"additionalProperties": false,
						"properties":           map[string]any{},
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
			test: "ObjectsAnnotations",
			src: `
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
						"properties":           map[string]any{},
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
		{
			test: "Null",
			src:  `foo: ~`,
			expected: schema.Schema{
				"type":                 "object",
				"additionalProperties": false,
				"properties": map[string]any{
					"foo": map[string]any{},
				},
			},
		},
		{
			test: "TypeAlreadySet",
			src: `
# @schema {"type": "foo"}
node: string
`,
			expected: schema.Schema{
				"type":                 "object",
				"additionalProperties": false,
				"properties": map[string]any{
					"node": map[string]any{"type": "foo"},
				},
			},
		},
		{
			test: "Enum",
			src: `
# @schema {"enum": []}
node: string
`,
			expected: schema.Schema{
				"type":                 "object",
				"additionalProperties": false,
				"properties": map[string]any{
					"node": map[string]any{"enum": []any{}},
				},
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			node, err := parser.Parse([]byte(test.src))
			s.Require().NoError(err)

			sch := schema.Schema{}

			inf := manifest.Inferrer{
				Schema: &sch,
			}
			err = inf.Infer(node)
			s.Require().NoError(err)

			s.Equal(test.expected, sch)
		})
	}
}

func (s *InferrerSuite) TestOptionErrors() {
	tests := []struct {
		test     string
		src      string
		expected errors.Assertion
	}{
		{
			test: "Syntax",
			src: `
# @option foo
foo: ~
`,
			expected: &parsing.FlattenErrorAssertion{
				Line:   2,
				Column: 12,
				Err: &serrors.Assertion{
					Message: "invalid character 'o' in literal false (expecting 'a')",
				},
			},
		},
		{
			test: "Type",
			src: `
# @option []
foo: ~
`,
			expected: &parsing.FlattenErrorAssertion{
				Line:   2,
				Column: 11,
				Err: &serrors.Assertion{
					Message: "wrong array value type",
				},
			},
		},
		{
			test: "InvalidType",
			src: `
# @option {"label": "Label", "type": 123}
foo: ~
`,
			expected: &parsing.FlattenErrorAssertion{
				Line:   2,
				Column: 40,
				Err: &serrors.Assertion{
					Message: "wrong number type for field \"type\"",
				},
			},
		},
		{
			test: "UnexpectedType",
			src: `
# @option {"label": "Label", "type": "unexpected"}
foo: ~
`,
			expected: &parsing.FlattenErrorAssertion{
				Line:   2,
				Column: 11,
				Err: &serrors.Assertion{
					Message: "unexpected \"unexpected\" option type",
				},
			},
		},
		{
			test: "Validation",
			src: `
# @option {"foo": "bar"}
foo: string
`,
			expected: &parsing.FlattenErrorAssertion{
				Line:   2,
				Column: 1,
				Err: &serrors.Assertion{
					Message: "missing option label property",
				},
			},
		},
		{
			test: "AutoDetection",
			src: `
# @option {"label": "Label"}
foo: ~
`,
			expected: &parsing.FlattenErrorAssertion{
				Line:   2,
				Column: 11,
				Err: &serrors.Assertion{
					Message: "unable to auto detect option type",
				},
			},
		},
		{
			test: "InvalidTextMissingType",
			src: `
# @option {"label": "Label", "type": "text"}
foo: ~
`,
			expected: &parsing.FlattenErrorAssertion{
				Line:   2,
				Column: 11,
				Err: &serrors.Assertion{
					Message: "invalid recipe option string type",
				},
			},
		},
		{
			test: "InvalidTextWrongType",
			src: `
# @option {"label": "Label", "type": "text"}
# @schema {"type": null}
foo: ~
`,
			expected: &parsing.FlattenErrorAssertion{
				Line:   2,
				Column: 11,
				Err: &serrors.Assertion{
					Message: "invalid recipe option string type",
				},
			},
		},
		{
			test: "InvalidSelectMissingEnum",
			src: `
# @option {"label": "Label", "type": "select"}
foo: ~
`,
			expected: &parsing.FlattenErrorAssertion{
				Line:   2,
				Column: 11,
				Err: &serrors.Assertion{
					Message: "invalid recipe option enum",
				},
			},
		},
		{
			test: "InvalidSelectWrongEnum",
			src: `
# @option {"label": "Label", "type": "select"}
# @schema {"enum": null}
foo: ~
`,
			expected: &parsing.FlattenErrorAssertion{
				Line:   2,
				Column: 11,
				Err: &serrors.Assertion{
					Message: "invalid recipe option enum",
				},
			},
		},
		{
			test: "InvalidSelectEmptyEnum",
			src: `
# @option {"label": "Label", "type": "select"}
# @schema {"enum": []}
foo: ~
`,
			expected: &parsing.FlattenErrorAssertion{
				Line:   2,
				Column: 11,
				Err: &serrors.Assertion{
					Message: "empty recipe option enum",
				},
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			node, err := parser.Parse([]byte(test.src))
			s.Require().NoError(err)

			inf := manifest.Inferrer{}
			err = inf.Infer(node)

			errors.Equal(s.T(), test.expected, err)
		})
	}
}

func (s *InferrerSuite) TestOptions() {
	tests := []struct {
		test     string
		src      string
		expected option.Assertions
	}{
		{
			test: "Text",
			src: `
# @option {"label": "Foo", "name": "bar", "type": "text"}
foo: bar
`,
			expected: option.Assertions{
				{
					Type:  &option.Text{},
					Label: "Foo",
					Name:  "bar",
					Path:  "foo",
				},
			},
		},
		{
			test: "TextNoName",
			src: `
# @option {"label": "Foo Bar", "type": "text"}
foo: bar
`,
			expected: option.Assertions{
				{
					Type:  &option.Text{},
					Label: "Foo Bar",
					Name:  "foo-bar",
					Path:  "foo",
				},
			},
		},
		{
			test: "TextTypeImplicit",
			src: `
# @option {"label": "Foo", "name": "bar"}
foo: bar
`,
			expected: option.Assertions{
				{
					Type:  &option.Text{},
					Label: "Foo",
					Name:  "bar",
					Path:  "foo",
				},
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			node, err := parser.Parse([]byte(test.src))
			s.Require().NoError(err)

			var options []app.RecipeOption

			inf := manifest.Inferrer{
				Schema:  &schema.Schema{},
				Options: &options,
			}
			err = inf.Infer(node)
			s.Require().NoError(err)

			option.Equals(s.T(), test.expected, options)
		})
	}
}
