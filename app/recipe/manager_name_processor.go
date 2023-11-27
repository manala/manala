package recipe

import (
	"log/slog"
	"manala/app"
	"manala/internal/watcher"
	"maps"
	"sort"
)

func NewNameProcessorManager(log *slog.Logger, cascadingManager app.RecipeManager) *NameProcessorManager {
	return &NameProcessorManager{
		log:              log.With("manager", "name_processor"),
		cascadingManager: cascadingManager,
		names:            map[int]string{},
	}
}

type NameProcessorManager struct {
	log              *slog.Logger
	cascadingManager app.RecipeManager
	names            map[int]string
}

func (manager *NameProcessorManager) AddName(name string, priority int) {
	manager.names[priority] = name
}

func (manager *NameProcessorManager) LoadRecipe(repository app.Repository, name string) (app.Recipe, error) {
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
	manager.log.Debug("cascading load recipe…",
		"name", name,
	)

	// Cascading manager
	return manager.cascadingManager.LoadRecipe(repository, name)
}

func (manager *NameProcessorManager) LoadPrecedingRecipe(repository app.Repository) (app.Recipe, error) {
	return manager.LoadRecipe(repository, "")
}

func (manager *NameProcessorManager) RepositoryRecipes(repository app.Repository) ([]app.Recipe, error) {
	// Log
	manager.log.Debug("repository recipes",
		"repository", repository.Url(),
	)

	// Log
	manager.log.Debug("cascading repository recipes…",
		"repository", repository.Url(),
	)

	return manager.cascadingManager.RepositoryRecipes(repository)
}

func (manager *NameProcessorManager) WatchRecipe(recipe app.Recipe, watcher *watcher.Watcher) error {
	// Log
	manager.log.Debug("watch recipe",
		"name", recipe.Name(),
	)

	// Log
	manager.log.Debug("cascading watch recipe…",
		"name", recipe.Name(),
	)

	return manager.cascadingManager.WatchRecipe(recipe, watcher)
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
		return "", &app.UnprocessableRecipeNameError{}
	}

	return name, nil
}
