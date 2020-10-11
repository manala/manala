package cmd

import (
	"fmt"
	"github.com/apex/log"
	"github.com/fsnotify/fsnotify"
	"github.com/gen2brain/beeep"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"manala/loaders"
	"manala/models"
	"manala/syncer"
	"manala/validator"
	"os"
	"path/filepath"
	"strings"
)

// WatchCmd represents the watch command
func WatchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "watch [dir]",
		Aliases: []string{"Watch project"},
		Short:   "Watch project",
		Long: `Watch (manala watch) will watch project, and launch update on changes.

Example: manala watch -> resulting in a watch in a directory (default to the current directory)`,
		Args:              cobra.MaximumNArgs(1),
		DisableAutoGenTag: true,
		RunE:              watchRun,
	}

	addRepositoryFlag(cmd, "force repository")
	addRecipeFlag(cmd, "force recipe")

	cmd.Flags().BoolP("all", "a", false, "watch recipe too")
	cmd.Flags().BoolP("notify", "n", false, "use system notifications")

	return cmd
}

func watchRun(cmd *cobra.Command, args []string) error {
	// Get flags
	watchAll, _ := cmd.Flags().GetBool("all")
	useNotify, _ := cmd.Flags().GetBool("notify")

	// New watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("error creating watcher: %v", err)
	}
	defer watcher.Close()

	// Loaders
	repoLoader := loaders.NewRepositoryLoader(
		viper.GetString("cache_dir"),
		viper.GetString("repository"),
	)
	recLoader := loaders.NewRecipeLoader()
	repoSrc, _ := cmd.Flags().GetString("repository")
	recName, _ := cmd.Flags().GetString("recipe")
	prjLoader := loaders.NewProjectLoader(repoLoader, recLoader, repoSrc, recName)

	// Directory
	dir := "."
	if len(args) != 0 {
		// Get directory from first command arg
		dir = args[0]
		if _, err := os.Stat(dir); err != nil {
			return fmt.Errorf("invalid directory: %s", dir)
		}
	}

	// Find project file
	prjFile, err := prjLoader.Find(dir, true)
	if err != nil {
		return err
	}

	if prjFile == nil {
		return fmt.Errorf("project not found: %s", dir)
	}

	var prj models.ProjectInterface

	// Get sync function
	syncProject := watchSyncProjectFunc(prjFile, &prj, prjLoader, watcher, watchAll)

	// Sync
	if err := syncProject(); err != nil {
		return err
	}

	// Watch project
	if err := watcher.Add(prj.Dir()); err != nil {
		return fmt.Errorf("error adding project watching: %v", err)
	}

	log.Info("Start watching...")

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				log.WithField("event", event).Debug("Watch event")

				if event.Op != fsnotify.Chmod {
					modified := false
					file := filepath.Clean(event.Name)
					dir := filepath.Dir(file)
					if file == prjFile.Name() {
						log.WithField("file", file).Info("Project config modified")
						modified = true
					} else if dir != prj.Dir() {
						// Modified directory is not project one. That could only means recipe's one
						log.WithField("dir", dir).Info("Recipe modified")
						modified = true
					}

					if modified {
						if err := syncProject(); err != nil {
							log.Error(err.Error())
							if useNotify {
								_ = beeep.Alert("Manala", strings.Replace(err.Error(), `"`, `\"`, -1), "")
							}
						} else {
							if useNotify {
								_ = beeep.Notify("Manala", "Project synced", "")
							}
						}
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.WithError(err).Error("Watching error")
			}
		}
	}()
	<-done

	return nil
}

func watchSyncProjectFunc(file *os.File, basePrj *models.ProjectInterface, prjLoader loaders.ProjectLoaderInterface, watcher *fsnotify.Watcher, watchAll bool) func() error {
	var baseRecDir string

	return func() error {
		// Load project
		prj, err := prjLoader.Load(file)
		if err != nil {
			return err
		}

		// Validate project
		if err := validator.ValidateProject(prj); err != nil {
			return err
		}

		log.Info("Project validated")

		*basePrj = prj

		if watchAll {
			// Initialize base recipe dir to the first synced recipe
			if baseRecDir == "" {
				baseRecDir = prj.Recipe().Dir()
			}

			// If recipe has changed, first, unwatch old one directories
			if baseRecDir != prj.Recipe().Dir() {
				if err := filepath.Walk(baseRecDir, func(path string, info os.FileInfo, err error) error {
					if info.Mode().IsDir() {
						if err := watcher.Remove(path); err != nil {
							return err
						}
					}
					return nil
				}); err != nil {
					return err
				}
			}

			// Watch all recipe directories; don't care if they are already watched
			if err := filepath.Walk(prj.Recipe().Dir(), func(path string, info os.FileInfo, err error) error {
				if info.Mode().IsDir() {
					if err := watcher.Add(path); err != nil {
						return err
					}
				}
				return nil
			}); err != nil {
				return err
			}
		}

		// Sync project
		if err := syncer.SyncProject(prj); err != nil {
			return err
		}

		log.Info("Project synced")

		return nil
	}
}
