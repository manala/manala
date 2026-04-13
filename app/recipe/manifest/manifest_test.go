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

type ManifestSuite struct{ suite.Suite }

func TestManifestSuite(t *testing.T) {
	suite.Run(t, new(ManifestSuite))
}

func (s *ManifestSuite) Test() {
	m := manifest.New()

	s.Empty(m.Description())
	s.Empty(m.Icon())
	s.Empty(m.Template())
	s.Empty(m.Partials())
	s.Equal(map[string]any{}, m.Vars())
	s.Equal([]sync.UnitInterface{}, m.Sync())
	s.Equal(schema.Schema{}, m.Schema())
}

func (s *ManifestSuite) TestUnmarshalYAMLErrors() {
	tests := []struct {
		test     string
		expected errors.Assertion
	}{
		{
			test: "Empty",
			expected: &parsing.ErrorAssertion{
				Err: &serrors.Assertion{
					Message: "empty yaml content",
				},
			},
		},
		{
			test: "Invalid",
			expected: &parsing.ErrorAssertion{
				Line:   1,
				Column: 1,
				Err: &serrors.Assertion{
					Message: "'@' is a reserved character",
				},
			},
		},
		{
			test: "IrregularType",
			expected: &parsing.ErrorAssertion{
				Line:   1,
				Column: 6,
				Err: &serrors.Assertion{
					Message: "irregular type",
				},
			},
		},
		{
			test: "IrregularMapKey",
			expected: &parsing.ErrorAssertion{
				Line:   1,
				Column: 2,
				Err: &serrors.Assertion{
					Message: "irregular map key",
				},
			},
		},
		{
			test: "NotMap",
			expected: &parsing.ErrorAssertion{
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
			expected: &parsing.ErrorAssertion{
				Err: &serrors.Assertion{
					Message: "missing manala property",
				},
			},
		},
		{
			test: "ConfigNotMap",
			expected: &parsing.ErrorAssertion{
				Line:   1,
				Column: 9,
				Err: &serrors.Assertion{
					Message: "string was used where mapping is expected",
				},
			},
		},
		{
			test: "ConfigEmpty",
			expected: &parsing.ErrorAssertion{
				Line:   1,
				Column: 1,
				Err: &serrors.Assertion{
					Message: "missing manala description property",
				},
			},
		},
		{
			test: "ConfigAdditionalProperties",
			expected: &parsing.ErrorAssertion{
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
			expected: &parsing.ErrorAssertion{
				Line:   1,
				Column: 7,
				Err: &serrors.Assertion{
					Message: "missing manala description property",
				},
			},
		},
		{
			test: "ConfigDescriptionNotString",
			expected: &parsing.ErrorAssertion{
				Line:   2,
				Column: 16,
				Err: &serrors.Assertion{
					Message: "field must be a string",
				},
			},
		},
		{
			test: "ConfigDescriptionEmpty",
			expected: &parsing.ErrorAssertion{
				Line:   2,
				Column: 16,
				Err: &serrors.Assertion{
					Message: "missing manala description property",
				},
			},
		},
		{
			test: "ConfigDescriptionTooLong",
			expected: &parsing.ErrorAssertion{
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
			expected: &parsing.ErrorAssertion{
				Line:   3,
				Column: 9,
				Err: &serrors.Assertion{
					Message: "field must be a string",
				},
			},
		},
		{
			test: "ConfigIconTooLong",
			expected: &parsing.ErrorAssertion{
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
			expected: &parsing.ErrorAssertion{
				Line:   3,
				Column: 13,
				Err: &serrors.Assertion{
					Message: "field must be a string",
				},
			},
		},
		{
			test: "ConfigTemplateTooLong",
			expected: &parsing.ErrorAssertion{
				Line:   3,
				Column: 13,
				Err: &serrors.Assertion{
					Message: "too long manala template field (max=100)",
				},
			},
		},
		// Config - Partials
		{
			test: "ConfigPartialsNotArray",
			expected: &parsing.ErrorAssertion{
				Line:   3,
				Column: 13,
				Err: &serrors.Assertion{
					Message: "string was used where sequence is expected",
				},
			},
		},
		// Config - Partials Item
		{
			test: "ConfigPartialsItemNotString",
			expected: &parsing.ErrorAssertion{
				Line:   4,
				Column: 7,
				Err: &serrors.Assertion{
					Message: "field must be a string",
				},
			},
		},
		{
			test: "ConfigPartialsItemEmpty",
			expected: &parsing.ErrorAssertion{
				Line:   3,
				Column: 11,
				Err: &serrors.Assertion{
					Message: "empty partials entry",
				},
			},
		},
		{
			test: "ConfigPartialsItemTooLong",
			expected: &parsing.ErrorAssertion{
				Line:   3,
				Column: 11,
				Err: &serrors.Assertion{
					Message: "too long partials entry (max=100)",
				},
			},
		},
		// Config - Sync
		{
			test: "ConfigSyncNotArray",
			expected: &parsing.ErrorAssertion{
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
			expected: &parsing.ErrorAssertion{
				Line:   4,
				Column: 5,
				Err: &serrors.Assertion{
					Message: "sync entry must be a string",
				},
			},
		},
		{
			test: "ConfigSyncItemEmpty",
			expected: &parsing.ErrorAssertion{
				Line:   4,
				Column: 5,
				Err: &serrors.Assertion{
					Message: "empty sync entry",
				},
			},
		},
		{
			test: "ConfigSyncItemTooLong",
			expected: &parsing.ErrorAssertion{
				Line:   4,
				Column: 5,
				Err: &serrors.Assertion{
					Message: "too long sync entry (max=256)",
				},
			},
		},
		// Schema
		{
			test: "SchemaInvalidJson",
			expected: &parsing.FlattenErrorAssertion{
				Line:   4,
				Column: 12,
				Err: &serrors.Assertion{
					Message: "invalid character 'o' in literal false (expecting 'a')",
				},
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			m := manifest.New()

			dir := filepath.FromSlash("testdata/ManifestSuite/TestUnmarshalYAMLErrors")

			reader, _ := os.Open(filepath.Join(dir, test.test+".yaml"))
			content, _ := io.ReadAll(reader)

			err := m.UnmarshalYAML(content)

			errors.Equal(s.T(), test.expected, err)
		})
	}
}

func (s *ManifestSuite) TestUnmarshalYAML() {
	tests := []struct {
		test                string
		expectedDescription string
		expectedIcon        string
		expectedTemplate    string
		expectedPartials    []string
		expectedVars        map[string]any
		expectedSync        *sync.UnitsAssertion
		expectedSchema      schema.Schema
	}{
		{
			test:                "All",
			expectedDescription: "description",
			expectedIcon:        "icon",
			expectedTemplate:    "template",
			expectedPartials:    []string{"partial.tmpl", "dir/partial.tmpl"},
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
						"type":                 "object",
						"additionalProperties": false,
						"properties":           map[string]any{},
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
			expectedSchema: schema.Schema{
				"type":                 "object",
				"additionalProperties": false,
				"properties":           map[string]any{},
			},
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

			dir := filepath.FromSlash("testdata/ManifestSuite/TestUnmarshalYAML")

			reader, _ := os.Open(filepath.Join(dir, test.test+".yaml"))
			content, _ := io.ReadAll(reader)

			err := m.UnmarshalYAML(content)

			s.Require().NoError(err)

			s.Equal(test.expectedDescription, m.Description())
			s.Equal(test.expectedIcon, m.Icon())
			s.Equal(test.expectedTemplate, m.Template())
			s.Equal(test.expectedPartials, m.Partials())
			s.Equal(test.expectedVars, m.Vars())
			sync.EqualUnits(s.T(), test.expectedSync, m.Sync())
			s.Equal(test.expectedSchema, m.Schema())
		})
	}
}

func (s *ManifestSuite) TestOptions() {
	m := manifest.New()

	dir := filepath.FromSlash("testdata/ManifestSuite/TestOptions")

	reader, _ := os.Open(filepath.Join(dir, "manifest.yaml"))
	content, _ := io.ReadAll(reader)

	err := m.UnmarshalYAML(content)

	opts := m.Options()

	s.Require().NoError(err)

	option.Equals(s.T(), option.Assertions{
		{
			Type:      &option.String{},
			Label:     "String",
			Name:      "string",
			MaxLength: 0,
		},
		{
			Type:      &option.String{},
			Label:     "String null",
			Name:      "string-null",
			MaxLength: 0,
		},
		{
			Type:      &option.String{},
			Label:     "String max length",
			Name:      "string-max-length",
			MaxLength: 123,
		},
		{
			Type:   &option.Enum{},
			Label:  "String float int",
			Name:   "string-float-int",
			Values: []any{"3.0"},
		},
		{
			Type:   &option.Enum{},
			Label:  "String asterisk",
			Name:   "string-asterisk",
			Values: []any{"*"},
		},
		{
			Type:      &option.String{},
			Label:     "Map single first",
			Name:      "map-single-first",
			MaxLength: 0,
		},
		{
			Type:      &option.String{},
			Label:     "Map multiple first",
			Name:      "map-multiple-first",
			MaxLength: 0,
		},
		{
			Type:      &option.String{},
			Label:     "Map multiple second",
			Name:      "map-multiple-second",
			MaxLength: 0,
		},
		{
			Type:  &option.Enum{},
			Label: "Enum null",
			Name:  "enum-null",
			Values: []any{
				nil,
				true,
				false,
				"string",
				int64(12),
				2.3,
				3.0,
				"3.0",
			},
		},
		{
			Type:      &option.String{},
			Label:     "Underscore key",
			Name:      "underscore-key",
			MaxLength: 0,
		},
		{
			Type:      &option.String{},
			Label:     "Hyphen key",
			Name:      "hyphen-key",
			MaxLength: 0,
		},
		{
			Type:      &option.String{},
			Label:     "Dot key",
			Name:      "dot-key",
			MaxLength: 0,
		},
		{
			Type:      &option.String{},
			Label:     "Custom name",
			Name:      "foo-bar",
			MaxLength: 0,
		},
	}, opts)
}
