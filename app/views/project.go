package views

import (
	"manala/app/interfaces"
)

func NormalizeProject(proj interfaces.Project) *ProjectView {
	return &ProjectView{
		Vars:   proj.Vars(),
		Recipe: NormalizeRecipe(proj.Recipe()),
	}
}

type ProjectView struct {
	Vars   map[string]interface{}
	Recipe *RecipeView
}
