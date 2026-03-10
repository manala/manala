package option_test

import (
	"github.com/manala/manala/app/recipe/option"
	"github.com/manala/manala/internal/path"
	"github.com/manala/manala/internal/schema"
	"github.com/manala/manala/internal/serrors"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

type Suite struct{ suite.Suite }

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) TestNewFromErrors() {
	tests := []struct {
		test     string
		data     string
		schema   schema.Schema
		expected *serrors.Assertion
	}{
		{
			test:   "Syntax",
			data:   `foo`,
			schema: schema.Schema{},
			expected: &serrors.Assertion{
				Message: "irregular recipe option",
				Errors: []*serrors.Assertion{
					{
						Message: "invalid character 'o' in literal false (expecting 'a')",
						Arguments: []any{
							"offset", int64(2),
						},
					},
				},
			},
		},
		{
			test:   "Type",
			data:   `[]`,
			schema: schema.Schema{},
			expected: &serrors.Assertion{
				Message: "irregular recipe option",
				Errors: []*serrors.Assertion{
					{
						Message: "cannot unmarshal into value",
						Arguments: []any{
							"offset", int64(1),
							"value", "array",
							"type", "map[string]interface {}",
						},
					},
				},
			},
		},
		{
			test:   "Validation",
			data:   `{"foo": "bar"}`,
			schema: schema.Schema{},
			expected: &serrors.Assertion{
				Message: "invalid recipe option",
				Errors: []*serrors.Assertion{
					{
						Message: "missing property",
						Arguments: []any{
							"property", "label",
						},
					},
					{
						Message: "additional property is not allowed",
						Arguments: []any{
							"path", "foo",
						},
					},
				},
			},
		},
		{
			test:   "AutoDetection",
			data:   `{"label": "Label"}`,
			schema: schema.Schema{},
			expected: &serrors.Assertion{
				Message: "unable to auto detect recipe option type",
				Arguments: []any{
					"label", "Label",
				},
			},
		},
		{
			test:   "InvalidTextMissingType",
			data:   `{"label": "Label", "type": "text"}`,
			schema: schema.Schema{},
			expected: &serrors.Assertion{
				Message: "invalid recipe option string type",
				Arguments: []any{
					"label", "Label",
				},
			},
		},
		{
			test:   "InvalidTextWrongType",
			data:   `{"label": "Label", "type": "text"}`,
			schema: schema.Schema{"type": nil},
			expected: &serrors.Assertion{
				Message: "invalid recipe option string type",
				Arguments: []any{
					"label", "Label",
				},
			},
		},
		{
			test:   "InvalidSelectMissingEnum",
			data:   `{"label": "Label", "type": "select"}`,
			schema: schema.Schema{},
			expected: &serrors.Assertion{
				Message: "invalid recipe option enum",
				Arguments: []any{
					"label", "Label",
				},
			},
		},
		{
			test:   "InvalidTextWrongEnum",
			data:   `{"label": "Label", "type": "select"}`,
			schema: schema.Schema{"enum": nil},

			expected: &serrors.Assertion{
				Message: "invalid recipe option enum",
				Arguments: []any{
					"label", "Label",
				},
			},
		},
		{
			test:   "InvalidTextEmptyEnum",
			data:   `{"label": "Label", "type": "select"}`,
			schema: schema.Schema{"enum": []any{}},
			expected: &serrors.Assertion{
				Message: "empty recipe option enum",
				Arguments: []any{
					"label", "Label",
				},
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			path := path.Path("")

			option, err := option.New(strings.NewReader(test.data), test.schema, path)

			s.Nil(option)
			serrors.Equal(s.T(), test.expected, err)
		})
	}
}

func (s *Suite) TestReadFrom() {
	tests := []struct {
		test          string
		data          string
		schema        schema.Schema
		expectedType  any
		expectedLabel string
		expectedName  string
	}{
		{
			test:          "Text",
			data:          `{"label": "Foo", "name": "bar", "type": "text"}`,
			schema:        schema.Schema{"type": "string"},
			expectedType:  &option.TextOption{},
			expectedLabel: "Foo",
			expectedName:  "bar",
		},
		{
			test:          "TextNoName",
			data:          `{"label": "Foo Bar", "type": "text"}`,
			schema:        schema.Schema{"type": "string"},
			expectedType:  &option.TextOption{},
			expectedLabel: "Foo Bar",
			expectedName:  "foo-bar",
		},
		{
			test:          "TextTypeImplicit",
			data:          `{"label": "Foo", "name": "bar"}`,
			schema:        schema.Schema{"type": "string"},
			expectedType:  &option.TextOption{},
			expectedLabel: "Foo",
			expectedName:  "bar",
		},
	}
	for _, test := range tests {
		s.Run(test.test, func() {
			path := path.Path("")

			option, err := option.New(strings.NewReader(test.data), test.schema, path)

			s.Require().NoError(err)
			s.IsType(test.expectedType, option)
			s.Equal(test.expectedLabel, option.Label())
			s.Equal(test.expectedName, option.Name())
		})
	}
}
