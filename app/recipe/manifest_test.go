package recipe

import (
	"encoding/json"
	"github.com/stretchr/testify/suite"
	"manala/app/recipe/option"
	"manala/internal/schema"
	"manala/internal/serrors"
	"manala/internal/syncer"
	"os"
	"path/filepath"
	"testing"
)

type ManifestSuite struct{ suite.Suite }

func TestManifestSuite(t *testing.T) {
	suite.Run(t, new(ManifestSuite))
}

func (s *ManifestSuite) Test() {
	manifest := NewManifest()

	s.Equal("", manifest.Description())
	s.Equal("", manifest.Icon())
	s.Equal("", manifest.Template())
	s.Equal(map[string]any{}, manifest.Vars())
	s.Equal([]syncer.UnitInterface{}, manifest.Sync())
	s.Equal(schema.Schema{}, manifest.Schema())
}

func (s *ManifestSuite) TestReadFromErrors() {
	tests := []struct {
		test     string
		expected *serrors.Assert
	}{
		{
			test: "Empty",
			expected: &serrors.Assert{
				Type:    serrors.Error{},
				Message: "irregular recipe manifest",
				Errors: []*serrors.Assert{
					{
						Type:    serrors.Error{},
						Message: "empty yaml file",
					},
				},
			},
		},
		{
			test: "Invalid",
			expected: &serrors.Assert{
				Type:    serrors.Error{},
				Message: "irregular recipe manifest",
				Errors: []*serrors.Assert{
					{
						Type:    serrors.Error{},
						Message: "unexpected mapping key",
						Arguments: []any{
							"line", 1,
							"column", 1,
						},
						Details: `
							>  1 | ::
							       ^
						`,
					},
				},
			},
		},
		{
			test: "IrregularType",
			expected: &serrors.Assert{
				Type:    serrors.Error{},
				Message: "irregular recipe manifest",
				Errors: []*serrors.Assert{
					{
						Type:    serrors.Error{},
						Message: "irregular type",
						Arguments: []any{
							"line", 1,
							"column", 6,
						},
						Details: `
							>  1 | foo: .inf
							            ^
						`,
					},
				},
			},
		},
		{
			test: "IrregularMapKey",
			expected: &serrors.Assert{
				Type:    serrors.Error{},
				Message: "irregular recipe manifest",
				Errors: []*serrors.Assert{
					{
						Type:    serrors.Error{},
						Message: "irregular map key",
						Arguments: []any{
							"line", 1,
							"column", 2,
						},
						Details: `
							>  1 | 0: bar
							        ^
						`,
					},
				},
			},
		},
		{
			test: "NotMap",
			expected: &serrors.Assert{
				Type:    serrors.Error{},
				Message: "invalid recipe manifest",
				Errors: []*serrors.Assert{
					{
						Type:    serrors.Error{},
						Message: "yaml document must be a map",
						Arguments: []any{
							"expected", "object",
							"actual", "string",
						},
					},
				},
			},
		},
		// Config
		{
			test: "ConfigAbsent",
			expected: &serrors.Assert{
				Type:    serrors.Error{},
				Message: "invalid recipe manifest",
				Errors: []*serrors.Assert{
					{
						Type:    serrors.Error{},
						Message: "missing manala property",
						Arguments: []any{
							"property", "manala",
						},
					},
				},
			},
		},
		{
			test: "ConfigNotMap",
			expected: &serrors.Assert{
				Type:    serrors.Error{},
				Message: "invalid recipe manifest",
				Errors: []*serrors.Assert{
					{
						Type:    serrors.Error{},
						Message: "manala field must be a map",
						Arguments: []any{
							"expected", "object",
							"actual", "string",
							"path", "manala",
							"line", 1,
							"column", 9,
						},
						Details: `
							>  1 | manala: foo
							               ^
						`,
					},
				},
			},
		},
		{
			test: "ConfigEmpty",
			expected: &serrors.Assert{
				Type:    serrors.Error{},
				Message: "invalid recipe manifest",
				Errors: []*serrors.Assert{
					{
						Type:    serrors.Error{},
						Message: "missing manala description property",
						Arguments: []any{
							"path", "manala",
							"property", "description",
							"line", 1,
							"column", 9,
						},
						Details: `
							>  1 | manala: {}
							               ^
						`,
					},
				},
			},
		},
		{
			test: "ConfigAdditionalProperties",
			expected: &serrors.Assert{
				Type:    serrors.Error{},
				Message: "invalid recipe manifest",
				Errors: []*serrors.Assert{
					{
						Type:    serrors.Error{},
						Message: "manala field don't support additional properties",
						Arguments: []any{
							"path", "manala.foo",
							"line", 3,
							"column", 8,
						},
						Details: `
							   1 | manala:
							   2 |   description: description
							>  3 |   foo: bar
							              ^
						`,
					},
				},
			},
		},
		// Config - Description
		{
			test: "ConfigDescriptionAbsent",
			expected: &serrors.Assert{
				Type:    serrors.Error{},
				Message: "invalid recipe manifest",
				Errors: []*serrors.Assert{
					{
						Type:    serrors.Error{},
						Message: "missing manala description property",
						Arguments: []any{
							"path", "manala",
							"property", "description",
							"line", 2,
							"column", 11,
						},
						Details: `
							   1 | manala:
							>  2 |   template: template
							                 ^
						`,
					},
				},
			},
		},
		{
			test: "ConfigDescriptionNotString",
			expected: &serrors.Assert{
				Type:    serrors.Error{},
				Message: "invalid recipe manifest",
				Errors: []*serrors.Assert{
					{
						Type:    serrors.Error{},
						Message: "manala description field must be a string",
						Arguments: []any{
							"expected", "string",
							"actual", "array",
							"path", "manala.description",
							"line", 2,
							"column", 16,
						},
						Details: `
							   1 | manala:
							>  2 |   description: []
							                      ^
							   3 |   template: template
						`,
					},
				},
			},
		},
		{
			test: "ConfigDescriptionEmpty",
			expected: &serrors.Assert{
				Type:    serrors.Error{},
				Message: "invalid recipe manifest",
				Errors: []*serrors.Assert{
					{
						Type:    serrors.Error{},
						Message: "empty manala description field",
						Arguments: []any{
							"minimum", 1,
							"path", "manala.description",
							"line", 2,
							"column", 16,
						},
						Details: `
							   1 | manala:
							>  2 |   description: ""
							                      ^
							   3 |   template: template
						`,
					},
				},
			},
		},
		{
			test: "ConfigDescriptionTooLong",
			expected: &serrors.Assert{
				Type:    serrors.Error{},
				Message: "invalid recipe manifest",
				Errors: []*serrors.Assert{
					{
						Type:    serrors.Error{},
						Message: "too long manala description field",
						Arguments: []any{
							"maximum", 256,
							"path", "manala.description",
							"line", 2,
							"column", 16,
						},
						Details: `
							   1 | manala:
							>  2 |   description: Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
							                      ^
							   3 |   template: template
						`,
					},
				},
			},
		},
		// Config - Icon
		{
			test: "ConfigIconNotString",
			expected: &serrors.Assert{
				Type:    serrors.Error{},
				Message: "invalid recipe manifest",
				Errors: []*serrors.Assert{
					{
						Type:    serrors.Error{},
						Message: "manala icon field must be a string",
						Arguments: []any{
							"expected", "string",
							"actual", "array",
							"path", "manala.icon",
							"line", 3,
							"column", 9,
						},
						Details: `
							   1 | manala:
							   2 |   description: description
							>  3 |   icon: []
							               ^
						`,
					},
				},
			},
		},
		{
			test: "ConfigIconEmpty",
			expected: &serrors.Assert{
				Type:    serrors.Error{},
				Message: "invalid recipe manifest",
				Errors: []*serrors.Assert{
					{
						Type:    serrors.Error{},
						Message: "empty manala icon field",
						Arguments: []any{
							"minimum", 1,
							"path", "manala.icon",
							"line", 3,
							"column", 9,
						},
						Details: `
							   1 | manala:
							   2 |   description: description
							>  3 |   icon: ""
							               ^
						`,
					},
				},
			},
		},
		{
			test: "ConfigIconTooLong",
			expected: &serrors.Assert{
				Type:    serrors.Error{},
				Message: "invalid recipe manifest",
				Errors: []*serrors.Assert{
					{
						Type:    serrors.Error{},
						Message: "too long manala icon field",
						Arguments: []any{
							"maximum", 100,
							"path", "manala.icon",
							"line", 3,
							"column", 9,
						},
						Details: `
							   1 | manala:
							   2 |   description: description
							>  3 |   icon: Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
							               ^
						`,
					},
				},
			},
		},
		// Config - Template
		{
			test: "ConfigTemplateNotString",
			expected: &serrors.Assert{
				Type:    serrors.Error{},
				Message: "invalid recipe manifest",
				Errors: []*serrors.Assert{
					{
						Type:    serrors.Error{},
						Message: "manala template field must be a string",
						Arguments: []any{
							"expected", "string",
							"actual", "array",
							"path", "manala.template",
							"line", 3,
							"column", 13,
						},
						Details: `
							   1 | manala:
							   2 |   description: description
							>  3 |   template: []
							                   ^
						`,
					},
				},
			},
		},
		{
			test: "ConfigTemplateEmpty",
			expected: &serrors.Assert{
				Type:    serrors.Error{},
				Message: "invalid recipe manifest",
				Errors: []*serrors.Assert{
					{
						Type:    serrors.Error{},
						Message: "empty manala template field",
						Arguments: []any{
							"minimum", 1,
							"path", "manala.template",
							"line", 3,
							"column", 13,
						},
						Details: `
							   1 | manala:
							   2 |   description: description
							>  3 |   template: ""
							                   ^
						`,
					},
				},
			},
		},
		{
			test: "ConfigTemplateTooLong",
			expected: &serrors.Assert{
				Type:    serrors.Error{},
				Message: "invalid recipe manifest",
				Errors: []*serrors.Assert{
					{
						Type:    serrors.Error{},
						Message: "too long manala template field",
						Arguments: []any{
							"maximum", 100,
							"path", "manala.template",
							"line", 3,
							"column", 13,
						},
						Details: `
							   1 | manala:
							   2 |   description: description
							>  3 |   template: Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
							                   ^
						`,
					},
				},
			},
		},
		// Config - Sync
		{
			test: "ConfigSyncNotArray",
			expected: &serrors.Assert{
				Type:    serrors.Error{},
				Message: "invalid recipe manifest",
				Errors: []*serrors.Assert{
					{
						Type:    serrors.Error{},
						Message: "manala sync field must be a sequence",
						Arguments: []any{
							"expected", "array",
							"actual", "string",
							"path", "manala.sync",
							"line", 3,
							"column", 9,
						},
						Details: `
							   1 | manala:
							   2 |   description: description
							>  3 |   sync: foo
							               ^
						`,
					},
				},
			},
		},
		// Config - Sync Item
		{
			test: "ConfigSyncItemNotString",
			expected: &serrors.Assert{
				Type:    serrors.Error{},
				Message: "invalid recipe manifest",
				Errors: []*serrors.Assert{
					{
						Type:    serrors.Error{},
						Message: "manala sync sequence entries must be strings",
						Arguments: []any{
							"expected", "string",
							"actual", "array",
							"path", "manala.sync[0]",
							"line", 4,
							"column", 7,
						},
						Details: `
							   1 | manala:
							   2 |   description: description
							   3 |   sync:
							>  4 |     - []
							             ^
						`,
					},
				},
			},
		},
		{
			test: "ConfigSyncItemEmpty",
			expected: &serrors.Assert{
				Type:    serrors.Error{},
				Message: "invalid recipe manifest",
				Errors: []*serrors.Assert{
					{
						Type:    serrors.Error{},
						Message: "empty manala sync sequence entry",
						Arguments: []any{
							"minimum", 1,
							"path", "manala.sync[0]",
							"line", 4,
							"column", 7,
						},
						Details: `
							   1 | manala:
							   2 |   description: description
							   3 |   sync:
							>  4 |     - ""
							             ^
						`,
					},
				},
			},
		},
		{
			test: "ConfigSyncItemTooLong",
			expected: &serrors.Assert{
				Type:    serrors.Error{},
				Message: "invalid recipe manifest",
				Errors: []*serrors.Assert{
					{
						Type:    serrors.Error{},
						Message: "too long manala sync sequence entry",
						Arguments: []any{
							"maximum", 256,
							"path", "manala.sync[0]",
							"line", 4,
							"column", 7,
						},
						Details: `
							   1 | manala:
							   2 |   description: description
							   3 |   sync:
							>  4 |     - Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
							             ^
						`,
					},
				},
			},
		},
		// Schema
		{
			test: "SchemaMisplacedTag",
			expected: &serrors.Assert{
				Type:    serrors.Error{},
				Message: "unable to infer recipe manifest schema",
				Errors: []*serrors.Assert{
					{
						Type:    serrors.Error{},
						Message: "misplaced schema tag",
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
			expected: &serrors.Assert{
				Type:    serrors.Error{},
				Message: "unable to infer recipe manifest schema",
				Errors: []*serrors.Assert{
					{
						Type:    serrors.Error{},
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
			manifest := NewManifest()

			dir := filepath.FromSlash("testdata/ManifestSuite/TestReadFromErrors/" + test.test)

			manifestFile, _ := os.Open(filepath.Join(dir, "manifest.yaml"))
			_, err := manifest.ReadFrom(manifestFile)

			serrors.Equal(s.Assert(), test.expected, err)
		})
	}
}

func (s *ManifestSuite) TestReadFrom() {
	tests := []struct {
		test                string
		expectedDescription string
		expectedIcon        string
		expectedTemplate    string
		expectedVars        map[string]any
		expectedSync        *syncer.UnitsAssert
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
			expectedSync: &syncer.UnitsAssert{
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
			expectedSync: &syncer.UnitsAssert{},
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
			expectedSync:        &syncer.UnitsAssert{},
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
			expectedSync: &syncer.UnitsAssert{},
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
			manifest := NewManifest()

			dir := filepath.FromSlash("testdata/ManifestSuite/TestReadFrom/" + test.test)

			manifestFile, _ := os.Open(filepath.Join(dir, "manifest.yaml"))
			_, err := manifest.ReadFrom(manifestFile)

			s.NoError(err)

			s.Equal(test.expectedDescription, manifest.Description())
			s.Equal(test.expectedIcon, manifest.Icon())
			s.Equal(test.expectedTemplate, manifest.Template())
			s.Equal(test.expectedVars, manifest.Vars())
			syncer.EqualUnits(s.Assert(), test.expectedSync, manifest.Sync())
			s.Equal(test.expectedSchema, manifest.Schema())
		})
	}
}

func (s *ManifestSuite) TestOptions() {
	manifest := NewManifest()

	dir := filepath.FromSlash("testdata/ManifestSuite/TestOptions")

	manifestFile, _ := os.Open(filepath.Join(dir, "manifest.yaml"))
	_, err := manifest.ReadFrom(manifestFile)

	options := manifest.Options()

	s.NoError(err)

	s.Len(options, 13)

	s.IsType(&option.TextOption{}, options[0])
	s.Equal("string", options[0].Name())
	s.Equal("String", options[0].Label())
	s.Equal("string", options[0].Path().String())
	s.Equal(0, options[0].(*option.TextOption).MaxLength)

	s.IsType(&option.TextOption{}, options[1])
	s.Equal("string-null", options[1].Name())
	s.Equal("String null", options[1].Label())
	s.Equal("string_null", options[1].Path().String())
	s.Equal(0, options[1].(*option.TextOption).MaxLength)

	s.IsType(&option.TextOption{}, options[2])
	s.Equal("string-max-length", options[2].Name())
	s.Equal("String max length", options[2].Label())
	s.Equal("string_max_length", options[2].Path().String())
	s.Equal(123, options[2].(*option.TextOption).MaxLength)

	s.IsType(&option.SelectOption{}, options[3])
	s.Equal("string-float-int", options[3].Name())
	s.Equal("String float int", options[3].Label())
	s.Equal("string_float_int", options[3].Path().String())
	s.Equal([]any{"3.0"}, options[3].(*option.SelectOption).Values)

	s.IsType(&option.SelectOption{}, options[4])
	s.Equal("string-asterisk", options[4].Name())
	s.Equal("String asterisk", options[4].Label())
	s.Equal("string_asterisk", options[4].Path().String())
	s.Equal([]any{"*"}, options[4].(*option.SelectOption).Values)

	s.IsType(&option.TextOption{}, options[5])
	s.Equal("map-single-first", options[5].Name())
	s.Equal("Map single first", options[5].Label())
	s.Equal("map_single.first", options[5].Path().String())
	s.Equal(0, options[5].(*option.TextOption).MaxLength)

	s.IsType(&option.TextOption{}, options[6])
	s.Equal("map-multiple-first", options[6].Name())
	s.Equal("Map multiple first", options[6].Label())
	s.Equal("map_multiple.first", options[6].Path().String())
	s.Equal(0, options[6].(*option.TextOption).MaxLength)

	s.IsType(&option.TextOption{}, options[7])
	s.Equal("map-multiple-second", options[7].Name())
	s.Equal("Map multiple second", options[7].Label())
	s.Equal("map_multiple.second", options[7].Path().String())
	s.Equal(0, options[7].(*option.TextOption).MaxLength)

	s.IsType(&option.SelectOption{}, options[8])
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

	s.IsType(&option.TextOption{}, options[9])
	s.Equal("underscore-key", options[9].Name())
	s.Equal("Underscore key", options[9].Label())
	s.Equal("underscore_key", options[9].Path().String())
	s.Equal(0, options[9].(*option.TextOption).MaxLength)

	s.IsType(&option.TextOption{}, options[10])
	s.Equal("hyphen-key", options[10].Name())
	s.Equal("Hyphen key", options[10].Label())
	s.Equal("hyphen-key", options[10].Path().String())
	s.Equal(0, options[10].(*option.TextOption).MaxLength)

	s.IsType(&option.TextOption{}, options[11])
	s.Equal("Dot key", options[11].Label())
	s.Equal("dot-key", options[11].Name())
	s.Equal("'dot.key'", options[11].Path().String())
	s.Equal(0, options[11].(*option.TextOption).MaxLength)

	s.IsType(&option.TextOption{}, options[12])
	s.Equal("foo-bar", options[12].Name())
	s.Equal("Custom name", options[12].Label())
	s.Equal(0, options[12].(*option.TextOption).MaxLength)
}
