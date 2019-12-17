package repository

import (
	"github.com/stretchr/testify/suite"
	"manala/pkg/recipe"
	"os"
	"testing"
)

/***************/
/* New - Suite */
/***************/

type NewTestSuite struct{ suite.Suite }

func TestNewTestSuite(t *testing.T) {
	// Run
	suite.Run(t, new(NewTestSuite))
}

/***************/
/* New - Tests */
/***************/

func (s *NewTestSuite) TestNew() {
	repo := New("foo")
	s.Implements((*Interface)(nil), repo)
	s.Equal("foo", repo.GetSrc())
}

/****************/
/* Load - Suite */
/****************/

type LoadTestSuite struct{ suite.Suite }

func TestLoadTestSuite(t *testing.T) {
	// Run
	suite.Run(t, new(LoadTestSuite))
}

func (s *LoadTestSuite) SetupTest() {
	dir := "testdata/load/cache"
	_ = os.RemoveAll(dir)
	_ = os.Mkdir(dir, 0755)
}

/****************/
/* Load - Tests */
/****************/

func (s *LoadTestSuite) TestLoadDir() {
	repo := New("foo")
	err := repo.Load("bar")
	s.NoError(err)
	s.Equal("foo", repo.GetDir())
}

func (s *LoadTestSuite) TestLoadGit() {
	repo := New("https://github.com/octocat/Hello-World.git")
	err := repo.Load("testdata/load/cache")
	s.NoError(err)
	s.Equal("https://github.com/octocat/Hello-World.git", repo.GetSrc())
	s.Equal("testdata/load/cache/repositories/1d60d0a17c4d14e9bda84ee53ee51311", repo.GetDir())

	s.DirExists("testdata/load/cache/repositories")
	stat, _ := os.Stat("testdata/load/cache/repositories")
	s.Equal(os.FileMode(0700), stat.Mode().Perm())

	s.DirExists("testdata/load/cache/repositories/1d60d0a17c4d14e9bda84ee53ee51311")
	stat, _ = os.Stat("testdata/load/cache/repositories/1d60d0a17c4d14e9bda84ee53ee51311")
	s.Equal(os.FileMode(0700), stat.Mode().Perm())
}

/************************/
/* Walk Recipes - Suite */
/************************/

type WalkRecipesTestSuite struct{ suite.Suite }

func TestWalkRecipesTestSuite(t *testing.T) {
	// Run
	suite.Run(t, new(WalkRecipesTestSuite))
}

/************************/
/* Walk Recipes - Tests */
/************************/

func (s *WalkRecipesTestSuite) TestWalkRecipes() {
	repo := New("testdata/walk_recipes")
	_ = repo.Load("")

	results := make(map[string]string)

	err := repo.WalkRecipes(func(rec recipe.Interface) {
		results[rec.GetName()] = rec.GetConfig().Description
	})

	s.NoError(err)
	s.Len(results, 3)
	s.Equal("Foo bar", results["foo"])
	s.Equal("Bar bar", results["bar"])
	s.Equal("Baz bar", results["baz"])
}

/***********************/
/* Load Recipe - Suite */
/***********************/

type LoadRecipeTestSuite struct{ suite.Suite }

func TestLoadRecipeTestSuite(t *testing.T) {
	// Run
	suite.Run(t, new(LoadRecipeTestSuite))
}

/***********************/
/* Load Recipe - Tests */
/***********************/

func (s *LoadRecipeTestSuite) TestLoadRecipe() {
	repo := New("testdata/load_recipe")
	_ = repo.Load("")

	rec, err := repo.LoadRecipe("foo")

	s.NoError(err)
	s.Implements((*recipe.Interface)(nil), rec)
}

func (s *LoadRecipeTestSuite) TestLoadRecipeNotFound() {
	repo := New("testdata/load_recipe")
	_ = repo.Load("")

	rec, err := repo.LoadRecipe("bar")

	s.Nil(rec)
	s.Error(err, "recipe not found")
}
