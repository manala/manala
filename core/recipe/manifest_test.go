package recipe

import (
	"github.com/stretchr/testify/suite"
	"manala/app/interfaces"
	"manala/internal/errors/serrors"
	"manala/internal/syncer"
	"manala/internal/testing/heredoc"
	"manala/internal/validation"
	"manala/internal/yaml"
	"os"
	"path/filepath"
	"testing"
)

type ManifestSuite struct{ suite.Suite }

func TestManifestSuite(t *testing.T) {
	suite.Run(t, new(ManifestSuite))
}

func (s *ManifestSuite) Test() {
	recMan := NewManifest()

	s.Equal("", recMan.Description())
	s.Equal("", recMan.Template())
	s.Equal(map[string]interface{}{}, recMan.Vars())
	s.Equal([]syncer.UnitInterface{}, recMan.Sync())
	s.Equal(map[string]interface{}{}, recMan.Schema())
}

func (s *ManifestSuite) TestReadFromErrors() {
	tests := []struct {
		test     string
		expected *serrors.Assert
	}{
		{
			test: "Empty",
			expected: &serrors.Assert{
				Type:    &serrors.WrapError{},
				Message: "irregular recipe manifest",
				Error: &serrors.Assert{
					Type:    &serrors.Error{},
					Message: "empty yaml file",
				},
			},
		},
		{
			test: "Invalid",
			expected: &serrors.Assert{
				Type:    &serrors.WrapError{},
				Message: "irregular recipe manifest",
				Error: &serrors.Assert{
					Type:    &yaml.Error{},
					Message: "unexpected mapping key",
					Arguments: []any{
						"line", 1,
						"column", 1,
					},
					Details: heredoc.Doc(`
						>  1 | ::
						       ^
					`),
				},
			},
		},
		{
			test: "IrregularType",
			expected: &serrors.Assert{
				Type:    &serrors.WrapError{},
				Message: "irregular recipe manifest",
				Error: &serrors.Assert{
					Type:    &yaml.NodeError{},
					Message: "irregular type",
					Arguments: []any{
						"line", 1,
						"column", 6,
					},
					Details: heredoc.Doc(`
						>  1 | foo: .inf
						            ^
					`),
				},
			},
		},
		{
			test: "IrregularMapKey",
			expected: &serrors.Assert{
				Type:    &serrors.WrapError{},
				Message: "irregular recipe manifest",
				Error: &serrors.Assert{
					Type:    &yaml.NodeError{},
					Message: "irregular map key",
					Arguments: []any{
						"line", 1,
						"column", 2,
					},
					Details: heredoc.Doc(`
						>  1 | 0: bar
						        ^
					`),
				},
			},
		},
		{
			test: "NotMap",
			expected: &serrors.Assert{
				Type:    &validation.Error{},
				Message: "invalid recipe manifest",
				Errors: []*serrors.Assert{
					{
						Type:    &yaml.NodeValidationResultError{},
						Message: "yaml document must be a map",
						Arguments: []any{
							"expected", "object",
							"given", "string",
						},
					},
				},
			},
		},
		// Config
		{
			test: "ConfigAbsent",
			expected: &serrors.Assert{
				Type:    &validation.Error{},
				Message: "invalid recipe manifest",
				Errors: []*serrors.Assert{
					{
						Type:    &yaml.NodeValidationResultError{},
						Message: "missing manala field",
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
				Type:    &validation.Error{},
				Message: "invalid recipe manifest",
				Errors: []*serrors.Assert{
					{
						Type:    &yaml.NodeValidationResultError{},
						Message: "manala field must be a map",
						Arguments: []any{
							"expected", "object",
							"given", "string",
							"line", 1,
							"column", 9,
						},
						Details: heredoc.Doc(`
							>  1 | manala: foo
							               ^
						`),
					},
				},
			},
		},
		{
			test: "ConfigEmpty",
			expected: &serrors.Assert{
				Type:    &validation.Error{},
				Message: "invalid recipe manifest",
				Errors: []*serrors.Assert{
					{
						Type:    &yaml.NodeValidationResultError{},
						Message: "missing manala description field",
						Arguments: []any{
							"property", "description",
							"line", 1,
							"column", 9,
						},
						Details: heredoc.Doc(`
							>  1 | manala: {}
							               ^
						`),
					},
				},
			},
		},
		{
			test: "ConfigAdditionalProperties",
			expected: &serrors.Assert{
				Type:    &validation.Error{},
				Message: "invalid recipe manifest",
				Errors: []*serrors.Assert{
					{
						Type:    &yaml.NodeValidationResultError{},
						Message: "manala field don't support additional properties",
						Arguments: []any{
							"property", "foo",
							"line", 2,
							"column", 14,
						},
						Details: heredoc.Doc(`
							   1 | manala:
							>  2 |   description: description
							                    ^
							   3 |   foo: bar
						`),
					},
				},
			},
		},
		// Config - Description
		{
			test: "ConfigDescriptionAbsent",
			expected: &serrors.Assert{
				Type:    &validation.Error{},
				Message: "invalid recipe manifest",
				Errors: []*serrors.Assert{
					{
						Type:    &yaml.NodeValidationResultError{},
						Message: "missing manala description field",
						Arguments: []any{
							"property", "description",
							"line", 2,
							"column", 11,
						},
						Details: heredoc.Doc(`
							   1 | manala:
							>  2 |   template: template
							                 ^
						`),
					},
				},
			},
		},
		{
			test: "ConfigDescriptionNotString",
			expected: &serrors.Assert{
				Type:    &validation.Error{},
				Message: "invalid recipe manifest",
				Errors: []*serrors.Assert{
					{
						Type:    &yaml.NodeValidationResultError{},
						Message: "manala description field must be a string",
						Arguments: []any{
							"expected", "string",
							"given", "array",
							"line", 2,
							"column", 16,
						},
						Details: heredoc.Doc(`
							   1 | manala:
							>  2 |   description: []
							                      ^
							   3 |   template: template
						`),
					},
				},
			},
		},
		{
			test: "ConfigDescriptionEmpty",
			expected: &serrors.Assert{
				Type:    &validation.Error{},
				Message: "invalid recipe manifest",
				Errors: []*serrors.Assert{
					{
						Type:    &yaml.NodeValidationResultError{},
						Message: "empty manala description field",
						Arguments: []any{
							"line", 2,
							"column", 16,
						},
						Details: heredoc.Doc(`
							   1 | manala:
							>  2 |   description: ""
							                      ^
							   3 |   template: template
						`),
					},
				},
			},
		},
		{
			test: "ConfigDescriptionTooLong",
			expected: &serrors.Assert{
				Type:    &validation.Error{},
				Message: "invalid recipe manifest",
				Errors: []*serrors.Assert{
					{
						Type:    &yaml.NodeValidationResultError{},
						Message: "too long manala description field",
						Arguments: []any{
							"line", 2,
							"column", 16,
						},
						Details: heredoc.Doc(`
							   1 | manala:
							>  2 |   description: Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
							                      ^
							   3 |   template: template
						`),
					},
				},
			},
		},
		// Config - Template
		{
			test: "ConfigTemplateNotString",
			expected: &serrors.Assert{
				Type:    &validation.Error{},
				Message: "invalid recipe manifest",
				Errors: []*serrors.Assert{
					{
						Type:    &yaml.NodeValidationResultError{},
						Message: "manala template field must be a string",
						Arguments: []any{
							"expected", "string",
							"given", "array",
							"line", 3,
							"column", 13,
						},
						Details: heredoc.Doc(`
							   1 | manala:
							   2 |   description: description
							>  3 |   template: []
							                   ^
						`),
					},
				},
			},
		},
		{
			test: "ConfigTemplateEmpty",
			expected: &serrors.Assert{
				Type:    &validation.Error{},
				Message: "invalid recipe manifest",
				Errors: []*serrors.Assert{
					{
						Type:    &yaml.NodeValidationResultError{},
						Message: "empty manala template field",
						Arguments: []any{
							"line", 3,
							"column", 13,
						},
						Details: heredoc.Doc(`
							   1 | manala:
							   2 |   description: description
							>  3 |   template: ""
							                   ^
						`),
					},
				},
			},
		},
		{
			test: "ConfigTemplateTooLong",
			expected: &serrors.Assert{
				Type:    &validation.Error{},
				Message: "invalid recipe manifest",
				Errors: []*serrors.Assert{
					{
						Type:    &yaml.NodeValidationResultError{},
						Message: "too long manala template field",
						Arguments: []any{
							"line", 3,
							"column", 13,
						},
						Details: heredoc.Doc(`
							   1 | manala:
							   2 |   description: description
							>  3 |   template: Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
							                   ^
						`),
					},
				},
			},
		},
		// Config - Sync
		{
			test: "ConfigSyncNotArray",
			expected: &serrors.Assert{
				Type:    &validation.Error{},
				Message: "invalid recipe manifest",
				Errors: []*serrors.Assert{
					{
						Type:    &yaml.NodeValidationResultError{},
						Message: "manala sync field must be a sequence",
						Arguments: []any{
							"expected", "array",
							"given", "string",
							"line", 3,
							"column", 9,
						},
						Details: heredoc.Doc(`
							   1 | manala:
							   2 |   description: description
							>  3 |   sync: foo
							               ^
						`),
					},
				},
			},
		},
		// Config - Sync Item
		{
			test: "ConfigSyncItemNotString",
			expected: &serrors.Assert{
				Type:    &validation.Error{},
				Message: "invalid recipe manifest",
				Errors: []*serrors.Assert{
					{
						Type:    &yaml.NodeValidationResultError{},
						Message: "manala sync sequence entries must be strings",
						Arguments: []any{
							"expected", "string",
							"given", "array",
							"line", 4,
							"column", 7,
						},
						Details: heredoc.Doc(`
							   1 | manala:
							   2 |   description: description
							   3 |   sync:
							>  4 |     - []
							             ^
						`),
					},
				},
			},
		},
		{
			test: "ConfigSyncItemEmpty",
			expected: &serrors.Assert{
				Type:    &validation.Error{},
				Message: "invalid recipe manifest",
				Errors: []*serrors.Assert{
					{
						Type:    &yaml.NodeValidationResultError{},
						Message: "empty manala sync sequence entry",
						Arguments: []any{
							"line", 4,
							"column", 7,
						},
						Details: heredoc.Doc(`
							   1 | manala:
							   2 |   description: description
							   3 |   sync:
							>  4 |     - ""
							             ^
						`),
					},
				},
			},
		},
		{
			test: "ConfigSyncItemTooLong",
			expected: &serrors.Assert{
				Type:    &validation.Error{},
				Message: "invalid recipe manifest",
				Errors: []*serrors.Assert{
					{
						Type:    &yaml.NodeValidationResultError{},
						Message: "too long manala sync sequence entry",
						Arguments: []any{
							"line", 4,
							"column", 7,
						},
						Details: heredoc.Doc(`
							   1 | manala:
							   2 |   description: description
							   3 |   sync:
							>  4 |     - Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
							             ^
						`),
					},
				},
			},
		},
		// Schema
		{
			test: "SchemaMisplacedTag",
			expected: &serrors.Assert{
				Type:    &serrors.WrapError{},
				Message: "unable to infer recipe manifest schema",
				Error: &serrors.Assert{
					Type:    &yaml.NodeError{},
					Message: "misplaced schema tag",
					Arguments: []any{
						"line", 4,
						"column", 9,
					},
					Details: heredoc.Doc(`
						   1 | manala:
						   2 |   description: description
						   3 |
						>  4 | foo: ~  # @schema {"type": "string", "minLength": 1}
						               ^
					`),
				},
			},
		},
		{
			test: "SchemaInvalidJson",
			expected: &serrors.Assert{
				Type:    &serrors.WrapError{},
				Message: "unable to infer recipe manifest schema",
				Error: &serrors.Assert{
					Type:    &yaml.NodeError{},
					Message: "invalid character 'o' in literal false (expecting 'a')",
					Arguments: []any{
						"line", 4,
						"column", 1,
					},
					Details: heredoc.Doc(`
						   1 | manala:
						   2 |   description: description
						   3 |
						>  4 | # @schema foo
						       ^
						   5 | foo: ~
					`),
				},
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			recMan := NewManifest()

			recDir := filepath.FromSlash("testdata/ManifestSuite/TestReadFromErrors/" + test.test)

			recManFile, _ := os.Open(filepath.Join(recDir, "manifest.yaml"))
			err := recMan.ReadFrom(recManFile)

			serrors.Equal(s.Assert(), test.expected, err)
		})
	}
}

func (s *ManifestSuite) TestReadFrom() {
	tests := []struct {
		test                string
		expectedDescription string
		expectedTemplate    string
		expectedVars        map[string]interface{}
		expectedSync        *syncer.UnitsAssert
		expectedSchema      map[string]interface{}
	}{
		{
			test:                "All",
			expectedDescription: "description",
			expectedTemplate:    "template",
			expectedVars: map[string]interface{}{
				"string":      "string",
				"string_null": nil,
				"sequence": []interface{}{
					"first",
				},
				"sequence_string_empty": []interface{}{},
				"boolean":               true,
				"integer":               uint64(123),
				"float":                 1.2,
				"map": map[string]interface{}{
					"string": "string",
					"map": map[string]interface{}{
						"string": "string",
					},
				},
				"map_empty": map[string]interface{}{},
				"map_single": map[string]interface{}{
					"first": "foo",
				},
				"map_multiple": map[string]interface{}{
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
			expectedSchema: map[string]interface{}{
				"type":                 "object",
				"additionalProperties": false,
				"properties": map[string]interface{}{
					"string": map[string]interface{}{
						"type": "string",
					},
					"string_null": map[string]interface{}{
						"type": "string",
					},
					"sequence": map[string]interface{}{
						"type": "array",
					},
					"sequence_string_empty": map[string]interface{}{
						"type": "array",
						"items": map[string]interface{}{
							"type": "string",
						},
					},
					"boolean": map[string]interface{}{
						"type": "boolean",
					},
					"integer": map[string]interface{}{
						"type": "integer",
					},
					"float": map[string]interface{}{
						"type": "number",
					},
					"map": map[string]interface{}{
						"type":                 "object",
						"additionalProperties": false,
						"properties": map[string]interface{}{
							"string": map[string]interface{}{
								"type": "string",
							},
							"map": map[string]interface{}{
								"type":                 "object",
								"additionalProperties": false,
								"properties": map[string]interface{}{
									"string": map[string]interface{}{
										"type": "string",
									},
								},
							},
						},
					},
					"map_empty": map[string]interface{}{
						"type": "object",
					},
					"map_single": map[string]interface{}{
						"type":                 "object",
						"additionalProperties": false,
						"properties": map[string]interface{}{
							"first": map[string]interface{}{
								"type":      "string",
								"minLength": float64(1),
							},
						},
					},
					"map_multiple": map[string]interface{}{
						"type":                 "object",
						"additionalProperties": false,
						"properties": map[string]interface{}{
							"first": map[string]interface{}{
								"type":      "string",
								"minLength": float64(1),
							},
							"second": map[string]interface{}{
								"type":      "string",
								"minLength": float64(1),
							},
						},
					},
					"enum": map[string]interface{}{
						"enum": []interface{}{
							nil,
							true,
							false,
							"string",
							float64(12),
							2.3,
							3.0,
							"3.0",
						},
					},
					"underscore_key": map[string]interface{}{
						"type": "string",
					},
					"hyphen-key": map[string]interface{}{
						"type": "string",
					},
					"dot.key": map[string]interface{}{
						"type": "string",
					},
				},
			},
		},
		{
			test:                "ConfigTemplateAbsent",
			expectedDescription: "description",
			expectedTemplate:    "",
			expectedVars: map[string]interface{}{
				"foo": "bar",
			},
			expectedSync: &syncer.UnitsAssert{},
			expectedSchema: map[string]interface{}{
				"type":                 "object",
				"additionalProperties": false,
				"properties": map[string]interface{}{
					"foo": map[string]interface{}{
						"type": "string",
					},
				},
			},
		},
		{
			test:                "VarsAbsent",
			expectedDescription: "description",
			expectedTemplate:    "template",
			expectedVars:        map[string]interface{}{},
			expectedSync:        &syncer.UnitsAssert{},
			expectedSchema:      map[string]interface{}{},
		},
		{
			test:                "VarsKeys",
			expectedDescription: "description",
			expectedTemplate:    "template",
			expectedVars: map[string]interface{}{
				"underscore_key": "ok",
				"hyphen-key":     "ok",
				"dot.key":        "ok",
			},
			expectedSync: &syncer.UnitsAssert{},
			expectedSchema: map[string]interface{}{
				"type":                 "object",
				"additionalProperties": false,
				"properties": map[string]interface{}{
					"underscore_key": map[string]interface{}{
						"type": "string",
					},
					"hyphen-key": map[string]interface{}{
						"type": "string",
					},
					"dot.key": map[string]interface{}{
						"type": "string",
					},
				},
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			recMan := NewManifest()

			recDir := filepath.FromSlash("testdata/ManifestSuite/TestReadFrom/" + test.test)

			recManFile, _ := os.Open(filepath.Join(recDir, "manifest.yaml"))
			err := recMan.ReadFrom(recManFile)

			s.NoError(err)

			s.Equal(test.expectedDescription, recMan.Description())
			s.Equal(test.expectedTemplate, recMan.Template())
			s.Equal(test.expectedVars, recMan.Vars())
			syncer.EqualUnits(s.Assert(), test.expectedSync, recMan.Sync())
			s.Equal(test.expectedSchema, recMan.Schema())
		})
	}
}

func (s *ManifestSuite) TestInitVars() {
	recMan := NewManifest()

	recDir := filepath.FromSlash("testdata/ManifestSuite/TestInitVars")

	recManFile, _ := os.Open(filepath.Join(recDir, "manifest.yaml"))
	_ = recMan.ReadFrom(recManFile)

	vars, err := recMan.InitVars(func(options []interfaces.RecipeOption) error {
		s.Len(options, 11)

		s.Equal("String", options[0].Label())
		s.Equal(map[string]interface{}{
			"type": "string",
		}, options[0].Schema())

		s.Equal("String null", options[1].Label())
		s.Equal(map[string]interface{}{
			"type": "string",
		}, options[1].Schema())

		s.Equal("String float int", options[2].Label())
		s.Equal(map[string]interface{}{
			"enum": []interface{}{"3.0"},
		}, options[2].Schema())
		_ = options[2].Set("3.0")

		s.Equal("String asterisk", options[3].Label())
		s.Equal(map[string]interface{}{
			"enum": []interface{}{"*"},
		}, options[3].Schema())
		_ = options[3].Set("*")

		s.Equal("Map single first", options[4].Label())
		s.Equal(map[string]interface{}{
			"type":      "string",
			"minLength": float64(1),
		}, options[4].Schema())

		s.Equal("Map multiple first", options[5].Label())
		s.Equal(map[string]interface{}{
			"type":      "string",
			"minLength": float64(1),
		}, options[5].Schema())

		s.Equal("Map multiple second", options[6].Label())
		s.Equal(map[string]interface{}{
			"type":      "string",
			"minLength": float64(1),
		}, options[6].Schema())

		s.Equal("Enum null", options[7].Label())
		s.Equal(map[string]interface{}{
			"enum": []interface{}{
				nil,
				true,
				false,
				"string",
				float64(12),
				2.3,
				3.0,
				"3.0",
			},
		}, options[7].Schema())

		s.Equal("Underscore key", options[8].Label())
		s.Equal(map[string]interface{}{
			"type": "string",
		}, options[8].Schema())

		s.Equal("Hyphen key", options[9].Label())
		s.Equal(map[string]interface{}{
			"type": "string",
		}, options[9].Schema())

		s.Equal("Dot key", options[10].Label())
		s.Equal(map[string]interface{}{
			"type": "string",
		}, options[10].Schema())

		return nil
	})

	s.NoError(err)

	s.Equal(map[string]interface{}{
		"string":                 "string",
		"string_null":            nil,
		"string_float_int":       "3.0",
		"string_float_int_value": "3.0",
		"string_asterisk":        "*",
		"string_asterisk_value":  "*",
		"sequence": []interface{}{
			"first",
		},
		"sequence_string_empty": []interface{}{},
		"boolean":               true,
		"integer":               uint64(123),
		"float":                 1.2,
		"map": map[string]interface{}{
			"string": "string",
			"map": map[string]interface{}{
				"string": "string",
			},
		},
		"map_empty": map[string]interface{}{},
		"map_single": map[string]interface{}{
			"first": "foo",
		},
		"map_multiple": map[string]interface{}{
			"first":  "foo",
			"second": "foo",
		},
		"enum":           nil,
		"underscore_key": "ok",
		"hyphen-key":     "ok",
		"dot.key":        "ok",
	}, vars)
}
