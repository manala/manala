package repository

import (
	"github.com/stretchr/testify/suite"
	"io"
	"manala/core"
	internalLog "manala/internal/log"
	internalTesting "manala/internal/testing"
	"testing"
)

type RepositorySuite struct{ suite.Suite }

func TestRepositorySuite(t *testing.T) {
	suite.Run(t, new(RepositorySuite))
}

func (s *RepositorySuite) Test() {
	log := internalLog.New(io.Discard)

	repo := NewRepository(
		log,
		"path",
		"dir",
	)

	s.Equal("path", repo.Path())
	s.Equal("path", repo.Source())
	s.Equal("dir", repo.Dir())

	s.Run("LoadRecipe", func() {
		path := internalTesting.DataPath(s, "repository")

		repo := NewRepository(
			log,
			"path",
			path,
		)

		rec, err := repo.LoadRecipe("recipe")

		s.NoError(err)
		s.Equal(internalTesting.DataPath(s, "repository", "recipe"), rec.Path())
		s.Equal("recipe", rec.Name())
		s.Equal("Recipe", rec.Description())
		s.Equal(map[string]interface{}{"foo": "bar"}, rec.Vars())
	})

	s.Run("WalkRecipes Empty", func() {
		path := internalTesting.DataPath(s, "repository")

		repo := NewRepository(
			log,
			"path",
			path,
		)

		err := repo.WalkRecipes(func(rec core.Recipe) {})

		s.EqualError(err, "empty repository")
	})

	s.Run("WalkRecipes", func() {
		path := internalTesting.DataPath(s, "repository")

		repo := NewRepository(
			log,
			"path",
			path,
		)

		count := 1

		err := repo.WalkRecipes(func(rec core.Recipe) {
			switch count {
			case 1:
				s.Equal(internalTesting.DataPath(s, "repository", "bar"), rec.Path())
				s.Equal("bar", rec.Name())
				s.Equal("Bar", rec.Description())
				s.Equal(map[string]interface{}{"bar": "bar"}, rec.Vars())
			case 2:
				s.Equal(internalTesting.DataPath(s, "repository", "foo"), rec.Path())
				s.Equal("foo", rec.Name())
				s.Equal("Foo", rec.Description())
				s.Equal(map[string]interface{}{"foo": "foo"}, rec.Vars())
			}

			count++
		})

		s.NoError(err)
	})
}
