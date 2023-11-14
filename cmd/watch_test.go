package cmd

import (
	"bytes"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/suite"
	"log/slog"
	"manala/app"
	"manala/app/config"
	"manala/internal/serrors"
	"manala/internal/testing/cmd"
	"manala/internal/ui/adapters/charm"
	"manala/internal/ui/log"
	"path/filepath"
	"testing"
)

type WatchSuite struct {
	suite.Suite
	configMock *config.Mock
	executor   *cmd.Executor
}

func TestWatchSuite(t *testing.T) {
	suite.Run(t, new(WatchSuite))
}

func (s *WatchSuite) SetupTest() {
	s.configMock = &config.Mock{}
	s.executor = cmd.NewExecutor(func(stdout *bytes.Buffer, stderr *bytes.Buffer) *cobra.Command {
		ui := charm.New(nil, stdout, stderr)
		return newWatchCmd(
			s.configMock,
			slog.New(log.NewSlogHandler(ui)),
			ui,
		)
	})
}

func (s *WatchSuite) TestProjectErrors() {
	s.configMock.
		On("CacheDir").Return("").
		On("Repository").Return("")

	s.Run("ProjectNotFound", func() {
		projectDir := filepath.FromSlash("testdata/WatchSuite/TestProjectErrors/ProjectNotFound/project")

		err := s.executor.Execute([]string{
			projectDir,
		})

		s.Empty(s.executor.Stdout)
		s.Empty(s.executor.Stderr)

		serrors.Equal(s.Assert(), &serrors.Assert{
			Type:    serrors.Error{},
			Message: "project not found",
			Arguments: []any{
				"dir", projectDir,
			},
		}, err)
	})

	s.Run("WrongProjectManifest", func() {
		projectDir := filepath.FromSlash("testdata/WatchSuite/TestProjectErrors/WrongProjectManifest/project")

		err := s.executor.Execute([]string{
			projectDir,
		})

		s.Empty(s.executor.Stdout)
		s.Empty(s.executor.Stderr)

		serrors.Equal(s.Assert(), &serrors.Assert{
			Type:    serrors.Error{},
			Message: "project manifest is a directory",
			Arguments: []any{
				"dir", filepath.Join(projectDir, ".manala.yaml"),
			},
		}, err)
	})

	s.Run("EmptyProjectManifest", func() {
		projectDir := filepath.FromSlash("testdata/WatchSuite/TestProjectErrors/EmptyProjectManifest/project")

		err := s.executor.Execute([]string{
			projectDir,
		})

		s.Empty(s.executor.Stdout)
		s.Empty(s.executor.Stderr)

		serrors.Equal(s.Assert(), &serrors.Assert{
			Type:    serrors.Error{},
			Message: "unable to read project manifest",
			Arguments: []any{
				"file", filepath.Join(projectDir, ".manala.yaml"),
			},
			Errors: []*serrors.Assert{
				{
					Type:    serrors.Error{},
					Message: "irregular project manifest",
					Errors: []*serrors.Assert{
						{
							Type:    serrors.Error{},
							Message: "empty yaml file",
						},
					},
				},
			},
		}, err)
	})

	s.Run("InvalidProjectManifest", func() {
		projectDir := filepath.FromSlash("testdata/WatchSuite/TestProjectErrors/InvalidProjectManifest/project")

		err := s.executor.Execute([]string{
			projectDir,
		})

		s.Empty(s.executor.Stdout)
		s.Empty(s.executor.Stderr)

		serrors.Equal(s.Assert(), &serrors.Assert{
			Type:    serrors.Error{},
			Message: "unable to read project manifest",
			Arguments: []any{
				"file", filepath.Join(projectDir, ".manala.yaml"),
			},
			Errors: []*serrors.Assert{
				{
					Type:    serrors.Error{},
					Message: "invalid project manifest",
					Errors: []*serrors.Assert{
						{
							Type:    serrors.Error{},
							Message: "missing manala recipe property",
							Arguments: []any{
								"path", "manala",
								"property", "recipe",
								"line", 1,
								"column", 9,
							},
							Details: `
								>  1 | manala: {}
								               ^
							`,
						},
					},
				},
			},
		}, err)
	})
}

