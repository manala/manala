package update_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/manala/manala/app"
	"github.com/manala/manala/app/api"
	"github.com/manala/manala/app/testing/errors"
	cmdUpdate "github.com/manala/manala/cmd/update"
	"github.com/manala/manala/internal/cache"
	"github.com/manala/manala/internal/errors/serror"
	"github.com/manala/manala/internal/errors/source"
	"github.com/manala/manala/internal/log"
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
			stdout, stderr, err := s.execute("",
				filepath.Join(dir, test.test, "project"),
			)

			s.Empty(stdout)

			s.Equal(test.expectedStderr, stderr.String())
			expectation.ExpectError(s.T(), test.expectedError, err)
		})
	}
}

func (s *CommandSuite) TestProjectRecursiveErrors() {
	dir := filepath.FromSlash("testdata/TestProjectRecursiveErrors")

	tests := []struct {
		test           string
		expectedStderr string
		expectedError  expectation.ErrorExpectation
	}{
		{
			test: "NotFound",
			expectedStderr: heredoc.Doc(`
				 ● loading projects recursive…
			`),
			expectedError: nil,
		},
		{
			test: "Invalid",
			expectedStderr: heredoc.Doc(`
				 ● loading projects recursive…
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
				 ● loading projects recursive…
			`),
			expectedError: serror.Expectation{
				Msg: "invalid project manifest vars",
				Err: expectation.Errors(
					source.Expectation(heredoc.Doc(`

						at %[1]s:5:6

						  2 │   recipe: recipe
						  3 │   repository: testdata/TestProjectRecursiveErrors/InvalidVars/repository
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
			stdout, stderr, err := s.execute("",
				filepath.Join(dir, test.test, "project"),
				"--recursive",
			)

			s.Empty(stdout)

			s.Equal(test.expectedStderr, stderr.String())
			expectation.ExpectError(s.T(), test.expectedError, err)
		})
	}
}

func (s *CommandSuite) TestRepository() {
	projectDir := filepath.FromSlash("testdata/TestRepository/project")
	repositoryURL := filepath.FromSlash("testdata/TestRepository/repository")

	_ = os.Remove(filepath.Join(projectDir, "file.txt"))
	_ = os.Remove(filepath.Join(projectDir, "template"))

	stdout, stderr, err := s.execute(repositoryURL,
		projectDir,
	)

	s.Require().NoError(err)
	heredoc.Equal(s.T(), `
		project successfully updated
	`, stdout)
	heredoc.Equal(s.T(), `
		 ● loading project…
		 ● syncing project…
		 ● file synced                      path=file.txt
		 ● file synced                      path=template
	`, stderr)

	heredoc.EqualFile(s.T(), `
		File
	`, filepath.Join(projectDir, "file.txt"))

	heredoc.EqualFile(s.T(), `
		Template

		foo: bar
	`, filepath.Join(projectDir, "template"))
}

func (s *CommandSuite) TestRepositoryArg() {
	projectDir := filepath.FromSlash("testdata/TestRepositoryArg/project")
	repositoryURL := filepath.FromSlash("testdata/TestRepositoryArg/repository")

	_ = os.Remove(filepath.Join(projectDir, "file.txt"))

	stdout, stderr, err := s.execute("",
		projectDir,
		"--repository", repositoryURL,
	)

	s.Require().NoError(err)
	heredoc.Equal(s.T(), `
		project successfully updated
	`, stdout)
	heredoc.Equal(s.T(), `
		 ● loading project…
		 ● syncing project…
		 ● file synced                      path=file.txt
	`, stderr)

	heredoc.EqualFile(s.T(), `
		File
	`, filepath.Join(projectDir, "file.txt"))
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
			stdout, stderr, err := s.execute("",
				filepath.Join(dir, test.test, "project"),
				"--repository", filepath.Join(dir, test.test, "repository"),
			)

			s.Empty(stdout)

			s.Equal(test.expectedStderr, stderr.String())
			expectation.ExpectError(s.T(), test.expectedError, err)
		})
	}
}

func (s *CommandSuite) TestRecipeArg() {
	projectDir := filepath.FromSlash("testdata/TestRecipeArg/project")
	repositoryURL := filepath.FromSlash("testdata/TestRecipeArg/repository")

	_ = os.Remove(filepath.Join(projectDir, "file.txt"))

	stdout, stderr, err := s.execute(repositoryURL,
		projectDir,
		"--recipe", "recipe",
	)

	s.Require().NoError(err)
	heredoc.Equal(s.T(), `
		project successfully updated
	`, stdout)
	heredoc.Equal(s.T(), `
		 ● loading project…
		 ● syncing project…
		 ● file synced                      path=file.txt
	`, stderr)

	heredoc.EqualFile(s.T(), `
		File
	`, filepath.Join(projectDir, "file.txt"))
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
			stdout, stderr, err := s.execute("",
				filepath.Join(dir, test.test, "project"),
				"--repository", filepath.Join(dir, test.test, "repository"),
			)

			s.Empty(stdout)

			s.Equal(test.expectedStderr, stderr.String())
			expectation.ExpectError(s.T(), test.expectedError, err)
		})
	}
}

func (s *CommandSuite) execute(defaultRepositoryURL string, args ...string) (*bytes.Buffer, *bytes.Buffer, error) {
	out := &bytes.Buffer{}
	err := &bytes.Buffer{}

	logger := log.New(output.NewDetached(err))
	logger.Verbose(1)

	command := cmdUpdate.NewCommand(
		logger,
		api.New(
			logger,
			cache.New(""),
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
