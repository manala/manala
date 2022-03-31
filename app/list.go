package app

import (
	"manala/models"
)

func (app *App) List() ([]models.RecipeInterface, error) {
	// Debug
	app.log.Debug("run list command")

	// Load repository
	repo, err := app.repositoryLoader.Load(
		app.config.GetString("repository"),
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