func (s *WatchSuite) TestRepositoryErrors() {
	s.configMock.
		On("CacheDir").Return("").
		On("Repository").Return("")

	s.Run("NoRepository", func() {
		projectDir := filepath.FromSlash("testdata/WatchSuite/TestRepositoryErrors/NoRepository/project")

		err := s.executor.Execute([]string{
			projectDir,
		})

		s.Empty(s.executor.Stdout)
		s.Empty(s.executor.Stderr)

		serrors.Equal(s.Assert(), &serrors.Assert{
			Type:    &app.UnprocessableRepositoryUrlError{},
			Message: "unable to process repository url",
		}, err)
	})

	s.Run("RepositoryNotFound", func() {
		projectDir := filepath.FromSlash("testdata/WatchSuite/TestRepositoryErrors/RepositoryNotFound/project")
		repositoryUrl := filepath.FromSlash("testdata/WatchSuite/TestRepositoryErrors/RepositoryNotFound/repository")

		err := s.executor.Execute([]string{
			projectDir,
			"--repository", repositoryUrl,
		})

		s.Empty(s.executor.Stdout)
		s.Empty(s.executor.Stderr)

		serrors.Equal(s.Assert(), &serrors.Assert{
			Type:    &app.UnsupportedRepositoryError{},
			Message: "unsupported repository url",
			Arguments: []any{
				"url", repositoryUrl,
			},
		}, err)
	})

	s.Run("WrongRepository", func() {
		projectDir := filepath.FromSlash("testdata/WatchSuite/TestRepositoryErrors/WrongRepository/project")
		repositoryUrl := filepath.FromSlash("testdata/WatchSuite/TestRepositoryErrors/WrongRepository/repository")

		err := s.executor.Execute([]string{
			projectDir,
			"--repository", repositoryUrl,
		})

		s.Empty(s.executor.Stdout)
		s.Empty(s.executor.Stderr)

		serrors.Equal(s.Assert(), &serrors.Assert{
			Type:    &app.UnsupportedRepositoryError{},
			Message: "unsupported repository url",
			Arguments: []any{
				"url", repositoryUrl,
			},
		}, err)
	})
}

func (s *WatchSuite) TestRecipeErrors() {
	s.configMock.
		On("CacheDir").Return("").
		On("Repository").Return("")

	s.Run("RecipeNotFound", func() {
		projectDir := filepath.FromSlash("testdata/WatchSuite/TestRecipeErrors/RecipeNotFound/project")
		repositoryUrl := filepath.FromSlash("testdata/WatchSuite/TestRecipeErrors/RecipeNotFound/repository")

		err := s.executor.Execute([]string{
			projectDir,
			"--repository", repositoryUrl,
			"--recipe", "recipe",
		})

		s.Empty(s.executor.Stdout)
		s.Empty(s.executor.Stderr)

		serrors.Equal(s.Assert(), &serrors.Assert{
			Type:    &app.NotFoundRecipeManifestError{},
			Message: "recipe manifest not found",
			Arguments: []any{
				"file", filepath.Join(repositoryUrl, "recipe", ".manala.yaml"),
			},
		}, err)
	})

	s.Run("WrongRecipeManifest", func() {
		projectDir := filepath.FromSlash("testdata/WatchSuite/TestRecipeErrors/WrongRecipeManifest/project")
		repositoryUrl := filepath.FromSlash("testdata/WatchSuite/TestRecipeErrors/WrongRecipeManifest/repository")

		err := s.executor.Execute([]string{
			projectDir,
			"--repository", repositoryUrl,
			"--recipe", "recipe",
		})

		s.Empty(s.executor.Stdout)
		s.Empty(s.executor.Stderr)

		serrors.Equal(s.Assert(), &serrors.Assert{
			Type:    serrors.Error{},
			Message: "recipe manifest is a directory",
			Arguments: []any{
				"dir", filepath.Join(repositoryUrl, "recipe", ".manala.yaml"),
			},
		}, err)
	})

	s.Run("InvalidRecipeManifest", func() {
		projectDir := filepath.FromSlash("testdata/WatchSuite/TestRecipeErrors/InvalidRecipeManifest/project")
		repositoryUrl := filepath.FromSlash("testdata/WatchSuite/TestRecipeErrors/InvalidRecipeManifest/repository")

		err := s.executor.Execute([]string{
			projectDir,
			"--repository", repositoryUrl,
			"--recipe", "recipe",
		})

		s.Empty(s.executor.Stdout)
		s.Empty(s.executor.Stderr)

		serrors.Equal(s.Assert(), &serrors.Assert{
			Type:    serrors.Error{},
			Message: "unable to read recipe manifest",
			Arguments: []any{
				"file", filepath.Join(repositoryUrl, "recipe", ".manala.yaml"),
			},
			Errors: []*serrors.Assert{
				{
					Type:    serrors.Error{},
					Message: "invalid recipe manifest",
					Errors: []*serrors.Assert{
						{
							Type:    serrors.Error{},
							Message: "missing manala description property",
							Arguments: []any{
								"path", "manala",
								"property", "description",
								"line", 1,
								"column", 9,
							},
							Details: `
								>  1 | manala: {}
								               ^
							`,
						},
					},
				},
			},
		}, err)
	})
}
