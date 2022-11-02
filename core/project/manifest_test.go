package project

import (
	"github.com/stretchr/testify/suite"
	internalReport "manala/internal/report"
	internalTesting "manala/internal/testing"
	"os"
	"testing"
)

type ManifestSuite struct{ suite.Suite }

func TestManifestSuite(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	suite.Run(t, new(ManifestSuite))
}

func (s *ManifestSuite) Test() {
	recMan := NewManifest()

	s.Equal("", recMan.Recipe())
	s.Equal("", recMan.Repository())
	s.Equal(map[string]interface{}{}, recMan.Vars())
}

func (s *ManifestSuite) TestReadFromErrors() {
	tests := []struct {
		name   string
		report *internalReport.Assert
	}{
		{
			name: "Empty",
			report: &internalReport.Assert{
				Message: "irregular project manifest",
				Err:     "empty yaml file",
			},
		},
		{
			name: "Invalid",
			report: &internalReport.Assert{
				Message: "irregular project manifest",
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
				Message: "irregular project manifest",
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
				Message: "irregular project manifest",
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
				Err: "invalid project manifest",
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
				Err: "invalid project manifest",
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
				Err: "invalid project manifest",
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
				Err: "invalid project manifest",
				Reports: []internalReport.Assert{
					{
						Message: "missing manala recipe field",
						Fields: map[string]interface{}{
							"line":     1,
							"column":   9,
							"property": "recipe",
						},
					},
				},
			},
		},
		{
			name: "Config Additional Properties",
			report: &internalReport.Assert{
				Err: "invalid project manifest",
				Reports: []internalReport.Assert{
					{
						Message: "manala field don't support additional properties",
						Fields: map[string]interface{}{
							"line":     2,
							"column":   9,
							"property": "foo",
						},
					},
				},
			},
		},
		// Config - Recipe
		{
			name: "Config Recipe Absent",
			report: &internalReport.Assert{
				Err: "invalid project manifest",
				Reports: []internalReport.Assert{
					{
						Message: "missing manala recipe field",
						Fields: map[string]interface{}{
							"line":     2,
							"column":   13,
							"property": "recipe",
						},
					},
				},
			},
		},
		{
			name: "Config Recipe Not String",
			report: &internalReport.Assert{
				Err: "invalid project manifest",
				Reports: []internalReport.Assert{
					{
						Message: "manala recipe field must be a string",
						Fields: map[string]interface{}{
							"line":     2,
							"column":   11,
							"expected": "string",
							"given":    "array",
						},
					},
				},
			},
		},
		{
			name: "Config Recipe Empty",
			report: &internalReport.Assert{
				Err: "invalid project manifest",
				Reports: []internalReport.Assert{
					{
						Message: "empty manala recipe field",
						Fields: map[string]interface{}{
							"line":   2,
							"column": 11,
						},
					},
				},
			},
		},
		{
			name: "Config Recipe Too Long",
			report: &internalReport.Assert{
				Err: "invalid project manifest",
				Reports: []internalReport.Assert{
					{
						Message: "too long manala recipe field",
						Fields: map[string]interface{}{
							"line":   2,
							"column": 11,
						},
					},
				},
			},
		},
		// Config - Repository
		{
			name: "Config Repository Not String",
			report: &internalReport.Assert{
				Err: "invalid project manifest",
				Reports: []internalReport.Assert{
					{
						Message: "manala repository field must be a string",
						Fields: map[string]interface{}{
							"line":     3,
							"column":   15,
							"expected": "string",
							"given":    "array",
						},
					},
				},
			},
		},
		{
			name: "Config Repository Empty",
			report: &internalReport.Assert{
				Err: "invalid project manifest",
				Reports: []internalReport.Assert{
					{
						Message: "empty manala repository field",
						Fields: map[string]interface{}{
							"line":   3,
							"column": 15,
						},
					},
				},
			},
		},
		{
			name: "Config Repository Too Long",
			report: &internalReport.Assert{
				Err: "invalid project manifest",
				Reports: []internalReport.Assert{
					{
						Message: "too long manala repository field",
						Fields: map[string]interface{}{
							"line":   3,
							"column": 15,
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			recMan := NewManifest()

			recManFile, _ := os.Open(internalTesting.DataPath(s, "manifest.yaml"))
			err := recMan.ReadFrom(recManFile)

			s.Error(err)

			report := internalReport.NewErrorReport(err)

			test.report.Equal(&s.Suite, report)
		})
	}
}

func (s *ManifestSuite) TestReadFrom() {
	tests := []struct {
		name       string
		recipe     string
		repository string
		vars       map[string]interface{}
	}{
		{
			name:       "All",
			recipe:     "recipe",
			repository: "repository",
			vars: map[string]interface{}{
				"foo": "bar",
			},
		},
		{
			name:       "Config Repository Absent",
			recipe:     "recipe",
			repository: "",
			vars: map[string]interface{}{
				"foo": "bar",
			},
		},
		{
			name:       "Vars Absent",
			recipe:     "recipe",
			repository: "repository",
			vars:       map[string]interface{}{},
		},
		{
			name:       "Vars Keys",
			recipe:     "recipe",
			repository: "repository",
			vars: map[string]interface{}{
				"underscore_key": "ok",
				"hyphen-key":     "ok",
				"dot.key":        "ok",
			},
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			recMan := NewManifest()

			recManFile, _ := os.Open(internalTesting.DataPath(s, "manifest.yaml"))
			err := recMan.ReadFrom(recManFile)

			s.NoError(err)
			s.Equal(test.recipe, recMan.Recipe())
			s.Equal(test.repository, recMan.Repository())
			s.Equal(test.vars, recMan.Vars())
		})
	}
}
