package recipe

import (
	"golang.org/x/exp/maps"
	"manala/core"
	internalLog "manala/internal/log"
	internalWatcher "manala/internal/watcher"
	"sort"
)

func NewNameProcessorManager(log *internalLog.Logger, cascadingManager core.RecipeManager) *NameProcessorManager {
	return &NameProcessorManager{
		log:              log,
		cascadingManager: cascadingManager,
		names:            map[int]string{},
	}
}

type NameProcessorManager struct {
	log              *internalLog.Logger
	cascadingManager core.RecipeManager
	names            map[int]string
}

func (manager *NameProcessorManager) AddName(name string, priority int) {
	manager.names[priority] = name
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
	// Clone manager names to add current processed one with priority 0 without touching them
	names := maps.Clone(manager.names)
	names[0] = name

	// Reverse order priorities
	priorities := maps.Keys(names)
	sort.Sort(sort.Reverse(sort.IntSlice(priorities)))

	var processedName string

	for _, priority := range priorities {
		_name := names[priority]

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
