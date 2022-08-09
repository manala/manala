package template

import (
	"bytes"
	"github.com/Masterminds/sprig/v3"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/suite"
	internalTesting "manala/internal/testing"
	"testing"
	textTemplate "text/template"
)

type FunctionsSuite struct {
	suite.Suite
	goldie *goldie.Goldie
	buffer *bytes.Buffer
}

func TestFunctionsSuite(t *testing.T) {
	suite.Run(t, new(FunctionsSuite))
}

func (s *FunctionsSuite) SetupTest() {
	s.goldie = goldie.New(s.T())
	s.buffer = &bytes.Buffer{}
}

func (s *FunctionsSuite) execute(content string, data interface{}) string {
	template := textTemplate.New("test")

	template.Funcs(sprig.TxtFuncMap())
	template.Funcs(FuncMap(template))

	_, err := template.Parse(content)
	s.NoError(err)

	buffer := &bytes.Buffer{}
	err = template.Execute(buffer, data)
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

		s.goldie.Assert(s.T(), internalTesting.Path(s, "content"), []byte(content))
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

		s.goldie.Assert(s.T(), internalTesting.Path(s, "content"), []byte(content))
	})

	s.Run("Mapping", func() {
		content := s.execute(`{{ omit .foo "baz" | toYaml }}`, map[string]interface{}{
			"foo": map[string]interface{}{
				"bar": true,
				"baz": true,
				"qux": true,
			},
		})

		s.goldie.Assert(s.T(), internalTesting.Path(s, "content"), []byte(content))
	})

	s.Run("Root Sequence", func() {
		content := s.execute(`{{ . | toYaml }}`, []string{
			"foo",
			"bar",
			"baz",
		})

		s.goldie.Assert(s.T(), internalTesting.Path(s, "content"), []byte(content))
	})

	s.Run("Nested Sequence", func() {
		content := s.execute(`{{ . | toYaml }}`, map[string]interface{}{
			"nested": []string{
				"foo",
				"bar",
				"baz",
			},
		})

		s.goldie.Assert(s.T(), internalTesting.Path(s, "content"), []byte(content))
	})

	s.Run("Quotes", func() {
		content := s.execute(`{{ . | toYaml }}`, `'single' "double"`)

		s.goldie.Assert(s.T(), internalTesting.Path(s, "content"), []byte(content))
	})

	s.Run("Block Scalar", func() {
		content := s.execute(`{{ . | toYaml }}`, map[string]interface{}{
			"scalar": `foo
bar\baz
`,
		})

		s.goldie.Assert(s.T(), internalTesting.Path(s, "content"), []byte(content))
	})

	s.Run("Indentation", func() {
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

		s.goldie.Assert(s.T(), internalTesting.Path(s, "content"), []byte(content))
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

	s.goldie.Assert(s.T(), internalTesting.Path(s, "content"), []byte(content))
}
