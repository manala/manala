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
		path := internalTesting.DataPath(s, "repository")

		manager := NewRepositoryManager(
			logger,
			core.NewRepositoryMock().
				WithPath(path).
				WithDir(path),
		)

		manifest, err := manager.LoadRecipeManifest("recipe")

		var _notFoundRecipeManifestError *core.NotFoundRecipeManifestError
		s.ErrorAs(err, &_notFoundRecipeManifestError)
		s.Nil(manifest)

		report := internalReport.NewErrorReport(err)

		reportAssert := &internalReport.Assert{
			Err: "recipe manifest not found",
			Fields: map[string]interface{}{
				"path": filepath.Join(path, "recipe"),
			},
		}
		reportAssert.Equal(&s.Suite, report)
	})

	s.Run("Directory", func() {
		path := internalTesting.DataPath(s, "repository")

		manager := NewRepositoryManager(
			logger,
			core.NewRepositoryMock().
				WithPath(path).
				WithDir(path),
		)

		manifest, err := manager.LoadRecipeManifest("recipe")

		s.Error(err)
		s.Nil(manifest)

		report := internalReport.NewErrorReport(err)

		reportAssert := &internalReport.Assert{
			Err: "recipe manifest is a directory",
			Fields: map[string]interface{}{
				"path": filepath.Join(path, "recipe", ".manala.yaml"),
			},
		}
		reportAssert.Equal(&s.Suite, report)
	})
}

func (s *ManagerSuite) TestLoadRecipeManifest() {
	logger := internalLog.New(io.Discard)

	manager := NewRepositoryManager(
		logger,
		core.NewRepositoryMock().
			WithPath(internalTesting.DataPath(s, "repository")).
			WithDir(internalTesting.DataPath(s, "repository")),
	)

	manifest, err := manager.LoadRecipeManifest("recipe")

	s.NoError(err)
	s.Equal("Recipe", manifest.Description())
	s.Equal(map[string]interface{}{"foo": "bar"}, manifest.Vars())
}

func (s *ManagerSuite) TestLoadRecipe() {
	logger := internalLog.New(io.Discard)

	manager := NewRepositoryManager(
		logger,
		core.NewRepositoryMock().
			WithPath(internalTesting.DataPath(s, "repository")).
			WithDir(internalTesting.DataPath(s, "repository")),
	)

	rec, err := manager.LoadRecipe("recipe")

	s.NoError(err)
	s.Equal(internalTesting.DataPath(s, "repository", "recipe"), rec.Path())
	s.Equal(internalTesting.DataPath(s, "repository"), rec.Repository().Path())
}
