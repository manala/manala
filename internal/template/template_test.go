package template

import (
	"bytes"
	"github.com/stretchr/testify/suite"
	"manala/internal/errors/serrors"
	"path/filepath"
	"testing"
)

type TemplateSuite struct {
	suite.Suite
	provider ProviderInterface
	buffer   *bytes.Buffer
}

func TestTemplateSuite(t *testing.T) {
	suite.Run(t, new(TemplateSuite))
}

func (s *TemplateSuite) SetupSuite() {
	s.provider = &Provider{}
}

func (s *TemplateSuite) SetupTest() {
	s.buffer = &bytes.Buffer{}
}

func (s *TemplateSuite) TestWriteTo() {
	template := s.provider.Template()
	err := template.WriteTo(s.buffer)

	s.NoError(err)
	s.Equal("", s.buffer.String())

	s.Run("DefaultFile", func() {
		s.buffer.Reset()

		dir := filepath.FromSlash("testdata/TemplateSuite/TestWriteTo/DefaultFile")

		template := s.provider.Template()
		template.WithDefaultFile(filepath.Join(dir, "template.tmpl"))
		template.WithDefaultContent(`{{ template "foo" . }}`)
		err := template.WriteTo(s.buffer)

		s.NoError(err)
		s.Equal("bar", s.buffer.String())
	})

	s.Run("File", func() {
		s.buffer.Reset()

		dir := filepath.FromSlash("testdata/TemplateSuite/TestWriteTo/File")

		template := s.provider.Template()
		template.WithFile(filepath.Join(dir, "template.tmpl"))
		err := template.WriteTo(s.buffer)

		s.NoError(err)
		s.Equal("bar", s.buffer.String())
	})

	s.Run("DefaultContent", func() {
		s.buffer.Reset()

		template := s.provider.Template()
		template.WithDefaultContent(`{{ "baz" }}`)
		err := template.WriteTo(s.buffer)

		s.NoError(err)
		s.Equal("baz", s.buffer.String())
	})

	s.Run("FileOverDefaultContent", func() {
		s.buffer.Reset()

		dir := filepath.FromSlash("testdata/TemplateSuite/TestWriteTo/FileOverDefaultContent")

		template := s.provider.Template()
		template.WithFile(filepath.Join(dir, "template.tmpl"))
		template.WithDefaultContent(`{{ "baz" }}`)
		err := template.WriteTo(s.buffer)

		s.NoError(err)
		s.Equal("bar", s.buffer.String())
	})

	s.Run("Data", func() {
		s.buffer.Reset()

		template := s.provider.Template()
		template.WithDefaultContent(`{{ . }}`)
		template.WithData("foo")
		err := template.WriteTo(s.buffer)

		s.NoError(err)
		s.Equal("foo", s.buffer.String())
	})

	s.Run("ParsingError", func() {
		s.buffer.Reset()

		template := s.provider.Template()
		template.WithDefaultContent(`{{ .foo }`)
		err := template.WriteTo(s.buffer)

		serrors.Equal(s.Assert(), &serrors.Assert{
			Type:    &Error{},
			Message: "unexpected \"}\" in operand",
			Arguments: []any{
				"line", 1,
			},
		}, err)
	})

	s.Run("ExecutionError", func() {
		s.buffer.Reset()

		template := s.provider.Template()
		template.WithDefaultContent(`{{ .foo }}`)
		err := template.WriteTo(s.buffer)

		serrors.Equal(s.Assert(), &serrors.Assert{
			Type:    &Error{},
			Message: "nil data; no entry for key \"foo\"",
			Arguments: []any{
				"context", ".foo",
				"line", 1,
				"column", 3,
			},
		}, err)
	})
}
