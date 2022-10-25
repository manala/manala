package recipe

import (
	"github.com/stretchr/testify/suite"
	"manala/core"
	internalReport "manala/internal/report"
	internalSyncer "manala/internal/syncer"
	internalTesting "manala/internal/testing"
	"os"
	"path/filepath"
	"testing"
)

type ManifestSuite struct{ suite.Suite }

func TestManifestSuite(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	suite.Run(t, new(ManifestSuite))
}

func (s *ManifestSuite) Test() {
	manifest := NewManifest("dir")

	s.Equal(filepath.Join("dir", ".manala.yaml"), manifest.Path())
	s.Equal("", manifest.Description())
	s.Equal("", manifest.Template())
	s.Equal(map[string]interface{}{}, manifest.Vars())
	s.Equal([]internalSyncer.UnitInterface{}, manifest.Sync())
	s.Equal(map[string]interface{}{}, manifest.Schema())
}

func (s *ManifestSuite) TestReadFromErrors() {
	tests := []struct {
		name   string
		report *internalReport.Assert
	}{
		{
			name: "Empty",
			report: &internalReport.Assert{
				Message: "irregular recipe manifest",
				Err:     "empty yaml file",
			},
		},
		{
			name: "Invalid",
			report: &internalReport.Assert{
				Message: "irregular recipe manifest",
				Fields: map[string]interface{}{
					"column": 1,
					"line":   1,
				},
				Err: "unexpected mapping key",
			},
		},
		{
			name: "Irregular Type",
			report: &internalReport.Assert{
				Message: "irregular recipe manifest",
				Err:     "irregular type",
				Fields: map[string]interface{}{
					"line":   1,
					"column": 6,
				},
			},
		},
		{
			name: "Irregular Map Key",
			report: &internalReport.Assert{
				Message: "irregular recipe manifest",
				Err:     "irregular map key",
				Fields: map[string]interface{}{
					"line":   1,
					"column": 2,
				},
			},
		},
		{
			name: "Not Map",
			report: &internalReport.Assert{
				Err: "invalid recipe manifest",
				Reports: []internalReport.Assert{
					{
						Message: "yaml document must be a map",
						Fields: map[string]interface{}{
							"expected": "object",
							"given":    "string",
						},
					},
				},
			},
		},
		// Config
		{
			name: "Config Absent",
			report: &internalReport.Assert{
				Err: "invalid recipe manifest",
				Reports: []internalReport.Assert{
					{
						Message: "missing manala field",
						Fields: map[string]interface{}{
							"property": "manala",
						},
					},
				},
			},
		},
		{
			name: "Config Not Map",
			report: &internalReport.Assert{
				Err: "invalid recipe manifest",
				Reports: []internalReport.Assert{
					{
						Message: "manala field must be a map",
						Fields: map[string]interface{}{
							"line":     1,
							"column":   9,
							"expected": "object",
							"given":    "string",
						},
					},
				},
			},
		},
		{
			name: "Config Empty",
			report: &internalReport.Assert{
				Err: "invalid recipe manifest",
				Reports: []internalReport.Assert{
					{
						Message: "missing manala description field",
						Fields: map[string]interface{}{
							"line":     1,
							"column":   9,
							"property": "description",
						},
					},
				},
			},
		},
		{
			name: "Config Additional Properties",
			report: &internalReport.Assert{
				Err: "invalid recipe manifest",
				Reports: []internalReport.Assert{
					{
						Message: "manala field don't support additional properties",
						Fields: map[string]interface{}{
							"line":     2,
							"column":   14,
							"property": "foo",
						},
					},
				},
			},
		},
		// Config - Description
		{
			name: "Config Description Absent",
			report: &internalReport.Assert{
				Err: "invalid recipe manifest",
				Reports: []internalReport.Assert{
					{
						Message: "missing manala description field",
						Fields: map[string]interface{}{
							"line":     2,
							"column":   11,
							"property": "description",
						},
					},
				},
			},
		},
		{
			name: "Config Description Not String",
			report: &internalReport.Assert{
				Err: "invalid recipe manifest",
				Reports: []internalReport.Assert{
					{
						Message: "manala description field must be a string",
						Fields: map[string]interface{}{
							"line":     2,
							"column":   16,
							"expected": "string",
							"given":    "array",
						},
					},
				},
			},
		},
		{
			name: "Config Description Empty",
			report: &internalReport.Assert{
				Err: "invalid recipe manifest",
				Reports: []internalReport.Assert{
					{
						Message: "empty manala description field",
						Fields: map[string]interface{}{
							"line":   2,
							"column": 16,
						},
					},
				},
			},
		},
		{
			name: "Config Description Too Long",
			report: &internalReport.Assert{
				Err: "invalid recipe manifest",
				Reports: []internalReport.Assert{
					{
						Message: "too long manala description field",
						Fields: map[string]interface{}{
							"line":   2,
							"column": 16,
						},
					},
				},
			},
		},
		// Config - Template
		{
			name: "Config Template Not String",
			report: &internalReport.Assert{
				Err: "invalid recipe manifest",
				Reports: []internalReport.Assert{
					{
						Message: "manala template field must be a string",
						Fields: map[string]interface{}{
							"line":     3,
							"column":   13,
							"expected": "string",
							"given":    "array",
						},
					},
				},
			},
		},
		{
			name: "Config Template Empty",
			report: &internalReport.Assert{
				Err: "invalid recipe manifest",
				Reports: []internalReport.Assert{
					{
						Message: "empty manala template field",
						Fields: map[string]interface{}{
							"line":   3,
							"column": 13,
						},
					},
				},
			},
		},
		{
			name: "Config Template Too Long",
			report: &internalReport.Assert{
				Err: "invalid recipe manifest",
				Reports: []internalReport.Assert{
					{
						Message: "too long manala template field",
						Fields: map[string]interface{}{
							"line":   3,
							"column": 13,
						},
					},
				},
			},
		},
		// Config - Sync
		{
			name: "Config Sync Not Array",
			report: &internalReport.Assert{
				Err: "invalid recipe manifest",
				Reports: []internalReport.Assert{
					{
						Message: "manala sync field must be a sequence",
						Fields: map[string]interface{}{
							"line":     3,
							"column":   9,
							"expected": "array",
							"given":    "string",
						},
					},
				},
			},
		},
		// Config - Sync Item
		{
			name: "Config Sync Item Not String",
			report: &internalReport.Assert{
				Err: "invalid recipe manifest",
				Reports: []internalReport.Assert{
					{
						Message: "manala sync sequence entries must be strings",
						Fields: map[string]interface{}{
							"line":     4,
							"column":   7,
							"expected": "string",
							"given":    "array",
						},
					},
				},
			},
		},
		{
			name: "Config Sync Item Empty",
			report: &internalReport.Assert{
				Err: "invalid recipe manifest",
				Reports: []internalReport.Assert{
					{
						Message: "empty manala sync sequence entry",
						Fields: map[string]interface{}{
							"line":   4,
							"column": 7,
						},
					},
				},
			},
		},
		{
			name: "Config Sync Item Too Long",
			report: &internalReport.Assert{
				Err: "invalid recipe manifest",
				Reports: []internalReport.Assert{
					{
						Message: "too long manala sync sequence entry",
						Fields: map[string]interface{}{
							"line":   4,
							"column": 7,
						},
					},
				},
			},
		},
		// Schema
		{
			name: "Schema Misplaced Tag",
			report: &internalReport.Assert{
				Message: "unable to infer recipe manifest schema",
				Err:     "misplaced schema tag",
				Fields: map[string]interface{}{
					"line":   4,
					"column": 6,
				},
			},
		},
		{
			name: "Schema Invalid Json",
			report: &internalReport.Assert{
				Message: "unable to unmarshal json",
				Err:     "invalid character 'o' in literal false (expecting 'a')",
				Fields: map[string]interface{}{
					"line":   5,
					"column": 4,
				},
			},
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			manifest := NewManifest("")

			manifestFile, _ := os.Open(internalTesting.DataPath(s, "manifest.yaml"))
			err := manifest.ReadFrom(manifestFile)

			s.Error(err)

			report := internalReport.NewErrorReport(err)

			test.report.Equal(&s.Suite, report)
		})
	}
}

