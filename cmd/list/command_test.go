package list_test

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/manala/manala/app"
	"github.com/manala/manala/app/api"
	"github.com/manala/manala/app/testing/errors"
	cmdList "github.com/manala/manala/cmd/list"
	"github.com/manala/manala/internal/caching"
	"github.com/manala/manala/internal/log"
	"github.com/manala/manala/internal/serrors"
	"github.com/manala/manala/internal/testing/expect"
	"github.com/manala/manala/internal/testing/heredoc"

	"github.com/stretchr/testify/suite"
)

type CommandSuite struct{ suite.Suite }

func TestCommandSuite(t *testing.T) {
	suite.Run(t, new(CommandSuite))
}

func (s *CommandSuite) TestRepositoryErrors() {
	s.Run("NoRepository", func() {
		stdOut, stdErr, err := s.execute("")

		s.Empty(stdOut)
		heredoc.Equal(s.T(), `
			 • loading repository…
		`, stdErr)

		expect.Error(s.T(), errors.Expectation{
			Type: &app.NotFoundRepositoryError{},
			Attrs: [][2]any{
				{"url", ""},
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

		expect.Error(s.T(), errors.Expectation{
			Type: &app.NotFoundRepositoryError{},
			Attrs: [][2]any{
				{"url", repositoryURL},
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

		expect.Error(s.T(), errors.Expectation{
			Type: &app.NotFoundRepositoryError{},
			Attrs: [][2]any{
				{"url", repositoryURL},
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

		expect.Error(s.T(), errors.Expectation{
			Type: &app.EmptyRepositoryError{},
			Attrs: [][2]any{
				{"url", repositoryURL},
			},
		}, err)
	})
}

func (s *CommandSuite) TestRepositoryCustom() {
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

func (s *CommandSuite) TestRepositoryConfig() {
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

func (s *CommandSuite) TestRecipeErrors() {
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

		expect.Error(s.T(), serrors.Expectation{
			Message: "recipe manifest is a directory",
			Attrs: [][2]any{
				{"dir", filepath.Join(repositoryURL, "recipe", ".manala.yaml")},
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

		expect.Error(s.T(), serrors.Expectation{
			Message: "unable to parse recipe manifest",
			Dump: heredoc.Doc(`
				at %[1]s:1:1

				▶ 1 │ manala: {}
				    ├─╯ missing manala description property
			`,
				filepath.Join(repositoryURL, "recipe", ".manala.yaml"),
			),
		}, err)
	})
}

func (s *CommandSuite) execute(defaultRepositoryURL string, args ...string) (*bytes.Buffer, *bytes.Buffer, error) {
	stdOut := &bytes.Buffer{}
	stdErr := &bytes.Buffer{}

	logger := log.New(stdErr)
	logger.Verbose(1)

	command := cmdList.NewCommand(
		logger,
		api.New(
			logger,
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
