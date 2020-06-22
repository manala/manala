package cmd

import (
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

	// Loaders
	repoLoader := loaders.NewRepositoryLoader(viper.GetString("cache_dir"))
	recLoader := loaders.NewRecipeLoader()
	prjLoader := loaders.NewProjectLoader(repoLoader, recLoader, viper.GetString("repository"))

	var prj models.ProjectInterface

	// Get sync function
	syncProject := watchSyncProjectFunc(&prj, prjLoader, watcher, watchRecipe)

	// Sync
	if err := syncProject(); err != nil {
		log.Fatal(err.Error())
	}

	// Get project config file
	prjConfigFile, _ := prjLoader.ConfigFile(prj.Dir())

	// Watch project
	if err := watcher.Add(prj.Dir()); err != nil {
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
					if file == prjConfigFile.Name() {
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
}

func watchSyncProjectFunc(basePrj *models.ProjectInterface, prjLoader loaders.ProjectLoaderInterface, watcher *fsnotify.Watcher, watchRecipe bool) func() error {
	var baseRecDir string

	return func() error {
		// Load project
		prj, err := prjLoader.Load(viper.GetString("dir"))
		if err != nil {
			return err
		}

		// Validate project
		if err := validator.ValidateProject(prj); err != nil {
			return err
		}

		log.Info("Project validated")

		*basePrj = prj

		if watchRecipe {
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
