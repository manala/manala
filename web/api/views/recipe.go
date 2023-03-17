package views

import (
	"manala/app/interfaces"
)

func NormalizeRecipe(rec interfaces.Recipe) *RecipeView {
	return &RecipeView{
		Name:        rec.Name(),
		Description: rec.Description(),
	}
}

type RecipeView struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

/**********/
/* Option */
/**********/

func NormalizeRecipeOption(option interfaces.RecipeOption) *RecipeOptionView {
	return &RecipeOptionView{
		Label: option.Label(),
	}
}

type RecipeOptionView struct {
	Label string `json:"label"`
}
