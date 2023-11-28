package recipe

import (
	"github.com/stretchr/testify/suite"
	"io"
	"log/slog"
	"manala/app"
	"manala/internal/serrors"
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
		repositoryUrl := filepath.FromSlash("testdata/DirManagerSuite/TestLoadManifestErrors/NotFound/repository")

		dir := filepath.Join(repositoryUrl, "recipe")
		manifestFile := filepath.Join(dir, ".manala.yaml")

		manager := NewDirManager(log)

		manifest, err := manager.loadManifest(manifestFile)

		s.Nil(manifest)

		serrors.Equal(s.Assert(), &serrors.Assert{
			Type:    &app.NotFoundRecipeManifestError{},
			Message: "recipe manifest not found",
			Arguments: []any{
				"file", filepath.Join(dir, ".manala.yaml"),
			},
		}, err)
	})

	s.Run("Directory", func() {
		repositoryUrl := filepath.FromSlash("testdata/DirManagerSuite/TestLoadManifestErrors/Directory/repository")

		dir := filepath.Join(repositoryUrl, "recipe")
		manifestFile := filepath.Join(dir, ".manala.yaml")

		manager := NewDirManager(log)

		manifest, err := manager.loadManifest(manifestFile)

		s.Nil(manifest)

		serrors.Equal(s.Assert(), &serrors.Assert{
			Type:    serrors.Error{},
			Message: "recipe manifest is a directory",
			Arguments: []any{
				"dir", manifestFile,
			},
		}, err)
	})
}

func (s *DirManagerSuite) TestLoadManifest() {
	repositoryUrl := filepath.FromSlash("testdata/DirManagerSuite/TestLoadManifest/repository")

	dir := filepath.Join(repositoryUrl, "recipe")
	manifestFile := filepath.Join(dir, ".manala.yaml")

	manager := NewDirManager(
		slog.New(slog.NewTextHandler(io.Discard, nil)),
	)

	manifest, err := manager.loadManifest(manifestFile)

	s.NoError(err)

	s.Equal("Recipe", manifest.Description())
	s.Equal("Icon", manifest.Icon())
	s.Equal(map[string]any{"foo": "bar"}, manifest.Vars())
}

func (s *DirManagerSuite) TestLoadRecipe() {
	repositoryUrl := filepath.FromSlash("testdata/DirManagerSuite/TestLoadRecipe/repository")

	repositoryMock := &app.RepositoryMock{}
	repositoryMock.
		On("Url").Return(repositoryUrl).
		On("Dir").Return(repositoryUrl)

	manager := NewDirManager(
		slog.New(slog.NewTextHandler(io.Discard, nil)),
	)

	recipe, err := manager.LoadRecipe(repositoryMock, "recipe")

	s.NoError(err)

	s.Equal(filepath.Join(repositoryUrl, "recipe"), recipe.Dir())
	s.Equal(repositoryUrl, recipe.Repository().Url())
}

func (s *DirManagerSuite) TestRepositoryRecipes() {
	manager := NewDirManager(
		slog.New(slog.NewTextHandler(io.Discard, nil)),
	)

	s.Run("Default", func() {
		repositoryUrl := filepath.FromSlash("testdata/DirManagerSuite/TestRepositoryRecipes/Default/repository")

		repositoryMock := &app.RepositoryMock{}
		repositoryMock.
			On("Url").Return(repositoryUrl).
			On("Dir").Return(repositoryUrl)

		recipes, err := manager.RepositoryRecipes(repositoryMock)

		s.NoError(err)

		s.Equal(filepath.Join(repositoryUrl, "bar"), recipes[0].Dir())
		s.Equal("bar", recipes[0].Name())
		s.Equal("Bar", recipes[0].Description())
		s.Equal("Bar Icon", recipes[0].Icon())
		s.Equal(map[string]any{"bar": "bar"}, recipes[0].Vars())

		s.Equal(filepath.Join(repositoryUrl, "foo"), recipes[1].Dir())
		s.Equal("foo", recipes[1].Name())
		s.Equal("Foo", recipes[1].Description())
		s.Equal("Foo Icon", recipes[1].Icon())
		s.Equal(map[string]any{"foo": "foo"}, recipes[1].Vars())
	})

	s.Run("Empty", func() {
		repositoryUrl := filepath.FromSlash("testdata/DirManagerSuite/TestRepositoryRecipes/Empty/repository")

		repositoryMock := &app.RepositoryMock{}
		repositoryMock.
			On("Url").Return(repositoryUrl).
			On("Dir").Return(repositoryUrl)

		recipes, err := manager.RepositoryRecipes(repositoryMock)

		s.Nil(recipes)

		serrors.Equal(s.Assert(), &serrors.Assert{
			Type:    &app.EmptyRepositoryError{},
			Message: "empty repository",
			Arguments: []any{
				"url", repositoryUrl,
			},
		}, err)
	})

}
