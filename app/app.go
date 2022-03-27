package app

import (
	"manala/loaders"
	"manala/logger"
	"manala/models"
	"manala/syncer"
)

type App struct {
	RepositoryLoader loaders.RepositoryLoaderInterface
	RecipeLoader     loaders.RecipeLoaderInterface
	ProjectLoader    loaders.ProjectLoaderInterface
	TemplateManager  models.TemplateManagerInterface
	WatcherManager   models.WatcherManagerInterface
	Sync             *syncer.Syncer
	Log              logger.Logger
}
