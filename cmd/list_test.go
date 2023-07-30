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

type ListSuite struct {
	suite.Suite
	configMock *mocks.ConfigMock
	executor   *cmd.Executor
}

func TestListSuite(t *testing.T) {
	suite.Run(t, new(ListSuite))
}

func (s *ListSuite) SetupTest() {
	s.configMock = &mocks.ConfigMock{}
	s.executor = cmd.NewExecutor(func(stdout *bytes.Buffer, stderr *bytes.Buffer) *cobra.Command {
		out := lipgloss.New(stdout, stderr)
		return newListCmd(
			s.configMock,
			slog.New(log.NewSlogHandler(out)),
			out,
		)
	})
}

func (s *ListSuite) TestRepositoryErrors() {
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
		repoUrl := filepath.FromSlash("testdata/ListSuite/TestRepositoryErrors/RepositoryNotFound/repository")

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
		repoUrl := filepath.FromSlash("testdata/ListSuite/TestRepositoryErrors/WrongRepository/repository")

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
		repoUrl := filepath.FromSlash("testdata/ListSuite/TestRepositoryErrors/EmptyRepository/repository")

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

func (s *ListSuite) TestRepositoryCustom() {
	repoUrl := filepath.FromSlash("testdata/ListSuite/TestRepositoryCustom/repository")

	s.configMock.
		On("CacheDir").Return("").
		On("Repository").Return("")

	err := s.executor.Execute([]string{
		"--repository", repoUrl,
	})

	s.NoError(err)

	s.Equal(heredoc.Docf(`
		bar  Bar
		foo  Foo
		`),
		s.executor.Stdout.String())
	s.Empty(s.executor.Stderr)
}

func (s *ListSuite) TestRepositoryConfig() {
	repoUrl := filepath.FromSlash("testdata/ListSuite/TestRepositoryConfig/repository")

	s.configMock.
		On("CacheDir").Return("").
		On("Repository").Return(repoUrl)

	err := s.executor.Execute([]string{})

	s.NoError(err)

	s.Equal(heredoc.Docf(`
		bar  Bar
		foo  Foo
		`),
		s.executor.Stdout.String())
	s.Empty(s.executor.Stderr)
}

func (s *ListSuite) TestRecipeErrors() {
	s.configMock.
		On("CacheDir").Return("").
		On("Repository").Return("")

	s.Run("WrongRecipeManifest", func() {
		repoUrl := filepath.FromSlash("testdata/ListSuite/TestRecipeErrors/WrongRecipeManifest/repository")

		err := s.executor.Execute([]string{
			"--repository", repoUrl,
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
		repoUrl := filepath.FromSlash("testdata/ListSuite/TestRecipeErrors/InvalidRecipeManifest/repository")

		err := s.executor.Execute([]string{
			"--repository", repoUrl,
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
