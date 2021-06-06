package template

import (
	"bytes"
	"github.com/stretchr/testify/suite"
	"manala/fs"
	"testing"
)

/*********/
/* Suite */
/*********/

type TemplateTestSuite struct {
	suite.Suite
	fsManager fs.ManagerInterface
	manager   ManagerInterface
}

func TestTemplateTestSuite(t *testing.T) {
	// Run
	suite.Run(t, new(TemplateTestSuite))
}

func (s *TemplateTestSuite) SetupTest() {
	s.fsManager = fs.NewManager()
	s.manager = NewManager()
}

/*********/
/* Tests */
/*********/

func (s *TemplateTestSuite) Test() {
	for _, t := range []struct {
		test    string
		file    string
		helpers string
		data    map[string]interface{}
		out     string
		err     string
	}{
		{
			test: "Base",
			file: "base.tmpl",
			out: `foo
`,
		},
		{
			test: "Invalid",
			file: "invalid.tmpl",
			err:  "template: :1:3: executing \"\" at <.foo>: map has no entry for key \"foo\"",
		},
		{
			test: "To Yaml",
			file: "to_yaml.tmpl",
			data: map[string]interface{}{
				"foo": map[string]interface{}{
					"bar": "string",
					"baz": struct {
						Foo string
						Bar int
					}{
						Foo: "foo",
						Bar: 123,
					},
					"qux":    123,
					"quux":   true,
					"corge":  false,
					"grault": 1.23,
					"garply": map[string]interface{}{},
					"waldo": map[string]interface{}{
						"foo": "bar",
						"bar": "baz",
					},
					"fred": []interface{}{},
					"plugh": []interface{}{
						"foo",
						"bar",
					},
					"xyzzy": nil,
					"thud":  "123",
				},
			},
			out: `foo:
    bar: string
    baz:
        foo: foo
        bar: 123
    corge: false
    fred: []
    garply: {}
    grault: 1.23
    plugh:
        - foo
        - bar
    quux: true
    qux: 123
    thud: "123"
    waldo:
        bar: baz
        foo: bar
    xyzzy: null
`,
		},
		{
			test: "Cases",
			file: "cases.tmpl",
			data: map[string]interface{}{
				"foo": map[string]interface{}{
					"bar":  true,
					"BAZ":  true,
					"qUx":  true,
					"QuuX": true,
				},
			},
			out: `foo:
    BAZ: true
    QuuX: true
    bar: true
    qUx: true
`,
		},
		{
			test: "Dict",
			file: "dict.tmpl",
			data: map[string]interface{}{
				"foo": map[string]interface{}{
					"bar": true,
					"baz": true,
					"qux": true,
				},
			},
			out: `bar: true
qux: true
`,
		},
		{
			test: "Include",
			file: "include.tmpl",
			out:  `foo: bar`,
		},
		{
			test:    "Helpers",
			file:    "helpers.tmpl",
			helpers: "_helpers.tmpl",
			out:     `bar: foo`,
		},
	} {
		s.Run(t.test, func() {
			fsys := s.fsManager.NewDirFs("testdata")
			template := s.manager.NewFsTemplate(fsys)
			if t.helpers != "" {
				_ = template.ParseFiles(t.helpers)
			}
			tmplContent, _ := fsys.ReadFile(t.file)
			_ = template.Parse(string(tmplContent))
			var tmplOut bytes.Buffer
			err := template.Execute(&tmplOut, t.data)
			if t.err != "" {
				s.Error(err)
				s.Equal(t.err, err.Error())
			} else {
				s.NoError(err)
			}
			s.Equal(t.out, tmplOut.String())
		})
	}
}
