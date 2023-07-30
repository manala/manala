package lipgloss

import (
	"bytes"
	"github.com/stretchr/testify/suite"
	"testing"
)

type Suite struct {
	suite.Suite
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) Test() {

	s.Run("writeOutString", func() {
		out := &bytes.Buffer{}
		err := &bytes.Buffer{}

		output := New(out, err)

		output.writeOutString("string")

		s.Equal("string", out.String())
		s.Empty(err)
	})

	s.Run("writeErrString", func() {
		out := &bytes.Buffer{}
		err := &bytes.Buffer{}

		output := New(out, err)

		output.writeErrString("string")

		s.Empty(out)
		s.Equal("string", err.String())
	})

	s.Run("outStyle", func() {
		output := New(nil, nil)

		s.Equal("string", output.outStyle().Render("string"))
	})

	s.Run("errStyle", func() {
		output := New(nil, nil)

		s.Equal("string", output.errStyle().Render("string"))
	})
}
