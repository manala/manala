package template

import (
	"bytes"
	"github.com/stretchr/testify/suite"
	"path/filepath"
	"testing"
)

type TemplateSuite struct {
	suite.Suite
	provider ProviderInterface
	buffer   bytes.Buffer
}

func TestTemplateSuite(t *testing.T) {
	suite.Run(t, new(TemplateSuite))
}

func (s *TemplateSuite) SetupSuite() {
	s.provider = &Provider{}
}

func (s *TemplateSuite) SetupTest() {
	s.buffer.Reset()
}

var templateTestPath = filepath.Join("testdata", "template")

func (s *TemplateSuite) TestWrite() {
	template := s.provider.Template()
	err := template.Write(&s.buffer)

	s.NoError(err)
	s.Equal("", s.buffer.String())

	s.Run("Default Files", func() {
		s.buffer.Reset()

		template := s.provider.Template()
		template.WithDefaultFile(filepath.Join(templateTestPath, "default_files", "foo.tmpl"))
		template.WithDefaultContent(`{{ template "foo" . }}`)
		err := template.Write(&s.buffer)

		s.NoError(err)
		s.Equal("bar", s.buffer.String())
	})

	s.Run("File", func() {
		s.buffer.Reset()

		template := s.provider.Template()
		template.WithFile(filepath.Join(templateTestPath, "file", "foo.tmpl"))
		err := template.Write(&s.buffer)

		s.NoError(err)
		s.Equal("bar", s.buffer.String())
	})

	s.Run("Default Content", func() {
		s.buffer.Reset()

		template := s.provider.Template()
		template.WithDefaultContent(`{{ "baz" }}`)
		err := template.Write(&s.buffer)

		s.NoError(err)
		s.Equal("baz", s.buffer.String())
	})

	s.Run("File Over Default Content", func() {
		s.buffer.Reset()

		template := s.provider.Template()
		template.WithFile(filepath.Join(templateTestPath, "file", "foo.tmpl"))
		template.WithDefaultContent(`{{ "baz" }}`)
		err := template.Write(&s.buffer)

		s.NoError(err)
		s.Equal("bar", s.buffer.String())
	})

	s.Run("Data", func() {
		s.buffer.Reset()

		template := s.provider.Template()
		template.WithDefaultContent(`{{ . }}`)
		template.WithData("foo")
		err := template.Write(&s.buffer)

		s.NoError(err)
		s.Equal("foo", s.buffer.String())
	})

	s.Run("Parsing Error", func() {
		s.buffer.Reset()

		template := s.provider.Template()
		template.WithDefaultContent(`{{ .foo }`)
		err := template.Write(&s.buffer)

		s.ErrorAs(err, &internalError)
		s.EqualError(internalError, "template error")
		s.Equal(1, internalError.Fields["line"])
		s.Equal(`unexpected "}" in operand`, internalError.Fields["message"])
	})

	s.Run("Execution Error", func() {
		s.buffer.Reset()

		template := s.provider.Template()
		template.WithDefaultContent(`{{ .foo }}`)
		err := template.Write(&s.buffer)

		s.ErrorAs(err, &internalError)
		s.EqualError(internalError, "template error")
		s.Equal(1, internalError.Fields["line"])
		s.Equal(3, internalError.Fields["column"])
		s.Equal(".foo", internalError.Fields["context"])
		s.Equal(`nil data; no entry for key "foo"`, internalError.Fields["message"])
	})
}
