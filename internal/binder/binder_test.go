package binder

import (
	"code.rocketnine.space/tslocum/cview"
	"github.com/stretchr/testify/suite"
	"manala/internal"
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
		options := []internal.RecipeManifestOption{
			{
				Schema: map[string]interface{}{"type": "string"},
			},
		}

		binder, err := NewRecipeFormBinder(options)

		s.NoError(err)
		s.Len(binder.Binds(), 1)

		bind := binder.Binds()[0]
		s.IsType((*cview.InputField)(nil), bind.Item)
		s.Equal(options[0].Label, bind.Item.GetLabel())

		item := bind.Item.(*cview.InputField)

		s.Equal("", item.GetText())
		s.Equal("", bind.Value)

		item.SetText("foo")
		s.Equal("foo", bind.Value)
	})

	s.Run("Enum", func() {
		options := []internal.RecipeManifestOption{
			{
				Schema: map[string]interface{}{
					"enum": []interface{}{true, false, nil, "foo", 123, "7.0", 7.1},
				},
			},
		}

		binder, err := NewRecipeFormBinder(options)

		s.NoError(err)
		s.Len(binder.Binds(), 1)

		bind := binder.Binds()[0]
		s.IsType((*cview.DropDown)(nil), bind.Item)
		s.Equal(options[0].Label, bind.Item.GetLabel())

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
	tests := []struct {
		name          string
		initialValue  interface{}
		actualValue   interface{}
		expectedValue interface{}
	}{
		{
			name:          "Nil",
			initialValue:  "",
			actualValue:   nil,
			expectedValue: nil,
		},
		{
			name:          "True",
			initialValue:  nil,
			actualValue:   true,
			expectedValue: true,
		},
		{
			name:          "False",
			initialValue:  nil,
			actualValue:   false,
			expectedValue: false,
		},
		{
			name:          "String",
			initialValue:  nil,
			actualValue:   "string",
			expectedValue: "string",
		},
		{
			name:          "String Asterisk",
			initialValue:  nil,
			actualValue:   "*",
			expectedValue: "*",
		},
		{
			name:          "String Int",
			initialValue:  nil,
			actualValue:   "12",
			expectedValue: "12",
		},
		{
			name:          "String Float",
			initialValue:  nil,
			actualValue:   "2.3",
			expectedValue: "2.3",
		},
		{
			name:          "String Float Int",
			initialValue:  nil,
			actualValue:   "3.0",
			expectedValue: "3.0",
		},
		{
			name:          "Integer Uint64",
			initialValue:  nil,
			actualValue:   uint64(12),
			expectedValue: uint64(12),
		},
		{
			name:          "Integer Float64",
			initialValue:  nil,
			actualValue:   float64(12),
			expectedValue: uint64(12),
		},
		{
			name:          "Float",
			initialValue:  nil,
			actualValue:   float64(2.3),
			expectedValue: float64(2.3),
		},
		{
			name:          "Float Int",
			initialValue:  nil,
			actualValue:   float64(3.0),
			expectedValue: uint64(3),
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			binder, _ := NewRecipeFormBinder([]internal.RecipeManifestOption{{
				Path:   "$.value",
				Schema: map[string]interface{}{"type": "string"},
			}})

			manifest := internal.NewProjectManifest("dir")
			manifest.Vars["value"] = test.initialValue

			binder.Binds()[0].Value = test.actualValue
			err := binder.Apply(manifest)

			s.NoError(err)
			s.Equal(test.expectedValue, manifest.Vars["value"])
		})
	}
}
