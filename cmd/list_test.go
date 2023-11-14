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
	"path/filepath"
	"testing"
)

type ListSuite struct {
	suite.Suite
	configMock *config.Mock
	executor   *cmd.Executor
}

func TestListSuite(t *testing.T) {
	suite.Run(t, new(ListSuite))
}

func (s *ListSuite) SetupTest() {
	s.configMock = &config.Mock{}
	s.executor = cmd.NewExecutor(func(stdout *bytes.Buffer, stderr *bytes.Buffer) *cobra.Command {
		ui := charm.New(nil, stdout, stderr)
		return newListCmd(
			s.configMock,
			slog.New(log.NewSlogHandler(ui)),
			ui,
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
			Type:    &app.UnprocessableRepositoryUrlError{},
			Message: "unable to process repository url",
		}, err)
	})

	s.Run("RepositoryNotFound", func() {
		repositoryUrl := filepath.FromSlash("testdata/ListSuite/TestRepositoryErrors/RepositoryNotFound/repository")

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
		repositoryUrl := filepath.FromSlash("testdata/ListSuite/TestRepositoryErrors/WrongRepository/repository")

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
		repositoryUrl := filepath.FromSlash("testdata/ListSuite/TestRepositoryErrors/EmptyRepository/repository")

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

func (s *ListSuite) TestRepositoryCustom() {
	repositoryUrl := filepath.FromSlash("testdata/ListSuite/TestRepositoryCustom/repository")

	s.configMock.
		On("CacheDir").Return("").
		On("Repository").Return("")

	err := s.executor.Execute([]string{
		"--repository", repositoryUrl,
	})

	s.NoError(err)

	heredoc.Equal(s.Assert(), `
		Recipes available in %[1]s
		───────────────────────────────────────────────────────────────────────
		 • bar
		   Bar
		 • foo
		   Foo
		`,
		s.executor.Stdout.String(),
		repositoryUrl,
	)
	s.Empty(s.executor.Stderr)
}

func (s *ListSuite) TestRepositoryConfig() {
	repositoryUrl := filepath.FromSlash("testdata/ListSuite/TestRepositoryConfig/repository")

	s.configMock.
		On("CacheDir").Return("").
		On("Repository").Return(repositoryUrl)

	err := s.executor.Execute([]string{})

	s.NoError(err)

	heredoc.Equal(s.Assert(), `
		Recipes available in %[1]s
		───────────────────────────────────────────────────────────────────────
		 • bar
		   Bar
		 • foo
		   Foo
		`,
		s.executor.Stdout.String(),
		repositoryUrl,
	)

	s.Empty(s.executor.Stderr)
}

func (s *ListSuite) TestRecipeErrors() {
	s.configMock.
		On("CacheDir").Return("").
		On("Repository").Return("")

	s.Run("WrongRecipeManifest", func() {
		repositoryUrl := filepath.FromSlash("testdata/ListSuite/TestRecipeErrors/WrongRecipeManifest/repository")

		err := s.executor.Execute([]string{
			"--repository", repositoryUrl,
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
		repositoryUrl := filepath.FromSlash("testdata/ListSuite/TestRecipeErrors/InvalidRecipeManifest/repository")

		err := s.executor.Execute([]string{
			"--repository", repositoryUrl,
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
