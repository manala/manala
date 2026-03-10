package app

/***********/
/* Project */
/***********/

// ProjectView is a secure and lightweight facade of a project, dedicated to template usage.
type ProjectView struct {
	Dir    string         `json:"dir"`
	Vars   map[string]any `json:"vars"`
	Recipe *RecipeView    `json:"recipe"`
}

// NewProjectView create a project view.
func NewProjectView(project Project) *ProjectView {
	return &ProjectView{
		Dir:    project.Dir(),
		Vars:   project.Vars(),
		Recipe: NewRecipeView(project.Recipe()),
	}
}

/**********/
/* Recipe */
/**********/

// RecipeView is a secure and lightweight facade of a recipe, dedicated to template usage.
type RecipeView struct {
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Icon        string          `json:"icon,omitempty"`
	Repository  *RepositoryView `json:"-"`
}

// NewRecipeView create a recipe view.
func NewRecipeView(recipe Recipe) *RecipeView {
	return &RecipeView{
		Name:        recipe.Name(),
		Description: recipe.Description(),
		Icon:        recipe.Icon(),
		Repository:  NewRepositoryView(recipe.Repository()),
	}
}

type RecipesView []*RecipeView

// NewRecipesView create a slice of recipes view.
func NewRecipesView(recipes []Recipe) RecipesView {
	views := make(RecipesView, len(recipes))
	for i := range recipes {
		views[i] = NewRecipeView(recipes[i])
	}

	return views
}

/**************/
/* Repository */
/**************/

// RepositoryView is a secure and lightweight facade of a repository, dedicated to template usage.
type RepositoryView struct {
	URL string
	// Path ensure backward compatibility, when "path" was used instead of "url" to define repository origin
	Path string
	// Source ensure backward compatibility, when "source" was used instead of "path" to define repository origin
	Source string
}

// NewRepositoryView create a repository view.
func NewRepositoryView(repository Repository) *RepositoryView {
	url := repository.URL()

	return &RepositoryView{
		URL:    url,
		Path:   url,
		Source: url,
	}
}
