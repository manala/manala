package list_test

import (
	"bytes"
	"log/slog"
	cmd "manala/cmd/list"
	"path/filepath"
	"testing"

	"manala/app"
	"manala/app/api"
	"manala/internal/caching"
	"manala/internal/serrors"
	"manala/internal/testing/heredoc"
	"manala/internal/ui/adapters/charm"
	"manala/internal/ui/log"

	"github.com/stretchr/testify/suite"
)

type Suite struct{ suite.Suite }

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) execute(defaultRepositoryURL string, args ...string) (*bytes.Buffer, *bytes.Buffer, error) {
	stdOut := &bytes.Buffer{}
	stdErr := &bytes.Buffer{}

	ui := charm.New(nil, stdOut, stdErr)
	log := slog.New(log.NewSlogHandler(ui))

	cmd := cmd.NewCmd(
		log,
		api.New(
			log,
			caching.NewCache(""),
			api.WithDefaultRepositoryURL(defaultRepositoryURL),
		),
		ui,
	)

	cmd.SilenceErrors = true
	cmd.SilenceUsage = true
	cmd.SetOut(stdOut)
	cmd.SetErr(stdErr)
	cmd.SetArgs(append([]string{}, args...))

	return stdOut, stdErr, cmd.Execute()
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
		Recipes available in %[1]s
		─────────────────────────────────────────────────────────────
		 • bar
		   Bar
		 • foo
		   Foo
	`, stdOut, repositoryURL)
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
		Recipes available in %[1]s
		─────────────────────────────────────────────────────────────
		 • bar
		   Bar
		 • foo
		   Foo
	`, stdOut, repositoryURL)
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
