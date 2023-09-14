package app

/***********/
/* Project */
/***********/

// NewProjectView create a project view
func NewProjectView(project Project) *ProjectView {
	return &ProjectView{
		Vars:   project.Vars(),
		Recipe: NewRecipeView(project.Recipe()),
	}
}

// ProjectView is a secure and lightweight facade of a project, dedicated to template usage
type ProjectView struct {
	Vars   map[string]any
	Recipe *RecipeView
}

/**********/
/* Recipe */
/**********/

// NewRecipeView create a recipe view
func NewRecipeView(recipe Recipe) *RecipeView {
	return &RecipeView{
		Name:       recipe.Name(),
		Repository: NewRepositoryView(recipe.Repository()),
	}
}

// RecipeView is a secure and lightweight facade of a recipe, dedicated to template usage
type RecipeView struct {
	Name       string
	Repository *RepositoryView
}

/**************/
/* Repository */
/**************/

// NewRepositoryView create a repository view
func NewRepositoryView(repository Repository) *RepositoryView {
	url := repository.Url()

	return &RepositoryView{
		Url:    url,
		Path:   url,
		Source: url,
	}
}

// RepositoryView is a secure and lightweight facade of a repository, dedicated to template usage
type RepositoryView struct {
	Url string
	// Path ensure backward compatibility, when "path" was used instead of "url" to define repository origin
	Path string
	// Source ensure backward compatibility, when "source" was used instead of "path" to define repository origin
	Source string
}
