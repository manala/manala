package recipe

import (
	"golang.org/x/exp/maps"
	"manala/app/interfaces"
	"manala/core"
	internalLog "manala/internal/log"
	internalWatcher "manala/internal/watcher"
	"sort"
)

func NewNameProcessorManager(log *internalLog.Logger, cascadingManager interfaces.RecipeManager) *NameProcessorManager {
	return &NameProcessorManager{
		log:              log,
		cascadingManager: cascadingManager,
		names:            map[int]string{},
	}
}

type NameProcessorManager struct {
	log              *internalLog.Logger
	cascadingManager interfaces.RecipeManager
	names            map[int]string
}

func (manager *NameProcessorManager) AddName(name string, priority int) {
	manager.names[priority] = name
}

func (manager *NameProcessorManager) LoadRecipe(repo interfaces.Repository, name string) (interfaces.Recipe, error) {
	// Log
	manager.log.
		WithField("manager", "name_processor").
		WithField("name", name).
		Debug("load recipe")
	manager.log.IncreasePadding()
	defer manager.log.DecreasePadding()

	// Process name
	name, err := manager.processName(name)
	if err != nil {
		return nil, err
	}

	// Cascading manager
	rec, err := manager.cascadingManager.LoadRecipe(repo, name)
	if err != nil {
		return nil, err
	}

	return rec, nil
}

func (manager *NameProcessorManager) LoadPrecedingRecipe(repo interfaces.Repository) (interfaces.Recipe, error) {
	return manager.LoadRecipe(repo, "")
}

func (manager *NameProcessorManager) WalkRecipes(repo interfaces.Repository, walker func(rec interfaces.Recipe) error) error {
	if err := manager.cascadingManager.WalkRecipes(repo, walker); err != nil {
		return err
	}

	return nil
}

func (manager *NameProcessorManager) WatchRecipe(rec interfaces.Recipe, watcher *internalWatcher.Watcher) error {
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

	for _, priority := range priorities {
		name = names[priority]

		manager.log.
			WithField("name", name).
			WithField("priority", priority).
			Debug("process name")

		if name != "" {
			break
		}
	}

	if name == "" {
		return "", core.NewUnprocessableRecipeNameError(
			"unable to process empty recipe name",
		)
	}

	return name, nil
}
