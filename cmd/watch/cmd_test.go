package watch

import (
	"bytes"
	"github.com/stretchr/testify/suite"
	"log/slog"
	"manala/app"
	"manala/app/api"
	"manala/internal/cache"
	"manala/internal/notifier"
	"manala/internal/serrors"
	"manala/internal/testing/heredoc"
	"manala/internal/ui/adapters/charm"
	"manala/internal/ui/log"
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
		notifier.NewNil(),
	)

	cmd.SilenceErrors = true
	cmd.SilenceUsage = true
	cmd.SetOut(stdOut)
	cmd.SetErr(stdErr)
	cmd.SetArgs(append([]string{}, args...))

	return stdOut, stdErr, cmd.Execute()
}

func (s *Suite) TestProjectErrors() {
	s.Run("ProjectNotFound", func() {
		projectDir := filepath.FromSlash("testdata/TestProjectErrors/ProjectNotFound/project")

		stdOut, stdErr, err := s.execute("",
			projectDir,
		)

		s.Empty(stdOut)
		heredoc.Equal(s.T(), `
			 • loading project…
		`, stdErr)

		serrors.Equal(s.T(), &serrors.Assertion{
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

		serrors.Equal(s.T(), &serrors.Assertion{
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

		serrors.Equal(s.T(), &serrors.Assertion{
			Message: "unable to read project manifest",
			Arguments: []any{
				"file", filepath.Join(projectDir, ".manala.yaml"),
			},
			Errors: []*serrors.Assertion{
				{
					Message: "irregular project manifest",
					Errors: []*serrors.Assertion{
						{
							Message: "empty yaml file",
						},
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

		serrors.Equal(s.T(), &serrors.Assertion{
			Message: "unable to read project manifest",
			Arguments: []any{
				"file", filepath.Join(projectDir, ".manala.yaml"),
			},
			Errors: []*serrors.Assertion{
				{
					Message: "invalid project manifest",
					Errors: []*serrors.Assertion{
						{
							Message: "missing manala recipe property",
							Arguments: []any{
								"path", "manala",
								"property", "recipe",
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

func (s *Suite) TestRepositoryErrors() {
	s.Run("NoRepository", func() {
		projectDir := filepath.FromSlash("testdata/TestRepositoryErrors/NoRepository/project")

		stdOut, stdErr, err := s.execute("",
			projectDir,
		)

		s.Empty(stdOut)
		heredoc.Equal(s.T(), `
			 • loading project…
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
		projectDir := filepath.FromSlash("testdata/TestRepositoryErrors/RepositoryNotFound/project")
		repositoryUrl := filepath.FromSlash("testdata/TestRepositoryErrors/RepositoryNotFound/repository")

		stdOut, stdErr, err := s.execute("",
			projectDir,
			"--repository", repositoryUrl,
		)

		s.Empty(stdOut)
		heredoc.Equal(s.T(), `
			 • loading project…
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
		projectDir := filepath.FromSlash("testdata/TestRepositoryErrors/WrongRepository/project")
		repositoryUrl := filepath.FromSlash("testdata/TestRepositoryErrors/WrongRepository/repository")

		stdOut, stdErr, err := s.execute("",
			projectDir,
			"--repository", repositoryUrl,
		)

		s.Empty(stdOut)
		heredoc.Equal(s.T(), `
			 • loading project…
		`, stdErr)

		serrors.Equal(s.T(), &serrors.Assertion{
			Type:    &app.NotFoundRepositoryError{},
			Message: "repository not found",
			Arguments: []any{
				"url", repositoryUrl,
			},
		}, err)
	})
}

func (s *Suite) TestRecipeErrors() {
	s.Run("RecipeNotFound", func() {
		projectDir := filepath.FromSlash("testdata/TestRecipeErrors/RecipeNotFound/project")
		repositoryUrl := filepath.FromSlash("testdata/TestRecipeErrors/RecipeNotFound/repository")

		stdOut, stdErr, err := s.execute("",
			projectDir,
			"--repository", repositoryUrl,
			"--recipe", "recipe",
		)

		s.Empty(stdOut)
		heredoc.Equal(s.T(), `
			 • loading project…
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
		projectDir := filepath.FromSlash("testdata/TestRecipeErrors/WrongRecipeManifest/project")
		repositoryUrl := filepath.FromSlash("testdata/TestRecipeErrors/WrongRecipeManifest/repository")

		stdOut, stdErr, err := s.execute("",
			projectDir,
			"--repository", repositoryUrl,
			"--recipe", "recipe",
		)

		s.Empty(stdOut)
		heredoc.Equal(s.T(), `
			 • loading project…
		`, stdErr)

		serrors.Equal(s.T(), &serrors.Assertion{
			Message: "recipe manifest is a directory",
			Arguments: []any{
				"dir", filepath.Join(repositoryUrl, "recipe", ".manala.yaml"),
			},
		}, err)
	})

	s.Run("InvalidRecipeManifest", func() {
		projectDir := filepath.FromSlash("testdata/TestRecipeErrors/InvalidRecipeManifest/project")
		repositoryUrl := filepath.FromSlash("testdata/TestRecipeErrors/InvalidRecipeManifest/repository")

		stdOut, stdErr, err := s.execute("",
			projectDir,
			"--repository", repositoryUrl,
			"--recipe", "recipe",
		)

		s.Empty(stdOut)
		heredoc.Equal(s.T(), `
			 • loading project…
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
