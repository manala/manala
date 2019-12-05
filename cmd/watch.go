package cmd

import (
	"github.com/apex/log"
	"github.com/fsnotify/fsnotify"
	"github.com/gen2brain/beeep"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"manala/pkg/project"
	"manala/pkg/recipe"
	"manala/pkg/repository"
	"manala/pkg/syncer"
	"os"
	"path/filepath"
	"strings"
)

// WatchCmd represents the watch command
func WatchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "watch",
		Aliases: []string{"Watch project"},
		Short:   "Watch project",
		Long: `Watch (manala watch) will watch project, and launch update on changes.

		Example: manala watch -> resulting in an watch in current directory`,
		Run:  watchRun,
		Args: cobra.NoArgs,
	}

	cmd.Flags().BoolP("recipe", "i", false, "watch recipe too")
	cmd.Flags().BoolP("notify", "n", false, "use system notifications")

	return cmd
}

func watchRun(cmd *cobra.Command, args []string) {
	// Get flags
	watchRecipe, _ := cmd.Flags().GetBool("recipe")
	useNotify, _ := cmd.Flags().GetBool("notify")

	// New watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.WithError(err).Fatal("Error creating watcher")
	}
	defer watcher.Close()

	var prj project.Project

	// Get sync function
	syncProject := watchSyncProjectFunc(&prj, watcher, watchRecipe)

	// Sync
	if err := syncProject(); err != nil {
		log.Fatal(err.Error())
	}

	// Watch project
	if err := watcher.Add(prj.Dir); err != nil {
		log.WithError(err).Fatal("Error adding project watching")
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
					dir := filepath.Dir(filepath.Clean(event.Name))
					// Modified directory is not project one. That could only means recipe's one
					if dir != prj.Dir {
						log.WithField("dir", dir).Info("Recipe modified")
						modified = true
					} else {
						file := filepath.Base(filepath.Clean(event.Name))
						if file == prj.ConfigFile {
							log.WithField("file", file).Info("Project config modified")
							modified = true
						}
					}

					if modified {
						if err := syncProject(); err != nil {
							log.Error(err.Error())
							if useNotify {
								_ = beeep.Alert("Manala", strings.Replace(err.Error(), `"`, `\"`, -1), "")
							}
						} else {
							log.Info("Project synced")
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
}

func watchSyncProjectFunc(basePrj *project.Project, watcher *fsnotify.Watcher, watchRecipe bool) func() error {
	var baseRecDir string

	return func() error {
		// Load project
		prj, err := project.Load(viper.GetString("dir"), viper.GetString("repository"))
		if err != nil {
			return err
		}
		*basePrj = *prj

		log.WithFields(log.Fields{
			"recipe":     prj.Config.Recipe,
			"repository": prj.Config.Repository,
		}).Info("Project loaded")

		// Load repository
		repo, err := repository.Load(prj.Config.Repository, viper.GetString("cache_dir"))
		if err != nil {
			return err
		}

		log.Info("Repository loaded")

		// Load recipe
		rec, err := recipe.Load(repo, prj.Config.Recipe)
		if err != nil {
			return err
		}

		log.Info("Recipe loaded")

		if watchRecipe {
			// Initialize base recipe dir to the first synced recipe
			if baseRecDir == "" {
				baseRecDir = rec.Dir
			}

			// If recipe has changed, first, unwatch old one directories
			if baseRecDir != rec.Dir {
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
			if err := filepath.Walk(rec.Dir, func(path string, info os.FileInfo, err error) error {
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
		if err := syncer.SyncProject(prj, rec); err != nil {
			return err
		}

		log.Info("Project synced")

		return nil
	}
}
