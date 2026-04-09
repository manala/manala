package manifest_test

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/manala/manala/app/recipe/manifest"
	"github.com/manala/manala/app/recipe/option"
	"github.com/manala/manala/internal/parsing"
	"github.com/manala/manala/internal/schema"
	"github.com/manala/manala/internal/serrors"
	"github.com/manala/manala/internal/sync"
	"github.com/manala/manala/internal/testing/errors"

	"github.com/stretchr/testify/suite"
)

type Suite struct{ suite.Suite }

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) Test() {
	m := manifest.New()

	s.Empty(m.Description())
	s.Empty(m.Icon())
	s.Empty(m.Template())
	s.Equal(map[string]any{}, m.Vars())
	s.Equal([]sync.UnitInterface{}, m.Sync())
	s.Equal(schema.Schema{}, m.Schema())
}

func (s *Suite) TestUnmarshalYAMLErrors() {
	tests := []struct {
		test     string
		expected errors.Assertion
	}{
		{
			test: "Empty",
			expected: &parsing.Assertion{
				Err: &serrors.Assertion{
					Message: "empty yaml content",
				},
			},
		},
		{
			test: "Invalid",
			expected: &parsing.Assertion{
				Line:   1,
				Column: 1,
				Err: &serrors.Assertion{
					Message: "'@' is a reserved character",
				},
			},
		},
		{
			test: "IrregularType",
			expected: &parsing.Assertion{
				Line:   1,
				Column: 6,
				Err: &serrors.Assertion{
					Message: "irregular type",
				},
			},
		},
		{
			test: "IrregularMapKey",
			expected: &parsing.Assertion{
				Line:   1,
				Column: 2,
				Err: &serrors.Assertion{
					Message: "irregular map key",
				},
			},
		},
		{
			test: "NotMap",
			expected: &parsing.Assertion{
				Line:   1,
				Column: 1,
				Err: &serrors.Assertion{
					Message: "yaml document must be a map",
				},
			},
		},
		// Config
		{
			test: "ConfigAbsent",
			expected: &parsing.Assertion{
				Err: &serrors.Assertion{
					Message: "missing manala property",
				},
			},
		},
		{
			test: "ConfigNotMap",
			expected: &parsing.Assertion{
				Line:   1,
				Column: 9,
				Err: &serrors.Assertion{
					Message: "string was used where mapping is expected",
				},
			},
		},
		{
			test: "ConfigEmpty",
			expected: &parsing.Assertion{
				Line:   1,
				Column: 1,
				Err: &serrors.Assertion{
					Message: "missing manala description property",
				},
			},
		},
		{
			test: "ConfigAdditionalProperties",
			expected: &parsing.Assertion{
				Line:   3,
				Column: 3,
				Err: &serrors.Assertion{
					Message: "unknown field \"foo\"",
				},
			},
		},
		// Config - Description
		{
			test: "ConfigDescriptionAbsent",
			expected: &parsing.Assertion{
				Line:   1,
				Column: 7,
				Err: &serrors.Assertion{
					Message: "missing manala description property",
				},
			},
		},
		{
			test: "ConfigDescriptionNotString",
			expected: &parsing.Assertion{
				Line:   2,
				Column: 16,
				Err: &serrors.Assertion{
					Message: "field must be a string",
				},
			},
		},
		{
			test: "ConfigDescriptionEmpty",
			expected: &parsing.Assertion{
				Line:   2,
				Column: 16,
				Err: &serrors.Assertion{
					Message: "missing manala description property",
				},
			},
		},
		{
			test: "ConfigDescriptionTooLong",
			expected: &parsing.Assertion{
				Line:   2,
				Column: 16,
				Err: &serrors.Assertion{
					Message: "too long manala description field (max=256)",
				},
			},
		},
		// Config - Icon
		{
			test: "ConfigIconNotString",
			expected: &parsing.Assertion{
				Line:   3,
				Column: 9,
				Err: &serrors.Assertion{
					Message: "field must be a string",
				},
			},
		},
		{
			test: "ConfigIconTooLong",
			expected: &parsing.Assertion{
				Line:   3,
				Column: 9,
				Err: &serrors.Assertion{
					Message: "too long manala icon field (max=100)",
				},
			},
		},
		// Config - Template
		{
			test: "ConfigTemplateNotString",
			expected: &parsing.Assertion{
				Line:   3,
				Column: 13,
				Err: &serrors.Assertion{
					Message: "field must be a string",
				},
			},
		},
		{
			test: "ConfigTemplateTooLong",
			expected: &parsing.Assertion{
				Line:   3,
				Column: 13,
				Err: &serrors.Assertion{
					Message: "too long manala template field (max=100)",
				},
			},
		},
		// Config - Sync
		{
			test: "ConfigSyncNotArray",
			expected: &parsing.Assertion{
				Line:   3,
				Column: 9,
				Err: &serrors.Assertion{
					Message: "sync field must be a sequence",
				},
			},
		},
		// Config - Sync Item
		{
			test: "ConfigSyncItemNotString",
			expected: &parsing.Assertion{
				Line:   4,
				Column: 5,
				Err: &serrors.Assertion{
					Message: "sync entry must be a string",
				},
			},
		},
		{
			test: "ConfigSyncItemEmpty",
			expected: &parsing.Assertion{
				Line:   4,
				Column: 5,
				Err: &serrors.Assertion{
					Message: "empty sync entry",
				},
			},
		},
		{
			test: "ConfigSyncItemTooLong",
			expected: &parsing.Assertion{
				Line:   4,
				Column: 5,
				Err: &serrors.Assertion{
					Message: "too long sync entry (max=256)",
				},
			},
		},
		// Schema
		{
			test: "SchemaMisplacedAnnotation",
			expected: &serrors.Assertion{
				Message: "unable to infer recipe manifest schema",
				Errors: []errors.Assertion{
					&serrors.Assertion{
						Message: "misplaced schema annotation",
						Arguments: []any{
							"line", 4,
							"column", 9,
						},
						Details: `
							   1 | manala:
							   2 |   description: description
							   3 |
							>  4 | foo: ~  # @schema {"type": "string", "minLength": 1}
							               ^
						`,
					},
				},
			},
		},
		{
			test: "SchemaInvalidJson",
			expected: &serrors.Assertion{
				Message: "unable to infer recipe manifest schema",
				Errors: []errors.Assertion{
					&serrors.Assertion{
						Message: "invalid character 'o' in literal false (expecting 'a')",
						Arguments: []any{
							"line", 4,
							"column", 1,
						},
						Details: `
							   1 | manala:
							   2 |   description: description
							   3 |
							>  4 | # @schema foo
							       ^
							   5 | foo: ~
						`,
					},
				},
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			m := manifest.New()

			dir := filepath.FromSlash("testdata/Suite/TestUnmarshalYAMLErrors")

			reader, _ := os.Open(filepath.Join(dir, test.test+".yaml"))
			content, _ := io.ReadAll(reader)

			err := m.UnmarshalYAML(content)

			errors.Equal(s.T(), test.expected, err)
		})
	}
}

