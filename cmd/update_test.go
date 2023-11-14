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

type UpdateSuite struct {
	suite.Suite
	configMock *config.Mock
	executor   *cmd.Executor
}

func TestUpdateSuite(t *testing.T) {
	suite.Run(t, new(UpdateSuite))
}

func (s *UpdateSuite) SetupTest() {
	s.configMock = &config.Mock{}
	s.executor = cmd.NewExecutor(func(stdout *bytes.Buffer, stderr *bytes.Buffer) *cobra.Command {
		ui := charm.New(nil, stdout, stderr)
		return newUpdateCmd(
			s.configMock,
			slog.New(log.NewSlogHandler(ui)),
			ui,
		)
	})
}

func (s *UpdateSuite) TestProjectErrors() {
	s.configMock.
		On("CacheDir").Return("").
		On("Repository").Return("")

	s.Run("ProjectNotFound", func() {
		projectDir := filepath.FromSlash("testdata/UpdateSuite/TestProjectErrors/ProjectNotFound/project")

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
		projectDir := filepath.FromSlash("testdata/UpdateSuite/TestProjectErrors/WrongProjectManifest/project")

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
		projectDir := filepath.FromSlash("testdata/UpdateSuite/TestProjectErrors/EmptyProjectManifest/project")

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
		projectDir := filepath.FromSlash("testdata/UpdateSuite/TestProjectErrors/InvalidProjectManifest/project")

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

func (s *UpdateSuite) TestRecursiveProjectErrors() {
	s.configMock.
		On("CacheDir").Return("").
		On("Repository").Return("")

	s.Run("ProjectNotFound", func() {
		projectDir := filepath.FromSlash("testdata/UpdateSuite/TestRecursiveProjectErrors/ProjectNotFound/project")

		err := s.executor.Execute([]string{
			projectDir,
			"--recursive",
		})

		s.NoError(err)

		s.Empty(s.executor.Stdout)
		heredoc.Equal(s.Assert(), `
			 • walk projects from…
			   dir=%[1]s
			`,
			s.executor.Stderr.String(),
			projectDir,
		)
	})

	s.Run("WrongProjectManifest", func() {
		projectDir := filepath.FromSlash("testdata/UpdateSuite/TestRecursiveProjectErrors/WrongProjectManifest/project")

		err := s.executor.Execute([]string{
			projectDir,
			"--recursive",
		})

		s.Empty(s.executor.Stdout)
		heredoc.Equal(s.Assert(), `
			 • walk projects from…
			   dir=%[1]s
			`,
			s.executor.Stderr.String(),
			projectDir,
		)

		serrors.Equal(s.Assert(), &serrors.Assert{
			Type:    serrors.Error{},
			Message: "project manifest is a directory",
			Arguments: []any{
				"dir", filepath.Join(projectDir, ".manala.yaml"),
			},
		}, err)
	})

	s.Run("EmptyProjectManifest", func() {
		projectDir := filepath.FromSlash("testdata/UpdateSuite/TestRecursiveProjectErrors/EmptyProjectManifest/project")

		err := s.executor.Execute([]string{
			projectDir,
			"--recursive",
		})

		s.Empty(s.executor.Stdout)
		heredoc.Equal(s.Assert(), `
			 • walk projects from…
			   dir=%[1]s
			`,
			s.executor.Stderr.String(),
			projectDir,
		)

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
		projectDir := filepath.FromSlash("testdata/UpdateSuite/TestRecursiveProjectErrors/InvalidProjectManifest/project")

		err := s.executor.Execute([]string{
			projectDir,
			"--recursive",
		})

		s.Empty(s.executor.Stdout)
		heredoc.Equal(s.Assert(), `
			 • walk projects from…
			   dir=%[1]s
			`,
			s.executor.Stderr.String(),
			projectDir,
		)

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

func (s *UpdateSuite) TestRepositoryErrors() {
	s.configMock.
		On("CacheDir").Return("").
		On("Repository").Return("")

	s.Run("NoRepository", func() {
		projectDir := filepath.FromSlash("testdata/UpdateSuite/TestRepositoryErrors/NoRepository/project")

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
		projectDir := filepath.FromSlash("testdata/UpdateSuite/TestRepositoryErrors/RepositoryNotFound/project")
		repositoryUrl := filepath.FromSlash("testdata/UpdateSuite/TestRepositoryErrors/RepositoryNotFound/repository")

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
		projectDir := filepath.FromSlash("testdata/UpdateSuite/TestRepositoryErrors/WrongRepository/project")
		repositoryUrl := filepath.FromSlash("testdata/UpdateSuite/TestRepositoryErrors/WrongRepository/repository")

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

func (s *UpdateSuite) TestRepositoryCustom() {
	projectDir := filepath.FromSlash("testdata/UpdateSuite/TestRepositoryCustom/project")
	repositoryUrl := filepath.FromSlash("testdata/UpdateSuite/TestRepositoryCustom/repository")

	_ = os.Remove(filepath.Join(projectDir, "file.txt"))

	s.configMock.
		On("CacheDir").Return("").
		On("Repository").Return("")

	err := s.executor.Execute([]string{
		projectDir,
		"--repository", repositoryUrl,
	})

	s.NoError(err)

	s.Empty(s.executor.Stdout)
	heredoc.Equal(s.Assert(), `
		 • sync project…
		   src=%[1]s dst=%[2]s
		 • file synced
		   path=file.txt
		`,
		s.executor.Stderr.String(),
		filepath.Join(repositoryUrl, "recipe"),
		projectDir,
	)

	heredoc.EqualFile(s.Assert(), `
		File
		`,
		filepath.Join(projectDir, "file.txt"),
	)
}

func (s *UpdateSuite) TestRepositoryConfig() {
	projectDir := filepath.FromSlash("testdata/UpdateSuite/TestRepositoryConfig/project")
	repositoryUrl := filepath.FromSlash("testdata/UpdateSuite/TestRepositoryConfig/repository")

	_ = os.Remove(filepath.Join(projectDir, "file.txt"))
	_ = os.Remove(filepath.Join(projectDir, "template"))

	s.configMock.
		On("CacheDir").Return("").
		On("Repository").Return(repositoryUrl)

	err := s.executor.Execute([]string{
		projectDir,
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

func (s *UpdateSuite) TestRecipeErrors() {
	s.configMock.
		On("CacheDir").Return("").
		On("Repository").Return("")

	s.Run("RecipeNotFound", func() {
		projectDir := filepath.FromSlash("testdata/UpdateSuite/TestRecipeErrors/RecipeNotFound/project")
		repositoryUrl := filepath.FromSlash("testdata/UpdateSuite/TestRecipeErrors/RecipeNotFound/repository")

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
		projectDir := filepath.FromSlash("testdata/UpdateSuite/TestRecipeErrors/WrongRecipeManifest/project")
		repositoryUrl := filepath.FromSlash("testdata/UpdateSuite/TestRecipeErrors/WrongRecipeManifest/repository")

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
		projectDir := filepath.FromSlash("testdata/UpdateSuite/TestRecipeErrors/InvalidRecipeManifest/project")
		repositoryUrl := filepath.FromSlash("testdata/UpdateSuite/TestRecipeErrors/InvalidRecipeManifest/repository")

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

func (s *UpdateSuite) TestRecipeCustom() {
	projectDir := filepath.FromSlash("testdata/UpdateSuite/TestRecipeCustom/project")
	repositoryUrl := filepath.FromSlash("testdata/UpdateSuite/TestRecipeCustom/repository")

	_ = os.Remove(filepath.Join(projectDir, "file.txt"))

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
		`,
		s.executor.Stderr.String(),
		filepath.Join(repositoryUrl, "recipe"),
		projectDir,
	)

	heredoc.EqualFile(s.Assert(), `
		File
		`,
		filepath.Join(projectDir, "file.txt"),
	)
}
