package recipe

import (
	"github.com/stretchr/testify/suite"
	"io"
	"manala/core"
	internalLog "manala/internal/log"
	internalReport "manala/internal/report"
	internalTesting "manala/internal/testing"
	"path/filepath"
	"testing"
)

type ManagerSuite struct{ suite.Suite }

func TestManagerSuite(t *testing.T) {
	suite.Run(t, new(ManagerSuite))
}

func (s *ManagerSuite) TestLoadRecipeManifestErrors() {
	logger := internalLog.New(io.Discard)

	s.Run("Not Found", func() {
		repoPath := internalTesting.DataPath(s, "repository")
		recPath := filepath.Join(repoPath, "recipe")

		repoMock := core.NewRepositoryMock()
		repoMock.
			On("Path").Return(repoPath).
			On("Dir").Return(repoPath)

		manager := NewRepositoryManager(
			logger,
			repoMock,
		)

		manifest, err := manager.LoadRecipeManifest("recipe")

		var _notFoundRecipeManifestError *core.NotFoundRecipeManifestError
		s.ErrorAs(err, &_notFoundRecipeManifestError)
		s.Nil(manifest)

		report := internalReport.NewErrorReport(err)

		reportAssert := &internalReport.Assert{
			Err: "recipe manifest not found",
			Fields: map[string]interface{}{
				"path": recPath,
			},
		}
		reportAssert.Equal(&s.Suite, report)
	})

	s.Run("Directory", func() {
		repoPath := internalTesting.DataPath(s, "repository")
		manifestPath := filepath.Join(repoPath, "recipe", ".manala.yaml")

		repoMock := core.NewRepositoryMock()
		repoMock.
			On("Path").Return(repoPath).
			On("Dir").Return(repoPath)

		manager := NewRepositoryManager(
			logger,
			repoMock,
		)

		manifest, err := manager.LoadRecipeManifest("recipe")

		s.Error(err)
		s.Nil(manifest)

		report := internalReport.NewErrorReport(err)

		reportAssert := &internalReport.Assert{
			Err: "recipe manifest is a directory",
			Fields: map[string]interface{}{
				"path": manifestPath,
			},
		}
		reportAssert.Equal(&s.Suite, report)
	})
}

func (s *ManagerSuite) TestLoadRecipeManifest() {
	logger := internalLog.New(io.Discard)

	repoPath := internalTesting.DataPath(s, "repository")

	repoMock := core.NewRepositoryMock()
	repoMock.
		On("Path").Return(repoPath).
		On("Dir").Return(repoPath)

	manager := NewRepositoryManager(
		logger,
		repoMock,
	)

	manifest, err := manager.LoadRecipeManifest("recipe")

	s.NoError(err)
	s.Equal("Recipe", manifest.Description())
	s.Equal(map[string]interface{}{"foo": "bar"}, manifest.Vars())
}

func (s *ManagerSuite) TestLoadRecipe() {
	logger := internalLog.New(io.Discard)

	repoPath := internalTesting.DataPath(s, "repository")
	recPath := filepath.Join(repoPath, "recipe")

	repoMock := core.NewRepositoryMock()
	repoMock.
		On("Path").Return(repoPath).
		On("Dir").Return(repoPath)

	manager := NewRepositoryManager(
		logger,
		repoMock,
	)

	rec, err := manager.LoadRecipe("recipe")

	s.NoError(err)
	s.Equal(recPath, rec.Path())
	s.Equal(repoPath, rec.Repository().Path())
}
