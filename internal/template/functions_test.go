package template

import (
	"bytes"
	"github.com/Masterminds/sprig/v3"
	"github.com/stretchr/testify/suite"
	"testing"
	textTemplate "text/template"
)

type FunctionsSuite struct {
	suite.Suite
	buffer bytes.Buffer
}

func TestFunctionsSuite(t *testing.T) {
	suite.Run(t, new(FunctionsSuite))
}

func (s *FunctionsSuite) SetupTest() {
	s.buffer.Reset()
}

func (s *FunctionsSuite) execute(content string, data interface{}) string {
	template := textTemplate.New("test")

	template.Funcs(sprig.TxtFuncMap())
	template.Funcs(FuncMap(template))

	_, err := template.Parse(content)
	s.NoError(err)

	var buffer bytes.Buffer
	err = template.Execute(&buffer, data)
	s.NoError(err)

	return buffer.String()
}

func (s *FunctionsSuite) TestToYaml() {

	s.Run("Default", func() {
		content := s.execute(`{{ . | toYaml }}`, map[string]interface{}{
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
    xyzzy: null`, content)
	})

	s.Run("Cases", func() {
		content := s.execute(`{{ . | toYaml }}`, map[string]interface{}{
			"foo": map[string]interface{}{
				"bar":  true,
				"BAZ":  true,
				"qUx":  true,
				"QuuX": true,
			},
		})

		s.Equal(`foo:
    BAZ: true
    QuuX: true
    bar: true
    qUx: true`, content)
	})

	s.Run("Mapping", func() {
		content := s.execute(`{{ omit .foo "baz" | toYaml }}`, map[string]interface{}{
			"foo": map[string]interface{}{
				"bar": true,
				"baz": true,
				"qux": true,
			},
		})

		s.Equal(`bar: true
qux: true`, content)
	})

	s.Run("Indent", func() {
		content := s.execute(`{{ . | toYaml }}`, map[string]interface{}{
			"mapping": map[string]interface{}{
				"foo": "bar",
				"bar": "baz",
			},
			"sequence": []string{
				"foo",
				"bar",
			},
		})

		s.Equal(`mapping:
    bar: baz
    foo: bar
sequence:
    - foo
    - bar`, content)
	})
}

func (s *FunctionsSuite) TestInclude() {
	content := s.execute(
		`{{- define "foo" -}}
	foo {{ . }}
{{- end -}}
{{- include "foo" . -}}
`,
		"bar",
	)

	s.Equal(`foo bar`, content)
}
