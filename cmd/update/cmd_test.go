package update_test

import (
	"bytes"
	"log/slog"
	cmd "manala/cmd/update"
	"os"
	"path/filepath"
	"testing"

	"manala/app"
	"manala/app/api"
	"manala/internal/caching"
	"manala/internal/serrors"
	"manala/internal/testing/heredoc"
	"manala/internal/ui/adapters/charm"
	"manala/internal/ui/log"

	"github.com/stretchr/testify/suite"
)

type Suite struct{ suite.Suite }

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) execute(defaultRepositoryURL string, args ...string) (*bytes.Buffer, *bytes.Buffer, error) {
	stdOut := &bytes.Buffer{}
	stdErr := &bytes.Buffer{}

	ui := charm.New(nil, stdOut, stdErr)
	log := slog.New(log.NewSlogHandler(ui))

	cmd := cmd.NewCmd(
		log,
		api.New(
			log,
			caching.NewCache(""),
			api.WithDefaultRepositoryURL(defaultRepositoryURL),
		),
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

func (s *Suite) TestRecursiveProjectErrors() {
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

		serrors.Equal(s.T(), &serrors.Assertion{
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
		projectDir := filepath.FromSlash("testdata/TestRecursiveProjectErrors/InvalidProjectManifest/project")

		stdOut, stdErr, err := s.execute("",
			projectDir,
			"--recursive",
		)

		s.Empty(stdOut)
		heredoc.Equal(s.T(), `
			 • loading projects recursive…
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
		repositoryURL := filepath.FromSlash("testdata/TestRepositoryErrors/RepositoryNotFound/repository")

		stdOut, stdErr, err := s.execute("",
			projectDir,
			"--repository", repositoryURL,
		)

		s.Empty(stdOut)
		heredoc.Equal(s.T(), `
			 • loading project…
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

		serrors.Equal(s.T(), &serrors.Assertion{
			Type:    &app.NotFoundRepositoryError{},
			Message: "repository not found",
			Arguments: []any{
				"url", repositoryURL,
			},
		}, err)
	})
}

func (s *Suite) TestRepositoryCustom() {
	projectDir := filepath.FromSlash("testdata/TestRepositoryCustom/project")
	repositoryURL := filepath.FromSlash("testdata/TestRepositoryCustom/repository")

	_ = os.Remove(filepath.Join(projectDir, "file.txt"))

	stdOut, stdErr, err := s.execute("",
		projectDir,
		"--repository", repositoryURL,
	)

	s.Require().NoError(err)
	s.Empty(stdOut)
	heredoc.Equal(s.T(), `
		 • loading project…
		 • file synced                      path=file.txt
	`, stdErr)

	heredoc.EqualFile(s.T(), `
		File
	`, filepath.Join(projectDir, "file.txt"))
}

func (s *Suite) TestRepositoryConfig() {
	projectDir := filepath.FromSlash("testdata/TestRepositoryConfig/project")
	repositoryURL := filepath.FromSlash("testdata/TestRepositoryConfig/repository")

	_ = os.Remove(filepath.Join(projectDir, "file.txt"))
	_ = os.Remove(filepath.Join(projectDir, "template"))

	stdOut, stdErr, err := s.execute(repositoryURL,
		projectDir,
	)

	s.Require().NoError(err)
	s.Empty(stdOut)
	heredoc.Equal(s.T(), `
		 • loading project…
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

func (s *Suite) TestRecipeErrors() {
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

		serrors.Equal(s.T(), &serrors.Assertion{
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

func (s *Suite) TestRecipeCustom() {
	projectDir := filepath.FromSlash("testdata/TestRecipeCustom/project")
	repositoryURL := filepath.FromSlash("testdata/TestRecipeCustom/repository")

	_ = os.Remove(filepath.Join(projectDir, "file.txt"))

	stdOut, stdErr, err := s.execute(repositoryURL,
		projectDir,
		"--recipe", "recipe",
	)

	s.Require().NoError(err)
	s.Empty(stdOut)
	heredoc.Equal(s.T(), `
		 • loading project…
		 • file synced                      path=file.txt
	`, stdErr)

	heredoc.EqualFile(s.T(), `
		File
	`, filepath.Join(projectDir, "file.txt"))
}
