package project

import (
	_ "embed"
	"github.com/caarlos0/log"
	"manala/core"
	internalLog "manala/internal/log"
	internalWatcher "manala/internal/watcher"
)

func NewLogManager(log *internalLog.Logger, cascadingManager core.ProjectManager) *LogManager {
	return &LogManager{
		log:              log,
		cascadingManager: cascadingManager,
	}
}

type LogManager struct {
	log              *internalLog.Logger
	cascadingManager core.ProjectManager
}

func (manager LogManager) IsProject(dir string) bool {
	ok := manager.cascadingManager.IsProject(dir)

	return ok
}

func (manager LogManager) CreateProject(dir string, rec core.Recipe, vars map[string]interface{}) (core.Project, error) {
	proj, err := manager.cascadingManager.CreateProject(dir, rec, vars)

	if err != nil {
		return nil, err
	}

	return proj, nil
}

func (manager LogManager) LoadProject(dir string) (core.Project, error) {
	// Log
	manager.log.
		WithField("dir", dir).
		Debug("load project")
	manager.log.IncreasePadding()

	proj, err := manager.cascadingManager.LoadProject(dir)

	// Log
	manager.log.DecreasePadding()

	if err != nil {
		return nil, err
	}

	manager.log.WithFields(log.Fields{
		"dir":        proj.Dir(),
		"repository": proj.Recipe().Repository().Url(),
		"recipe":     proj.Recipe().Name(),
	}).Info("project loaded")

	return proj, nil
}

func (manager LogManager) WatchProject(proj core.Project, watcher *internalWatcher.Watcher) error {
	err := manager.cascadingManager.WatchProject(proj, watcher)

	if err != nil {
		return err
	}

	return nil
}
