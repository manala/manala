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

type InitSuite struct {
	suite.Suite
	configMock *mocks.ConfigMock
	executor   *cmd.Executor
}

func TestInitSuite(t *testing.T) {
	suite.Run(t, new(InitSuite))
}

func (s *InitSuite) SetupTest() {
	s.configMock = &mocks.ConfigMock{}
	s.executor = cmd.NewExecutor(func(stdout *bytes.Buffer, stderr *bytes.Buffer) *cobra.Command {
		out := lipgloss.New(stdout, stderr)
		return newInitCmd(
			s.configMock,
			slog.New(log.NewSlogHandler(out)),
			out,
		)
	})
}

func (s *InitSuite) TestProjectErrors() {
	s.configMock.
		On("CacheDir").Return("").
		On("Repository").Return("")

	s.Run("AlreadyExistingProject", func() {
		projDir := filepath.FromSlash("testdata/InitSuite/TestProjectErrors/AlreadyExistingProject/project")

		err := s.executor.Execute([]string{
			projDir,
		})

		s.Empty(s.executor.Stdout)
		s.Empty(s.executor.Stderr)

		serrors.Equal(s.Assert(), &serrors.Assert{
			Type:    &core.AlreadyExistingProjectError{},
			Message: "already existing project",
			Arguments: []any{
				"dir", projDir,
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
			Type:    &core.UnprocessableRepositoryUrlError{},
			Message: "unable to process repository url",
		}, err)
	})

	s.Run("RepositoryNotFound", func() {
		repoUrl := filepath.FromSlash("testdata/InitSuite/TestRepositoryErrors/RepositoryNotFound/repository")

		err := s.executor.Execute([]string{
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
		repoUrl := filepath.FromSlash("testdata/InitSuite/TestRepositoryErrors/WrongRepository/repository")

		err := s.executor.Execute([]string{
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

	s.Run("EmptyRepository", func() {
		repoUrl := filepath.FromSlash("testdata/InitSuite/TestRepositoryErrors/EmptyRepository/repository")

		err := s.executor.Execute([]string{
			"--repository", repoUrl,
		})

		s.Empty(s.executor.Stdout)
		s.Empty(s.executor.Stderr)

		serrors.Equal(s.Assert(), &serrors.Assert{
			Type:    &serrors.Error{},
			Message: "empty repository",
			Arguments: []any{
				"dir", repoUrl,
			},
		}, err)
	})
}

func (s *InitSuite) TestRepositoryCustom() {
	projDir := filepath.FromSlash("testdata/InitSuite/TestRepositoryCustom/project")
	repoUrl := filepath.FromSlash("testdata/InitSuite/TestRepositoryCustom/repository")

	_ = os.RemoveAll(projDir)

	s.configMock.
		On("CacheDir").Return("").
		On("Repository").Return("")

	err := s.executor.Execute([]string{
		projDir,
		"--repository", repoUrl,
		"--recipe", "recipe",
	})

	s.NoError(err)

	s.Empty(s.executor.Stdout)
	s.Equal(heredoc.Docf(`
		  • sync project…                      src=%[1]s dst=%[2]s
		`, filepath.Join(repoUrl, "recipe"), projDir,
	), s.executor.Stderr.String())

	s.DirExists(projDir)

	file.EqualContent(s.Assert(), heredoc.Docf(`
		####################################################################
		#                         !!! REMINDER !!!                         #
		# Don't forget to run `+"`"+`manala up`+"`"+` each time you update this file ! #
		####################################################################

		manala:
		    recipe: recipe
		    repository: %[1]s
		`, repoUrl),
		filepath.Join(projDir, ".manala.yaml"),
	)
}

func (s *InitSuite) TestRepositoryConfig() {
	projDir := filepath.FromSlash("testdata/InitSuite/TestRepositoryConfig/project")
	repoUrl := filepath.FromSlash("testdata/InitSuite/TestRepositoryConfig/repository")

	_ = os.RemoveAll(projDir)

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
		  • file synced                        path=template
		`, filepath.Join(repoUrl, "recipe"), projDir,
	), s.executor.Stderr.String())

	s.DirExists(projDir)

	file.EqualContent(s.Assert(), heredoc.Docf(`
			manala:
			    recipe: recipe
			    repository: %[1]s

			foo: bar
		`, repoUrl),
		filepath.Join(projDir, ".manala.yaml"),
	)

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

func (s *InitSuite) TestRecipeErrors() {
	s.configMock.
		On("CacheDir").Return("").
		On("Repository").Return("")

	s.Run("RecipeNotFound", func() {
		repoUrl := filepath.FromSlash("testdata/InitSuite/TestRecipeErrors/RecipeNotFound/repository")

		err := s.executor.Execute([]string{
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
		repoUrl := filepath.FromSlash("testdata/InitSuite/TestRecipeErrors/WrongRecipeManifest/repository")

		err := s.executor.Execute([]string{
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
		repoUrl := filepath.FromSlash("testdata/InitSuite/TestRecipeErrors/InvalidRecipeManifest/repository")

		err := s.executor.Execute([]string{
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
