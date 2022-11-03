package recipe

import (
	"github.com/caarlos0/log"
	"manala/core"
	internalLog "manala/internal/log"
	internalWatcher "manala/internal/watcher"
)

func NewLogManager(log *internalLog.Logger, cascadingManager core.RecipeManager) *LogManager {
	return &LogManager{
		log:              log,
		cascadingManager: cascadingManager,
	}
}

type LogManager struct {
	log              *internalLog.Logger
	cascadingManager core.RecipeManager
}

func (manager *LogManager) LoadRecipe(repo core.Repository, name string) (core.Recipe, error) {
	// Log
	manager.log.WithFields(log.Fields{
		"name": name,
	}).Debug("load recipe")
	manager.log.IncreasePadding()

	rec, err := manager.cascadingManager.LoadRecipe(repo, name)

	// Log
	manager.log.DecreasePadding()

	if err != nil {
		return nil, err
	}

	return rec, nil
}

func (manager *LogManager) WalkRecipes(repo core.Repository, walker func(rec core.Recipe) error) error {
	err := manager.cascadingManager.WalkRecipes(repo, walker)

	if err != nil {
		return err
	}

	return nil
}

func (manager *LogManager) WatchRecipe(rec core.Recipe, watcher *internalWatcher.Watcher) error {
	err := manager.cascadingManager.WatchRecipe(rec, watcher)

	if err != nil {
		return err
	}

	return nil
}
