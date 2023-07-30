package application

import (
	"bytes"
	"github.com/stretchr/testify/suite"
	"log/slog"
	"manala/app/interfaces"
	"manala/app/mocks"
	"manala/internal/testing/file"
	"manala/internal/testing/heredoc"
	"manala/internal/ui/log"
	"manala/internal/ui/output/lipgloss"
	"os"
	"path/filepath"
	"testing"
)

type ApplicationSuite struct {
	suite.Suite
}

func TestApplicationSuite(t *testing.T) {
	suite.Run(t, new(ApplicationSuite))
}

func (s *ApplicationSuite) TestCreateProject() {
	projDir := filepath.FromSlash("testdata/ApplicationSuite/TestCreateProject/project")
	repoUrl := filepath.FromSlash("testdata/ApplicationSuite/TestCreateProject/repository")

	_ = os.RemoveAll(projDir)

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	confMock := &mocks.ConfigMock{}
	confMock.
		On("CacheDir").Return("").
		On("Repository").Return("")

	out := lipgloss.New(stdout, stderr)

	app := NewApplication(
		confMock,
		slog.New(log.NewSlogHandler(out)),
		out,
		WithRepositoryUrl(repoUrl),
	)

	proj, err := app.CreateProject(
		projDir,
		// Recipe selector
		func(recWalker func(walker func(rec interfaces.Recipe) error) error) (interfaces.Recipe, error) {
			var rec interfaces.Recipe
			_ = recWalker(func(_rec interfaces.Recipe) error {
				rec = _rec
				return nil
			})
			return rec, nil
		},
		// Options selector
		func(rec interfaces.Recipe, options []interfaces.RecipeOption) error {
			// String float int
			_ = options[0].Set("3.0")
			// String asterisk
			_ = options[1].Set("*")
			return nil
		},
	)

	s.NotNil(proj)

	s.NoError(err)

	s.Empty(stdout)
	s.Empty(stderr)

	file.EqualContent(s.Assert(), heredoc.Docf(`
		manala:
		    recipe: recipe

		string_float_int: '3.0'
		string_float_int_value: '3.0'

		string_asterisk: '*'
		string_asterisk_value: '*'
		`),
		filepath.Join(projDir, ".manala.yaml"),
	)
}
