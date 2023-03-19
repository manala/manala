package binder

import (
	"code.rocketnine.space/tslocum/cview"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"manala/app/interfaces"
	"manala/app/mocks"
	"testing"
)

type RecipeFormBinderSuite struct {
	suite.Suite
}

func TestRecipeFormBinderSuite(t *testing.T) {
	suite.Run(t, new(RecipeFormBinderSuite))
}

func (s *RecipeFormBinderSuite) TestNew() {

	s.Run("String", func() {
		optionMock := mocks.MockRecipeOption()
		optionMock.
			On("Label").Return("Option").
			On("Schema").Return(map[string]interface{}{
			"type": "string",
		})

		options := []interfaces.RecipeOption{
			optionMock,
		}

		binder, err := NewRecipeFormBinder(options)

		s.NoError(err)
		s.Len(binder.Binds(), 1)

		bind := binder.Binds()[0]
		s.IsType((*cview.InputField)(nil), bind.Item)
		s.Equal(options[0].Label(), bind.Item.GetLabel())

		item := bind.Item.(*cview.InputField)

		s.Equal("", item.GetText())
		s.Equal("", bind.Value)

		item.SetText("foo")
		s.Equal("foo", bind.Value)
	})

	s.Run("Enum", func() {
		optionMock := mocks.MockRecipeOption()
		optionMock.
			On("Label").Return("Option").
			On("Schema").Return(map[string]interface{}{
			"enum": []interface{}{true, false, nil, "foo", 123, "7.0", 7.1},
		})

		options := []interfaces.RecipeOption{
			optionMock,
		}

		binder, err := NewRecipeFormBinder(options)

		s.NoError(err)
		s.Len(binder.Binds(), 1)

		bind := binder.Binds()[0]
		s.IsType((*cview.DropDown)(nil), bind.Item)
		s.Equal(options[0].Label(), bind.Item.GetLabel())

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
	})

}

func (s *RecipeFormBinderSuite) TestApply() {
	var value interface{}

	optionMock := mocks.MockRecipeOption()
	optionMock.
		On("Label").Return("Option").
		On("Schema").Return(map[string]interface{}{
		"type": "string",
	}).
		On("Set", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		value = args.Get(0)
	})

	binder, _ := NewRecipeFormBinder([]interfaces.RecipeOption{
		optionMock,
	})

	binder.Binds()[0].Value = "bar"
	err := binder.Apply()

	s.NoError(err)
	s.Equal("bar", value)
}
