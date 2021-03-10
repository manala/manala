package template

import (
	"bytes"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

/*********/
/* Suite */
/*********/

type TemplateTestSuite struct {
	suite.Suite
}

func TestTemplateTestSuite(t *testing.T) {
	// Run
	suite.Run(t, new(TemplateTestSuite))
}

/*********/
/* Tests */
/*********/

func (s *TemplateTestSuite) Test() {
	for _, t := range []struct {
		test    string
		file    string
		helpers string
		data    interface{}
		out     string
		err     string
	}{
		{
			test: "Base",
			file: "testdata/base.tmpl",
			out: `foo
`,
		},
		{
			test: "Invalid",
			file: "testdata/invalid.tmpl",
			err:  "template: :1:3: executing \"\" at <.foo>: nil data; no entry for key \"foo\"",
		},
		{
			test: "To Yaml",
			file: "testdata/to_yaml.tmpl",
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
			file: "testdata/cases.tmpl",
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
			file: "testdata/dict.tmpl",
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
			file: "testdata/include.tmpl",
			out:  `foo: bar`,
		},
		{
			test:    "Helpers",
			file:    "testdata/helpers.tmpl",
			helpers: "testdata/_helpers.tmpl",
			out:     `bar: foo`,
		},
	} {
		s.Run(t.test, func() {
			tmpl := New()
			if t.helpers != "" {
				_ = tmpl.ParseFiles(t.helpers)
			}
			tmplContent, _ := os.ReadFile(t.file)
			_ = tmpl.Parse(string(tmplContent))
			var tmplOut bytes.Buffer
			err := tmpl.Execute(&tmplOut, t.data)
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
