package template

import (
	"github.com/manala/manala/app"
	"github.com/manala/manala/internal/template/engine"
)

type Engine struct {
	engine *engine.Engine
}

func NewEngine() *Engine {
	return &Engine{
		engine: engine.New(),
	}
}

func (e *Engine) Executor(vars map[string]any, recipe app.Recipe, dir string) (*engine.Executor, error) {
	return e.engine.Executor(
		map[string]any{
			"Vars":   vars,
			"Recipe": NewRecipeView(recipe),
			// Legacy: switch to RepositoryView when it has been cleaned of its own legacy fields
			"Repository": map[string]any{
				"URL": recipe.Repository().URL(),
			},
			// Legacy: remove
			"Dir": dir,
		},
		recipe.Partials()...,
	)
}
