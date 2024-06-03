package template_test

import (
	"bytes"
	"manala/internal/template"
	"manala/internal/testing/heredoc"
	"testing"
	gotemplate "text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/stretchr/testify/suite"
)

type FunctionsSuite struct{ suite.Suite }

func TestFunctionsSuite(t *testing.T) {
	suite.Run(t, new(FunctionsSuite))
}

func (s *FunctionsSuite) execute(content string, data any) string {
	tmpl := gotemplate.New("test")

	tmpl.Funcs(sprig.TxtFuncMap())
	tmpl.Funcs(template.FuncMap(tmpl))

	_, err := tmpl.Parse(content)
	s.Require().NoError(err)

	buffer := &bytes.Buffer{}
	err = tmpl.Execute(buffer, data)
	s.Require().NoError(err)

	return buffer.String()
}

func (s *FunctionsSuite) TestToYaml() {
	s.Run("Default", func() {
		content := s.execute(`{{ . | toYaml }}`, map[string]any{
			"foo": map[string]any{
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
				"garply": map[string]any{},
				"waldo": map[string]any{
					"foo": "bar",
					"bar": "baz",
				},
				"fred": []any{},
				"plugh": []any{
					"foo",
					"bar",
				},
				"xyzzy": nil,
				"thud":  "123",
			},
		})

		heredoc.Equal(s.T(), `
			foo:
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
			    thud: '123'
			    waldo:
			        bar: baz
			        foo: bar
			    xyzzy: null`,
			content,
		)
	})

	s.Run("Cases", func() {
		content := s.execute(`{{ . | toYaml }}`, map[string]any{
			"foo": map[string]any{
				"bar":  true,
				"BAZ":  true,
				"qUx":  true,
				"QuuX": true,
			},
		})

		heredoc.Equal(s.T(), `
			foo:
			    BAZ: true
			    QuuX: true
			    bar: true
			    qUx: true`,
			content,
		)
	})

	s.Run("Mapping", func() {
		content := s.execute(`{{ omit .foo "baz" | toYaml }}`, map[string]any{
			"foo": map[string]any{
				"bar": true,
				"baz": true,
				"qux": true,
			},
		})

		heredoc.Equal(s.T(), `
			bar: true
			qux: true`,
			content,
		)
	})

	s.Run("RootSequence", func() {
		content := s.execute(`{{ . | toYaml }}`, []string{
			"foo",
			"bar",
			"baz",
		})

		heredoc.Equal(s.T(), `
			- foo
			- bar
			- baz`,
			content,
		)
	})

	s.Run("NestedSequence", func() {
		content := s.execute(`{{ . | toYaml }}`, map[string]any{
			"nested": []string{
				"foo",
				"bar",
				"baz",
			},
		})

		heredoc.Equal(s.T(), `
			nested:
			    - foo
			    - bar
			    - baz`,
			content,
		)
	})

	s.Run("Quotes", func() {
		content := s.execute(`{{ . | toYaml }}`, `'single' "double"`)

		heredoc.Equal(s.T(),
			`'\'single\' "double"'`,
			content,
		)
	})

	s.Run("BlockScalar", func() {
		content := s.execute(`{{ . | toYaml }}`, map[string]any{
			"scalar": `foo
bar\baz
`,
		})

		heredoc.Equal(s.T(), `
			scalar: |
			  foo
			  bar\baz`,
			content,
		)
	})

	s.Run("Indentation", func() {
		content := s.execute(`{{ . | toYaml }}`, map[string]any{
			"mapping": map[string]any{
				"foo": "bar",
				"bar": "baz",
			},
			"sequence": []string{
				"foo",
				"bar",
			},
		})

		heredoc.Equal(s.T(), `
			mapping:
			    bar: baz
			    foo: bar
			sequence:
			    - foo
			    - bar`,
			content,
		)
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

	heredoc.Equal(s.T(),
		`foo bar`,
		content,
	)
}
