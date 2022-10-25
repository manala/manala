package core

/***********/
/* Project */
/***********/

func NewProjectView(proj Project) *ProjectView {
	return &ProjectView{
		Vars:   proj.Vars(),
		Recipe: NewRecipeView(proj.Recipe()),
	}
}

type ProjectView struct {
	Vars   map[string]interface{}
	Recipe *RecipeView
}

/**********/
/* Recipe */
/**********/

func NewRecipeView(rec Recipe) *RecipeView {
	return &RecipeView{
		Name:       rec.Name(),
		Repository: NewRepositoryView(rec.Repository()),
	}
}

type RecipeView struct {
	Name       string
	Repository *RepositoryView
}

/**************/
/* Repository */
/**************/

func NewRepositoryView(repo Repository) *RepositoryView {
	return &RepositoryView{
		Path:   repo.Path(),
		Source: repo.Source(),
	}
}

type RepositoryView struct {
	Path   string
	Source string
}
