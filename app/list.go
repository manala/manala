package app

import (
	"fmt"
	"io"
	"manala/models"
)

func (app *App) List(
	repository string,
	out io.Writer,
) error {
	// Load repository
	repo, err := app.RepositoryLoader.Load(repository)
	if err != nil {
		return err
	}

	// Walk into recipes
	if err := app.RecipeLoader.Walk(repo, func(rec models.RecipeInterface) {
		_, _ = fmt.Fprintf(out, "%s: %s\n", rec.Name(), rec.Description())
	}); err != nil {
		return err
	}

	return nil
}