func (s *ManifestSuite) TestReadFrom() {
	tests := []struct {
		name        string
		description string
		template    string
		vars        map[string]interface{}
		sync        *syncAssert
		schema      map[string]interface{}
	}{
		{
			name:        "All",
			description: "description",
			template:    "template",
			vars: map[string]interface{}{
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
			sync: &syncAssert{
				{Source: "file", Destination: "file"},
				{Source: "dir/file", Destination: "dir/file"},
				{Source: "file", Destination: "dir/file"},
				{Source: "dir/file", Destination: "file"},
				{Source: "src_file", Destination: "dst_file"},
				{Source: "src_dir/file", Destination: "dst_dir/file"},
			},
			schema: map[string]interface{}{
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
			name:        "Config Template Absent",
			description: "description",
			template:    "",
			vars: map[string]interface{}{
				"foo": "bar",
			},
			sync: &syncAssert{},
			schema: map[string]interface{}{
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
			name:        "Vars Absent",
			description: "description",
			template:    "template",
			vars:        map[string]interface{}{},
			sync:        &syncAssert{},
			schema:      map[string]interface{}{},
		},
		{
			name:        "Vars Keys",
			description: "description",
			template:    "template",
			vars: map[string]interface{}{
				"underscore_key": "ok",
				"hyphen-key":     "ok",
				"dot.key":        "ok",
			},
			sync: &syncAssert{},
			schema: map[string]interface{}{
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
		s.Run(test.name, func() {
			manifest := NewManifest("")

			manifestFile, _ := os.Open(internalTesting.DataPath(s, "manifest.yaml"))
			err := manifest.ReadFrom(manifestFile)

			s.NoError(err)
			s.Equal(test.description, manifest.Description())
			s.Equal(test.template, manifest.Template())
			s.Equal(test.vars, manifest.Vars())
			test.sync.Equal(&s.Suite, manifest.Sync())
			s.Equal(test.schema, manifest.Schema())
		})
	}
}

func (s *ManifestSuite) TestInitVars() {
	manifest := NewManifest("")

	manifestFile, _ := os.Open(internalTesting.DataPath(s, "manifest.yaml"))
	_ = manifest.ReadFrom(manifestFile)

	_, err := manifest.InitVars(func(options []core.RecipeOption) error {
		s.Len(options, 9)

		s.Equal("String", options[0].Label())
		s.Equal(map[string]interface{}{
			"type": "string",
		}, options[0].Schema())

		s.Equal("String null", options[1].Label())
		s.Equal(map[string]interface{}{
			"type": "string",
		}, options[1].Schema())

		s.Equal("Map single first", options[2].Label())
		s.Equal(map[string]interface{}{
			"type":      "string",
			"minLength": float64(1),
		}, options[2].Schema())

		s.Equal("Map multiple first", options[3].Label())
		s.Equal(map[string]interface{}{
			"type":      "string",
			"minLength": float64(1),
		}, options[3].Schema())

		s.Equal("Map multiple second", options[4].Label())
		s.Equal(map[string]interface{}{
			"type":      "string",
			"minLength": float64(1),
		}, options[4].Schema())

		s.Equal("Enum null", options[5].Label())
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
		}, options[5].Schema())

		s.Equal("Underscore key", options[6].Label())
		s.Equal(map[string]interface{}{
			"type": "string",
		}, options[6].Schema())

		s.Equal("Hyphen key", options[7].Label())
		s.Equal(map[string]interface{}{
			"type": "string",
		}, options[7].Schema())

		s.Equal("Dot key", options[8].Label())
		s.Equal(map[string]interface{}{
			"type": "string",
		}, options[8].Schema())

		return nil
	})

	s.NoError(err)
}
