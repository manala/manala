package binder

import (
	"fmt"
	"github.com/xeipuuv/gojsonpointer"
	"gitlab.com/tslocum/cview"
	"manala/models"
)

func NewRecipeFormBinder(rec models.RecipeInterface) (*RecipeFormBinder, error) {
	bndr := &RecipeFormBinder{
		recipe: &rec,
	}

	for _, option := range rec.Options() {

		// Bind
		bind := &recipeFormBind{
			Option: option,
		}

		if _, ok := option.Schema["enum"]; ok {
			// Dropdown item based on enum schema
			item := cview.NewDropDown()

			// Item label
			item.SetLabel(option.Label)

			// Item options
			values, ok := option.Schema["enum"].([]interface{})
			if !ok {
				return nil, fmt.Errorf("invalid recipe option enum type: " + option.Label)
			}
			if len(values) == 0 {
				return nil, fmt.Errorf("empty recipe option enum: " + option.Label)
			}
			for _, value := range values {
				var text string
				switch value {
				case nil:
					text = "<None>"
				case true:
					text = "<True>"
				case false:
					text = "<False>"
				default:
					text = fmt.Sprintf("%v", value)
				}
				itemValue := value
				item.AddOption(text, func() {
					bind.Value = itemValue
				})
			}

			// Item current option
			item.SetCurrentOption(0)

			// Bind
			bind.Item = item
		} else if t, ok := option.Schema["type"]; ok && t == "string" {
			// Input field item based on string type
			item := cview.NewInputField()

			// Item label
			item.SetLabel(option.Label)

			item.SetChangedFunc(func(text string) {
				bind.Value = text
			})

			// Item text
			item.SetText("")

			// Bind
			bind.Item = item
		} else {
			return nil, fmt.Errorf("unable to bind recipe option into a form item: " + option.Label)
		}

		bndr.binds = append(bndr.binds, bind)
	}

	return bndr, nil
}

type RecipeFormBinder struct {
	recipe *models.RecipeInterface
	binds  []*recipeFormBind
}

func (bndr *RecipeFormBinder) Binds() []*recipeFormBind {
	return bndr.binds
}

func (bndr *RecipeFormBinder) BindForm(form *cview.Form) {
	for i, bind := range bndr.binds {
		form.AddFormItem(bind.Item)
		bind.ItemIndex = i
	}
}

func (bndr *RecipeFormBinder) ApplyValues(values map[string]interface{}) error {
	for _, bind := range bndr.binds {
		// Json pointer
		pointer, err := gojsonpointer.NewJsonPointer(bind.Option.Path)
		if err != nil {
			return err
		}
		_, err = pointer.Set(values, bind.Value)
		if err != nil {
			return err
		}
	}

	return nil
}

type recipeFormBind struct {
	Option    models.RecipeOption
	Item      cview.FormItem
	ItemIndex int
	Value     interface{}
}
