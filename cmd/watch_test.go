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
	"manala/internal/testing/heredoc"
	"manala/internal/ui/log"
	"manala/internal/ui/output/lipgloss"
	"manala/internal/validation"
	"manala/internal/yaml"
	"path/filepath"
	"testing"
)

type WatchSuite struct {
	suite.Suite
	configMock *mocks.ConfigMock
	executor   *cmd.Executor
}

func TestWatchSuite(t *testing.T) {
	suite.Run(t, new(WatchSuite))
}

func (s *WatchSuite) SetupTest() {
	s.configMock = &mocks.ConfigMock{}
	s.executor = cmd.NewExecutor(func(stdout *bytes.Buffer, stderr *bytes.Buffer) *cobra.Command {
		out := lipgloss.New(stdout, stderr)
		return newWatchCmd(
			s.configMock,
			slog.New(log.NewSlogHandler(out)),
			out,
		)
	})
}

func (s *WatchSuite) TestProjectErrors() {
	s.configMock.
		On("CacheDir").Return("").
		On("Repository").Return("")

	s.Run("ProjectNotFound", func() {
		projDir := filepath.FromSlash("testdata/WatchSuite/TestProjectErrors/ProjectNotFound/project")

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
		projDir := filepath.FromSlash("testdata/WatchSuite/TestProjectErrors/WrongProjectManifest/project")

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
		projDir := filepath.FromSlash("testdata/WatchSuite/TestProjectErrors/EmptyProjectManifest/project")

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
		projDir := filepath.FromSlash("testdata/WatchSuite/TestProjectErrors/InvalidProjectManifest/project")

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

func (s *WatchSuite) TestRepositoryErrors() {
	s.configMock.
		On("CacheDir").Return("").
		On("Repository").Return("")

	s.Run("NoRepository", func() {
		projDir := filepath.FromSlash("testdata/WatchSuite/TestRepositoryErrors/NoRepository/project")

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
		projDir := filepath.FromSlash("testdata/WatchSuite/TestRepositoryErrors/RepositoryNotFound/project")
		repoUrl := filepath.FromSlash("testdata/WatchSuite/TestRepositoryErrors/RepositoryNotFound/repository")

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
		projDir := filepath.FromSlash("testdata/WatchSuite/TestRepositoryErrors/WrongRepository/project")
		repoUrl := filepath.FromSlash("testdata/WatchSuite/TestRepositoryErrors/WrongRepository/repository")

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

func (s *WatchSuite) TestRecipeErrors() {
	s.configMock.
		On("CacheDir").Return("").
		On("Repository").Return("")

	s.Run("RecipeNotFound", func() {
		projDir := filepath.FromSlash("testdata/WatchSuite/TestRecipeErrors/RecipeNotFound/project")
		repoUrl := filepath.FromSlash("testdata/WatchSuite/TestRecipeErrors/RecipeNotFound/repository")

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
		projDir := filepath.FromSlash("testdata/WatchSuite/TestRecipeErrors/WrongRecipeManifest/project")
		repoUrl := filepath.FromSlash("testdata/WatchSuite/TestRecipeErrors/WrongRecipeManifest/repository")

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
		projDir := filepath.FromSlash("testdata/WatchSuite/TestRecipeErrors/InvalidRecipeManifest/project")
		repoUrl := filepath.FromSlash("testdata/WatchSuite/TestRecipeErrors/InvalidRecipeManifest/repository")

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
