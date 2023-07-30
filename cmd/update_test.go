package cmd

import (
	"bytes"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/suite"
	"log/slog"
	"manala/app/mocks"
	"manala/core"
	"manala/internal/errors/serrors"
	"manala/internal/testing/cmd"
	"manala/internal/testing/file"
	"manala/internal/testing/heredoc"
	"manala/internal/ui/log"
	"manala/internal/ui/output/lipgloss"
	"manala/internal/validation"
	"manala/internal/yaml"
	"os"
	"path/filepath"
	"testing"
)

type UpdateSuite struct {
	suite.Suite
	configMock *mocks.ConfigMock
	executor   *cmd.Executor
}

func TestUpdateSuite(t *testing.T) {
	suite.Run(t, new(UpdateSuite))
}

func (s *UpdateSuite) SetupTest() {
	s.configMock = &mocks.ConfigMock{}
	s.executor = cmd.NewExecutor(func(stdout *bytes.Buffer, stderr *bytes.Buffer) *cobra.Command {
		out := lipgloss.New(stdout, stderr)
		return newUpdateCmd(
			s.configMock,
			slog.New(log.NewSlogHandler(out)),
			out,
		)
	})
}

func (s *UpdateSuite) TestProjectErrors() {
	s.configMock.
		On("CacheDir").Return("").
		On("Repository").Return("")

	s.Run("ProjectNotFound", func() {
		projDir := filepath.FromSlash("testdata/UpdateSuite/TestProjectErrors/ProjectNotFound/project")

		err := s.executor.Execute([]string{
			projDir,
		})

		s.Empty(s.executor.Stdout)
		s.Empty(s.executor.Stderr)

		serrors.Equal(s.Assert(), &serrors.Assert{
			Type:    &serrors.Error{},
			Message: "project not found",
			Arguments: []any{
				"dir", projDir,
			},
		}, err)
	})

	s.Run("WrongProjectManifest", func() {
		projDir := filepath.FromSlash("testdata/UpdateSuite/TestProjectErrors/WrongProjectManifest/project")

		err := s.executor.Execute([]string{
			projDir,
		})

		s.Empty(s.executor.Stdout)
		s.Empty(s.executor.Stderr)

		serrors.Equal(s.Assert(), &serrors.Assert{
			Type:    &serrors.Error{},
			Message: "project manifest is a directory",
			Arguments: []any{
				"dir", filepath.Join(projDir, ".manala.yaml"),
			},
		}, err)
	})

	s.Run("EmptyProjectManifest", func() {
		projDir := filepath.FromSlash("testdata/UpdateSuite/TestProjectErrors/EmptyProjectManifest/project")

		err := s.executor.Execute([]string{
			projDir,
		})

		s.Empty(s.executor.Stdout)
		s.Empty(s.executor.Stderr)

		serrors.Equal(s.Assert(), &serrors.Assert{
			Type:    &serrors.WrapError{},
			Message: "unable to read project manifest",
			Arguments: []any{
				"file", filepath.Join(projDir, ".manala.yaml"),
			},
			Error: &serrors.Assert{
				Type:    &serrors.WrapError{},
				Message: "irregular project manifest",
				Error: &serrors.Assert{
					Type:    &serrors.Error{},
					Message: "empty yaml file",
				},
			},
		}, err)
	})

	s.Run("InvalidProjectManifest", func() {
		projDir := filepath.FromSlash("testdata/UpdateSuite/TestProjectErrors/InvalidProjectManifest/project")

		err := s.executor.Execute([]string{
			projDir,
		})

		s.Empty(s.executor.Stdout)
		s.Empty(s.executor.Stderr)

		serrors.Equal(s.Assert(), &serrors.Assert{
			Type:    &serrors.WrapError{},
			Message: "unable to read project manifest",
			Arguments: []any{
				"file", filepath.Join(projDir, ".manala.yaml"),
			},
			Error: &serrors.Assert{
				Type:    &validation.Error{},
				Message: "invalid project manifest",
				Errors: []*serrors.Assert{
					{
						Type:    &yaml.NodeValidationResultError{},
						Message: "missing manala recipe field",
						Arguments: []any{
							"property", "recipe",
							"line", 1,
							"column", 9,
						},
						Details: heredoc.Doc(`
							>  1 | manala: {}
							               ^
						`),
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
		projDir := filepath.FromSlash("testdata/UpdateSuite/TestRecursiveProjectErrors/ProjectNotFound/project")

		err := s.executor.Execute([]string{
			projDir,
			"--recursive",
		})

		s.NoError(err)

		s.Empty(s.executor.Stdout)
		s.Equal(heredoc.Docf(`
			  • walk projects from…                dir=%[1]s
		`, projDir,
		), s.executor.Stderr.String())
	})

	s.Run("WrongProjectManifest", func() {
		projDir := filepath.FromSlash("testdata/UpdateSuite/TestRecursiveProjectErrors/WrongProjectManifest/project")

		err := s.executor.Execute([]string{
			projDir,
			"--recursive",
		})

		s.Empty(s.executor.Stdout)
		s.Equal(heredoc.Docf(`
			  • walk projects from…                dir=%[1]s
		`, projDir,
		), s.executor.Stderr.String())

		serrors.Equal(s.Assert(), &serrors.Assert{
			Type:    &serrors.Error{},
			Message: "project manifest is a directory",
			Arguments: []any{
				"dir", filepath.Join(projDir, ".manala.yaml"),
			},
		}, err)
	})

	s.Run("EmptyProjectManifest", func() {
		projDir := filepath.FromSlash("testdata/UpdateSuite/TestRecursiveProjectErrors/EmptyProjectManifest/project")

		err := s.executor.Execute([]string{
			projDir,
			"--recursive",
		})

		s.Empty(s.executor.Stdout)
		s.Equal(heredoc.Docf(`
			  • walk projects from…                dir=%[1]s
		`, projDir,
		), s.executor.Stderr.String())

		serrors.Equal(s.Assert(), &serrors.Assert{
			Type:    &serrors.WrapError{},
			Message: "unable to read project manifest",
			Arguments: []any{
				"file", filepath.Join(projDir, ".manala.yaml"),
			},
			Error: &serrors.Assert{
				Type:    &serrors.WrapError{},
				Message: "irregular project manifest",
				Error: &serrors.Assert{
					Type:    &serrors.Error{},
					Message: "empty yaml file",
				},
			},
		}, err)
	})

	s.Run("InvalidProjectManifest", func() {
		projDir := filepath.FromSlash("testdata/UpdateSuite/TestRecursiveProjectErrors/InvalidProjectManifest/project")

		err := s.executor.Execute([]string{
			projDir,
			"--recursive",
		})

		s.Empty(s.executor.Stdout)
		s.Equal(heredoc.Docf(`
			  • walk projects from…                dir=%[1]s
		`, projDir,
		), s.executor.Stderr.String())

		serrors.Equal(s.Assert(), &serrors.Assert{
			Type:    &serrors.WrapError{},
			Message: "unable to read project manifest",
			Arguments: []any{
				"file", filepath.Join(projDir, ".manala.yaml"),
			},
			Error: &serrors.Assert{
				Type:    &validation.Error{},
				Message: "invalid project manifest",
				Errors: []*serrors.Assert{
					{
						Type:    &yaml.NodeValidationResultError{},
						Message: "missing manala recipe field",
						Arguments: []any{
							"property", "recipe",
							"line", 1,
							"column", 9,
						},
						Details: heredoc.Doc(`
							>  1 | manala: {}
							               ^
						`),
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
		projDir := filepath.FromSlash("testdata/UpdateSuite/TestRepositoryErrors/NoRepository/project")

		err := s.executor.Execute([]string{
			projDir,
		})

		s.Empty(s.executor.Stdout)
		s.Empty(s.executor.Stderr)

		serrors.Equal(s.Assert(), &serrors.Assert{
			Type:    &core.UnprocessableRepositoryUrlError{},
			Message: "unable to process repository url",
		}, err)
	})

	s.Run("RepositoryNotFound", func() {
		projDir := filepath.FromSlash("testdata/UpdateSuite/TestRepositoryErrors/RepositoryNotFound/project")
		repoUrl := filepath.FromSlash("testdata/UpdateSuite/TestRepositoryErrors/RepositoryNotFound/repository")

		err := s.executor.Execute([]string{
			projDir,
			"--repository", repoUrl,
		})

		s.Empty(s.executor.Stdout)
		s.Empty(s.executor.Stderr)

		serrors.Equal(s.Assert(), &serrors.Assert{
			Type:    &core.UnsupportedRepositoryError{},
			Message: "unsupported repository url",
			Arguments: []any{
				"url", repoUrl,
			},
		}, err)
	})

	s.Run("WrongRepository", func() {
		projDir := filepath.FromSlash("testdata/UpdateSuite/TestRepositoryErrors/WrongRepository/project")
		repoUrl := filepath.FromSlash("testdata/UpdateSuite/TestRepositoryErrors/WrongRepository/repository")

		err := s.executor.Execute([]string{
			projDir,
			"--repository", repoUrl,
		})

		s.Empty(s.executor.Stdout)
		s.Empty(s.executor.Stderr)

		serrors.Equal(s.Assert(), &serrors.Assert{
			Type:    &core.UnsupportedRepositoryError{},
			Message: "unsupported repository url",
			Arguments: []any{
				"url", repoUrl,
			},
		}, err)
	})
}

func (s *UpdateSuite) TestRepositoryCustom() {
	projDir := filepath.FromSlash("testdata/UpdateSuite/TestRepositoryCustom/project")
	repoUrl := filepath.FromSlash("testdata/UpdateSuite/TestRepositoryCustom/repository")

	_ = os.Remove(filepath.Join(projDir, "file.txt"))

	s.configMock.
		On("CacheDir").Return("").
		On("Repository").Return("")

	err := s.executor.Execute([]string{
		projDir,
		"--repository", repoUrl,
	})

	s.NoError(err)

	s.Empty(s.executor.Stdout)
	s.Equal(heredoc.Docf(`
		  • sync project…                      src=%[1]s dst=%[2]s
		  • file synced                        path=file.txt
		`, filepath.Join(repoUrl, "recipe"), projDir,
	), s.executor.Stderr.String())

	file.EqualContent(s.Assert(), heredoc.Docf(`
		File
		`),
		filepath.Join(projDir, "file.txt"),
	)
}

func (s *UpdateSuite) TestRepositoryConfig() {
	projDir := filepath.FromSlash("testdata/UpdateSuite/TestRepositoryConfig/project")
	repoUrl := filepath.FromSlash("testdata/UpdateSuite/TestRepositoryConfig/repository")

	_ = os.Remove(filepath.Join(projDir, "file.txt"))
	_ = os.Remove(filepath.Join(projDir, "template"))

	s.configMock.
		On("CacheDir").Return("").
		On("Repository").Return(repoUrl)

	err := s.executor.Execute([]string{
		projDir,
	})

	s.NoError(err)

	s.Empty(s.executor.Stdout)
	s.Equal(heredoc.Docf(`
		  • sync project…                      src=%[1]s dst=%[2]s
		  • file synced                        path=file.txt
		  • file synced                        path=template
		`, filepath.Join(repoUrl, "recipe"), projDir,
	), s.executor.Stderr.String())

	file.EqualContent(s.Assert(), heredoc.Docf(`
		File
		`),
		filepath.Join(projDir, "file.txt"),
	)

	file.EqualContent(s.Assert(), heredoc.Doc(`
		Template

		foo: bar
		`),
		filepath.Join(projDir, "template"),
	)
}

func (s *UpdateSuite) TestRecipeErrors() {
	s.configMock.
		On("CacheDir").Return("").
		On("Repository").Return("")

	s.Run("RecipeNotFound", func() {
		projDir := filepath.FromSlash("testdata/UpdateSuite/TestRecipeErrors/RecipeNotFound/project")
		repoUrl := filepath.FromSlash("testdata/UpdateSuite/TestRecipeErrors/RecipeNotFound/repository")

		err := s.executor.Execute([]string{
			projDir,
			"--repository", repoUrl,
			"--recipe", "recipe",
		})

		s.Empty(s.executor.Stdout)
		s.Empty(s.executor.Stderr)

		serrors.Equal(s.Assert(), &serrors.Assert{
			Type:    &core.NotFoundRecipeManifestError{},
			Message: "recipe manifest not found",
			Arguments: []any{
				"file", filepath.Join(repoUrl, "recipe", ".manala.yaml"),
			},
		}, err)
	})

	s.Run("WrongRecipeManifest", func() {
		projDir := filepath.FromSlash("testdata/UpdateSuite/TestRecipeErrors/WrongRecipeManifest/project")
		repoUrl := filepath.FromSlash("testdata/UpdateSuite/TestRecipeErrors/WrongRecipeManifest/repository")

		err := s.executor.Execute([]string{
			projDir,
			"--repository", repoUrl,
			"--recipe", "recipe",
		})

		s.Empty(s.executor.Stdout)
		s.Empty(s.executor.Stderr)

		serrors.Equal(s.Assert(), &serrors.Assert{
			Type:    &serrors.Error{},
			Message: "recipe manifest is a directory",
			Arguments: []any{
				"dir", filepath.Join(repoUrl, "recipe", ".manala.yaml"),
			},
		}, err)
	})

	s.Run("InvalidRecipeManifest", func() {
		projDir := filepath.FromSlash("testdata/UpdateSuite/TestRecipeErrors/InvalidRecipeManifest/project")
		repoUrl := filepath.FromSlash("testdata/UpdateSuite/TestRecipeErrors/InvalidRecipeManifest/repository")

		err := s.executor.Execute([]string{
			projDir,
			"--repository", repoUrl,
			"--recipe", "recipe",
		})

		s.Empty(s.executor.Stdout)
		s.Empty(s.executor.Stderr)

		serrors.Equal(s.Assert(), &serrors.Assert{
			Type:    &serrors.WrapError{},
			Message: "unable to read recipe manifest",
			Arguments: []any{
				"file", filepath.Join(repoUrl, "recipe", ".manala.yaml"),
			},
			Error: &serrors.Assert{
				Type:    &validation.Error{},
				Message: "invalid recipe manifest",
				Errors: []*serrors.Assert{
					{
						Type:    &yaml.NodeValidationResultError{},
						Message: "missing manala description field",
						Arguments: []any{
							"property", "description",
							"line", 1,
							"column", 9,
						},
						Details: heredoc.Doc(`
							>  1 | manala: {}
							               ^
						`),
					},
				},
			},
		}, err)
	})
}

func (s *UpdateSuite) TestRecipeCustom() {
	projDir := filepath.FromSlash("testdata/UpdateSuite/TestRecipeCustom/project")
	repoUrl := filepath.FromSlash("testdata/UpdateSuite/TestRecipeCustom/repository")

	_ = os.Remove(filepath.Join(projDir, "file.txt"))

	s.configMock.
		On("CacheDir").Return("").
		On("Repository").Return(repoUrl)

	err := s.executor.Execute([]string{
		projDir,
		"--recipe", "recipe",
	})

	s.NoError(err)

	s.Empty(s.executor.Stdout)
	s.Equal(heredoc.Docf(`
		  • sync project…                      src=%[1]s dst=%[2]s
		  • file synced                        path=file.txt
		`, filepath.Join(repoUrl, "recipe"), projDir,
	), s.executor.Stderr.String())

	file.EqualContent(s.Assert(), heredoc.Docf(`
		File
		`),
		filepath.Join(projDir, "file.txt"),
	)
}
