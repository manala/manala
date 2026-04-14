package manifest_test

import (
	"encoding/json"
	"testing"

	"github.com/manala/manala/app/recipe/manifest"
	"github.com/manala/manala/app/recipe/option"
	"github.com/manala/manala/internal/parsing"
	"github.com/manala/manala/internal/schema"
	"github.com/manala/manala/internal/serrors"
	"github.com/manala/manala/internal/sync"
	"github.com/manala/manala/internal/testing/errors"
	"github.com/manala/manala/internal/testing/heredoc"

	"github.com/stretchr/testify/suite"
)

type ManifestSuite struct{ suite.Suite }

func TestManifestSuite(t *testing.T) {
	suite.Run(t, new(ManifestSuite))
}

func (s *ManifestSuite) TestNew() {
	m := manifest.New()

	s.Empty(m.Description)
	s.Empty(m.Icon)
	s.Empty(m.Template)
	s.Empty(m.Partials)
	sync.EqualUnits(s.T(), &sync.UnitsAssertion{}, m.Sync)
	s.Equal(map[string]any{}, m.Vars())
	s.Equal(schema.Schema{}, m.Schema())
	option.Equals(s.T(), option.Assertions{}, m.Options())
}

func (s *ManifestSuite) TestUnmarshalRequired() {
	m := manifest.New()

	err := m.Unmarshal([]byte(heredoc.Doc(`
		manala:
		  description: description
	`)))
	s.Require().NoError(err)

	s.Equal("description", m.Description)
	s.Empty(m.Icon)
	s.Empty(m.Template)
	s.Empty(m.Partials)
	sync.EqualUnits(s.T(), &sync.UnitsAssertion{}, m.Sync)
	s.Equal(map[string]any{}, m.Vars())
	s.Equal(schema.Schema{
		"type":                 "object",
		"additionalProperties": false,
		"properties":           map[string]any{},
	}, m.Schema())
	option.Equals(s.T(), option.Assertions{}, m.Options())
}

func (s *ManifestSuite) TestUnmarshal() {
	m := manifest.New()

	err := m.Unmarshal([]byte(heredoc.Doc(`
		manala:
		  description: description
		  icon: icon
		  template: template
		  partials:
		    - partial.tmpl
		    - dir/partial.tmpl
		  sync:
		    - file
		    - dir/file
		    - file dir/file
		    - dir/file file
		    - src_file dst_file
		    - src_dir/file dst_dir/file
		# @option {"label": "String"}
		string: string
		# @option {"label": "String null"}
		# @schema {"type": "string"}
		string_null: ~
		# @option {"label": "String max length"}
		# @schema {"maxLength": 123}
		string_max_length: string
		# @option {"label": "String float int"}
		# @schema {"enum": ["3.0"]}
		string_float_int: ~
		string_float_int_value: "3.0"
		# @option {"label": "String asterisk"}
		# @schema {"enum": ["*"]}
		string_asterisk: ~
		string_asterisk_value: "*"
		sequence:
		  - first
		# @schema {"items": {"type": "string"}}
		sequence_string_empty: []
		boolean: true
		integer: 123
		float: 1.2
		map:
		  string: string
		  map:
		    string: string
		map_empty: {}
		map_single:
		  # @option {"label": "Map single first"}
		  # @schema {"minLength": 1}
		  first: foo
		map_multiple:
		  # @option {"label": "Map multiple first"}
		  # @schema {"minLength": 1}
		  first: foo
		  # @option {"label": "Map multiple second"}
		  # @schema {"minLength": 1}
		  second: foo
		# @option {"label": "Enum null"}
		# @schema {"enum": [null, true, false, "string", 12, 2.3, 3.0, "3.0"]}
		enum: ~
		# @option {"label": "Underscore key"}
		underscore_key: ok
		# @option {"label": "Hyphen key"}
		hyphen-key: ok
		# @option {"label": "Dot key"}
		dot.key: ok
		# @option {"label": "Custom name", "name": "foo-bar"}
		custom_name: ok
	`)))
	s.Require().NoError(err)

	// Description
	s.Equal("description", m.Description)

	// Icon
	s.Equal("icon", m.Icon)

	// Template
	s.Equal("template", m.Template)

	// Partials
	s.Equal([]string{
		"partial.tmpl",
		"dir/partial.tmpl",
	}, m.Partials)

	// Sync
	sync.EqualUnits(s.T(), &sync.UnitsAssertion{
		{Source: "file", Destination: "file"},
		{Source: "dir/file", Destination: "dir/file"},
		{Source: "file", Destination: "dir/file"},
		{Source: "dir/file", Destination: "file"},
		{Source: "src_file", Destination: "dst_file"},
		{Source: "src_dir/file", Destination: "dst_dir/file"},
	}, m.Sync)

	// Vars
	s.Equal(map[string]any{
		"string":                 "string",
		"string_null":            nil,
		"string_max_length":      "string",
		"string_float_int":       nil,
		"string_float_int_value": "3.0",
		"string_asterisk":        nil,
		"string_asterisk_value":  "*",
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
		"custom_name":    "ok",
	}, m.Vars())

	// Schema
	s.Equal(schema.Schema{
		"type":                 "object",
		"additionalProperties": false,
		"properties": map[string]any{
			"string": map[string]any{
				"type": "string",
			},
			"string_null": map[string]any{
				"type": "string",
			},
			"string_max_length": map[string]any{
				"type":      "string",
				"maxLength": json.Number("123"),
			},
			"string_float_int": map[string]any{
				"enum": []any{"3.0"},
			},
			"string_float_int_value": map[string]any{
				"type": "string",
			},
			"string_asterisk": map[string]any{
				"enum": []any{"*"},
			},
			"string_asterisk_value": map[string]any{
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
			"custom_name": map[string]any{
				"type": "string",
			},
		},
	}, m.Schema())

	// Options
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
	}, m.Options())
}

