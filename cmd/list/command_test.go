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
	"github.com/manala/manala/internal/output"
	"github.com/manala/manala/internal/serrors"
	"github.com/manala/manala/internal/testing/expect"
	"github.com/manala/manala/internal/testing/heredoc"

	"github.com/stretchr/testify/suite"
)

type CommandSuite struct{ suite.Suite }

func TestCommandSuite(t *testing.T) {
	suite.Run(t, new(CommandSuite))
}

func (s *CommandSuite) TestRepository() {
	repositoryURL := filepath.FromSlash("testdata/TestRepository/repository")

	stdout, stderr, err := s.execute(repositoryURL)

	s.Require().NoError(err)
	heredoc.Equal(s.T(), `
		bar
		  Bar
		foo
		  Foo
	`, stdout)
	heredoc.Equal(s.T(), `
		 ● loading repository…
		 ● loading recipes…
	`, stderr)
}

func (s *CommandSuite) TestRepositoryArg() {
	repositoryURL := filepath.FromSlash("testdata/TestRepositoryArg/repository")

	stdout, stderr, err := s.execute("",
		"--repository", repositoryURL,
	)

	s.Require().NoError(err)
	heredoc.Equal(s.T(), `
		bar
		  Bar
		foo
		  Foo
	`, stdout)
	heredoc.Equal(s.T(), `
		 ● loading repository…
		 ● loading recipes…
	`, stderr)
}

func (s *CommandSuite) TestRepositoryErrors() {
	dir := filepath.FromSlash("testdata/TestRepositoryErrors")

	tests := []struct {
		test           string
		expectedStderr string
		expectedError  expect.ErrorExpectation
	}{
		{
			test: "NotFound",
			expectedStderr: heredoc.Doc(`
				 ● loading repository…
			`),
			expectedError: errors.Expectation{
				Type: &app.NotFoundRepositoryError{},
				Attrs: [][2]any{
					{"url", filepath.Join(dir, "NotFound", "repository")},
				},
			},
		},
		{
			test: "Empty",
			expectedStderr: heredoc.Doc(`
				 ● loading repository…
				 ● loading recipes…
			`),
			expectedError: errors.Expectation{
				Type: &app.EmptyRepositoryError{},
				Attrs: [][2]any{
					{"url", filepath.Join(dir, "Empty", "repository")},
				},
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			stdout, stderr, err := s.execute("",
				"--repository", filepath.Join(dir, test.test, "repository"),
			)

			s.Empty(stdout)

			s.Equal(test.expectedStderr, stderr.String())
			expect.Error(s.T(), test.expectedError, err)
		})
	}
}

func (s *CommandSuite) TestRecipeErrors() {
	dir := filepath.FromSlash("testdata/TestRecipeErrors")

	tests := []struct {
		test           string
		expectedStderr string
		expectedError  expect.ErrorExpectation
	}{
		{
			test: "Unparsable",
			expectedStderr: heredoc.Doc(`
				 ● loading repository…
				 ● loading recipes…
			`),
			expectedError: serrors.Expectation{
				Message: "unable to parse recipe manifest",
				Dump: heredoc.Doc(`
					at %[1]s:1:1

					▶ 1 │ manala: {}
					    ├─╯ missing manala "description" property
				`,
					filepath.Join(dir, "Unparsable", "repository", "recipe", ".manala.yaml"),
				),
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			stdout, stderr, err := s.execute("",
				"--repository", filepath.Join(dir, test.test, "repository"),
			)

			s.Empty(stdout)

			s.Equal(test.expectedStderr, stderr.String())
			expect.Error(s.T(), test.expectedError, err)
		})
	}
}

func (s *CommandSuite) execute(defaultRepositoryURL string, args ...string) (*bytes.Buffer, *bytes.Buffer, error) {
	out := &bytes.Buffer{}
	err := &bytes.Buffer{}

	logger := log.New(output.NewDetached(err))
	logger.Verbose(1)

	command := cmdList.NewCommand(
		logger,
		api.New(
			logger,
			caching.NewCache(""),
			api.WithDefaultRepositoryURL(defaultRepositoryURL),
		),
		output.NewDetached(out),
	)

	command.SilenceErrors = true
	command.SilenceUsage = true
	command.SetOut(out)
	command.SetErr(err)
	command.SetArgs(append([]string{}, args...))

	return out, err, command.Execute()
}
