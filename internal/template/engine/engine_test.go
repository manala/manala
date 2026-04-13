package engine_test

import (
	"path/filepath"
	"testing"

	"github.com/manala/manala/internal/serrors"
	"github.com/manala/manala/internal/template/engine"
	"github.com/manala/manala/internal/testing/errors"

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
		errors.Equal(s.T(), &serrors.Assertion{
			Message: "unable to read template file",
			Arguments: []any{
				"path", filepath.Join(dir, "not_found.tmpl"),
			},
			Errors: []errors.Assertion{
				&serrors.Assertion{
					Message: "file does not exist",
					Arguments: []any{
						"operation", "open",
						"path", filepath.Join(dir, "not_found.tmpl"),
					},
				},
			},
		}, err)
	})

	s.Run("FilesInvalid", func() {
		dir := filepath.FromSlash("testdata/EngineSuite/TestExecutor/FilesInvalid")

		executor, err := e.Executor("data", filepath.Join(dir, "partial.tmpl"))

		s.Nil(executor)
		errors.Equal(s.T(), &serrors.Assertion{
			Message: "unable to parse template file",
			Arguments: []any{
				"path", filepath.Join(dir, "partial.tmpl"),
				"line", 2, "column", 0,
			},
			Dump: `
				  1 | foo
				> 2 |   {{ .bar }
				  3 | baz
				* unexpected "}" in operand
			`,
		}, err)
	})
}
