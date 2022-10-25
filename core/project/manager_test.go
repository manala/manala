package project

import (
	"github.com/stretchr/testify/suite"
	"io"
	"manala/core"
	internalLog "manala/internal/log"
	internalReport "manala/internal/report"
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
	manager := NewManager(
		logger,
		core.NewRepositoryManagerMock(),
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

		manifest, err := manager.LoadProjectManifest(path)

		s.Error(err)
		s.Nil(manifest)

		report := internalReport.NewErrorReport(err)

		reportAssert := &internalReport.Assert{
			Err: "project manifest is a directory",
			Fields: map[string]interface{}{
				"path": filepath.Join(path, ".manala.yaml"),
			},
		}
		reportAssert.Equal(&s.Suite, report)
	})
}

func (s *ManagerSuite) TestLoadProjectManifest() {
	logger := internalLog.New(io.Discard)
	manager := NewManager(
		logger,
		core.NewRepositoryManagerMock(),
	)

	manifest, err := manager.LoadProjectManifest(internalTesting.DataPath(s, "project"))

	s.NoError(err)
	s.Equal("recipe", manifest.Recipe())
	s.Equal("repository", manifest.Repository())
	s.Equal(map[string]interface{}{"foo": "bar"}, manifest.Vars())
}

func (s *ManagerSuite) TestCreateProjectErrors() {
	logger := internalLog.New(io.Discard)
	manager := NewManager(
		logger,
		core.NewRepositoryManagerMock(),
	)

	s.Run("File", func() {
		path := internalTesting.DataPath(s)

		rec := core.NewRecipeMock().
			WithName("recipe").
			WithRepository(
				core.NewRepositoryMock().
					WithPath("repository"),
			)

		proj, err := manager.CreateProject(
			filepath.Join(path, "project"),
			rec,
			nil,
		)

		s.Nil(proj)
		s.Error(err)

		report := internalReport.NewErrorReport(err)

		reportAssert := &internalReport.Assert{
			Err: "project is not a directory",
			Fields: map[string]interface{}{
				"path": filepath.Join(path, "project"),
			},
		}
		reportAssert.Equal(&s.Suite, report)
	})
}

func (s *ManagerSuite) TestCreateProject() {
	logger := internalLog.New(io.Discard)
	manager := NewManager(
		logger,
		core.NewRepositoryManagerMock(),
	)

	s.Run("Root", func() {
		path := internalTesting.DataPath(s)
		manifestPath := filepath.Join(path, ".manala.yaml")

		_ = os.Remove(manifestPath)

		rec := core.NewRecipeMock().
			WithName("recipe").
			WithRepository(
				core.NewRepositoryMock().
					WithPath("repository"),
			)

		proj, err := manager.CreateProject(
			path,
			rec,
			nil,
		)

		s.NotNil(proj)
		s.NoError(err)

		s.FileExists(manifestPath)
	})

	s.Run("Directory", func() {
		path := internalTesting.DataPath(s)
		directoryPath := filepath.Join(path, "directory")
		manifestPath := filepath.Join(directoryPath, ".manala.yaml")

		_ = os.RemoveAll(directoryPath)

		rec := core.NewRecipeMock().
			WithName("recipe").
			WithRepository(
				core.NewRepositoryMock().
					WithPath("repository"),
			)

		proj, err := manager.CreateProject(
			filepath.Join(path, "directory"),
			rec,
			nil,
		)

		s.NotNil(proj)
		s.NoError(err)

		s.DirExists(directoryPath)
		s.FileExists(manifestPath)
	})
}

func (s *ManagerSuite) TestLoadProject() {
	logger := internalLog.New(io.Discard)

	repo := core.NewRepositoryMock()

	rec := core.NewRecipeMock().
		WithRepository(repo)

	manager := NewManager(
		logger,
		core.NewRepositoryManagerMock().
			WithLoadRepository(
				repo.WithLoadRecipe(rec),
			),
	)

	project, err := manager.LoadProject(
		internalTesting.DataPath(s, "project"),
		"repository",
		"recipe",
	)

	s.NoError(err)

	s.Equal(internalTesting.DataPath(s, "project"), project.Path())
}
