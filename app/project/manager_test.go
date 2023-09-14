package project

import (
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"io"
	"log/slog"
	"manala/app"
	"manala/internal/schema"
	"manala/internal/serrors"
	"manala/internal/template"
	"manala/internal/validator"
	"os"
	"path/filepath"
	"testing"
)

type ManagerSuite struct{ suite.Suite }

func TestManagerSuite(t *testing.T) {
	suite.Run(t, new(ManagerSuite))
}

func (s *ManagerSuite) TestIsProject() {
	repositoryManagerMock := &app.RepositoryManagerMock{}

	recipeManagerMock := &app.RecipeManagerMock{}

	manager := NewManager(
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		repositoryManagerMock,
		recipeManagerMock,
	)

	s.Run("True", func() {
		dir := filepath.FromSlash("testdata/ManagerSuite/TestIsProject/True")

		s.True(manager.IsProject(dir))
	})

	s.Run("False", func() {
		dir := filepath.FromSlash("testdata/ManagerSuite/TestIsProject/False")

		s.False(manager.IsProject(dir))
	})
}

func (s *ManagerSuite) TestLoadManifestErrors() {
	repositoryManagerMock := &app.RepositoryManagerMock{}

	recipeManagerMock := &app.RecipeManagerMock{}

	manager := NewManager(
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		repositoryManagerMock,
		recipeManagerMock,
	)

	s.Run("NotFound", func() {
		dir := filepath.FromSlash("testdata/ManagerSuite/TestLoadManifestErrors/NotFound/project")

		manifest, err := manager.loadManifest(filepath.Join(dir, ".manala.yaml"))

		s.Nil(manifest)

		serrors.Equal(s.Assert(), &serrors.Assert{
			Type:    &app.NotFoundProjectManifestError{},
			Message: "project manifest not found",
			Arguments: []any{
				"file", filepath.Join(dir, ".manala.yaml"),
			},
		}, err)
	})

	s.Run("Directory", func() {
		dir := filepath.FromSlash("testdata/ManagerSuite/TestLoadManifestErrors/Directory/project")

		manifest, err := manager.loadManifest(filepath.Join(dir, ".manala.yaml"))

		s.Nil(manifest)

		serrors.Equal(s.Assert(), &serrors.Assert{
			Type:    serrors.Error{},
			Message: "project manifest is a directory",
			Arguments: []any{
				"dir", filepath.Join(dir, ".manala.yaml"),
			},
		}, err)
	})
}

func (s *ManagerSuite) TestLoadManifest() {
	dir := filepath.FromSlash("testdata/ManagerSuite/TestLoadManifest/project")

	repositoryManagerMock := &app.RepositoryManagerMock{}

	recipeManagerMock := &app.RecipeManagerMock{}

	manager := NewManager(
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		repositoryManagerMock,
		recipeManagerMock,
	)

	manifest, err := manager.loadManifest(filepath.Join(dir, ".manala.yaml"))

	s.NoError(err)

	s.Equal("recipe", manifest.Recipe())
	s.Equal("repository", manifest.Repository())
	s.Equal(map[string]any{"foo": "bar"}, manifest.Vars())
}

func (s *ManagerSuite) TestCreateProjectErrors() {
	repositoryManagerMock := &app.RepositoryManagerMock{}

	recipeManagerMock := &app.RecipeManagerMock{}

	manager := NewManager(
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		repositoryManagerMock,
		recipeManagerMock,
	)

	s.Run("File", func() {
		dir := filepath.FromSlash("testdata/ManagerSuite/TestCreateProjectErrors/File/project")

		repositoryMock := &app.RepositoryMock{}
		repositoryMock.
			On("Url").Return("repository")

		recipeMock := &app.RecipeMock{}
		recipeMock.
			On("Name").Return("recipe").
			On("Repository").Return(repositoryMock).
			On("ProjectManifestTemplate").Return(template.NewTemplate())

		project, err := manager.CreateProject(
			dir,
			recipeMock,
			nil,
		)

		s.Nil(project)

		serrors.Equal(s.Assert(), &serrors.Assert{
			Type:    serrors.Error{},
			Message: "project is not a directory",
			Arguments: []any{
				"path", dir,
			},
		}, err)
	})
}

