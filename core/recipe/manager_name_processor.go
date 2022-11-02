package recipe

import (
	"manala/core"
	internalLog "manala/internal/log"
	internalWatcher "manala/internal/watcher"
)

func NewNameProcessorManager(log *internalLog.Logger, cascadingManager core.RecipeManager) *NameProcessorManager {
	return &NameProcessorManager{
		log:              log,
		cascadingManager: cascadingManager,
	}
}

type NameProcessorManager struct {
	log              *internalLog.Logger
	uppermostName    string
	cascadingManager core.RecipeManager
}

func (manager *NameProcessorManager) WithUppermostName(name string) {
	manager.uppermostName = name
}

func (manager *NameProcessorManager) LoadRecipe(repo core.Repository, name string) (core.Recipe, error) {
	name, err := manager.processName(name)
	if err != nil {
		return nil, err
	}

	rec, err := manager.cascadingManager.LoadRecipe(repo, name)
	if err != nil {
		return nil, err
	}

	return rec, nil
}

func (manager *NameProcessorManager) LoadPrecedingRecipe(repo core.Repository) (core.Recipe, error) {
	return manager.LoadRecipe(repo, "")
}

func (manager *NameProcessorManager) WalkRecipes(repo core.Repository, walker func(rec core.Recipe) error) error {
	if err := manager.cascadingManager.WalkRecipes(repo, walker); err != nil {
		return err
	}

	return nil
}

func (manager *NameProcessorManager) WatchRecipe(rec core.Recipe, watcher *internalWatcher.Watcher) error {
	if err := manager.cascadingManager.WatchRecipe(rec, watcher); err != nil {
		return err
	}

	return nil
}

func (manager *NameProcessorManager) processName(name string) (string, error) {
	var processedName string

	for _, _name := range []string{manager.uppermostName, name} {
		if _name == "" {
			continue
		}

		processedName = _name
		break
	}

	if processedName == "" {
		return "", core.NewUnprocessableRecipeNameError(
			"unable to process empty recipe name",
		)
	}

	return processedName, nil
}