func (s *ManifestSuite) TestUnmarshalErrors() {
	tests := []struct {
		test     string
		content  string
		expected errors.Assertion
	}{
		{
			test:    "Empty",
			content: "",
			expected: &parsing.ErrorAssertion{
				Err: &serrors.Assertion{
					Message: "empty yaml content",
				},
			},
		},
		{
			test: "Invalid",
			content: heredoc.Doc(`
				@
			`),
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
			content: heredoc.Doc(`
				foo: .inf
			`),
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
			content: heredoc.Doc(`
				0: foo
			`),
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
			content: heredoc.Doc(`
				foo
			`),
			expected: &parsing.ErrorAssertion{
				Line:   1,
				Column: 1,
				Err: &serrors.Assertion{
					Message: "yaml document must be a map",
				},
			},
		},
		// Manala
		{
			test: "ManalaAbsent",
			content: heredoc.Doc(`
				foo: bar
			`),
			expected: &parsing.ErrorAssertion{
				Err: &serrors.Assertion{
					Message: "missing manala property",
				},
			},
		},
		{
			test: "ManalaNotMap",
			content: heredoc.Doc(`
				manala: foo
			`),
			expected: &parsing.ErrorAssertion{
				Line:   1,
				Column: 9,
				Err: &serrors.Assertion{
					Message: "string was used where mapping is expected",
				},
			},
		},
		{
			test: "ManalaEmpty",
			content: heredoc.Doc(`
				manala: {}
			`),
			expected: &parsing.ErrorAssertion{
				Line:   1,
				Column: 1,
				Err: &serrors.Assertion{
					Message: "missing manala description property",
				},
			},
		},
		{
			test: "ManalaAdditionalProperties",
			content: heredoc.Doc(`
				manala:
				  description: description
				  foo: bar
			`),
			expected: &parsing.ErrorAssertion{
				Line:   3,
				Column: 3,
				Err: &serrors.Assertion{
					Message: "unknown field \"foo\"",
				},
			},
		},
		// Description
		{
			test: "DescriptionAbsent",
			content: heredoc.Doc(`
				manala:
				  template: template
			`),
			expected: &parsing.ErrorAssertion{
				Line:   1,
				Column: 7,
				Err: &serrors.Assertion{
					Message: "missing manala description property",
				},
			},
		},
		{
			test: "DescriptionNotString",
			content: heredoc.Doc(`
				manala:
				  description: []
				  template: template
			`),
			expected: &parsing.ErrorAssertion{
				Line:   2,
				Column: 16,
				Err: &serrors.Assertion{
					Message: "field must be a string",
				},
			},
		},
		{
			test: "DescriptionEmpty",
			content: heredoc.Doc(`
				manala:
				  description: ""
				  template: template
			`),
			expected: &parsing.ErrorAssertion{
				Line:   2,
				Column: 16,
				Err: &serrors.Assertion{
					Message: "missing manala description property",
				},
			},
		},
		{
			test: "DescriptionTooLong",
			content: heredoc.Doc(`
				manala:
				  description: Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
				  template: template
			`),
			expected: &parsing.ErrorAssertion{
				Line:   2,
				Column: 16,
				Err: &serrors.Assertion{
					Message: "too long manala description field (max=256)",
				},
			},
		},
		// Icon
		{
			test: "IconNotString",
			content: heredoc.Doc(`
				manala:
				  description: description
				  icon: []
			`),
			expected: &parsing.ErrorAssertion{
				Line:   3,
				Column: 9,
				Err: &serrors.Assertion{
					Message: "field must be a string",
				},
			},
		},
		{
			test: "IconTooLong",
			content: heredoc.Doc(`
				manala:
				  description: description
				  icon: Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
			`),
			expected: &parsing.ErrorAssertion{
				Line:   3,
				Column: 9,
				Err: &serrors.Assertion{
					Message: "too long manala icon field (max=100)",
				},
			},
		},
		// Template
		{
			test: "TemplateNotString",
			content: heredoc.Doc(`
				manala:
				  description: description
				  template: []
			`),
			expected: &parsing.ErrorAssertion{
				Line:   3,
				Column: 13,
				Err: &serrors.Assertion{
					Message: "field must be a string",
				},
			},
		},
		{
			test: "TemplateTooLong",
			content: heredoc.Doc(`
				manala:
				  description: description
				  template: Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
			`),
			expected: &parsing.ErrorAssertion{
				Line:   3,
				Column: 13,
				Err: &serrors.Assertion{
					Message: "too long manala template field (max=100)",
				},
			},
		},
		// Partials
		{
			test: "PartialsNotArray",
			content: heredoc.Doc(`
				manala:
				  description: description
				  partials: foo
			`),
			expected: &parsing.ErrorAssertion{
				Line:   3,
				Column: 13,
				Err: &serrors.Assertion{
					Message: "string was used where sequence is expected",
				},
			},
		},
		// Partials Item
		{
			test: "PartialsItemNotString",
			content: heredoc.Doc(`
				manala:
				  description: description
				  partials:
				    - []
			`),
			expected: &parsing.ErrorAssertion{
				Line:   4,
				Column: 7,
				Err: &serrors.Assertion{
					Message: "field must be a string",
				},
			},
		},
		{
			test: "PartialsItemEmpty",
			content: heredoc.Doc(`
				manala:
				  description: description
				  partials:
				    - ""
			`),
			expected: &parsing.ErrorAssertion{
				Line:   3,
				Column: 11,
				Err: &serrors.Assertion{
					Message: "empty partials entry",
				},
			},
		},
		{
			test: "PartialsItemTooLong",
			content: heredoc.Doc(`
				manala:
				  description: description
				  partials:
				    - Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
			`),
			expected: &parsing.ErrorAssertion{
				Line:   3,
				Column: 11,
				Err: &serrors.Assertion{
					Message: "too long partials entry (max=100)",
				},
			},
		},
		// Sync
		{
			test: "SyncNotArray",
			content: heredoc.Doc(`
				manala:
				  description: description
				  sync: foo
			`),
			expected: &parsing.ErrorAssertion{
				Line:   3,
				Column: 9,
				Err: &serrors.Assertion{
					Message: "sync field must be a sequence",
				},
			},
		},
		// Sync Item
		{
			test: "SyncItemNotString",
			content: heredoc.Doc(`
				manala:
				  description: description
				  sync:
				    - []
			`),
			expected: &parsing.ErrorAssertion{
				Line:   4,
				Column: 5,
				Err: &serrors.Assertion{
					Message: "sync entry must be a string",
				},
			},
		},
		{
			test: "SyncItemEmpty",
			content: heredoc.Doc(`
				manala:
				  description: description
				  sync:
				    - ""
			`),
			expected: &parsing.ErrorAssertion{
				Line:   4,
				Column: 5,
				Err: &serrors.Assertion{
					Message: "empty sync entry",
				},
			},
		},
		{
			test: "SyncItemTooLong",
			content: heredoc.Doc(`
				manala:
				  description: description
				  sync:
				    - Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
			`),
			expected: &parsing.ErrorAssertion{
				Line:   4,
				Column: 5,
				Err: &serrors.Assertion{
					Message: "too long sync entry (max=256)",
				},
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			m := manifest.New()

			err := m.Unmarshal([]byte(test.content))

			errors.Equal(s.T(), test.expected, err)
		})
	}
}
