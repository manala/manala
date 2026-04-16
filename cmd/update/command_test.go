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
	s.Run("ProjectNotFound", func() {
		projectDir := filepath.FromSlash("testdata/TestProjectErrors/ProjectNotFound/project")

		stdOut, stdErr, err := s.execute("",
			projectDir,
		)

		s.Empty(stdOut)
		heredoc.Equal(s.T(), `
			 • loading project…
		`, stdErr)

		expect.Error(s.T(), errors.Expectation{
			Type: &app.NotFoundProjectError{},
			Attrs: [][2]any{
				{"dir", projectDir},
			},
		}, err)
	})

	s.Run("WrongProjectManifest", func() {
		projectDir := filepath.FromSlash("testdata/TestProjectErrors/WrongProjectManifest/project")

		stdOut, stdErr, err := s.execute("",
			projectDir,
		)

		s.Empty(stdOut)
		heredoc.Equal(s.T(), `
			 • loading project…
		`, stdErr)

		expect.Error(s.T(), serrors.Expectation{
			Message: "project manifest is a directory",
			Attrs: [][2]any{
				{"dir", filepath.Join(projectDir, ".manala.yaml")},
			},
		}, err)
	})

	s.Run("EmptyProjectManifest", func() {
		projectDir := filepath.FromSlash("testdata/TestProjectErrors/EmptyProjectManifest/project")

		stdOut, stdErr, err := s.execute("",
			projectDir,
		)

		s.Empty(stdOut)
		heredoc.Equal(s.T(), `
			 • loading project…
		`, stdErr)

		expect.Error(s.T(), serrors.Expectation{
			Message: "unable to parse project manifest",
			Dump: heredoc.Doc(`
				at %[1]s:0

			`,
				filepath.Join(projectDir, ".manala.yaml"),
			),
		}, err)
	})

	s.Run("InvalidProjectManifest", func() {
		projectDir := filepath.FromSlash("testdata/TestProjectErrors/InvalidProjectManifest/project")

		stdOut, stdErr, err := s.execute("",
			projectDir,
		)

		s.Empty(stdOut)
		heredoc.Equal(s.T(), `
			 • loading project…
		`, stdErr)

		expect.Error(s.T(), serrors.Expectation{
			Message: "unable to parse project manifest",
			Dump: heredoc.Doc(`
				at %[1]s:1:1

				▶ 1 │ manala: {}
				    ├─╯ missing manala recipe property
			`,
				filepath.Join(projectDir, ".manala.yaml"),
			),
		}, err)
	})
}

func (s *CommandSuite) TestRecursiveProjectErrors() {
	s.Run("ProjectNotFound", func() {
		projectDir := filepath.FromSlash("testdata/TestRecursiveProjectErrors/ProjectNotFound/project")

		stdOut, stdErr, err := s.execute("",
			projectDir,
			"--recursive",
		)

		s.Require().NoError(err)
		s.Empty(stdOut)
		heredoc.Equal(s.T(), `
			 • loading projects recursive…
		`, stdErr)
	})

	s.Run("WrongProjectManifest", func() {
		projectDir := filepath.FromSlash("testdata/TestRecursiveProjectErrors/WrongProjectManifest/project")

		stdOut, stdErr, err := s.execute("",
			projectDir,
			"--recursive",
		)

		s.Empty(stdOut)
		heredoc.Equal(s.T(), `
			 • loading projects recursive…
		`, stdErr)

		expect.Error(s.T(), serrors.Expectation{
			Message: "project manifest is a directory",
			Attrs: [][2]any{
				{"dir", filepath.Join(projectDir, ".manala.yaml")},
			},
		}, err)
	})

	s.Run("EmptyProjectManifest", func() {
		projectDir := filepath.FromSlash("testdata/TestRecursiveProjectErrors/EmptyProjectManifest/project")

		stdOut, stdErr, err := s.execute("",
			projectDir,
			"--recursive",
		)

		s.Empty(stdOut)
		heredoc.Equal(s.T(), `
			 • loading projects recursive…
		`, stdErr)

		expect.Error(s.T(), serrors.Expectation{
			Message: "unable to parse project manifest",
			Dump: heredoc.Doc(`
				at %[1]s:0

			`,
				filepath.Join(projectDir, ".manala.yaml"),
			),
		}, err)
	})

	s.Run("InvalidProjectManifest", func() {
		projectDir := filepath.FromSlash("testdata/TestRecursiveProjectErrors/InvalidProjectManifest/project")

		stdOut, stdErr, err := s.execute("",
			projectDir,
			"--recursive",
		)

		s.Empty(stdOut)
		heredoc.Equal(s.T(), `
			 • loading projects recursive…
		`, stdErr)

		expect.Error(s.T(), serrors.Expectation{
			Message: "unable to parse project manifest",
			Dump: heredoc.Doc(`
				at %[1]s:1:1

				▶ 1 │ manala: {}
				    ├─╯ missing manala recipe property
			`,
				filepath.Join(projectDir, ".manala.yaml"),
			),
		}, err)
	})
}

