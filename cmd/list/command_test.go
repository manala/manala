package list_test

import (
	"bytes"
	"log/slog"
	"path/filepath"
	"testing"

	"github.com/manala/manala/app"
	"github.com/manala/manala/app/api"
	cmdList "github.com/manala/manala/cmd/list"
	"github.com/manala/manala/internal/caching"
	"github.com/manala/manala/internal/serrors"
	"github.com/manala/manala/internal/testing/heredoc"
	"github.com/manala/manala/internal/ui/adapters/charm"
	"github.com/manala/manala/internal/ui/log"

	"github.com/stretchr/testify/suite"
)

type Suite struct{ suite.Suite }

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) TestRepositoryErrors() {
	s.Run("NoRepository", func() {
		stdOut, stdErr, err := s.execute("")

		s.Empty(stdOut)
		heredoc.Equal(s.T(), `
			 • loading repository…
		`, stdErr)

		serrors.Equal(s.T(), &serrors.Assertion{
			Type:    &app.NotFoundRepositoryError{},
			Message: "repository not found",
			Arguments: []any{
				"url", "",
			},
		}, err)
	})

	s.Run("RepositoryNotFound", func() {
		repositoryURL := filepath.FromSlash("testdata/TestRepositoryErrors/RepositoryNotFound/repository")

		stdOut, stdErr, err := s.execute("",
			"--repository", repositoryURL,
		)

		s.Empty(stdOut)
		heredoc.Equal(s.T(), `
			 • loading repository…
		`, stdErr)

		serrors.Equal(s.T(), &serrors.Assertion{
			Type:    &app.NotFoundRepositoryError{},
			Message: "repository not found",
			Arguments: []any{
				"url", repositoryURL,
			},
		}, err)
	})

	s.Run("WrongRepository", func() {
		repositoryURL := filepath.FromSlash("testdata/TestRepositoryErrors/WrongRepository/repository")

		stdOut, stdErr, err := s.execute("",
			"--repository", repositoryURL,
		)

		s.Empty(stdOut)
		heredoc.Equal(s.T(), `
			 • loading repository…
		`, stdErr)

		serrors.Equal(s.T(), &serrors.Assertion{
			Type:    &app.NotFoundRepositoryError{},
			Message: "repository not found",
			Arguments: []any{
				"url", repositoryURL,
			},
		}, err)
	})

	s.Run("EmptyRepository", func() {
		repositoryURL := filepath.FromSlash("testdata/TestRepositoryErrors/EmptyRepository/repository")

		stdOut, stdErr, err := s.execute("",
			"--repository", repositoryURL,
		)

		s.Empty(stdOut)
		heredoc.Equal(s.T(), `
			 • loading repository…
			 • loading recipes…
		`, stdErr)

		serrors.Equal(s.T(), &serrors.Assertion{
			Type:    &app.EmptyRepositoryError{},
			Message: "empty repository",
			Arguments: []any{
				"url", repositoryURL,
			},
		}, err)
	})
}

func (s *Suite) TestRepositoryCustom() {
	repositoryURL := filepath.FromSlash("testdata/TestRepositoryCustom/repository")

	stdOut, stdErr, err := s.execute("",
		"--repository", repositoryURL,
	)

	s.Require().NoError(err)
	heredoc.Equal(s.T(), `
		bar
		  Bar
		foo
		  Foo
	`, stdOut)
	heredoc.Equal(s.T(), `
		 • loading repository…
		 • loading recipes…
	`, stdErr)
}

func (s *Suite) TestRepositoryConfig() {
	repositoryURL := filepath.FromSlash("testdata/TestRepositoryConfig/repository")

	stdOut, stdErr, err := s.execute(repositoryURL)

	s.Require().NoError(err)
	heredoc.Equal(s.T(), `
		bar
		  Bar
		foo
		  Foo
	`, stdOut)
	heredoc.Equal(s.T(), `
		 • loading repository…
		 • loading recipes…
	`, stdErr)
}

func (s *Suite) TestRecipeErrors() {
	s.Run("WrongRecipeManifest", func() {
		repositoryURL := filepath.FromSlash("testdata/TestRecipeErrors/WrongRecipeManifest/repository")

		stdOut, stdErr, err := s.execute("",
			"--repository", repositoryURL,
		)

		s.Empty(stdOut)
		heredoc.Equal(s.T(), `
			 • loading repository…
			 • loading recipes…
		`, stdErr)

		serrors.Equal(s.T(), &serrors.Assertion{
			Message: "recipe manifest is a directory",
			Arguments: []any{
				"dir", filepath.Join(repositoryURL, "recipe", ".manala.yaml"),
			},
		}, err)
	})

	s.Run("InvalidRecipeManifest", func() {
		repositoryURL := filepath.FromSlash("testdata/TestRecipeErrors/InvalidRecipeManifest/repository")

		stdOut, stdErr, err := s.execute("",
			"--repository", repositoryURL,
		)

		s.Empty(stdOut)
		heredoc.Equal(s.T(), `
			 • loading repository…
			 • loading recipes…
		`, stdErr)

		serrors.Equal(s.T(), &serrors.Assertion{
			Message: "unable to read recipe manifest",
			Arguments: []any{
				"file", filepath.Join(repositoryURL, "recipe", ".manala.yaml"),
			},
			Errors: []*serrors.Assertion{
				{
					Message: "invalid recipe manifest",
					Errors: []*serrors.Assertion{
						{
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

func (s *Suite) execute(defaultRepositoryURL string, args ...string) (*bytes.Buffer, *bytes.Buffer, error) {
	stdOut := &bytes.Buffer{}
	stdErr := &bytes.Buffer{}

	ui := charm.New(stdErr)
	log := slog.New(log.NewSlogHandler(ui))

	command := cmdList.NewCommand(
		log,
		api.New(
			log,
			caching.NewCache(""),
			api.WithDefaultRepositoryURL(defaultRepositoryURL),
		),
		stdOut,
	)

	command.SilenceErrors = true
	command.SilenceUsage = true
	command.SetOut(stdOut)
	command.SetErr(stdErr)
	command.SetArgs(append([]string{}, args...))

	return stdOut, stdErr, command.Execute()
}
