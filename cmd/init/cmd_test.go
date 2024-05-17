package init

import (
	"bytes"
	"github.com/stretchr/testify/suite"
	"log/slog"
	"manala/app"
	"manala/app/api"
	"manala/internal/cache"
	"manala/internal/serrors"
	"manala/internal/testing/heredoc"
	"manala/internal/ui/adapters/charm"
	"manala/internal/ui/log"
	"os"
	"path/filepath"
	"testing"
)

type Suite struct{ suite.Suite }

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) execute(defaultRepositoryUrl string, args ...string) (*bytes.Buffer, *bytes.Buffer, error) {
	stdOut := &bytes.Buffer{}
	stdErr := &bytes.Buffer{}

	ui := charm.New(nil, stdOut, stdErr)
	log := slog.New(log.NewSlogHandler(ui))

	cmd := NewCmd(
		log,
		api.New(
			log,
			cache.New(""),
			api.WithDefaultRepositoryUrl(defaultRepositoryUrl),
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
		repositoryUrl := filepath.FromSlash("testdata/TestRepositoryErrors/RepositoryNotFound/repository")

		stdOut, stdErr, err := s.execute("",
			"--repository", repositoryUrl,
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
				"url", repositoryUrl,
			},
		}, err)
	})

	s.Run("WrongRepository", func() {
		repositoryUrl := filepath.FromSlash("testdata/TestRepositoryErrors/WrongRepository/repository")

		stdOut, stdErr, err := s.execute("",
			"--repository", repositoryUrl,
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
				"url", repositoryUrl,
			},
		}, err)
	})

	s.Run("EmptyRepository", func() {
		repositoryUrl := filepath.FromSlash("testdata/TestRepositoryErrors/EmptyRepository/repository")

		stdOut, stdErr, err := s.execute("",
			"--repository", repositoryUrl,
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
				"url", repositoryUrl,
			},
		}, err)
	})
}

func (s *Suite) TestRepositoryCustom() {
	projectDir := filepath.FromSlash("testdata/TestRepositoryCustom/project")
	repositoryUrl := filepath.FromSlash("testdata/TestRepositoryCustom/repository")

	_ = os.RemoveAll(projectDir)

	stdOut, stdErr, err := s.execute("",
		projectDir,
		"--repository", repositoryUrl,
		"--recipe", "recipe",
	)

	s.NoError(err)

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
	`, filepath.Join(projectDir, ".manala.yaml"), repositoryUrl)
}

func (s *Suite) TestRepositoryConfig() {
	projectDir := filepath.FromSlash("testdata/TestRepositoryConfig/project")
	repositoryUrl := filepath.FromSlash("testdata/TestRepositoryConfig/repository")

	_ = os.RemoveAll(projectDir)

	stdOut, stdErr, err := s.execute(repositoryUrl,
		projectDir,
		"--recipe", "recipe",
	)

	s.NoError(err)

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
	`, filepath.Join(projectDir, ".manala.yaml"), repositoryUrl)

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
		repositoryUrl := filepath.FromSlash("testdata/TestRecipeErrors/RecipeNotFound/repository")

		stdOut, stdErr, err := s.execute("",
			"--repository", repositoryUrl,
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
				"repository", repositoryUrl,
				"name", "recipe",
			},
		}, err)
	})

	s.Run("WrongRecipeManifest", func() {
		repositoryUrl := filepath.FromSlash("testdata/TestRecipeErrors/WrongRecipeManifest/repository")

		stdOut, stdErr, err := s.execute("",
			"--repository", repositoryUrl,
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
				"dir", filepath.Join(repositoryUrl, "recipe", ".manala.yaml"),
			},
		}, err)
	})

	s.Run("InvalidRecipeManifest", func() {
		repositoryUrl := filepath.FromSlash("testdata/TestRecipeErrors/InvalidRecipeManifest/repository")

		stdOut, stdErr, err := s.execute("",
			"--repository", repositoryUrl,
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
				"file", filepath.Join(repositoryUrl, "recipe", ".manala.yaml"),
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
