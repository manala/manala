package app

import (
	"manala/models"
)

func (app *App) List() ([]models.RecipeInterface, error) {
	// Load repository
	repo, err := app.repositoryLoader.Load(
		app.Config.GetString("repository"),
		app.Config.GetString("cache-dir"),
	)
	if err != nil {
		return nil, err
	}

	// Walk into recipes
	var recipes []models.RecipeInterface
	if err := app.recipeLoader.Walk(repo, func(rec models.RecipeInterface) {
		recipes = append(recipes, rec)
	}); err != nil {
		return nil, err
	}

	return recipes, nil
}
