package views

import (
	"manala/app/interfaces"
)

func NormalizeRecipe(rec interfaces.Recipe) *RecipeView {
	return &RecipeView{
		Name:       rec.Name(),
		Repository: NormalizeRepository(rec.Repository()),
	}
}

type RecipeView struct {
	Name       string
	Repository *RepositoryView
}
