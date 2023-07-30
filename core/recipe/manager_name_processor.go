package recipe

import (
	"log/slog"
	"manala/app/interfaces"
	"manala/core"
	"manala/internal/watcher"
	"maps"
	"sort"
)

func NewNameProcessorManager(log *slog.Logger, cascadingManager interfaces.RecipeManager) *NameProcessorManager {
	return &NameProcessorManager{
		log:              log.With("manager", "name_processor"),
		cascadingManager: cascadingManager,
		names:            map[int]string{},
	}
}

type NameProcessorManager struct {
	log              *slog.Logger
	cascadingManager interfaces.RecipeManager
	names            map[int]string
}

func (manager *NameProcessorManager) AddName(name string, priority int) {
	manager.names[priority] = name
}

func (manager *NameProcessorManager) LoadRecipe(repo interfaces.Repository, name string) (interfaces.Recipe, error) {
	// Log
	manager.log.Debug("load recipe",
		"name", name,
	)

	// Process name
	name, err := manager.processName(name)
	if err != nil {
		return nil, err
	}

	// Log
	manager.log.Debug("cascade recipe loading…",
		"name", name,
	)

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
	// Log
	manager.log.Debug("walk recipes",
		"repository", repo.Url(),
	)

	// Log
	manager.log.Debug("cascade recipes walking…",
		"repository", repo.Url(),
	)

	if err := manager.cascadingManager.WalkRecipes(repo, walker); err != nil {
		return err
	}

	return nil
}

func (manager *NameProcessorManager) WatchRecipe(rec interfaces.Recipe, watcher *watcher.Watcher) error {
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
	priorities := make([]int, 0, len(names))
	for priority := range names {
		priorities = append(priorities, priority)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(priorities)))

	for _, priority := range priorities {
		name = names[priority]

		// Log
		manager.log.Debug("process recipe name",
			"name", name,
			"priority", priority,
		)

		if name != "" {
			break
		}
	}

	if name == "" {
		return "", &core.UnprocessableRecipeNameError{}
	}

	return name, nil
}
