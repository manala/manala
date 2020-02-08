package binder

import (
	"github.com/stretchr/testify/suite"
	"gitlab.com/tslocum/cview"
	"manala/models"
	"testing"
)

/******************************/
/* Recipe form Binder - Suite */
/******************************/

type RecipeFormBinderTestSuite struct {
	suite.Suite
	recipe models.RecipeInterface
}

func TestRecipeOptionFormItemBindTestSuite(t *testing.T) {
	// Run
	suite.Run(t, new(RecipeFormBinderTestSuite))
}

func (s *RecipeFormBinderTestSuite) SetupTest() {
	s.recipe = models.NewRecipe(
		"foo",
		"bar",
		"baz",
		models.NewRepository(
			"foo",
			"bar",
		),
	)
}

/*****************************/
/* Recipe form Binder- Tests */
/*****************************/

func (s *RecipeFormBinderTestSuite) TestNewEnum() {
	s.recipe.AddOptions([]models.RecipeOption{
		{
			Label:  "Foo bar",
			Path:   "/foo",
			Schema: map[string]interface{}{"enum": []interface{}{true, false, nil, "foo", 123, "7.0", 7.1}},
		},
	})

	bndr, err := NewRecipeFormBinder(s.recipe)
	s.NoError(err)
	s.Len(bndr.Binds(), 1)

	bind := bndr.Binds()[0]
	s.IsType((*cview.DropDown)(nil), bind.Item)
	s.Equal(s.recipe.Options()[0].Label, bind.Item.GetLabel())

	item := bind.Item.(*cview.DropDown)

	itemIndex, itemValue := item.GetCurrentOption()
	s.Equal(0, itemIndex)
	s.Equal("<True>", itemValue)
	s.Equal(true, bind.Value)

	item.SetCurrentOption(1)
	itemIndex, itemValue = item.GetCurrentOption()
	s.Equal("<False>", itemValue)
	s.Equal(false, bind.Value)

	item.SetCurrentOption(2)
	itemIndex, itemValue = item.GetCurrentOption()
	s.Equal("<None>", itemValue)
	s.Equal(nil, bind.Value)

	item.SetCurrentOption(3)
	itemIndex, itemValue = item.GetCurrentOption()
	s.Equal("foo", itemValue)
	s.Equal("foo", bind.Value)

	item.SetCurrentOption(4)
	itemIndex, itemValue = item.GetCurrentOption()
	s.Equal("123", itemValue)
	s.Equal(123, bind.Value)

	item.SetCurrentOption(5)
	itemIndex, itemValue = item.GetCurrentOption()
	s.Equal("7.0", itemValue)
	s.Equal("7.0", bind.Value)

	item.SetCurrentOption(6)
	itemIndex, itemValue = item.GetCurrentOption()
	s.Equal("7.1", itemValue)
	s.Equal(7.1, bind.Value)
}

func (s *RecipeFormBinderTestSuite) TestNewTypeString() {
	s.recipe.AddOptions([]models.RecipeOption{
		{
			Label:  "Foo bar",
			Path:   "/foo",
			Schema: map[string]interface{}{"type": "string"},
		},
	})

	bndr, err := NewRecipeFormBinder(s.recipe)
	s.NoError(err)
	s.Len(bndr.Binds(), 1)

	bind := bndr.Binds()[0]
	s.IsType((*cview.InputField)(nil), bind.Item)
	s.Equal(s.recipe.Options()[0].Label, bind.Item.GetLabel())

	item := bind.Item.(*cview.InputField)

	s.Equal("", item.GetText())
	s.Equal("", bind.Value)

	item.SetText("foo")
	s.Equal("foo", bind.Value)
}
