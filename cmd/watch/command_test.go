package watch_test

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/manala/manala/app"
	"github.com/manala/manala/app/api"
	"github.com/manala/manala/app/testing/errors"
	cmdWatch "github.com/manala/manala/cmd/watch"
	"github.com/manala/manala/internal/cache"
	"github.com/manala/manala/internal/errors/serror"
	"github.com/manala/manala/internal/errors/source"
	"github.com/manala/manala/internal/log"
	"github.com/manala/manala/internal/notify"
	"github.com/manala/manala/internal/output"
	"github.com/manala/manala/internal/testing/expectation"
	"github.com/manala/manala/internal/testing/heredoc"

	"github.com/stretchr/testify/suite"
)

type CommandSuite struct{ suite.Suite }

func TestCommandSuite(t *testing.T) {
	suite.Run(t, new(CommandSuite))
}

func (s *CommandSuite) TestProjectErrors() {
	dir := filepath.FromSlash("testdata/TestProjectErrors")

	tests := []struct {
		test           string
		expectedStderr string
		expectedError  expectation.ErrorExpectation
	}{
		{
			test: "NotFound",
			expectedStderr: heredoc.Doc(`
				 ● loading project…
			`),
			expectedError: errors.Expectation{
				Type: &app.NotFoundProjectError{},
				Attrs: [][2]any{
					{"dir", filepath.Join(dir, "NotFound", "project")},
				},
			},
		},
		{
			test: "Invalid",
			expectedStderr: heredoc.Doc(`
				 ● loading project…
			`),
			expectedError: serror.Expectation{
				Msg: "invalid project manifest",
				Err: expectation.Errors(
					source.Expectation(heredoc.Doc(`

						at %[1]s:1:1

						▶ 1 │ manala: {}
						    ├─╯ missing property 'recipe'
					`,
						filepath.Join(dir, "Invalid", "project", ".manala.yaml"),
					)),
				),
			},
		},
		{
			test: "InvalidVars",
			expectedStderr: heredoc.Doc(`
				 ● loading project…
			`),
			expectedError: serror.Expectation{
				Msg: "invalid project manifest vars",
				Err: expectation.Errors(
					source.Expectation(heredoc.Doc(`

						at %[1]s:5:6

						  2 │   recipe: recipe
						  3 │   repository: testdata/TestProjectErrors/InvalidVars/repository
						  4 │
						▶ 5 │ foo: bar
						    ├──────╯ got string, want integer
					`,
						filepath.Join(dir, "InvalidVars", "project", ".manala.yaml"),
					)),
				),
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			stdout, stderr, err := s.execute(
				filepath.Join(dir, test.test, "project"),
			)

			s.Empty(stdout)

			s.Equal(test.expectedStderr, stderr.String())
			expectation.ExpectError(s.T(), test.expectedError, err)
		})
	}
}

func (s *CommandSuite) TestRepositoryErrors() {
	dir := filepath.FromSlash("testdata/TestRepositoryErrors")

	tests := []struct {
		test           string
		expectedStderr string
		expectedError  expectation.ErrorExpectation
	}{
		{
			test: "NotFound",
			expectedStderr: heredoc.Doc(`
				 ● loading project…
			`),
			expectedError: errors.Expectation{
				Type: &app.NotFoundRepositoryError{},
				Attrs: [][2]any{
					{"url", filepath.Join(dir, "NotFound", "repository")},
				},
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			stdout, stderr, err := s.execute(
				filepath.Join(dir, test.test, "project"),
				"--repository", filepath.Join(dir, test.test, "repository"),
			)

			s.Empty(stdout)

			s.Equal(test.expectedStderr, stderr.String())
			expectation.ExpectError(s.T(), test.expectedError, err)
		})
	}
}

func (s *CommandSuite) TestRecipeErrors() {
	dir := filepath.FromSlash("testdata/TestRecipeErrors")

	tests := []struct {
		test           string
		expectedStderr string
		expectedError  expectation.ErrorExpectation
	}{
		{
			test: "NotFound",
			expectedStderr: heredoc.Doc(`
				 ● loading project…
			`),
			expectedError: errors.Expectation{
				Type: &app.NotFoundRecipeError{},
				Attrs: [][2]any{
					{"repository", filepath.Join(dir, "NotFound", "repository")},
					{"name", "recipe"},
				},
			},
		},
		{
			test: "Invalid",
			expectedStderr: heredoc.Doc(`
				 ● loading project…
			`),
			expectedError: serror.Expectation{
				Msg: "invalid recipe manifest",
				Err: expectation.Errors(
					source.Expectation(heredoc.Doc(`

						at %[1]s:1:1

						▶ 1 │ manala: {}
						    ├─╯ missing property 'description'
					`,
						filepath.Join(dir, "Invalid", "repository", "recipe", ".manala.yaml"),
					)),
				),
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			stdout, stderr, err := s.execute(
				filepath.Join(dir, test.test, "project"),
				"--repository", filepath.Join(dir, test.test, "repository"),
				"--recipe", "recipe",
			)

			s.Empty(stdout)

			s.Equal(test.expectedStderr, stderr.String())
			expectation.ExpectError(s.T(), test.expectedError, err)
		})
	}
}

func (s *CommandSuite) execute(args ...string) (*bytes.Buffer, *bytes.Buffer, error) {
	out := &bytes.Buffer{}
	err := &bytes.Buffer{}

	logger := log.New(output.NewDetached(err))
	logger.Verbose(1)

	command := cmdWatch.NewCommand(
		logger,
		api.New(
			logger,
			cache.New(""),
		),
		output.NewDetached(out),
		notify.New(notify.DiscardHandler),
	)

	command.SilenceErrors = true
	command.SilenceUsage = true
	command.SetOut(out)
	command.SetErr(err)
	command.SetArgs(append([]string{}, args...))

	return out, err, command.Execute()
}
