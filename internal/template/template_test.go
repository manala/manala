package template

import (
	"bytes"
	"github.com/stretchr/testify/suite"
	internalReport "manala/internal/report"
	internalTesting "manala/internal/testing"
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

	s.Run("Default File", func() {
		s.buffer.Reset()

		template := s.provider.Template()
		template.WithDefaultFile(internalTesting.DataPath(s, "template.tmpl"))
		template.WithDefaultContent(`{{ template "foo" . }}`)
		err := template.WriteTo(s.buffer)

		s.NoError(err)
		s.Equal("bar", s.buffer.String())
	})

	s.Run("File", func() {
		s.buffer.Reset()

		template := s.provider.Template()
		template.WithFile(internalTesting.DataPath(s, "template.tmpl"))
		err := template.WriteTo(s.buffer)

		s.NoError(err)
		s.Equal("bar", s.buffer.String())
	})

	s.Run("Default Content", func() {
		s.buffer.Reset()

		template := s.provider.Template()
		template.WithDefaultContent(`{{ "baz" }}`)
		err := template.WriteTo(s.buffer)

		s.NoError(err)
		s.Equal("baz", s.buffer.String())
	})

	s.Run("File Over Default Content", func() {
		s.buffer.Reset()

		template := s.provider.Template()
		template.WithFile(internalTesting.DataPath(s, "template.tmpl"))
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

	s.Run("Parsing Error", func() {
		s.buffer.Reset()

		template := s.provider.Template()
		template.WithDefaultContent(`{{ .foo }`)
		err := template.WriteTo(s.buffer)

		s.EqualError(err, "unexpected \"}\" in operand")

		report := internalReport.NewErrorReport(err)

		reportAssert := &internalReport.Assert{
			Message: "template error",
			Err:     "unexpected \"}\" in operand",
			Fields: map[string]interface{}{
				"line": 1,
			},
		}
		reportAssert.Equal(&s.Suite, report)
	})

	s.Run("Execution Error", func() {
		s.buffer.Reset()

		template := s.provider.Template()
		template.WithDefaultContent(`{{ .foo }}`)
		err := template.WriteTo(s.buffer)

		s.EqualError(err, "nil data; no entry for key \"foo\"")

		report := internalReport.NewErrorReport(err)

		reportAssert := &internalReport.Assert{
			Message: "template error",
			Err:     "nil data; no entry for key \"foo\"",
			Fields: map[string]interface{}{
				"line":    1,
				"column":  3,
				"context": ".foo",
			},
		}
		reportAssert.Equal(&s.Suite, report)
	})
}