func (s *Suite) TestUnmarshalYAML() {
	tests := []struct {
		test                string
		expectedDescription string
		expectedIcon        string
		expectedTemplate    string
		expectedVars        map[string]any
		expectedSync        *sync.UnitsAssertion
		expectedSchema      schema.Schema
	}{
		{
			test:                "All",
			expectedDescription: "description",
			expectedIcon:        "icon",
			expectedTemplate:    "template",
			expectedVars: map[string]any{
				"string":      "string",
				"string_null": nil,
				"sequence": []any{
					"first",
				},
				"sequence_string_empty": []any{},
				"boolean":               true,
				"integer":               uint64(123),
				"float":                 1.2,
				"map": map[string]any{
					"string": "string",
					"map": map[string]any{
						"string": "string",
					},
				},
				"map_empty": map[string]any{},
				"map_single": map[string]any{
					"first": "foo",
				},
				"map_multiple": map[string]any{
					"first":  "foo",
					"second": "foo",
				},
				"enum":           nil,
				"underscore_key": "ok",
				"hyphen-key":     "ok",
				"dot.key":        "ok",
			},
			expectedSync: &sync.UnitsAssertion{
				{Source: "file", Destination: "file"},
				{Source: "dir/file", Destination: "dir/file"},
				{Source: "file", Destination: "dir/file"},
				{Source: "dir/file", Destination: "file"},
				{Source: "src_file", Destination: "dst_file"},
				{Source: "src_dir/file", Destination: "dst_dir/file"},
			},
			expectedSchema: schema.Schema{
				"type":                 "object",
				"additionalProperties": false,
				"properties": map[string]any{
					"string": map[string]any{
						"type": "string",
					},
					"string_null": map[string]any{
						"type": "string",
					},
					"sequence": map[string]any{
						"type": "array",
					},
					"sequence_string_empty": map[string]any{
						"type": "array",
						"items": map[string]any{
							"type": "string",
						},
					},
					"boolean": map[string]any{
						"type": "boolean",
					},
					"integer": map[string]any{
						"type": "integer",
					},
					"float": map[string]any{
						"type": "number",
					},
					"map": map[string]any{
						"type":                 "object",
						"additionalProperties": false,
						"properties": map[string]any{
							"string": map[string]any{
								"type": "string",
							},
							"map": map[string]any{
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
					"map_empty": map[string]any{
						"type": "object",
					},
					"map_single": map[string]any{
						"type":                 "object",
						"additionalProperties": false,
						"properties": map[string]any{
							"first": map[string]any{
								"type":      "string",
								"minLength": json.Number("1"),
							},
						},
					},
					"map_multiple": map[string]any{
						"type":                 "object",
						"additionalProperties": false,
						"properties": map[string]any{
							"first": map[string]any{
								"type":      "string",
								"minLength": json.Number("1"),
							},
							"second": map[string]any{
								"type":      "string",
								"minLength": json.Number("1"),
							},
						},
					},
					"enum": map[string]any{
						"enum": []any{
							nil,
							true,
							false,
							"string",
							json.Number("12"),
							json.Number("2.3"),
							json.Number("3.0"),
							"3.0",
						},
					},
					"underscore_key": map[string]any{
						"type": "string",
					},
					"hyphen-key": map[string]any{
						"type": "string",
					},
					"dot.key": map[string]any{
						"type": "string",
					},
				},
			},
		},
		{
			test:                "ConfigTemplateAbsent",
			expectedDescription: "description",
			expectedIcon:        "",
			expectedTemplate:    "",
			expectedVars: map[string]any{
				"foo": "bar",
			},
			expectedSync: &sync.UnitsAssertion{},
			expectedSchema: schema.Schema{
				"type":                 "object",
				"additionalProperties": false,
				"properties": map[string]any{
					"foo": map[string]any{
						"type": "string",
					},
				},
			},
		},
		{
			test:                "VarsAbsent",
			expectedDescription: "description",
			expectedIcon:        "icon",
			expectedTemplate:    "template",
			expectedVars:        map[string]any{},
			expectedSync:        &sync.UnitsAssertion{},
			expectedSchema:      schema.Schema{},
		},
		{
			test:                "VarsKeys",
			expectedDescription: "description",
			expectedIcon:        "icon",
			expectedTemplate:    "template",
			expectedVars: map[string]any{
				"underscore_key": "ok",
				"hyphen-key":     "ok",
				"dot.key":        "ok",
			},
			expectedSync: &sync.UnitsAssertion{},
			expectedSchema: schema.Schema{
				"type":                 "object",
				"additionalProperties": false,
				"properties": map[string]any{
					"underscore_key": map[string]any{
						"type": "string",
					},
					"hyphen-key": map[string]any{
						"type": "string",
					},
					"dot.key": map[string]any{
						"type": "string",
					},
				},
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			m := manifest.New()

			dir := filepath.FromSlash("testdata/Suite/TestUnmarshalYAML")

			reader, _ := os.Open(filepath.Join(dir, test.test+".yaml"))
			content, _ := io.ReadAll(reader)

			err := m.UnmarshalYAML(content)

			s.Require().NoError(err)

			s.Equal(test.expectedDescription, m.Description())
			s.Equal(test.expectedIcon, m.Icon())
			s.Equal(test.expectedTemplate, m.Template())
			s.Equal(test.expectedVars, m.Vars())
			sync.EqualUnits(s.T(), test.expectedSync, m.Sync())
			s.Equal(test.expectedSchema, m.Schema())
		})
	}
}

func (s *Suite) TestOptions() {
	m := manifest.New()

	dir := filepath.FromSlash("testdata/Suite/TestOptions")

	reader, _ := os.Open(filepath.Join(dir, "manifest.yaml"))
	content, _ := io.ReadAll(reader)

	err := m.UnmarshalYAML(content)

	options := m.Options()

	s.Require().NoError(err)

	s.Require().Len(options, 13)

	s.Require().IsType((*option.TextOption)(nil), options[0])
	s.Equal("string", options[0].Name())
	s.Equal("String", options[0].Label())
	s.Equal("string", options[0].Path().String())
	s.Equal(0, options[0].(*option.TextOption).MaxLength)

	s.Require().IsType((*option.TextOption)(nil), options[1])
	s.Equal("string-null", options[1].Name())
	s.Equal("String null", options[1].Label())
	s.Equal("string_null", options[1].Path().String())
	s.Equal(0, options[1].(*option.TextOption).MaxLength)

	s.Require().IsType((*option.TextOption)(nil), options[2])
	s.Equal("string-max-length", options[2].Name())
	s.Equal("String max length", options[2].Label())
	s.Equal("string_max_length", options[2].Path().String())
	s.Equal(123, options[2].(*option.TextOption).MaxLength)

	s.Require().IsType((*option.SelectOption)(nil), options[3])
	s.Equal("string-float-int", options[3].Name())
	s.Equal("String float int", options[3].Label())
	s.Equal("string_float_int", options[3].Path().String())
	s.Equal([]any{"3.0"}, options[3].(*option.SelectOption).Values)

	s.Require().IsType((*option.SelectOption)(nil), options[4])
	s.Equal("string-asterisk", options[4].Name())
	s.Equal("String asterisk", options[4].Label())
	s.Equal("string_asterisk", options[4].Path().String())
	s.Equal([]any{"*"}, options[4].(*option.SelectOption).Values)

	s.Require().IsType((*option.TextOption)(nil), options[5])
	s.Equal("map-single-first", options[5].Name())
	s.Equal("Map single first", options[5].Label())
	s.Equal("map_single.first", options[5].Path().String())
	s.Equal(0, options[5].(*option.TextOption).MaxLength)

	s.Require().IsType((*option.TextOption)(nil), options[6])
	s.Equal("map-multiple-first", options[6].Name())
	s.Equal("Map multiple first", options[6].Label())
	s.Equal("map_multiple.first", options[6].Path().String())
	s.Equal(0, options[6].(*option.TextOption).MaxLength)

	s.Require().IsType((*option.TextOption)(nil), options[7])
	s.Equal("map-multiple-second", options[7].Name())
	s.Equal("Map multiple second", options[7].Label())
	s.Equal("map_multiple.second", options[7].Path().String())
	s.Equal(0, options[7].(*option.TextOption).MaxLength)

	s.Require().IsType((*option.SelectOption)(nil), options[8])
	s.Equal("Enum null", options[8].Label())
	s.Equal("enum-null", options[8].Name())
	s.Equal("enum", options[8].Path().String())
	s.Equal([]any{
		nil,
		true,
		false,
		"string",
		int64(12),
		2.3,
		3.0,
		"3.0",
	}, options[8].(*option.SelectOption).Values)

	s.Require().IsType((*option.TextOption)(nil), options[9])
	s.Equal("underscore-key", options[9].Name())
	s.Equal("Underscore key", options[9].Label())
	s.Equal("underscore_key", options[9].Path().String())
	s.Equal(0, options[9].(*option.TextOption).MaxLength)

	s.Require().IsType((*option.TextOption)(nil), options[10])
	s.Equal("hyphen-key", options[10].Name())
	s.Equal("Hyphen key", options[10].Label())
	s.Equal("hyphen-key", options[10].Path().String())
	s.Equal(0, options[10].(*option.TextOption).MaxLength)

	s.Require().IsType((*option.TextOption)(nil), options[11])
	s.Equal("Dot key", options[11].Label())
	s.Equal("dot-key", options[11].Name())
	s.Equal("'dot.key'", options[11].Path().String())
	s.Equal(0, options[11].(*option.TextOption).MaxLength)

	s.Require().IsType((*option.TextOption)(nil), options[12])
	s.Equal("foo-bar", options[12].Name())
	s.Equal("Custom name", options[12].Label())
	s.Equal(0, options[12].(*option.TextOption).MaxLength)
}
