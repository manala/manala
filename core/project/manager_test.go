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

	manager := NewManager(
		logger,
		core.NewRepositoryManagerMock(),
	)

	manifest, err := manager.LoadProjectManifest(path)

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
		path := internalTesting.DataPath(s, "project")

		proj, err := manager.CreateProject(
			path,
			core.NewRecipeMock().
				WithName("recipe").
				WithRepository(
					core.NewRepositoryMock().
						WithPath("repository"),
				),
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

	manager := NewManager(
		logger,
		core.NewRepositoryManagerMock(),
	)

	s.Run("Root", func() {
		path := internalTesting.DataPath(s, "project")
		manifestPath := filepath.Join(path, ".manala.yaml")

		_ = os.Remove(manifestPath)

		proj, err := manager.CreateProject(
			path,
			core.NewRecipeMock().
				WithName("recipe").
				WithRepository(
					core.NewRepositoryMock().
						WithPath("repository"),
				),
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

		proj, err := manager.CreateProject(
			path,
			core.NewRecipeMock().
				WithName("recipe").
				WithRepository(
					core.NewRepositoryMock().
						WithPath("repository"),
				),
			nil,
		)

		s.NotNil(proj)
		s.NoError(err)

		s.DirExists(path)
		s.FileExists(manifestPath)
	})
}

func (s *ManagerSuite) TestLoadProject() {
	logger := internalLog.New(io.Discard)

	path := internalTesting.DataPath(s, "project")

	repo := core.NewRepositoryMock()

	manager := NewManager(
		logger,
		core.NewRepositoryManagerMock().WithLoadRepository(
			repo.WithLoadRecipe(
				core.NewRecipeMock().
					WithRepository(repo),
			),
		),
	)

	project, err := manager.LoadProject(
		path,
		"repository",
		"recipe",
	)

	s.NoError(err)

	s.Equal(path, project.Path())
}
