package init_test

import (
	"bytes"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/manala/manala/app"
	"github.com/manala/manala/app/api"
	cmdInit "github.com/manala/manala/cmd/init"
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

func (s *Suite) TestProjectErrors() {
	s.Run("AlreadyExistingProject", func() {
		projectDir := filepath.FromSlash("testdata/TestProjectErrors/AlreadyExistingProject/project")

		stdOut, stdErr, err := s.execute("",
			projectDir,
		)

		s.Empty(stdOut)
		heredoc.Equal(s.T(), `
			 • finding project…
		`, stdErr)

		serrors.Equal(s.T(), &serrors.Assertion{
			Type:    &app.AlreadyExistingProjectError{},
			Message: "already existing project",
			Arguments: []any{
				"dir", projectDir,
			},
		}, err)
	})
}

func (s *Suite) TestRepositoryErrors() {
	s.Run("NoRepository", func() {
		stdOut, stdErr, err := s.execute("")

		s.Empty(stdOut)
		heredoc.Equal(s.T(), `
			 • finding project…
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
			 • finding project…
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
			 • finding project…
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
			 • finding project…
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
	projectDir := filepath.FromSlash("testdata/TestRepositoryCustom/project")
	repositoryURL := filepath.FromSlash("testdata/TestRepositoryCustom/repository")

	_ = os.RemoveAll(projectDir)

	stdOut, stdErr, err := s.execute("",
		projectDir,
		"--repository", repositoryURL,
		"--recipe", "recipe",
	)

	s.Require().NoError(err)
	s.Empty(stdOut)
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

func (s *Suite) TestRepositoryConfig() {
	projectDir := filepath.FromSlash("testdata/TestRepositoryConfig/project")
	repositoryURL := filepath.FromSlash("testdata/TestRepositoryConfig/repository")

	_ = os.RemoveAll(projectDir)

	stdOut, stdErr, err := s.execute(repositoryURL,
		projectDir,
		"--recipe", "recipe",
	)

	s.Require().NoError(err)
	s.Empty(stdOut)
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

func (s *Suite) TestRecipeErrors() {
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

		serrors.Equal(s.T(), &serrors.Assertion{
			Type:    &app.NotFoundRecipeError{},
			Message: "recipe not found",
			Arguments: []any{
				"repository", repositoryURL,
				"name", "recipe",
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
			"--recipe", "recipe",
		)

		s.Empty(stdOut)
		heredoc.Equal(s.T(), `
			 • finding project…
			 • loading repository…
			 • loading recipe…
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

	ui := charm.New(nil, stdOut, stdErr)
	log := slog.New(log.NewSlogHandler(ui))

	command := cmdInit.NewCommand(
		log,
		api.New(
			log,
			caching.NewCache(""),
			api.WithDefaultRepositoryURL(defaultRepositoryURL),
		),
		ui,
	)

	command.SilenceErrors = true
	command.SilenceUsage = true
	command.SetOut(stdOut)
	command.SetErr(stdErr)
	command.SetArgs(append([]string{}, args...))

	return stdOut, stdErr, command.Execute()
}
