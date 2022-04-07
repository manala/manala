package binder

import (
	"code.rocketnine.space/tslocum/cview"
	"fmt"
	"github.com/ohler55/ojg/jp"
	"manala/internal"
)

func NewRecipeFormBinder(options []internal.RecipeManifestOption) (*RecipeFormBinder, error) {
	binder := &RecipeFormBinder{}

	for _, option := range options {

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
				itemOption := cview.NewDropDownOption(text)
				itemOption.SetSelectedFunc(func(index int, option *cview.DropDownOption) {
					bind.Value = itemValue
				})
				item.AddOptions(itemOption)
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

		binder.binds = append(binder.binds, bind)
	}

	return binder, nil
}

type RecipeFormBinder struct {
	binds []*recipeFormBind
}

func (binder *RecipeFormBinder) Binds() []*recipeFormBind {
	return binder.binds
}

func (binder *RecipeFormBinder) BindForm(form *cview.Form) {
	for i, bind := range binder.binds {
		form.AddFormItem(bind.Item)
		bind.ItemIndex = i
	}
}

func (binder *RecipeFormBinder) Apply(manifest *internal.ProjectManifest) error {
	for _, bind := range binder.binds {
		jsonPath, err := jp.ParseString(bind.Option.Path)
		if err != nil {
			return err
		}
		if err := jsonPath.SetOne(manifest.Vars, bind.Value); err != nil {
			return err
		}
	}

	return nil
}

type recipeFormBind struct {
	Option    internal.RecipeManifestOption
	Item      cview.FormItem
	ItemIndex int
	Value     interface{}
}