func (s *ManagerSuite) TestCreateProject() {
	repositoryManagerMock := &app.RepositoryManagerMock{}

	recipeManagerMock := &app.RecipeManagerMock{}

	manager := NewManager(
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		repositoryManagerMock,
		recipeManagerMock,
	)

	s.Run("Root", func() {
		dir := filepath.FromSlash("testdata/ManagerSuite/TestCreateProject/Root/project")

		_ = os.Remove(filepath.Join(dir, ".manala.yaml"))

		repositoryMock := &app.RepositoryMock{}
		repositoryMock.
			On("Url").Return("repository")

		recipeMock := &app.RecipeMock{}
		recipeMock.
			On("Name").Return("recipe").
			On("Repository").Return(repositoryMock).
			On("ProjectManifestTemplate").Return(template.NewTemplate())

		project, err := manager.CreateProject(
			dir,
			recipeMock,
			nil,
		)

		s.NotNil(project)

		s.NoError(err)

		s.FileExists(filepath.Join(dir, ".manala.yaml"))
	})

	s.Run("Directory", func() {
		dir := filepath.FromSlash("testdata/ManagerSuite/TestCreateProject/Directory/project")

		_ = os.RemoveAll(dir)

		repositoryMock := &app.RepositoryMock{}
		repositoryMock.
			On("Url").Return("repository")

		recipeMock := &app.RecipeMock{}
		recipeMock.
			On("Name").Return("recipe").
			On("Repository").Return(repositoryMock).
			On("ProjectManifestTemplate").Return(template.NewTemplate())

		project, err := manager.CreateProject(
			dir,
			recipeMock,
			nil,
		)

		s.NotNil(project)

		s.NoError(err)

		s.DirExists(dir)
		s.FileExists(filepath.Join(dir, ".manala.yaml"))
	})
}

func (s *ManagerSuite) TestLoadProjectErrors() {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))

	s.Run("Vars", func() {
		dir := filepath.FromSlash("testdata/ManagerSuite/TestLoadProjectErrors/Vars/project")

		repositoryMock := &app.RepositoryMock{}

		repositoryManagerMock := &app.RepositoryManagerMock{}
		repositoryManagerMock.
			On("LoadRepository", mock.Anything).Return(repositoryMock, nil)

		recipeMock := &app.RecipeMock{}
		recipeMock.
			On("Vars").Return(map[string]any{}).
			On("ProjectValidator", mock.Anything).Return(
			schema.NewValidator(
				schema.Schema{
					"type": "object",
					"properties": map[string]any{
						"foo": map[string]any{
							"type": "integer",
						},
					},
				},
			))

		recipeManagerMock := &app.RecipeManagerMock{}
		recipeManagerMock.
			On("LoadRecipe", mock.Anything, mock.Anything).Return(recipeMock, nil)

		manager := NewManager(
			log,
			repositoryManagerMock,
			recipeManagerMock,
		)

		project, err := manager.LoadProject(dir)

		s.Nil(project)

		serrors.Equal(s.Assert(), &serrors.Assert{
			Type:    serrors.Error{},
			Message: "invalid project manifest vars",
			Arguments: []any{
				"file", filepath.Join(dir, ".manala.yaml"),
			},
			Errors: []*serrors.Assert{
				{
					Type:    serrors.Error{},
					Message: "invalid type",
					Arguments: []any{
						"expected", "integer",
						"actual", "string",
						"path", "foo",
						"line", 5,
						"column", 6,
					},
					Details: `
						   2 |   recipe: recipe
						   3 |   repository: repository
						   4 |
						>  5 | foo: bar
						            ^
					`,
				},
			},
		}, err)
	})
}

func (s *ManagerSuite) TestLoadProject() {
	dir := filepath.FromSlash("testdata/ManagerSuite/TestLoadProject/project")

	repositoryMock := &app.RepositoryMock{}

	repositoryManagerMock := &app.RepositoryManagerMock{}
	repositoryManagerMock.
		On("LoadRepository", mock.Anything).Return(repositoryMock, nil)

	recipeMock := &app.RecipeMock{}
	recipeMock.
		On("Vars").Return(map[string]any{}).
		On("Repository").Return(repositoryMock).
		On("ProjectValidator", mock.Anything).Return(validator.New())

	recipeManagerMock := &app.RecipeManagerMock{}
	recipeManagerMock.
		On("LoadRecipe", mock.Anything, mock.Anything).Return(recipeMock, nil)

	manager := NewManager(
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		repositoryManagerMock,
		recipeManagerMock,
	)

	project, err := manager.LoadProject(dir)

	s.NoError(err)

	s.Equal(dir, project.Dir())
}
