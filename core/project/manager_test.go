package project

import (
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"io"
	"log/slog"
	"manala/app/mocks"
	"manala/core"
	"manala/internal/errors/serrors"
	"manala/internal/template"
	"manala/internal/testing/heredoc"
	"manala/internal/validation"
	"manala/internal/yaml"
	"os"
	"path/filepath"
	"testing"
)

type ManagerSuite struct{ suite.Suite }

func TestManagerSuite(t *testing.T) {
	suite.Run(t, new(ManagerSuite))
}

func (s *ManagerSuite) TestIsProject() {
	repoManagerMock := &mocks.RepositoryManagerMock{}

	recManagerMock := &mocks.RecipeManagerMock{}

	manager := NewManager(
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		repoManagerMock,
		recManagerMock,
	)

	s.Run("True", func() {
		projDir := filepath.FromSlash("testdata/ManagerSuite/TestIsProject/True")

		s.True(manager.IsProject(projDir))
	})

	s.Run("False", func() {
		projDir := filepath.FromSlash("testdata/ManagerSuite/TestIsProject/False")

		s.False(manager.IsProject(projDir))
	})
}

func (s *ManagerSuite) TestLoadManifestErrors() {
	repoManagerMock := &mocks.RepositoryManagerMock{}

	recManagerMock := &mocks.RecipeManagerMock{}

	manager := NewManager(
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		repoManagerMock,
		recManagerMock,
	)

	s.Run("NotFound", func() {
		projDir := filepath.FromSlash("testdata/ManagerSuite/TestLoadManifestErrors/NotFound/project")

		projMan, err := manager.loadManifest(filepath.Join(projDir, ".manala.yaml"))

		s.Nil(projMan)

		serrors.Equal(s.Assert(), &serrors.Assert{
			Type:    &core.NotFoundProjectManifestError{},
			Message: "project manifest not found",
			Arguments: []any{
				"file", filepath.Join(projDir, ".manala.yaml"),
			},
		}, err)
	})

	s.Run("Directory", func() {
		projDir := filepath.FromSlash("testdata/ManagerSuite/TestLoadManifestErrors/Directory/project")

		projMan, err := manager.loadManifest(filepath.Join(projDir, ".manala.yaml"))

		s.Nil(projMan)

		serrors.Equal(s.Assert(), &serrors.Assert{
			Type:    &serrors.Error{},
			Message: "project manifest is a directory",
			Arguments: []any{
				"dir", filepath.Join(projDir, ".manala.yaml"),
			},
		}, err)
	})
}

func (s *ManagerSuite) TestLoadManifest() {
	projDir := filepath.FromSlash("testdata/ManagerSuite/TestLoadManifest/project")

	repoManagerMock := &mocks.RepositoryManagerMock{}

	recManagerMock := &mocks.RecipeManagerMock{}

	manager := NewManager(
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		repoManagerMock,
		recManagerMock,
	)

	projMan, err := manager.loadManifest(filepath.Join(projDir, ".manala.yaml"))

	s.NoError(err)

	s.Equal("recipe", projMan.Recipe())
	s.Equal("repository", projMan.Repository())
	s.Equal(map[string]interface{}{"foo": "bar"}, projMan.Vars())
}

func (s *ManagerSuite) TestCreateProjectErrors() {
	repoManagerMock := &mocks.RepositoryManagerMock{}

	recManagerMock := &mocks.RecipeManagerMock{}

	manager := NewManager(
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		repoManagerMock,
		recManagerMock,
	)

	s.Run("File", func() {
		projDir := filepath.FromSlash("testdata/ManagerSuite/TestCreateProjectErrors/File/project")

		repoMock := &mocks.RepositoryMock{}
		repoMock.
			On("Url").Return("repository")

		recMock := &mocks.RecipeMock{}
		recMock.
			On("Name").Return("recipe").
			On("Repository").Return(repoMock).
			On("ProjectManifestTemplate").Return(template.NewTemplate())

		proj, err := manager.CreateProject(
			projDir,
			recMock,
			nil,
		)

		s.Nil(proj)

		serrors.Equal(s.Assert(), &serrors.Assert{
			Type:    &serrors.Error{},
			Message: "project is not a directory",
			Arguments: []any{
				"path", projDir,
			},
		}, err)
	})
}

