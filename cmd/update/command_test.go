package update_test

import (
	"bytes"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/manala/manala/app"
	"github.com/manala/manala/app/api"
	cmdUpdate "github.com/manala/manala/cmd/update"
	"github.com/manala/manala/internal/caching"
	"github.com/manala/manala/internal/parsing"
	"github.com/manala/manala/internal/serrors"
	"github.com/manala/manala/internal/testing/errors"
	"github.com/manala/manala/internal/testing/heredoc"
	"github.com/manala/manala/internal/ui/adapters/charm"
	"github.com/manala/manala/internal/ui/log"

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

		errors.Equal(s.T(), &serrors.Assertion{
			Type:    &app.NotFoundProjectError{},
			Message: "project not found",
			Arguments: []any{
				"dir", projectDir,
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

		errors.Equal(s.T(), &serrors.Assertion{
			Message: "project manifest is a directory",
			Arguments: []any{
				"dir", filepath.Join(projectDir, ".manala.yaml"),
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

		errors.Equal(s.T(), &serrors.Assertion{
			Message: "unable to parse project manifest",
			Arguments: []any{
				"file", filepath.Join(projectDir, ".manala.yaml"),
			},
			Errors: []errors.Assertion{
				&parsing.ErrorAssertion{
					Err: &serrors.Assertion{
						Message: "empty yaml content",
					},
				},
			},
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

		errors.Equal(s.T(), &serrors.Assertion{
			Message: "unable to parse project manifest",
			Arguments: []any{
				"file", filepath.Join(projectDir, ".manala.yaml"),
				"line", 1, "column", 1,
			},
			Dump: `
				> 1 | manala: {}
				      ^
				* missing manala recipe property
			`,
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

		errors.Equal(s.T(), &serrors.Assertion{
			Message: "project manifest is a directory",
			Arguments: []any{
				"dir", filepath.Join(projectDir, ".manala.yaml"),
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

		errors.Equal(s.T(), &serrors.Assertion{
			Message: "unable to parse project manifest",
			Arguments: []any{
				"file", filepath.Join(projectDir, ".manala.yaml"),
			},
			Errors: []errors.Assertion{
				&parsing.ErrorAssertion{
					Err: &serrors.Assertion{
						Message: "empty yaml content",
					},
				},
			},
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

		errors.Equal(s.T(), &serrors.Assertion{
			Message: "unable to parse project manifest",
			Arguments: []any{
				"file", filepath.Join(projectDir, ".manala.yaml"),
				"line", 1, "column", 1,
			},
			Dump: `
				> 1 | manala: {}
				      ^
				* missing manala recipe property
			`,
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

		errors.Equal(s.T(), &serrors.Assertion{
			Type:    &app.NotFoundRepositoryError{},
			Message: "repository not found",
			Arguments: []any{
				"url", "",
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

		errors.Equal(s.T(), &serrors.Assertion{
			Type:    &app.NotFoundRepositoryError{},
			Message: "repository not found",
			Arguments: []any{
				"url", repositoryURL,
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

		errors.Equal(s.T(), &serrors.Assertion{
			Type:    &app.NotFoundRepositoryError{},
			Message: "repository not found",
			Arguments: []any{
				"url", repositoryURL,
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

		errors.Equal(s.T(), &serrors.Assertion{
			Type:    &app.NotFoundRecipeError{},
			Message: "recipe not found",
			Arguments: []any{
				"repository", repositoryURL,
				"name", "recipe",
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

		errors.Equal(s.T(), &serrors.Assertion{
			Message: "recipe manifest is a directory",
			Arguments: []any{
				"dir", filepath.Join(repositoryURL, "recipe", ".manala.yaml"),
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

		errors.Equal(s.T(), &serrors.Assertion{
			Message: "unable to parse recipe manifest",
			Arguments: []any{
				"file", filepath.Join(repositoryURL, "recipe", ".manala.yaml"),
				"line", 1, "column", 1,
			},
			Dump: `
				> 1 | manala: {}
				      ^
				* missing manala description property
			`,
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

	ui := charm.New(stdErr)
	log := slog.New(log.NewSlogHandler(ui))

	command := cmdUpdate.NewCommand(
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
