package binder

import (
	"code.rocketnine.space/tslocum/cview"
	"github.com/stretchr/testify/suite"
	"manala/models"
	"testing"
)

/******************************/
/* Recipe form Binder - Suite */
/******************************/

type RecipeFormBinderTestSuite struct {
	suite.Suite
	repository models.RepositoryInterface
}

func TestRecipeOptionFormItemBindTestSuite(t *testing.T) {
	// Run
	suite.Run(t, new(RecipeFormBinderTestSuite))
}

func (s *RecipeFormBinderTestSuite) SetupTest() {
	s.repository = models.NewRepository(
		"foo",
		"bar",
		false,
	)
}

/*****************************/
/* Recipe form Binder- Tests */
/*****************************/

func (s *RecipeFormBinderTestSuite) TestNewEnum() {
	recipe := models.NewRecipe(
		"foo",
		"bar",
		"",
		"baz",
		s.repository,
		nil,
		nil,
		nil,
		[]models.RecipeOption{
			{
				Label:  "Foo bar",
				Path:   "/foo",
				Schema: map[string]interface{}{"enum": []interface{}{true, false, nil, "foo", 123, "7.0", 7.1}},
			},
		},
	)

	bndr, err := NewRecipeFormBinder(recipe)
	s.NoError(err)
	s.Len(bndr.Binds(), 1)

	bind := bndr.Binds()[0]
	s.IsType((*cview.DropDown)(nil), bind.Item)
	s.Equal(recipe.Options()[0].Label, bind.Item.GetLabel())

	item := bind.Item.(*cview.DropDown)

	itemIndex, itemOption := item.GetCurrentOption()
	s.Equal(0, itemIndex)
	s.Equal("<True>", itemOption.GetText())
	s.Equal(true, bind.Value)

	item.SetCurrentOption(1)
	_, itemOption = item.GetCurrentOption()
	s.Equal("<False>", itemOption.GetText())
	s.Equal(false, bind.Value)

	item.SetCurrentOption(2)
	_, itemOption = item.GetCurrentOption()
	s.Equal("<None>", itemOption.GetText())
	s.Equal(nil, bind.Value)

	item.SetCurrentOption(3)
	_, itemOption = item.GetCurrentOption()
	s.Equal("foo", itemOption.GetText())
	s.Equal("foo", bind.Value)

	item.SetCurrentOption(4)
	_, itemOption = item.GetCurrentOption()
	s.Equal("123", itemOption.GetText())
	s.Equal(123, bind.Value)

	item.SetCurrentOption(5)
	_, itemOption = item.GetCurrentOption()
	s.Equal("7.0", itemOption.GetText())
	s.Equal("7.0", bind.Value)

	item.SetCurrentOption(6)
	_, itemOption = item.GetCurrentOption()
	s.Equal("7.1", itemOption.GetText())
	s.Equal(7.1, bind.Value)
}

func (s *RecipeFormBinderTestSuite) TestNewTypeString() {
	recipe := models.NewRecipe(
		"foo",
		"bar",
		"",
		"baz",
		s.repository,
		nil,
		nil,
		nil,
		[]models.RecipeOption{
			{
				Label:  "Foo bar",
				Path:   "/foo",
				Schema: map[string]interface{}{"type": "string"},
			},
		},
	)

	bndr, err := NewRecipeFormBinder(recipe)
	s.NoError(err)
	s.Len(bndr.Binds(), 1)

	bind := bndr.Binds()[0]
	s.IsType((*cview.InputField)(nil), bind.Item)
	s.Equal(recipe.Options()[0].Label, bind.Item.GetLabel())

	item := bind.Item.(*cview.InputField)

	s.Equal("", item.GetText())
	s.Equal("", bind.Value)

	item.SetText("foo")
	s.Equal("foo", bind.Value)
}
