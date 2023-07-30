package recipe

import (
	"github.com/stretchr/testify/suite"
	"io"
	"log/slog"
	"manala/app/interfaces"
	"manala/app/mocks"
	"manala/core"
	"manala/internal/errors/serrors"
	"path/filepath"
	"testing"
)

type DirManagerSuite struct{ suite.Suite }

func TestDirManagerSuite(t *testing.T) {
	suite.Run(t, new(DirManagerSuite))
}

func (s *DirManagerSuite) TestLoadManifestErrors() {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))

	s.Run("NotFound", func() {
		repoUrl := filepath.FromSlash("testdata/DirManagerSuite/TestLoadManifestErrors/NotFound/repository")

		recDir := filepath.Join(repoUrl, "recipe")
		recManFile := filepath.Join(recDir, ".manala.yaml")

		manager := NewDirManager(log)

		recMan, err := manager.loadManifest(recManFile)

		s.Nil(recMan)

		serrors.Equal(s.Assert(), &serrors.Assert{
			Type:    &core.NotFoundRecipeManifestError{},
			Message: "recipe manifest not found",
			Arguments: []any{
				"file", filepath.Join(recDir, ".manala.yaml"),
			},
		}, err)
	})

	s.Run("Directory", func() {
		repoUrl := filepath.FromSlash("testdata/DirManagerSuite/TestLoadManifestErrors/Directory/repository")

		recDir := filepath.Join(repoUrl, "recipe")
		recManFile := filepath.Join(recDir, ".manala.yaml")

		manager := NewDirManager(log)

		recMan, err := manager.loadManifest(recManFile)

		s.Nil(recMan)

		serrors.Equal(s.Assert(), &serrors.Assert{
			Type:    &serrors.Error{},
			Message: "recipe manifest is a directory",
			Arguments: []any{
				"dir", recManFile,
			},
		}, err)
	})
}

func (s *DirManagerSuite) TestLoadManifest() {
	repoUrl := filepath.FromSlash("testdata/DirManagerSuite/TestLoadManifest/repository")

	recDir := filepath.Join(repoUrl, "recipe")
	recManFile := filepath.Join(recDir, ".manala.yaml")

	manager := NewDirManager(
		slog.New(slog.NewTextHandler(io.Discard, nil)),
	)

	recMan, err := manager.loadManifest(recManFile)

	s.NoError(err)

	s.Equal("Recipe", recMan.Description())
	s.Equal(map[string]interface{}{"foo": "bar"}, recMan.Vars())
}

func (s *DirManagerSuite) TestLoadRecipe() {
	repoUrl := filepath.FromSlash("testdata/DirManagerSuite/TestLoadRecipe/repository")

	repoMock := &mocks.RepositoryMock{}
	repoMock.
		On("Url").Return(repoUrl).
		On("Dir").Return(repoUrl)

	manager := NewDirManager(
		slog.New(slog.NewTextHandler(io.Discard, nil)),
	)

	rec, err := manager.LoadRecipe(repoMock, "recipe")

	s.NoError(err)

	s.Equal(filepath.Join(repoUrl, "recipe"), rec.Dir())
	s.Equal(repoUrl, rec.Repository().Url())
}

func (s *DirManagerSuite) TestWalkRecipes() {
	manager := NewDirManager(
		slog.New(slog.NewTextHandler(io.Discard, nil)),
	)

	s.Run("Default", func() {
		repoUrl := filepath.FromSlash("testdata/DirManagerSuite/TestWalkRecipes/Default/repository")

		repoMock := &mocks.RepositoryMock{}
		repoMock.
			On("Url").Return(repoUrl).
			On("Dir").Return(repoUrl)

		count := 1

		err := manager.WalkRecipes(repoMock, func(rec interfaces.Recipe) error {
			switch count {
			case 1:
				s.Equal(filepath.Join(repoUrl, "bar"), rec.Dir())
				s.Equal("bar", rec.Name())
				s.Equal("Bar", rec.Description())
				s.Equal(map[string]interface{}{"bar": "bar"}, rec.Vars())
			case 2:
				s.Equal(filepath.Join(repoUrl, "foo"), rec.Dir())
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
		repoUrl := filepath.FromSlash("testdata/DirManagerSuite/TestWalkRecipes/Empty/repository")

		repoMock := &mocks.RepositoryMock{}
		repoMock.
			On("Url").Return(repoUrl).
			On("Dir").Return(repoUrl)

		err := manager.WalkRecipes(repoMock, func(rec interfaces.Recipe) error {
			return nil
		})

		serrors.Equal(s.Assert(), &serrors.Assert{
			Type:    &serrors.Error{},
			Message: "empty repository",
			Arguments: []any{
				"dir", repoUrl,
			},
		}, err)
	})

}
