package recipe

import (
	"github.com/stretchr/testify/suite"
	"io"
	"manala/app/interfaces"
	"manala/app/mocks"
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

func (s *ManagerSuite) TestLoadManifestErrors() {
	log := internalLog.New(io.Discard)

	s.Run("Not Found", func() {
		repoUrl := internalTesting.DataPath(s, "repository")

		recDir := filepath.Join(repoUrl, "recipe")
		recManFile := filepath.Join(recDir, ".manala.yaml")

		manager := NewManager(log)

		recMan, err := manager.loadManifest(recManFile)

		var _notFoundRecipeManifestError *core.NotFoundRecipeManifestError
		s.ErrorAs(err, &_notFoundRecipeManifestError)
		s.Nil(recMan)

		report := internalReport.NewErrorReport(err)

		reportAssert := &internalReport.Assert{
			Err: "recipe manifest not found",
			Fields: map[string]interface{}{
				"file": filepath.Join(recDir, ".manala.yaml"),
			},
		}
		reportAssert.Equal(&s.Suite, report)
	})

	s.Run("Directory", func() {
		repoUrl := internalTesting.DataPath(s, "repository")

		recDir := filepath.Join(repoUrl, "recipe")
		recManFile := filepath.Join(recDir, ".manala.yaml")

		manager := NewManager(log)

		recMan, err := manager.loadManifest(recManFile)

		s.Error(err)
		s.Nil(recMan)

		report := internalReport.NewErrorReport(err)

		reportAssert := &internalReport.Assert{
			Err: "recipe manifest is a directory",
			Fields: map[string]interface{}{
				"dir": recManFile,
			},
		}
		reportAssert.Equal(&s.Suite, report)
	})
}

func (s *ManagerSuite) TestLoadManifest() {
	log := internalLog.New(io.Discard)

	repoUrl := internalTesting.DataPath(s, "repository")

	recDir := filepath.Join(repoUrl, "recipe")
	recManFile := filepath.Join(recDir, ".manala.yaml")

	manager := NewManager(log)

	recMan, err := manager.loadManifest(recManFile)

	s.NoError(err)
	s.Equal("Recipe", recMan.Description())
	s.Equal(map[string]interface{}{"foo": "bar"}, recMan.Vars())
}

func (s *ManagerSuite) TestLoadRecipe() {
	log := internalLog.New(io.Discard)

	repoUrl := internalTesting.DataPath(s, "repository")
	recDir := filepath.Join(repoUrl, "recipe")

	repoMock := mocks.MockRepository()
	repoMock.
		On("Url").Return(repoUrl).
		On("Dir").Return(repoUrl)

	manager := NewManager(log)

	rec, err := manager.LoadRecipe(repoMock, "recipe")

	s.NoError(err)
	s.Equal(recDir, rec.Dir())
	s.Equal(repoUrl, rec.Repository().Url())
}

func (s *ManagerSuite) TestWalkRecipes() {
	log := internalLog.New(io.Discard)

	manager := NewManager(log)

	s.Run("Default", func() {
		repoUrl := internalTesting.DataPath(s, "repository")

		repoMock := mocks.MockRepository()
		repoMock.
			On("Url").Return(repoUrl).
			On("Dir").Return(repoUrl)

		count := 1

		err := manager.WalkRecipes(repoMock, func(rec interfaces.Recipe) error {
			switch count {
			case 1:
				s.Equal(internalTesting.DataPath(s, "repository", "bar"), rec.Dir())
				s.Equal("bar", rec.Name())
				s.Equal("Bar", rec.Description())
				s.Equal(map[string]interface{}{"bar": "bar"}, rec.Vars())
			case 2:
				s.Equal(internalTesting.DataPath(s, "repository", "foo"), rec.Dir())
				s.Equal("foo", rec.Name())
				s.Equal("Foo", rec.Description())
				s.Equal(map[string]interface{}{"foo": "foo"}, rec.Vars())
			}

			count++

			return nil
		})

		s.NoError(err)
	})

	s.Run("Empty", func() {
		repoUrl := internalTesting.DataPath(s, "repository")

		repoMock := mocks.MockRepository()
		repoMock.
			On("Url").Return(repoUrl).
			On("Dir").Return(repoUrl)

		err := manager.WalkRecipes(repoMock, func(rec interfaces.Recipe) error {
			return nil
		})

		s.EqualError(err, "empty repository")
	})

}
