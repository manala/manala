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
	"manala/internal/testing/heredoc"
	"manala/internal/ui/adapters/charm"
	"manala/internal/ui/log"
	"os"
	"path/filepath"
	"testing"
)

type InitSuite struct {
	suite.Suite
	configMock *config.Mock
	executor   *cmd.Executor
}

func TestInitSuite(t *testing.T) {
	suite.Run(t, new(InitSuite))
}

func (s *InitSuite) SetupTest() {
	s.configMock = &config.Mock{}
	s.executor = cmd.NewExecutor(func(stdout *bytes.Buffer, stderr *bytes.Buffer) *cobra.Command {
		ui := charm.New(nil, stdout, stderr)
		return newInitCmd(
			s.configMock,
			slog.New(log.NewSlogHandler(ui)),
			ui,
			ui,
		)
	})
}

func (s *InitSuite) TestProjectErrors() {
	s.configMock.
		On("CacheDir").Return("").
		On("Repository").Return("")

	s.Run("AlreadyExistingProject", func() {
		projectDir := filepath.FromSlash("testdata/InitSuite/TestProjectErrors/AlreadyExistingProject/project")

		err := s.executor.Execute([]string{
			projectDir,
		})

		s.Empty(s.executor.Stdout)
		s.Empty(s.executor.Stderr)

		serrors.Equal(s.Assert(), &serrors.Assert{
			Type:    &app.AlreadyExistingProjectError{},
			Message: "already existing project",
			Arguments: []any{
				"dir", projectDir,
			},
		}, err)
	})
}

func (s *InitSuite) TestRepositoryErrors() {
	s.configMock.
		On("CacheDir").Return("").
		On("Repository").Return("")

	s.Run("NoRepository", func() {
		err := s.executor.Execute([]string{})

		s.Empty(s.executor.Stdout)
		s.Empty(s.executor.Stderr)

		serrors.Equal(s.Assert(), &serrors.Assert{
			Type:    &app.UnprocessableRepositoryUrlError{},
			Message: "unable to process repository url",
		}, err)
	})

	s.Run("RepositoryNotFound", func() {
		repositoryUrl := filepath.FromSlash("testdata/InitSuite/TestRepositoryErrors/RepositoryNotFound/repository")

		err := s.executor.Execute([]string{
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
		repositoryUrl := filepath.FromSlash("testdata/InitSuite/TestRepositoryErrors/WrongRepository/repository")

		err := s.executor.Execute([]string{
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

	s.Run("EmptyRepository", func() {
		repositoryUrl := filepath.FromSlash("testdata/InitSuite/TestRepositoryErrors/EmptyRepository/repository")

		err := s.executor.Execute([]string{
			"--repository", repositoryUrl,
		})

		s.Empty(s.executor.Stdout)
		s.Empty(s.executor.Stderr)

		serrors.Equal(s.Assert(), &serrors.Assert{
			Type:    &app.EmptyRepositoryError{},
			Message: "empty repository",
			Arguments: []any{
				"url", repositoryUrl,
			},
		}, err)
	})
}

func (s *InitSuite) TestRepositoryCustom() {
	projectDir := filepath.FromSlash("testdata/InitSuite/TestRepositoryCustom/project")
	repositoryUrl := filepath.FromSlash("testdata/InitSuite/TestRepositoryCustom/repository")

	_ = os.RemoveAll(projectDir)

	s.configMock.
		On("CacheDir").Return("").
		On("Repository").Return("")

	err := s.executor.Execute([]string{
		projectDir,
		"--repository", repositoryUrl,
		"--recipe", "recipe",
	})

	s.NoError(err)

	s.Empty(s.executor.Stdout)
	heredoc.Equal(s.Assert(), `
		 • sync project…
		   src=%[1]s dst=%[2]s
		`,
		s.executor.Stderr.String(),
		filepath.Join(repositoryUrl, "recipe"),
		projectDir,
	)

	s.DirExists(projectDir)

	heredoc.EqualFile(s.Assert(), `
		####################################################################
		#                         !!! REMINDER !!!                         #
		# Don't forget to run `+"`"+`manala up`+"`"+` each time you update this file ! #
		####################################################################

		manala:
		    recipe: recipe
		    repository: %[1]s
		`,
		filepath.Join(projectDir, ".manala.yaml"),
		repositoryUrl,
	)
}

func (s *InitSuite) TestRepositoryConfig() {
	projectDir := filepath.FromSlash("testdata/InitSuite/TestRepositoryConfig/project")
	repositoryUrl := filepath.FromSlash("testdata/InitSuite/TestRepositoryConfig/repository")

	_ = os.RemoveAll(projectDir)

	s.configMock.
		On("CacheDir").Return("").
		On("Repository").Return(repositoryUrl)

	err := s.executor.Execute([]string{
		projectDir,
		"--recipe", "recipe",
	})

	s.NoError(err)

	s.Empty(s.executor.Stdout)
	heredoc.Equal(s.Assert(), `
		 • sync project…
		   src=%[1]s dst=%[2]s
		 • file synced
		   path=file.txt
		 • file synced
		   path=template
		`,
		s.executor.Stderr.String(),
		filepath.Join(repositoryUrl, "recipe"),
		projectDir,
	)

	s.DirExists(projectDir)

	heredoc.EqualFile(s.Assert(), `
		manala:
		    recipe: recipe
		    repository: %[1]s

		foo: bar
		`,
		filepath.Join(projectDir, ".manala.yaml"),
		repositoryUrl,
	)

	heredoc.EqualFile(s.Assert(), `
		File
		`,
		filepath.Join(projectDir, "file.txt"),
	)

	heredoc.EqualFile(s.Assert(), `
		Template

		foo: bar
		`,
		filepath.Join(projectDir, "template"),
	)
}

func (s *InitSuite) TestRecipeErrors() {
	s.configMock.
		On("CacheDir").Return("").
		On("Repository").Return("")

	s.Run("RecipeNotFound", func() {
		repositoryUrl := filepath.FromSlash("testdata/InitSuite/TestRecipeErrors/RecipeNotFound/repository")

		err := s.executor.Execute([]string{
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
		repositoryUrl := filepath.FromSlash("testdata/InitSuite/TestRecipeErrors/WrongRecipeManifest/repository")

		err := s.executor.Execute([]string{
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
		repositoryUrl := filepath.FromSlash("testdata/InitSuite/TestRecipeErrors/InvalidRecipeManifest/repository")

		err := s.executor.Execute([]string{
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