func (s *CommandSuite) TestRepositoryErrors() {
	s.Run("NoRepository", func() {
		projectDir := filepath.FromSlash("testdata/TestRepositoryErrors/NoRepository/project")

		stdOut, stdErr, err := s.execute("",
			projectDir,
		)

		s.Empty(stdOut)
		heredoc.Equal(s.T(), `
			 • loading project…
		`, stdErr)

		expect.Error(s.T(), errors.Expectation{
			Type: &app.NotFoundRepositoryError{},
			Attrs: [][2]any{
				{"url", ""},
			},
		}, err)
	})

	s.Run("RepositoryNotFound", func() {
		projectDir := filepath.FromSlash("testdata/TestRepositoryErrors/RepositoryNotFound/project")
		repositoryURL := filepath.FromSlash("testdata/TestRepositoryErrors/RepositoryNotFound/repository")

		stdOut, stdErr, err := s.execute("",
			projectDir,
			"--repository", repositoryURL,
		)

		s.Empty(stdOut)
		heredoc.Equal(s.T(), `
			 • loading project…
		`, stdErr)

		expect.Error(s.T(), errors.Expectation{
			Type: &app.NotFoundRepositoryError{},
			Attrs: [][2]any{
				{"url", repositoryURL},
			},
		}, err)
	})

	s.Run("WrongRepository", func() {
		projectDir := filepath.FromSlash("testdata/TestRepositoryErrors/WrongRepository/project")
		repositoryURL := filepath.FromSlash("testdata/TestRepositoryErrors/WrongRepository/repository")

		stdOut, stdErr, err := s.execute("",
			projectDir,
			"--repository", repositoryURL,
		)

		s.Empty(stdOut)
		heredoc.Equal(s.T(), `
			 • loading project…
		`, stdErr)

		expect.Error(s.T(), errors.Expectation{
			Type: &app.NotFoundRepositoryError{},
			Attrs: [][2]any{
				{"url", repositoryURL},
			},
		}, err)
	})
}

func (s *CommandSuite) TestRepositoryCustom() {
	projectDir := filepath.FromSlash("testdata/TestRepositoryCustom/project")
	repositoryURL := filepath.FromSlash("testdata/TestRepositoryCustom/repository")

	_ = os.Remove(filepath.Join(projectDir, "file.txt"))

	stdOut, stdErr, err := s.execute("",
		projectDir,
		"--repository", repositoryURL,
	)

	s.Require().NoError(err)
	heredoc.Equal(s.T(), `
		project successfully updated
	`, stdOut)
	heredoc.Equal(s.T(), `
		 • loading project…
		 • syncing project…
		 • file synced                      path=file.txt
	`, stdErr)

	heredoc.EqualFile(s.T(), `
		File
	`, filepath.Join(projectDir, "file.txt"))
}