func (s *ManagerSuite) TestCreateProject() {
	repoManagerMock := &mocks.RepositoryManagerMock{}

	recManagerMock := &mocks.RecipeManagerMock{}

	manager := NewManager(
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		repoManagerMock,
		recManagerMock,
	)

	s.Run("Root", func() {
		projDir := filepath.FromSlash("testdata/ManagerSuite/TestCreateProject/Root/project")

		_ = os.Remove(filepath.Join(projDir, ".manala.yaml"))

		repoMock := &mocks.RepositoryMock{}
		repoMock.
			On("Url").Return("repository")

		recMock := &mocks.RecipeMock{}
		recMock.
			On("Name").Return("recipe").
			On("Repository").Return(repoMock).
			On("ProjectManifestTemplate").Return(template.NewTemplate())

		proj, err := manager.CreateProject(
			projDir,
			recMock,
			nil,
		)

		s.NotNil(proj)

		s.NoError(err)

		s.FileExists(filepath.Join(projDir, ".manala.yaml"))
	})

	s.Run("Directory", func() {
		projDir := filepath.FromSlash("testdata/ManagerSuite/TestCreateProject/Directory/project")

		_ = os.RemoveAll(projDir)

		repoMock := &mocks.RepositoryMock{}
		repoMock.
			On("Url").Return("repository")

		recMock := &mocks.RecipeMock{}
		recMock.
			On("Name").Return("recipe").
			On("Repository").Return(repoMock).
			On("ProjectManifestTemplate").Return(template.NewTemplate())

		proj, err := manager.CreateProject(
			projDir,
			recMock,
			nil,
		)

		s.NotNil(proj)

		s.NoError(err)

		s.DirExists(projDir)
		s.FileExists(filepath.Join(projDir, ".manala.yaml"))
	})
}

func (s *ManagerSuite) TestLoadProjectErrors() {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))

	s.Run("Vars", func() {
		projDir := filepath.FromSlash("testdata/ManagerSuite/TestLoadProjectErrors/Vars/project")

		repoMock := &mocks.RepositoryMock{}

		repoManagerMock := &mocks.RepositoryManagerMock{}
		repoManagerMock.
			On("LoadRepository", mock.Anything).Return(repoMock, nil)

		recMock := &mocks.RecipeMock{}
		recMock.
			On("Vars").Return(map[string]interface{}{}).
			On("Schema").Return(map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"foo": map[string]interface{}{
					"type": "integer",
				},
			},
		})

		recManagerMock := &mocks.RecipeManagerMock{}
		recManagerMock.
			On("LoadRecipe", mock.Anything, mock.Anything).Return(recMock, nil)

		manager := NewManager(
			log,
			repoManagerMock,
			recManagerMock,
		)

		project, err := manager.LoadProject(projDir)

		s.Nil(project)

		serrors.Equal(s.Assert(), &serrors.Assert{
			Type:    &validation.Error{},
			Message: "invalid project manifest vars",
			Arguments: []any{
				"file", filepath.Join(projDir, ".manala.yaml"),
			},
			Errors: []*serrors.Assert{
				{
					Type:    &yaml.NodeValidationResultError{},
					Message: "invalid type",
					Arguments: []any{
						"expected", "integer",
						"given", "string",
						"line", 5,
						"column", 6,
					},
					Details: heredoc.Doc(`
						   2 |   recipe: recipe
						   3 |   repository: repository
						   4 |
						>  5 | foo: bar
						            ^
					`),
				},
			},
		}, err)
	})
}

func (s *ManagerSuite) TestLoadProject() {
	projDir := filepath.FromSlash("testdata/ManagerSuite/TestLoadProject/project")

	repoMock := &mocks.RepositoryMock{}

	repoManagerMock := &mocks.RepositoryManagerMock{}
	repoManagerMock.
		On("LoadRepository", mock.Anything).Return(repoMock, nil)

	recMock := &mocks.RecipeMock{}
	recMock.
		On("Vars").Return(map[string]interface{}{}).
		On("Schema").Return(map[string]interface{}{}).
		On("Repository").Return(repoMock)

	recManagerMock := &mocks.RecipeManagerMock{}
	recManagerMock.
		On("LoadRecipe", mock.Anything, mock.Anything).Return(recMock, nil)

	manager := NewManager(
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		repoManagerMock,
		recManagerMock,
	)

	project, err := manager.LoadProject(projDir)

	s.NoError(err)

	s.Equal(projDir, project.Dir())
}
