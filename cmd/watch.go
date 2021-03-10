package cmd

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/gen2brain/beeep"
	"github.com/spf13/cobra"
	"manala/loaders"
	"manala/logger"
	"manala/models"
	"manala/syncer"
	"manala/validator"
	"os"
	"path/filepath"
	"strings"
)

type WatchCmd struct {
	Log           *logger.Logger
	ProjectLoader loaders.ProjectLoaderInterface
	Sync          *syncer.Syncer
}

func (cmd *WatchCmd) Command() *cobra.Command {
	command := &cobra.Command{
		Use:     "watch [dir]",
		Aliases: []string{"Watch project"},
		Short:   "Watch project",
		Long: `Watch (manala watch) will watch project, and launch update on changes.

Example: manala watch -> resulting in a watch in a directory (default to the current directory)`,
		Args:              cobra.MaximumNArgs(1),
		DisableAutoGenTag: true,
		RunE: func(command *cobra.Command, args []string) error {
			// Get directory from first command arg
			dir := "."
			if len(args) != 0 {
				dir = args[0]
			}

			flags := command.Flags()

			repoSrc, _ := flags.GetString("repository")
			recName, _ := flags.GetString("recipe")

			watchAll, _ := flags.GetBool("all")
			useNotify, _ := flags.GetBool("notify")

			return cmd.Run(dir, repoSrc, recName, watchAll, useNotify)
		},
	}

	flags := command.Flags()

	flags.StringP("repository", "o", "", "with repository source")
	flags.StringP("recipe", "i", "", "with recipe name")

	flags.BoolP("all", "a", false, "watch recipe too")
	flags.BoolP("notify", "n", false, "use system notifications")

	return command
}

func (cmd *WatchCmd) Run(dir string, repoSrc string, recName string, watchAll bool, useNotify bool) error {
	// New watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("error creating watcher: %v", err)
	}
	defer watcher.Close()

	// Check directory
	if dir != "." {
		if _, err := os.Stat(dir); err != nil {
			return fmt.Errorf("invalid directory: %s", dir)
		}
	}

	// Find project file
	prjFile, err := cmd.ProjectLoader.Find(dir, true)
	if err != nil {
		return err
	}

	if prjFile == nil {
		return fmt.Errorf("project not found: %s", dir)
	}

	var prj models.ProjectInterface

	// Get sync function
	syncProject := cmd.runProjectSync(prjFile, &prj, repoSrc, recName, watcher, watchAll)

	// Sync
	if err := syncProject(); err != nil {
		return err
	}

	// Watch project
	if err := watcher.Add(prj.Dir()); err != nil {
		return fmt.Errorf("error adding project watching: %v", err)
	}

	cmd.Log.Info("Start watching...")

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				cmd.Log.DebugWithField("Watch event", "event", event)

				if event.Op != fsnotify.Chmod {
					modified := false
					file := filepath.Clean(event.Name)
					dir := filepath.Dir(file)
					if file == prjFile.Name() {
						cmd.Log.InfoWithField("Project config modified", "file", file)
						modified = true
					} else if dir != prj.Dir() {
						// Modified directory is not project one. That could only means recipe's one
						cmd.Log.InfoWithField("Recipe modified", "dir", dir)
						modified = true
					}

					if modified {
						if err := syncProject(); err != nil {
							cmd.Log.Error(err.Error())
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
				cmd.Log.ErrorWithError("Watching error", err)
			}
		}
	}()
	<-done

	return nil
}

func (cmd *WatchCmd) runProjectSync(prjFile *os.File, basePrj *models.ProjectInterface, repoSrc string, recName string, watcher *fsnotify.Watcher, watchAll bool) func() error {
	var baseRecDir string

	return func() error {
		// Load project
		prj, err := cmd.ProjectLoader.Load(prjFile, repoSrc, recName)
		if err != nil {
			return err
		}

		// Validate project
		if err := validator.ValidateProject(prj); err != nil {
			return err
		}

		cmd.Log.Info("Project validated")

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
		if err := cmd.Sync.SyncProject(prj); err != nil {
			return err
		}

		cmd.Log.Info("Project synced")

		return nil
	}
}
