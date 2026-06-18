package validation_test

import (
	"testing"

	"github.com/manala/manala/internal/testing/expectation"
	"github.com/manala/manala/internal/validation"
	"github.com/manala/manala/internal/validation/validationtest"

	"github.com/stretchr/testify/suite"
)

type ValidatorLocatorSuite struct{ suite.Suite }

func TestValidatorSuite(t *testing.T) {
	suite.Run(t, new(ValidatorLocatorSuite))
}

func (s *ValidatorLocatorSuite) TestValidateViolations() {
	tests := []struct {
		test     string
		schema   map[string]any
		value    any
		expected expectation.ErrorExpectation
	}{
		{
			test: "AdditionalProperties",
			schema: map[string]any{
				"type":                 "object",
				"additionalProperties": false,
				"properties": map[string]any{
					"root_foo": map[string]any{"type": "string"},
					"root_bar": map[string]any{
						"type":                 "object",
						"additionalProperties": false,
						"properties": map[string]any{
							"nested_foo": map[string]any{"type": "string"},
						},
					},
				},
			},
			value: map[string]any{
				"root_foo": "string",
				"root_bar": map[string]any{
					"nested_foo": "string",
					"nested_bar": "string",
				},
				"root_baz": "string",
			},
			expected: expectation.Errors(
				validationtest.ViolationExpectation{
					Location: "/root_bar/nested_bar",
					Position: [2]int{0, 0},
					Err:      expectation.ErrorMessage("additional property 'nested_bar' not allowed"),
				},
				validationtest.ViolationExpectation{
					Location: "/root_baz",
					Position: [2]int{0, 0},
					Err:      expectation.ErrorMessage("additional property 'root_baz' not allowed"),
				},
			),
		},
		{
			test: "PropertyNames",
			schema: map[string]any{
				"type":          "object",
				"propertyNames": map[string]any{"pattern": "^[a-z_]+$"},
				"properties": map[string]any{
					"root_bar": map[string]any{
						"type":          "object",
						"propertyNames": map[string]any{"pattern": "^[a-z_]+$"},
					},
				},
			},
			value: map[string]any{
				"root_foo": "string",
				"root_bar": map[string]any{
					"nested_foo": "string",
					"nested_BAR": "string",
				},
				"root_BAZ": "string",
			},
			expected: expectation.Errors(
				validationtest.ViolationExpectation{
					Location: "",
					Position: [2]int{0, 0},
					Err:      expectation.ErrorMessage("'nested_BAR' does not match pattern '^[a-z_]+$'"),
				},
				validationtest.ViolationExpectation{
					Location: "",
					Position: [2]int{0, 0},
					Err:      expectation.ErrorMessage("'root_BAZ' does not match pattern '^[a-z_]+$'"),
				},
				validationtest.ViolationExpectation{
					Location: "",
					Position: [2]int{0, 0},
					Err:      expectation.ErrorMessage("invalid propertyName 'nested_BAR'"),
				},
				validationtest.ViolationExpectation{
					Location: "",
					Position: [2]int{0, 0},
					Err:      expectation.ErrorMessage("invalid propertyName 'root_BAZ'"),
				},
			),
		},
		{
			test: "PatternProperties",
			schema: map[string]any{
				"type":              "object",
				"patternProperties": map[string]any{"^root_(foo|baz)": map[string]any{"type": "integer"}},
				"properties": map[string]any{
					"root_bar": map[string]any{
						"type":              "object",
						"patternProperties": map[string]any{"^nested_foo": map[string]any{"type": "integer"}},
					},
				},
			},
			value: map[string]any{
				"root_foo": 123,
				"root_bar": map[string]any{
					"nested_foo": "string",
				},
				"root_baz": "string",
			},
			expected: expectation.Errors(
				validationtest.ViolationExpectation{
					Location: "/root_bar/nested_foo",
					Position: [2]int{0, 0},
					Err:      expectation.ErrorMessage("got string, want integer"),
				},
				validationtest.ViolationExpectation{
					Location: "/root_baz",
					Position: [2]int{0, 0},
					Err:      expectation.ErrorMessage("got string, want integer"),
				},
			),
		},
		{
			test: "Required",
			schema: map[string]any{
				"type":     "object",
				"required": []any{"root_foo"},
				"properties": map[string]any{
					"root_bar": map[string]any{
						"type":     "object",
						"required": []any{"nested_foo"},
					},
				},
			},
			value: map[string]any{
				"root_bar": map[string]any{},
			},
			expected: expectation.Errors(
				validationtest.ViolationExpectation{
					Location: "",
					Position: [2]int{0, 0},
					Err:      expectation.ErrorMessage("missing property 'root_foo'"),
				},
				validationtest.ViolationExpectation{
					Location: "/root_bar",
					Position: [2]int{0, 0},
					Err:      expectation.ErrorMessage("missing property 'nested_foo'"),
				},
			),
		},
		{
			test: "DependentRequired",
			schema: map[string]any{
				"type":              "object",
				"dependentRequired": map[string]any{"root_foo": []any{"root_baz"}},
				"properties": map[string]any{
					"root_bar": map[string]any{
						"type":              "object",
						"dependentRequired": map[string]any{"nested_foo": []any{"nested_baz"}},
					},
				},
			},
			value: map[string]any{
				"root_foo": "string",
				"root_bar": map[string]any{
					"nested_foo": "string",
				},
			},
			expected: expectation.Errors(
				validationtest.ViolationExpectation{
					Location: "",
					Position: [2]int{0, 0},
					Err:      expectation.ErrorMessage("properties 'root_baz' required, if 'root_foo' exists"),
				},
				validationtest.ViolationExpectation{
					Location: "/root_bar",
					Position: [2]int{0, 0},
					Err:      expectation.ErrorMessage("properties 'nested_baz' required, if 'nested_foo' exists"),
				},
			),
		},
		{
			test: "MinProperties",
			schema: map[string]any{
				"type":          "object",
				"minProperties": 3,
				"properties": map[string]any{
					"root_bar": map[string]any{
						"type":          "object",
						"minProperties": 2,
					},
				},
			},
			value: map[string]any{
				"root_foo": "string",
				"root_bar": map[string]any{
					"nested_foo": "string",
				},
			},
			expected: expectation.Errors(
				validationtest.ViolationExpectation{
					Location: "",
					Position: [2]int{0, 0},
					Err:      expectation.ErrorMessage("minProperties: got 2, want 3"),
				},
				validationtest.ViolationExpectation{
					Location: "/root_bar",
					Position: [2]int{0, 0},
					Err:      expectation.ErrorMessage("minProperties: got 1, want 2"),
				},
			),
		},
		{
			test: "MaxProperties",
			schema: map[string]any{
				"type":          "object",
				"maxProperties": 1,
				"properties": map[string]any{
					"root_bar": map[string]any{
						"type":          "object",
						"maxProperties": 1,
					},
				},
			},
			value: map[string]any{
				"root_foo": "string",
				"root_bar": map[string]any{
					"nested_foo": "string",
					"nested_bar": "string",
				},
			},
			expected: expectation.Errors(
				validationtest.ViolationExpectation{
					Location: "",
					Position: [2]int{0, 0},
					Err:      expectation.ErrorMessage("maxProperties: got 2, want 1"),
				},
				validationtest.ViolationExpectation{
					Location: "/root_bar",
					Position: [2]int{0, 0},
					Err:      expectation.ErrorMessage("maxProperties: got 2, want 1"),
				},
			),
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			validator, err := validation.NewValidator(test.schema)
			s.Require().NoError(err)

			err = validator.Validate(test.value)
			expectation.ExpectError(s.T(), test.expected, err)
		})
	}
}
