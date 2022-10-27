package project

import (
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"io"
	"manala/core"
	internalLog "manala/internal/log"
	internalReport "manala/internal/report"
	internalTemplate "manala/internal/template"
	internalTesting "manala/internal/testing"
	"os"
	"path/filepath"
	"testing"
)

type ManagerSuite struct{ suite.Suite }

func TestManagerSuite(t *testing.T) {
	suite.Run(t, new(ManagerSuite))
}

func (s *ManagerSuite) TestLoadProjectManifestErrors() {
	logger := internalLog.New(io.Discard)

	repoManagerMock := core.NewRepositoryManagerMock()

	manager := NewManager(
		logger,
		repoManagerMock,
	)

	s.Run("Not Found", func() {
		path := internalTesting.DataPath(s, "project")

		manifest, err := manager.LoadProjectManifest(path)

		var _notFoundProjectManifestError *core.NotFoundProjectManifestError
		s.ErrorAs(err, &_notFoundProjectManifestError)
		s.Nil(manifest)

		report := internalReport.NewErrorReport(err)

		reportAssert := &internalReport.Assert{
			Err: "project manifest not found",
			Fields: map[string]interface{}{
				"path": path,
			},
		}
		reportAssert.Equal(&s.Suite, report)
	})

	s.Run("Directory", func() {
		path := internalTesting.DataPath(s, "project")
		manifestPath := filepath.Join(path, ".manala.yaml")

		manifest, err := manager.LoadProjectManifest(path)

		s.Error(err)
		s.Nil(manifest)

		report := internalReport.NewErrorReport(err)

		reportAssert := &internalReport.Assert{
			Err: "project manifest is a directory",
			Fields: map[string]interface{}{
				"path": manifestPath,
			},
		}
		reportAssert.Equal(&s.Suite, report)
	})
}

func (s *ManagerSuite) TestLoadProjectManifest() {
	logger := internalLog.New(io.Discard)

	path := internalTesting.DataPath(s, "project")

	repoManagerMock := core.NewRepositoryManagerMock()

	manager := NewManager(
		logger,
		repoManagerMock,
	)

	manifest, err := manager.LoadProjectManifest(path)

	s.NoError(err)
	s.Equal("recipe", manifest.Recipe())
	s.Equal("repository", manifest.Repository())
	s.Equal(map[string]interface{}{"foo": "bar"}, manifest.Vars())
}

func (s *ManagerSuite) TestCreateProjectErrors() {
	logger := internalLog.New(io.Discard)

	repoManagerMock := core.NewRepositoryManagerMock()

	manager := NewManager(
		logger,
		repoManagerMock,
	)

	s.Run("File", func() {
		path := internalTesting.DataPath(s, "project")

		repoMock := core.NewRepositoryMock()
		repoMock.
			On("Path").Return("repository").
			On("Source").Return("repository")

		recMock := core.NewRecipeMock()
		recMock.
			On("Name").Return("recipe").
			On("Repository").Return(repoMock).
			On("ProjectManifestTemplate").Return(internalTemplate.NewTemplate())

		proj, err := manager.CreateProject(
			path,
			recMock,
			nil,
		)

		s.Nil(proj)
		s.Error(err)

		report := internalReport.NewErrorReport(err)

		reportAssert := &internalReport.Assert{
			Err: "project is not a directory",
			Fields: map[string]interface{}{
				"path": path,
			},
		}
		reportAssert.Equal(&s.Suite, report)
	})
}

func (s *ManagerSuite) TestCreateProject() {
	logger := internalLog.New(io.Discard)

	repoManagerMock := core.NewRepositoryManagerMock()

	manager := NewManager(
		logger,
		repoManagerMock,
	)

	s.Run("Root", func() {
		path := internalTesting.DataPath(s, "project")
		manifestPath := filepath.Join(path, ".manala.yaml")

		_ = os.Remove(manifestPath)

		repoMock := core.NewRepositoryMock()
		repoMock.
			On("Path").Return("repository").
			On("Source").Return("repository")

		recMock := core.NewRecipeMock()
		recMock.
			On("Name").Return("recipe").
			On("Repository").Return(repoMock).
			On("ProjectManifestTemplate").Return(internalTemplate.NewTemplate())

		proj, err := manager.CreateProject(
			path,
			recMock,
			nil,
		)

		s.NotNil(proj)
		s.NoError(err)

		s.FileExists(manifestPath)
	})

	s.Run("Directory", func() {
		path := internalTesting.DataPath(s, "project")
		manifestPath := filepath.Join(path, ".manala.yaml")

		_ = os.RemoveAll(path)

		repoMock := core.NewRepositoryMock()
		repoMock.
			On("Path").Return("repository").
			On("Source").Return("repository")

		recMock := core.NewRecipeMock()
		recMock.
			On("Name").Return("recipe").
			On("Repository").Return(repoMock).
			On("ProjectManifestTemplate").Return(internalTemplate.NewTemplate())

		proj, err := manager.CreateProject(
			path,
			recMock,
			nil,
		)

		s.NotNil(proj)
		s.NoError(err)

		s.DirExists(path)
		s.FileExists(manifestPath)
	})
}

func (s *ManagerSuite) TestLoadProjectErrors() {
	logger := internalLog.New(io.Discard)

	s.Run("Vars", func() {
		path := internalTesting.DataPath(s, "project")
		manifestPath := filepath.Join(path, ".manala.yaml")

		repoMock := core.NewRepositoryMock()

		repoManagerMock := core.NewRepositoryManagerMock()
		repoManagerMock.
			On("LoadRepository", mock.Anything).Return(repoMock, nil)

		recMock := core.NewRecipeMock()
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

		repoMock.
			On("LoadRecipe", mock.Anything).Return(recMock, nil)

		manager := NewManager(
			logger,
			repoManagerMock,
		)

		project, err := manager.LoadProject(
			path,
			"repository",
			"recipe",
		)

		s.Nil(project)
		s.Error(err)

		report := internalReport.NewErrorReport(err)

		reportAssert := &internalReport.Assert{
			Err: "invalid project manifest vars",
			Fields: map[string]interface{}{
				"path": manifestPath,
			},
			Reports: []internalReport.Assert{
				{
					Message: "invalid type",
					Fields: map[string]interface{}{
						"line":     5,
						"column":   6,
						"expected": "integer",
						"given":    "string",
					},
				},
			},
		}
		reportAssert.Equal(&s.Suite, report)
	})
}

func (s *ManagerSuite) TestLoadProject() {
	logger := internalLog.New(io.Discard)

	path := internalTesting.DataPath(s, "project")

	repoMock := core.NewRepositoryMock()

	repoManagerMock := core.NewRepositoryManagerMock()
	repoManagerMock.
		On("LoadRepository", mock.Anything).Return(repoMock, nil)

	recMock := core.NewRecipeMock()
	recMock.
		On("Vars").Return(map[string]interface{}{}).
		On("Schema").Return(map[string]interface{}{}).
		On("Repository").Return(repoMock)

	repoMock.
		On("LoadRecipe", mock.Anything).Return(recMock, nil)

	manager := NewManager(
		logger,
		repoManagerMock,
	)

	project, err := manager.LoadProject(
		path,
		"repository",
		"recipe",
	)

	s.NoError(err)

	s.Equal(path, project.Path())
}
