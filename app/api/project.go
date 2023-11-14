package api

import (
	"errors"
	"github.com/gen2brain/beeep"
	"manala/app"
	"manala/internal/filepath/backwalk"
	"manala/internal/serrors"
	"manala/internal/watcher"
	"os"
	"path/filepath"
	"slices"
)

func (api *Api) IsProject(dir string) bool {
	return api.projectManager.IsProject(dir)
}

func (api *Api) CreateProject(dir string, recipe app.Recipe, vars map[string]any) (app.Project, error) {
	// Log
	api.log.Debug("create project…",
		"dir", dir,
	)

	// Create project
	return api.projectManager.CreateProject(dir, recipe, vars)
}

func (api *Api) LoadProjectFrom(dir string) (app.Project, error) {
	var (
		project app.Project
		err     error
	)

	// Log
	api.log.Debug("backwalk projects from…",
		"dir", dir,
	)

	// Backwalk projects from dir
	err = backwalk.Backwalk(
		dir,
		func(path string, entry os.DirEntry, err error) error {
			if err != nil {
				return serrors.New("file system error").
					WithArguments("path", path).
					WithErrors(serrors.NewOs(err))
			}

			// Log
			api.log.Debug("try to load project…",
				"dir", path,
			)

			// Load project
			project, err = api.projectManager.LoadProject(path)
			if err != nil {
				var _notFoundProjectManifestError *app.NotFoundProjectManifestError
				if errors.As(err, &_notFoundProjectManifestError) {
					err = nil
				}
				return err
			}

			// Stop backwalk
			return filepath.SkipDir
		})

	if err != nil {
		return nil, err
	}

	if project == nil {
		return nil, serrors.New("project not found").
			WithArguments("dir", dir)
	}

	return project, nil
}

func (api *Api) WalkProjects(dir string, walker func(project app.Project) error) error {
	// Log
	api.log.Info("walk projects from…",
		"dir", dir,
	)

	err := filepath.WalkDir(
		dir,
		func(path string, entry os.DirEntry, err error) error {
			if err != nil {
				return serrors.New("file system error").
					WithArguments("path", path).
					WithErrors(serrors.NewOs(err))
			}

			// Only directories
			if !entry.IsDir() {
				return nil
			}

			// Exclusions
			if slices.Contains(api.exclusionPaths, filepath.Base(path)) {
				// Log
				api.log.Debug("exclude path",
					"path", path,
				)

				return filepath.SkipDir
			}

			// Log
			api.log.Debug("try to load project…",
				"dir", path,
			)

			// Load project
			project, err := api.projectManager.LoadProject(path)
			if err != nil {
				var _notFoundProjectManifestError *app.NotFoundProjectManifestError
				if errors.As(err, &_notFoundProjectManifestError) {
					err = nil
				}
				return err
			}

			// Walk function
			return walker(project)
		},
	)

	return err
}

func (api *Api) SyncProject(project app.Project) error {
	// Log
	api.log.Info("sync project…",
		"src", project.Recipe().Dir(),
		"dst", project.Dir(),
	)

	// Loop over project recipe sync units
	for _, unit := range project.Recipe().Sync() {
		if err := api.syncer.Sync(
			project.Recipe().Dir(),
			unit.Source(),
			project.Dir(),
			unit.Destination(),
			project,
		); err != nil {

			return err
		}
	}

	return nil
}

func (api *Api) WatchProject(project app.Project, all bool, notify bool) error {
	// Log
	api.log.Info("watch project…",
		"src", project.Recipe().Dir(),
		"dst", project.Dir(),
	)

	dir := project.Dir()

	watcher, err := api.watcherManager.NewWatcher(

		// On start
		func(watcher *watcher.Watcher) {
			// Watch project
			_ = api.projectManager.WatchProject(project, watcher)
		},

		// On change
		func(watcher *watcher.Watcher) {
			// Load project
			var err error
			project, err := api.projectManager.LoadProject(dir)
			if err != nil {
				if notify {
					_ = beeep.Alert("Manala", err.Error(), "")
				}
				api.out.Error(err)
				return
			}

			// Sync project
			err = api.SyncProject(project)

			if err != nil {
				if notify {
					_ = beeep.Alert("Manala", err.Error(), "")
				}
				api.out.Error(err)
				return
			}

			if notify {
				_ = beeep.Notify("Manala", "Project synced", "")
			}
		},

		// On all
		func(watcher *watcher.Watcher) {
			if all && project != nil {
				_ = api.recipeManager.WatchRecipe(project.Recipe(), watcher)
			}
		},
	)
	if err != nil {
		return nil
	}

	//goland:noinspection GoUnhandledErrorResult
	defer watcher.Close()

	watcher.Watch()

	return nil
}
