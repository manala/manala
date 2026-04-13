package template

import "github.com/manala/manala/app"

/**********/
/* Recipe */
/**********/

// RecipeView is a secure and lightweight facade of a Recipe, dedicated to template usage.
type RecipeView struct {
	Name        string
	Description string
	Icon        string
	// Legacy: remove
	Repository *RepositoryView
}

// NewRecipeView create a RecipeView from a Recipe.
func NewRecipeView(recipe app.Recipe) *RecipeView {
	return &RecipeView{
		Name:        recipe.Name(),
		Description: recipe.Description(),
		Icon:        recipe.Icon(),
		Repository:  NewRepositoryView(recipe.Repository()),
	}
}

/**************/
/* Repository */
/**************/

// RepositoryView is a secure and lightweight facade of a Repository, dedicated to template usage.
type RepositoryView struct {
	URL string
	// Legacy: remove
	Path string
	// Legacy: remove
	Source string
}

// NewRepositoryView create a RepositoryView from a Repository.
func NewRepositoryView(repository app.Repository) *RepositoryView {
	url := repository.URL()

	return &RepositoryView{
		URL:    url,
		Path:   url,
		Source: url,
	}
}
