package project

import (
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"io"
	"manala/app/mocks"
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

func (s *ManagerSuite) TestIsProject() {
	log := internalLog.New(io.Discard)

	repoManagerMock := mocks.MockRepositoryManager()

	recManagerMock := mocks.MockRecipeManager()

	manager := NewManager(
		log,
		repoManagerMock,
		recManagerMock,
	)

	s.Run("True", func() {
		projDir := internalTesting.DataPath(s)

		s.True(manager.IsProject(projDir))
	})

	s.Run("False", func() {
		projDir := internalTesting.DataPath(s)

		s.False(manager.IsProject(projDir))
	})
}

func (s *ManagerSuite) TestLoadManifestErrors() {
	log := internalLog.New(io.Discard)

	repoManagerMock := mocks.MockRepositoryManager()

	recManagerMock := mocks.MockRecipeManager()

	manager := NewManager(
		log,
		repoManagerMock,
		recManagerMock,
	)

	s.Run("Not Found", func() {
		projDir := internalTesting.DataPath(s, "project")
		projManFile := filepath.Join(projDir, ".manala.yaml")

		projMan, err := manager.loadManifest(projManFile)

		var _notFoundProjectManifestError *core.NotFoundProjectManifestError
		s.ErrorAs(err, &_notFoundProjectManifestError)
		s.Nil(projMan)

		report := internalReport.NewErrorReport(err)

		reportAssert := &internalReport.Assert{
			Err: "project manifest not found",
			Fields: map[string]interface{}{
				"file": projManFile,
			},
		}
		reportAssert.Equal(&s.Suite, report)
	})

	s.Run("Directory", func() {
		projDir := internalTesting.DataPath(s, "project")
		projManFile := filepath.Join(projDir, ".manala.yaml")

		projMan, err := manager.loadManifest(projManFile)

		s.Error(err)
		s.Nil(projMan)

		report := internalReport.NewErrorReport(err)

		reportAssert := &internalReport.Assert{
			Err: "project manifest is a directory",
			Fields: map[string]interface{}{
				"dir": projManFile,
			},
		}
		reportAssert.Equal(&s.Suite, report)
	})
}

func (s *ManagerSuite) TestLoadManifest() {
	log := internalLog.New(io.Discard)

	projDir := internalTesting.DataPath(s, "project")
	projManFile := filepath.Join(projDir, ".manala.yaml")

	repoManagerMock := mocks.MockRepositoryManager()

	recManagerMock := mocks.MockRecipeManager()

	manager := NewManager(
		log,
		repoManagerMock,
		recManagerMock,
	)

	projMan, err := manager.loadManifest(projManFile)

	s.NoError(err)
	s.Equal("recipe", projMan.Recipe())
	s.Equal("repository", projMan.Repository())
	s.Equal(map[string]interface{}{"foo": "bar"}, projMan.Vars())
}

func (s *ManagerSuite) TestCreateProjectErrors() {
	log := internalLog.New(io.Discard)

	repoManagerMock := mocks.MockRepositoryManager()

	recManagerMock := mocks.MockRecipeManager()

	manager := NewManager(
		log,
		repoManagerMock,
		recManagerMock,
	)

	s.Run("File", func() {
		projDir := internalTesting.DataPath(s, "project")

		repoMock := mocks.MockRepository()
		repoMock.
			On("Url").Return("repository")

		recMock := mocks.MockRecipe()
		recMock.
			On("Name").Return("recipe").
			On("Repository").Return(repoMock).
			On("ProjectManifestTemplate").Return(internalTemplate.NewTemplate())

		proj, err := manager.CreateProject(
			projDir,
			recMock,
			nil,
		)

		s.Nil(proj)
		s.Error(err)

		report := internalReport.NewErrorReport(err)

		reportAssert := &internalReport.Assert{
			Err: "project is not a directory",
			Fields: map[string]interface{}{
				"file": projDir,
			},
		}
		reportAssert.Equal(&s.Suite, report)
	})
}

func (s *ManagerSuite) TestCreateProject() {
	log := internalLog.New(io.Discard)

	repoManagerMock := mocks.MockRepositoryManager()

	recManagerMock := mocks.MockRecipeManager()

	manager := NewManager(
		log,
		repoManagerMock,
		recManagerMock,
	)

	s.Run("Root", func() {
		projDir := internalTesting.DataPath(s, "project")
		projManFile := filepath.Join(projDir, ".manala.yaml")

		_ = os.Remove(projManFile)

		repoMock := mocks.MockRepository()
		repoMock.
			On("Url").Return("repository")

		recMock := mocks.MockRecipe()
		recMock.
			On("Name").Return("recipe").
			On("Repository").Return(repoMock).
			On("ProjectManifestTemplate").Return(internalTemplate.NewTemplate())

		proj, err := manager.CreateProject(
			projDir,
			recMock,
			nil,
		)

		s.NotNil(proj)
		s.NoError(err)

		s.FileExists(projManFile)
	})

	s.Run("Directory", func() {
		projDir := internalTesting.DataPath(s, "project")
		projManFile := filepath.Join(projDir, ".manala.yaml")

		_ = os.RemoveAll(projDir)

		repoMock := mocks.MockRepository()
		repoMock.
			On("Url").Return("repository")

		recMock := mocks.MockRecipe()
		recMock.
			On("Name").Return("recipe").
			On("Repository").Return(repoMock).
			On("ProjectManifestTemplate").Return(internalTemplate.NewTemplate())

		proj, err := manager.CreateProject(
			projDir,
			recMock,
			nil,
		)

		s.NotNil(proj)
		s.NoError(err)

		s.DirExists(projDir)
		s.FileExists(projManFile)
	})
}

func (s *ManagerSuite) TestLoadProjectErrors() {
	log := internalLog.New(io.Discard)

	s.Run("Vars", func() {
		projDir := internalTesting.DataPath(s, "project")
		projManFile := filepath.Join(projDir, ".manala.yaml")

		repoMock := mocks.MockRepository()

		repoManagerMock := mocks.MockRepositoryManager()
		repoManagerMock.
			On("LoadRepository", mock.Anything).Return(repoMock, nil)

		recMock := mocks.MockRecipe()
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

		recManagerMock := mocks.MockRecipeManager()
		recManagerMock.
			On("LoadRecipe", mock.Anything, mock.Anything).Return(recMock, nil)

		manager := NewManager(
			log,
			repoManagerMock,
			recManagerMock,
		)

		project, err := manager.LoadProject(projDir)

		s.Nil(project)
		s.Error(err)

		report := internalReport.NewErrorReport(err)

		reportAssert := &internalReport.Assert{
			Err: "invalid project manifest vars",
			Fields: map[string]interface{}{
				"file": projManFile,
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
	log := internalLog.New(io.Discard)

	projDir := internalTesting.DataPath(s, "project")

	repoMock := mocks.MockRepository()

	repoManagerMock := mocks.MockRepositoryManager()
	repoManagerMock.
		On("LoadRepository", mock.Anything).Return(repoMock, nil)

	recMock := mocks.MockRecipe()
	recMock.
		On("Vars").Return(map[string]interface{}{}).
		On("Schema").Return(map[string]interface{}{}).
		On("Repository").Return(repoMock)

	recManagerMock := mocks.MockRecipeManager()
	recManagerMock.
		On("LoadRecipe", mock.Anything, mock.Anything).Return(recMock, nil)

	manager := NewManager(
		log,
		repoManagerMock,
		recManagerMock,
	)

	project, err := manager.LoadProject(projDir)

	s.NoError(err)

	s.Equal(projDir, project.Dir())
}
