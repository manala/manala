package api

import (
	"bytes"
	"log/slog"
	"manala/app/config"
	"manala/internal/testing/heredoc"
	"manala/internal/ui/adapters/charm"
	"manala/internal/ui/log"
	"os"
	"path/filepath"
)

func (s *Suite) TestCreateProject() {
	projectDir := filepath.FromSlash("testdata/TestCreateProject/project")
	repositoryUrl := filepath.FromSlash("testdata/TestCreateProject/repository")

	_ = os.RemoveAll(projectDir)

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	configMock := &config.Mock{}
	configMock.
		On("CacheDir").Return("").
		On("Repository").Return("")

	ui := charm.New(nil, stdout, stderr)

	api := New(
		configMock,
		slog.New(log.NewSlogHandler(ui)),
		ui,
		WithRepositoryUrl(repositoryUrl),
		WithRecipeName("recipe"),
	)

	repository, err := api.LoadPrecedingRepository()
	s.NotNil(repository)
	s.NoError(err)

	recipe, err := api.LoadPrecedingRecipe(repository)
	s.NotNil(recipe)
	s.NoError(err)

	vars := recipe.Vars()
	vars["string_float_int"] = "3.0"
	vars["string_asterisk"] = "*"

	project, err := api.CreateProject(projectDir, recipe, vars)

	s.NotNil(project)
	s.NoError(err)

	s.Empty(stdout)
	s.Empty(stderr)

	heredoc.EqualFile(s.Assert(), `
		manala:
		    recipe: recipe

		string_float_int: '3.0'
		string_float_int_value: '3.0'

		string_asterisk: '*'
		string_asterisk_value: '*'
		`,
		filepath.Join(projectDir, ".manala.yaml"),
	)
}
