package api

import "manala/app"

func (api *Api) LoadPrecedingRepository() (app.Repository, error) {
	// Log
	api.log.Debug("load preceding repository…")

	// Load preceding repository
	return api.repositoryManager.LoadPrecedingRepository()
}

func (api *Api) RepositoryRecipes(repository app.Repository) ([]app.Recipe, error) {
	// Log
	api.log.Debug("repository recipes…")

	// Walk repository recipes
	return api.recipeManager.RepositoryRecipes(repository)
}
