package option

import (
	"github.com/stretchr/testify/suite"
	"manala/internal/path"
	"manala/internal/schema"
	"manala/internal/serrors"
	"strings"
	"testing"
)

type Suite struct{ suite.Suite }

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) Test() {
	name := "name"
	label := "Label"
	schema := schema.Schema{"foo": "bar"}

	option := &option{
		name:   name,
		label:  label,
		schema: schema,
	}

	s.Equal(name, option.Name())
	s.Equal(label, option.Label())
	s.Equal(schema, option.Schema())
}

func (s *Suite) TestNewFromErrors() {
	tests := []struct {
		test     string
		data     string
		schema   schema.Schema
		expected *serrors.Assert
	}{
		{
			test:   "Syntax",
			data:   `foo`,
			schema: schema.Schema{},
			expected: &serrors.Assert{
				Type:    serrors.Error{},
				Message: "irregular recipe option",
				Errors: []*serrors.Assert{
					{
						Type:    serrors.Error{},
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
			expected: &serrors.Assert{
				Type:    serrors.Error{},
				Message: "irregular recipe option",
				Errors: []*serrors.Assert{
					{
						Type:    serrors.Error{},
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
			expected: &serrors.Assert{
				Type:    serrors.Error{},
				Message: "invalid recipe option",
				Errors: []*serrors.Assert{
					{
						Type:    serrors.Error{},
						Message: "missing property",
						Arguments: []any{
							"property", "label",
						},
					},
					{
						Type:    serrors.Error{},
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
			expected: &serrors.Assert{
				Type:    serrors.Error{},
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
			expected: &serrors.Assert{
				Type:    serrors.Error{},
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
			expected: &serrors.Assert{
				Type:    serrors.Error{},
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
			expected: &serrors.Assert{
				Type:    serrors.Error{},
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

			expected: &serrors.Assert{
				Type:    serrors.Error{},
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
			expected: &serrors.Assert{
				Type:    serrors.Error{},
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

			option, err := New(strings.NewReader(test.data), test.schema, path)

			s.Nil(option)
			serrors.Equal(s.Assert(), test.expected, err)
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
			expectedType:  &TextOption{},
			expectedLabel: "Foo",
			expectedName:  "bar",
		},
		{
			test:          "TextNoName",
			data:          `{"label": "Foo Bar", "type": "text"}`,
			schema:        schema.Schema{"type": "string"},
			expectedType:  &TextOption{},
			expectedLabel: "Foo Bar",
			expectedName:  "foo-bar",
		},
		{
			test:          "TextTypeImplicit",
			data:          `{"label": "Foo", "name": "bar"}`,
			schema:        schema.Schema{"type": "string"},
			expectedType:  &TextOption{},
			expectedLabel: "Foo",
			expectedName:  "bar",
		},
	}
	for _, test := range tests {
		s.Run(test.test, func() {
			path := path.Path("")

			option, err := New(strings.NewReader(test.data), test.schema, path)

			s.NoError(err)
			s.IsType(test.expectedType, option)
			s.Equal(test.expectedLabel, option.Label())
			s.Equal(test.expectedName, option.Name())
		})
	}
}
