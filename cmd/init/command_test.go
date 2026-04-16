package init_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/manala/manala/app"
	"github.com/manala/manala/app/api"
	"github.com/manala/manala/app/testing/errors"
	cmdInit "github.com/manala/manala/cmd/init"
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
	s.Run("AlreadyExistingProject", func() {
		projectDir := filepath.FromSlash("testdata/TestProjectErrors/AlreadyExistingProject/project")

		stdOut, stdErr, err := s.execute("",
			projectDir,
		)

		s.Empty(stdOut)
		heredoc.Equal(s.T(), `
			 • finding project…
		`, stdErr)

		expect.Error(s.T(), errors.Expectation{
			Type: &app.AlreadyExistingProjectError{},
			Attrs: [][2]any{
				{"dir", projectDir},
			},
		}, err)
	})
}

func (s *CommandSuite) TestRepositoryErrors() {
	s.Run("NoRepository", func() {
		stdOut, stdErr, err := s.execute("")

		s.Empty(stdOut)
		heredoc.Equal(s.T(), `
			 • finding project…
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
			 • finding project…
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
			 • finding project…
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
			 • finding project…
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
	projectDir := filepath.FromSlash("testdata/TestRepositoryCustom/project")
	repositoryURL := filepath.FromSlash("testdata/TestRepositoryCustom/repository")

	_ = os.RemoveAll(projectDir)

	stdOut, stdErr, err := s.execute("",
		projectDir,
		"--repository", repositoryURL,
		"--recipe", "recipe",
	)

	s.Require().NoError(err)
	heredoc.Equal(s.T(), `
		project successfully initialized
	`, stdOut)
	heredoc.Equal(s.T(), `
		 • finding project…
		 • loading repository…
		 • loading recipe…
		 • creating project…
		 • syncing project…
	`, stdErr)

	s.DirExists(projectDir)

	heredoc.EqualFile(s.T(), `
		####################################################################
		#                         !!! REMINDER !!!                         #
		# Don't forget to run `+"`"+`manala up`+"`"+` each time you update this file ! #
		####################################################################

		manala:
		    recipe: recipe
		    repository: %[1]s
	`, filepath.Join(projectDir, ".manala.yaml"), repositoryURL)
}

func (s *CommandSuite) TestRepositoryConfig() {
	projectDir := filepath.FromSlash("testdata/TestRepositoryConfig/project")
	repositoryURL := filepath.FromSlash("testdata/TestRepositoryConfig/repository")

	_ = os.RemoveAll(projectDir)

	stdOut, stdErr, err := s.execute(repositoryURL,
		projectDir,
		"--recipe", "recipe",
	)

	s.Require().NoError(err)
	heredoc.Equal(s.T(), `
		project successfully initialized
	`, stdOut)
	heredoc.Equal(s.T(), `
		 • finding project…
		 • loading repository…
		 • loading recipe…
		 • creating project…
		 • syncing project…
		 • file synced                      path=file.txt
		 • file synced                      path=template
	`, stdErr)

	s.DirExists(projectDir)

	heredoc.EqualFile(s.T(), `
		manala:
		    recipe: recipe
		    repository: %[1]s

		foo: bar
	`, filepath.Join(projectDir, ".manala.yaml"), repositoryURL)

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
		repositoryURL := filepath.FromSlash("testdata/TestRecipeErrors/RecipeNotFound/repository")

		stdOut, stdErr, err := s.execute("",
			"--repository", repositoryURL,
			"--recipe", "recipe",
		)

		s.Empty(stdOut)
		heredoc.Equal(s.T(), `
			 • finding project…
			 • loading repository…
			 • loading recipe…
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
		repositoryURL := filepath.FromSlash("testdata/TestRecipeErrors/WrongRecipeManifest/repository")

		stdOut, stdErr, err := s.execute("",
			"--repository", repositoryURL,
			"--recipe", "recipe",
		)

		s.Empty(stdOut)
		heredoc.Equal(s.T(), `
			 • finding project…
			 • loading repository…
			 • loading recipe…
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
			"--recipe", "recipe",
		)

		s.Empty(stdOut)
		heredoc.Equal(s.T(), `
			 • finding project…
			 • loading repository…
			 • loading recipe…
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

	command := cmdInit.NewCommand(
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
