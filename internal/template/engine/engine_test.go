package engine_test

import (
	"path/filepath"
	"testing"

	"github.com/manala/manala/internal/serrors"
	"github.com/manala/manala/internal/template/engine"
	"github.com/manala/manala/internal/testing/expect"
	"github.com/manala/manala/internal/testing/heredoc"

	"github.com/stretchr/testify/suite"
)

type EngineSuite struct{ suite.Suite }

func TestEngineSuite(t *testing.T) {
	suite.Run(t, new(EngineSuite))
}

func (s *EngineSuite) TestExecutor() {
	e := engine.New()

	s.Run("Default", func() {
		executor, err := e.Executor("data")

		s.Require().NoError(err)
		s.NotNil(executor)
	})

	s.Run("Files", func() {
		dir := filepath.FromSlash("testdata/EngineSuite/TestExecutor/Files")

		executor, err := e.Executor("data", filepath.Join(dir, "partial.tmpl"))

		s.Require().NoError(err)
		s.NotNil(executor)
	})

	s.Run("FilesNotFound", func() {
		dir := filepath.FromSlash("testdata/EngineSuite/TestExecutor/FilesNotFound")

		executor, err := e.Executor("data", filepath.Join(dir, "not_found.tmpl"))

		s.Nil(executor)
		expect.Error(s.T(), serrors.Expectation{
			Message: "unable to read template file",
			Attrs: [][2]any{
				{"path", filepath.Join(dir, "not_found.tmpl")},
			},
			Errors: []expect.ErrorExpectation{
				serrors.Expectation{
					Message: "file does not exist",
					Attrs: [][2]any{
						{"operation", "open"},
						{"path", filepath.Join(dir, "not_found.tmpl")},
					},
				},
			},
		}, err)
	})

	s.Run("FilesInvalid", func() {
		dir := filepath.FromSlash("testdata/EngineSuite/TestExecutor/FilesInvalid")

		executor, err := e.Executor("data", filepath.Join(dir, "partial.tmpl"))

		s.Nil(executor)
		expect.Error(s.T(), serrors.Expectation{
			Message: "unable to parse template file",
			Dump: heredoc.Doc(`
				in %[1]s:2
				  1 | foo
				> 2 |   {{ .bar }
				  3 | baz
				* unexpected "}" in operand
			`,
				filepath.Join(dir, "partial.tmpl"),
			),
		}, err)
	})
}
