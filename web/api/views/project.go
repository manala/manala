package views

import (
	"manala/app/interfaces"
)

func NormalizeProject(proj interfaces.Project) *ProjectView {
	return &ProjectView{
		Dir:    proj.Dir(),
		Recipe: NormalizeRecipe(proj.Recipe()),
		Vars:   proj.Vars(),
	}
}

type ProjectView struct {
	Dir    string                 `json:"dir"`
	Recipe *RecipeView            `json:"recipe"`
	Vars   map[string]interface{} `json:"vars"`
}
