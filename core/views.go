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
	url := repo.Url()
	return &RepositoryView{
		Url:    url,
		Path:   url,
		Source: url,
	}
}

type RepositoryView struct {
	Url string
	// Path ensure backward compatibility, when "path" was used instead of "url" to define repository origin
	Path string
	// Source ensure backward compatibility, when "source" was used instead of "path" to define repository origin
	Source string
}
