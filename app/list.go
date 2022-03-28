package app

import (
	"fmt"
	"io"
	"manala/models"
)

func (app *App) List(
	out io.Writer,
) error {
	// Load repository
	repo, err := app.repositoryLoader.Load(
		app.Config.GetString("repository"),
		app.Config.GetString("cache-dir"),
	)
	if err != nil {
		return err
	}

	// Walk into recipes
	if err := app.recipeLoader.Walk(repo, func(rec models.RecipeInterface) {
		_, _ = fmt.Fprintf(out, "%s: %s\n", rec.Name(), rec.Description())
	}); err != nil {
		return err
	}

	return nil
}
