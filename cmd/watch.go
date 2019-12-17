package cmd

import (
	"github.com/apex/log"
	"github.com/fsnotify/fsnotify"
	"github.com/gen2brain/beeep"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"manala/pkg/project"
	"manala/pkg/repository"
	"manala/pkg/sync"
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

	var prj project.Interface

	// Get sync function
	syncProject := watchSyncProjectFunc(&prj, watcher, watchRecipe)

	// Sync
	if err := syncProject(); err != nil {
		log.Fatal(err.Error())
	}

	// Watch project
	if err := watcher.Add(prj.GetDir()); err != nil {
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
					file := filepath.Clean(event.Name)
					dir := filepath.Dir(file)
					if file == prj.GetConfigFile() {
						log.WithField("file", file).Info("Project config modified")
						modified = true
					} else if dir != prj.GetDir() {
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
}

func watchSyncProjectFunc(basePrj *project.Interface, watcher *fsnotify.Watcher, watchRecipe bool) func() error {
	var baseRecDir string

	return func() error {
		// Create project
		prj := project.New(viper.GetString("dir"))

		// Load project
		if err := prj.Load(project.Config{
			Repository: viper.GetString("repository"),
		}); err != nil {
			return err
		}

		*basePrj = prj

		log.WithFields(log.Fields{
			"recipe":     prj.GetConfig().Recipe,
			"repository": prj.GetConfig().Repository,
		}).Info("Project loaded")

		// Load repository
		repo := repository.New(prj.GetConfig().Repository)
		if err := repo.Load(viper.GetString("cache_dir")); err != nil {
			return err
		}

		log.Info("Repository loaded")

		// Load recipe
		rec, err := repo.LoadRecipe(prj.GetConfig().Recipe)
		if err != nil {
			log.Fatal(err.Error())
		}

		log.Info("Recipe loaded")

		if watchRecipe {
			// Initialize base recipe dir to the first synced recipe
			if baseRecDir == "" {
				baseRecDir = rec.GetDir()
			}

			// If recipe has changed, first, unwatch old one directories
			if baseRecDir != rec.GetDir() {
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
			if err := filepath.Walk(rec.GetDir(), func(path string, info os.FileInfo, err error) error {
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
		if err := sync.SyncProject(prj, rec); err != nil {
			return err
		}

		log.Info("Project synced")

		return nil
	}
}
