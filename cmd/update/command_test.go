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

func (s *CommandSuite) TestProjectErrors() {
	dir := filepath.FromSlash("testdata/TestProjectErrors")

	tests := []struct {
		test           string
		expectedStderr string
		expectedError  expect.ErrorExpectation
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
			test: "Unparsable",
			expectedStderr: heredoc.Doc(`
				 ● loading project…
			`),
			expectedError: serrors.Expectation{
				Message: "unable to parse project manifest",
				Dump: heredoc.Doc(`
				at %[1]s:1:1

				▶ 1 │ manala: {}
				    ├─╯ missing manala "recipe" property
			`,
					filepath.Join(dir, "Unparsable", "project", ".manala.yaml"),
				),
			},
		},
		{
			test: "Invalid",
			expectedStderr: heredoc.Doc(`
				 ● loading project…
			`),
			expectedError: serrors.Expectation{
				Message: "invalid project manifest vars",
				Attrs: [][2]any{
					{"file", filepath.Join(dir, "Invalid", "project", ".manala.yaml")},
				},
				Errors: []expect.ErrorExpectation{
					serrors.Expectation{
						Message: "invalid type",
						Attrs: [][2]any{
							{"expected", "integer"},
							{"actual", "string"},
							{"path", "foo"},
						},
					},
				},
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
			expect.Error(s.T(), test.expectedError, err)
		})
	}
}

func (s *CommandSuite) TestProjectRecursiveErrors() {
	dir := filepath.FromSlash("testdata/TestProjectRecursiveErrors")

	tests := []struct {
		test           string
		expectedStderr string
		expectedError  expect.ErrorExpectation
	}{
		{
			test: "NotFound",
			expectedStderr: heredoc.Doc(`
				 ● loading projects recursive…
			`),
			expectedError: nil,
		},
		{
			test: "Unparsable",
			expectedStderr: heredoc.Doc(`
				 ● loading projects recursive…
			`),
			expectedError: serrors.Expectation{
				Message: "unable to parse project manifest",
				Dump: heredoc.Doc(`
				at %[1]s:1:1

				▶ 1 │ manala: {}
				    ├─╯ missing manala "recipe" property
			`,
					filepath.Join(dir, "Unparsable", "project", ".manala.yaml"),
				),
			},
		},
		{
			test: "Invalid",
			expectedStderr: heredoc.Doc(`
				 ● loading projects recursive…
			`),
			expectedError: serrors.Expectation{
				Message: "invalid project manifest vars",
				Attrs: [][2]any{
					{"file", filepath.Join(dir, "Invalid", "project", ".manala.yaml")},
				},
				Errors: []expect.ErrorExpectation{
					serrors.Expectation{
						Message: "invalid type",
						Attrs: [][2]any{
							{"expected", "integer"},
							{"actual", "string"},
							{"path", "foo"},
						},
					},
				},
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
			expect.Error(s.T(), test.expectedError, err)
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
		expectedError  expect.ErrorExpectation
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
			expect.Error(s.T(), test.expectedError, err)
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
		expectedError  expect.ErrorExpectation
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
			test: "Unparsable",
			expectedStderr: heredoc.Doc(`
				 ● loading project…
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
				filepath.Join(dir, test.test, "project"),
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

	command := cmdUpdate.NewCommand(
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
