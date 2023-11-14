package api

import "manala/app"

func (api *Api) LoadPrecedingRepository() (app.Repository, error) {
	// Log
	api.log.Debug("load preceding repository…")

	// Load preceding repository
	return api.repositoryManager.LoadPrecedingRepository()
}

func (api *Api) WalkRepositoryRecipes(repository app.Repository, walker func(recipe app.Recipe) error) error {
	// Log
	api.log.Debug("walk repository recipes…")

	// Walk repository recipes
	return api.recipeManager.WalkRecipes(repository, walker)
}
