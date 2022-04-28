package template

import (
	"bytes"
	"github.com/Masterminds/sprig/v3"
	"github.com/stretchr/testify/suite"
	"testing"
	textTemplate "text/template"
)

type FuncsSuite struct {
	suite.Suite
	buffer bytes.Buffer
}

func TestFuncsSuite(t *testing.T) {
	suite.Run(t, new(FuncsSuite))
}

func (s *FuncsSuite) SetupTest() {
	s.buffer.Reset()
}

func (s *FuncsSuite) TestToYaml() {
	template := textTemplate.New("")
	template.Funcs(textTemplate.FuncMap{
		"toYaml": funcToYaml(template),
	})

	_, err := template.Parse(`{{ . | toYaml }}`)
	s.NoError(err)

	err = template.Execute(&s.buffer, map[string]interface{}{
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
	})

	s.NoError(err)
	s.Equal(`foo:
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
    xyzzy: null`, s.buffer.String())

	s.Run("Cases", func() {
		s.buffer.Reset()
		template := textTemplate.New("")
		template.Funcs(textTemplate.FuncMap{
			"toYaml": funcToYaml(template),
		})

		_, err := template.Parse(`{{ . | toYaml }}`)
		s.NoError(err)

		err = template.Execute(&s.buffer, map[string]interface{}{
			"foo": map[string]interface{}{
				"bar":  true,
				"BAZ":  true,
				"qUx":  true,
				"QuuX": true,
			},
		})

		s.NoError(err)
		s.Equal(`foo:
    BAZ: true
    QuuX: true
    bar: true
    qUx: true`, s.buffer.String())
	})

	s.Run("Dict", func() {
		s.buffer.Reset()
		template := textTemplate.New("")
		template.Funcs(sprig.TxtFuncMap())
		template.Funcs(textTemplate.FuncMap{
			"toYaml": funcToYaml(template),
		})

		_, err := template.Parse(`{{ omit .foo "baz" | toYaml }}`)
		s.NoError(err)

		err = template.Execute(&s.buffer, map[string]interface{}{
			"foo": map[string]interface{}{
				"bar": true,
				"baz": true,
				"qux": true,
			},
		})

		s.NoError(err)
		s.Equal(`bar: true
qux: true`, s.buffer.String())
	})
}

func (s *FuncsSuite) TestInclude() {
	template := textTemplate.New("")
	template.Funcs(textTemplate.FuncMap{
		"include": funcInclude(template),
	})

	_, err := template.Parse(`
{{- define "foo" -}}
	foo {{ . }}
{{- end -}}
{{- include "foo" . -}}
`)
	s.NoError(err)

	err = template.Execute(&s.buffer, "bar")

	s.NoError(err)
	s.Equal(`foo bar`, s.buffer.String())
}