func (s *CommandSuite) TestRepositoryConfig() {
	projectDir := filepath.FromSlash("testdata/TestRepositoryConfig/project")
	repositoryURL := filepath.FromSlash("testdata/TestRepositoryConfig/repository")

	_ = os.Remove(filepath.Join(projectDir, "file.txt"))
	_ = os.Remove(filepath.Join(projectDir, "template"))

	stdOut, stdErr, err := s.execute(repositoryURL,
		projectDir,
	)

	s.Require().NoError(err)
	heredoc.Equal(s.T(), `
		project successfully updated
	`, stdOut)
	heredoc.Equal(s.T(), `
		 • loading project…
		 • syncing project…
		 • file synced                      path=file.txt
		 • file synced                      path=template
	`, stdErr)

	heredoc.EqualFile(s.T(), `
		File
	`, filepath.Join(projectDir, "file.txt"))

	heredoc.EqualFile(s.T(), `
		Template

		foo: bar
	`, filepath.Join(projectDir, "template"))
}

func (s *CommandSuite) TestRecipeErrors() {
	s.Run("RecipeNotFound", func() {
		projectDir := filepath.FromSlash("testdata/TestRecipeErrors/RecipeNotFound/project")
		repositoryURL := filepath.FromSlash("testdata/TestRecipeErrors/RecipeNotFound/repository")

		stdOut, stdErr, err := s.execute("",
			projectDir,
			"--repository", repositoryURL,
			"--recipe", "recipe",
		)

		s.Empty(stdOut)
		heredoc.Equal(s.T(), `
			 • loading project…
		`, stdErr)

		expect.Error(s.T(), errors.Expectation{
			Type: &app.NotFoundRecipeError{},
			Attrs: [][2]any{
				{"repository", repositoryURL},
				{"name", "recipe"},
			},
		}, err)
	})

	s.Run("WrongRecipeManifest", func() {
		projectDir := filepath.FromSlash("testdata/TestRecipeErrors/WrongRecipeManifest/project")
		repositoryURL := filepath.FromSlash("testdata/TestRecipeErrors/WrongRecipeManifest/repository")

		stdOut, stdErr, err := s.execute("",
			projectDir,
			"--repository", repositoryURL,
			"--recipe", "recipe",
		)

		s.Empty(stdOut)
		heredoc.Equal(s.T(), `
			 • loading project…
		`, stdErr)

		expect.Error(s.T(), serrors.Expectation{
			Message: "recipe manifest is a directory",
			Attrs: [][2]any{
				{"dir", filepath.Join(repositoryURL, "recipe", ".manala.yaml")},
			},
		}, err)
	})

	s.Run("InvalidRecipeManifest", func() {
		projectDir := filepath.FromSlash("testdata/TestRecipeErrors/InvalidRecipeManifest/project")
		repositoryURL := filepath.FromSlash("testdata/TestRecipeErrors/InvalidRecipeManifest/repository")

		stdOut, stdErr, err := s.execute("",
			projectDir,
			"--repository", repositoryURL,
			"--recipe", "recipe",
		)

		s.Empty(stdOut)
		heredoc.Equal(s.T(), `
			 • loading project…
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

func (s *CommandSuite) TestRecipeCustom() {
	projectDir := filepath.FromSlash("testdata/TestRecipeCustom/project")
	repositoryURL := filepath.FromSlash("testdata/TestRecipeCustom/repository")

	_ = os.Remove(filepath.Join(projectDir, "file.txt"))

	stdOut, stdErr, err := s.execute(repositoryURL,
		projectDir,
		"--recipe", "recipe",
	)

	s.Require().NoError(err)
	heredoc.Equal(s.T(), `
		project successfully updated
	`, stdOut)
	heredoc.Equal(s.T(), `
		 • loading project…
		 • syncing project…
		 • file synced                      path=file.txt
	`, stdErr)

	heredoc.EqualFile(s.T(), `
		File
	`, filepath.Join(projectDir, "file.txt"))
}

func (s *CommandSuite) execute(defaultRepositoryURL string, args ...string) (*bytes.Buffer, *bytes.Buffer, error) {
	stdOut := &bytes.Buffer{}
	stdErr := &bytes.Buffer{}

	logger := log.New(stdErr)
	logger.Verbose(1)

	command := cmdUpdate.NewCommand(
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
