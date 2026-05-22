package manifest_test

import (
	"encoding/json"
	"testing"

	"github.com/manala/manala/app"
	"github.com/manala/manala/app/recipe/manifest"
	"github.com/manala/manala/app/recipe/option"
	jsonerrors "github.com/manala/manala/internal/json/errors"
	"github.com/manala/manala/internal/testing/expectation"
	"github.com/manala/manala/internal/testing/heredoc"
	"github.com/manala/manala/internal/validation"
	yamlannotation "github.com/manala/manala/internal/yaml/annotation"
	yamlerrors "github.com/manala/manala/internal/yaml/errors"
	yamlparser "github.com/manala/manala/internal/yaml/parser"

	"github.com/stretchr/testify/suite"
)

type InferrerSuite struct{ suite.Suite }

func TestInferrerSuite(t *testing.T) {
	suite.Run(t, new(InferrerSuite))
}

func (s *InferrerSuite) TestSchema() {
	tests := []struct {
		test     string
		src      string
		expected map[string]any
	}{
		{
			test: "Scalars",
			src: heredoc.Doc(`
				string: string
				integer: 12
				number: 2.3
				boolean: true
			`),
			expected: map[string]any{
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
			src: heredoc.Doc(`
				# @schema {"type": "string"}
				string_null: ~
				# @schema {"maxLength": 123}
				string_max_length: string
				# @schema {"enum": ["3.0"]}
				string_float_int: ~
				string_float_int_value: "3.0"
				# @schema {"enum": ["*"]}
				string_asterisk: ~
				string_asterisk_value: "*"
				# @schema {"minimum": 10}
				integer: 12
			`),
			expected: map[string]any{
				"type":                 "object",
				"additionalProperties": false,
				"properties": map[string]any{
					"string_null":            map[string]any{"type": "string"},
					"string_max_length":      map[string]any{"type": "string", "maxLength": json.Number("123")},
					"string_float_int":       map[string]any{"enum": []any{"3.0"}},
					"string_float_int_value": map[string]any{"type": "string"},
					"string_asterisk":        map[string]any{"enum": []any{"*"}},
					"string_asterisk_value":  map[string]any{"type": "string"},
					"integer":                map[string]any{"type": "integer", "minimum": json.Number("10")},
				},
			},
		},
		{
			test: "Arrays",
			src: heredoc.Doc(`
				array_empty: []
				array_single:
				  - alone
				array_multiple:
				  - first
				  - second
			`),
			expected: map[string]any{
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
			src: heredoc.Doc(`
				# @schema {"items": {"type": "string"}}
				array_empty: []
			`),
			expected: map[string]any{
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
			src: heredoc.Doc(`
				object_empty: {}
				object_single:
				  alone: foo
				object_multiple:
				  first: foo
				  second: bar
				object_nested:
				  string: string
				  object:
				    string: string
			`),
			expected: map[string]any{
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
					"object_nested": map[string]any{
						"type":                 "object",
						"additionalProperties": false,
						"properties": map[string]any{
							"string": map[string]any{
								"type": "string",
							},
							"object": map[string]any{
								"type":                 "object",
								"additionalProperties": false,
								"properties": map[string]any{
									"string": map[string]any{
										"type": "string",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			test: "ObjectsAnnotations",
			src: heredoc.Doc(`
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
			`),
			expected: map[string]any{
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
			src: heredoc.Doc(`
				foo: ~
			`),
			expected: map[string]any{
				"type":                 "object",
				"additionalProperties": false,
				"properties": map[string]any{
					"foo": map[string]any{},
				},
			},
		},
		{
			test: "TypeAlreadySet",
			src: heredoc.Doc(`
				# @schema {"type": "foo"}
				node: string
			`),
			expected: map[string]any{
				"type":                 "object",
				"additionalProperties": false,
				"properties": map[string]any{
					"node": map[string]any{"type": "foo"},
				},
			},
		},
		{
			test: "Enums",
			src: heredoc.Doc(`
				# @schema {"enum": [null, true, false, "string", 12, 2.3, 3.0, "3.0"]}
				enum: ~
				# @schema {"enum": []}
				enum_empty: string
			`),
			expected: map[string]any{
				"type":                 "object",
				"additionalProperties": false,
				"properties": map[string]any{
					"enum": map[string]any{
						"enum": []any{nil, true, false, "string", json.Number("12"), json.Number("2.3"), json.Number("3.0"), "3.0"},
					},
					"enum_empty": map[string]any{"enum": []any{}},
				},
			},
		},
		{
			test: "Keys",
			src: heredoc.Doc(`
				underscore_key: ok
				hyphen-key: ok
				dot.key: ok
				custom_name: ok
			`),
			expected: map[string]any{
				"type":                 "object",
				"additionalProperties": false,
				"properties": map[string]any{
					"underscore_key": map[string]any{"type": "string"},
					"hyphen-key":     map[string]any{"type": "string"},
					"dot.key":        map[string]any{"type": "string"},
					"custom_name":    map[string]any{"type": "string"},
				},
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			node, err := yamlparser.Parse([]byte(test.src))
			s.Require().NoError(err)

			sch := map[string]any{}

			inf := manifest.Inferrer{
				Schema: &sch,
			}
			err = inf.Infer(node)
			s.Require().NoError(err)

			s.Equal(test.expected, sch)
		})
	}
}

func (s *InferrerSuite) TestSchemaErrors() {
	tests := []struct {
		test     string
		src      string
		expected expectation.ErrorExpectation
	}{
		{
			test: "Annotation",
			src: heredoc.Doc(`
				# @schema foo
				node: ~
			`),
			expected: yamlerrors.Expectation{
				Position: [2]int{1, 1},
				Err: jsonerrors.Expectation{
					Position: [2]int{1, 12},
					Err:      expectation.ErrorMessage("invalid character 'o' in literal false (expecting 'a')"),
				},
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			node, err := yamlparser.Parse([]byte(test.src))
			s.Require().NoError(err)

			inf := manifest.Inferrer{}
			err = inf.Infer(node)

			expectation.ExpectError(s.T(), test.expected, err)
		})
	}
}

func (s *InferrerSuite) TestOptions() {
	tests := []struct {
		test     string
		src      string
		expected option.Expectations
	}{
		{
			test: "String",
			src: heredoc.Doc(`
				# @option {"label": "Foo", "name": "bar", "type": "string"}
				foo: bar
			`),
			expected: option.Expectations{
				{
					Type:  &option.String{},
					Label: "Foo",
					Name:  "bar",
				},
			},
		},
		{
			test: "StringNoName",
			src: heredoc.Doc(`
				# @option {"label": "Foo Bar", "type": "string"}
				foo: bar
			`),
			expected: option.Expectations{
				{
					Type:  &option.String{},
					Label: "Foo Bar",
					Name:  "foo-bar",
				},
			},
		},
		{
			test: "StringTypeImplicit",
			src: heredoc.Doc(`
				# @option {"label": "Foo", "name": "bar"}
				foo: bar
			`),
			expected: option.Expectations{
				{
					Type:  &option.String{},
					Label: "Foo",
					Name:  "bar",
				},
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			node, err := yamlparser.Parse([]byte(test.src))
			s.Require().NoError(err)

			var options []app.RecipeOption

			inf := manifest.Inferrer{
				Schema:  &map[string]any{},
				Options: &options,
			}
			err = inf.Infer(node)
			s.Require().NoError(err)

			option.ExpectOptions(s.T(), test.expected, options)
		})
	}
}

func (s *InferrerSuite) TestOptionErrors() {
	tests := []struct {
		test     string
		src      string
		expected expectation.ErrorExpectation
	}{
		{
			test: "Syntax",
			src: heredoc.Doc(`
				# @option foo
				foo: ~
			`),
			expected: yamlerrors.Expectation{
				Position: [2]int{1, 1},
				Err: jsonerrors.Expectation{
					Position: [2]int{1, 12},
					Err:      expectation.ErrorMessage("invalid character 'o' in literal false (expecting 'a')"),
				},
			},
		},
		{
			test: "Type",
			src: heredoc.Doc(`
				# @option []
				foo: ~
			`),
			expected: yamlerrors.Expectation{
				Position: [2]int{1, 1},
				Err: jsonerrors.Expectation{
					Position: [2]int{1, 11},
					Err:      expectation.ErrorMessage("wrong array value type"),
				},
			},
		},
		{
			test: "InvalidType",
			src: heredoc.Doc(`
				# @option {"label": "Label", "type": 123}
				foo: ~
			`),
			expected: yamlerrors.Expectation{
				Position: [2]int{1, 1},
				Err: jsonerrors.Expectation{
					Position: [2]int{1, 40},
					Err:      expectation.ErrorMessage("wrong number type for field \"type\""),
				},
			},
		},
		{
			test: "UnexpectedType",
			src: heredoc.Doc(`
				# @option {"label": "Label", "type": "unexpected"}
				foo: ~
			`),
			expected: yamlerrors.Expectation{
				Position: [2]int{1, 1},
				Err: yamlannotation.ErrorExpectation{
					Position: [2]int{1, 11},
					Err:      expectation.ErrorMessage("unexpected \"unexpected\" option type"),
				},
			},
		},
		{
			test: "Validation",
			src: heredoc.Doc(`
				# @option {"foo": "bar"}
				foo: string
			`),
			expected: yamlerrors.Expectation{
				Position: [2]int{1, 1},
				Err: expectation.Errors(
					validation.ViolationExpectation{
						Location: "",
						Position: [2]int{0, 0},
						Err:      expectation.ErrorMessage("missing property 'label'"),
					},
					validation.ViolationExpectation{
						Location: "/foo",
						Position: [2]int{1, 12},
						Err:      expectation.ErrorMessage("additional property 'foo' not allowed"),
					},
				),
			},
		},
		{
			test: "AutoDetection",
			src: heredoc.Doc(`
				# @option {"label": "Label"}
				foo: ~
			`),
			expected: yamlerrors.Expectation{
				Position: [2]int{1, 1},
				Err: yamlannotation.ErrorExpectation{
					Position: [2]int{1, 11},
					Err:      expectation.ErrorMessage("unable to auto detect option type"),
				},
			},
		},
		{
			test: "InvalidStringMissingType",
			src: heredoc.Doc(`
				# @option {"label": "Label", "type": "string"}
				foo: ~
			`),
			expected: yamlerrors.Expectation{
				Position: [2]int{1, 1},
				Err: yamlannotation.ErrorExpectation{
					Position: [2]int{1, 11},
					Err:      expectation.ErrorMessage("invalid recipe option string type"),
				},
			},
		},
		{
			test: "InvalidStringWrongType",
			src: heredoc.Doc(`
				# @option {"label": "Label", "type": "string"}
				# @schema {"type": null}
				foo: ~
			`),
			expected: yamlerrors.Expectation{
				Position: [2]int{1, 1},
				Err: yamlannotation.ErrorExpectation{
					Position: [2]int{1, 11},
					Err:      expectation.ErrorMessage("invalid recipe option string type"),
				},
			},
		},
		{
			test: "InvalidEnumMissingValues",
			src: heredoc.Doc(`
				# @option {"label": "Label", "type": "enum"}
				foo: ~
			`),
			expected: yamlerrors.Expectation{
				Position: [2]int{1, 1},
				Err: yamlannotation.ErrorExpectation{
					Position: [2]int{1, 11},
					Err:      expectation.ErrorMessage("invalid recipe option enum"),
				},
			},
		},
		{
			test: "InvalidEnumWrongValues",
			src: heredoc.Doc(`
				# @option {"label": "Label", "type": "enum"}
				# @schema {"enum": null}
				foo: ~
			`),
			expected: yamlerrors.Expectation{
				Position: [2]int{1, 1},
				Err: yamlannotation.ErrorExpectation{
					Position: [2]int{1, 11},
					Err:      expectation.ErrorMessage("invalid recipe option enum"),
				},
			},
		},
		{
			test: "InvalidEnumEmptyValues",
			src: heredoc.Doc(`
				# @option {"label": "Label", "type": "enum"}
				# @schema {"enum": []}
				foo: ~
			`),
			expected: yamlerrors.Expectation{
				Position: [2]int{1, 1},
				Err: yamlannotation.ErrorExpectation{
					Position: [2]int{1, 11},
					Err:      expectation.ErrorMessage("empty recipe option enum"),
				},
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			node, err := yamlparser.Parse([]byte(test.src))
			s.Require().NoError(err)

			inf := manifest.Inferrer{}
			err = inf.Infer(node)

			expectation.ExpectError(s.T(), test.expected, err)
		})
	}
}
