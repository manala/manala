package template

import (
	"bytes"
	"github.com/Masterminds/sprig/v3"
	"manala/internal/testing/heredoc"
	"text/template"
)

func (s *Suite) execute(content string, data any) string {
	template := template.New("test")

	template.Funcs(sprig.TxtFuncMap())
	template.Funcs(FuncMap(template))

	_, err := template.Parse(content)
	s.NoError(err)

	buffer := &bytes.Buffer{}
	err = template.Execute(buffer, data)
	s.NoError(err)

	return buffer.String()
}

func (s *Suite) TestFunctionToYaml() {

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

		heredoc.Equal(s.Assert(), `
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

		heredoc.Equal(s.Assert(), `
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

		heredoc.Equal(s.Assert(), `
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

		heredoc.Equal(s.Assert(), `
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

		heredoc.Equal(s.Assert(), `
			nested:
			    - foo
			    - bar
			    - baz`,
			content,
		)
	})

	s.Run("Quotes", func() {
		content := s.execute(`{{ . | toYaml }}`, `'single' "double"`)

		heredoc.Equal(s.Assert(),
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

		heredoc.Equal(s.Assert(), `
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

		heredoc.Equal(s.Assert(), `
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

func (s *Suite) TestFunctionInclude() {
	content := s.execute(
		`{{- define "foo" -}}
	foo {{ . }}
{{- end -}}
{{- include "foo" . -}}
`,
		"bar",
	)

	heredoc.Equal(s.Assert(),
		`foo bar`,
		content,
	)
}
