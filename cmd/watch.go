package cmd

import (
	"fmt"
	"github.com/gen2brain/beeep"
	"github.com/spf13/cobra"
	"manala/loaders"
	"manala/logger"
	"manala/models"
	"manala/syncer"
	"manala/validator"
	"os"
	"strings"
)

type WatchCmd struct {
	Log            logger.Logger
	ProjectLoader  loaders.ProjectLoaderInterface
	WatcherManager models.WatcherManagerInterface
	Sync           *syncer.Syncer
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

			withRepositorySource, _ := flags.GetString("repository")
			withRecipeName, _ := flags.GetString("recipe")

			watchAll, _ := flags.GetBool("all")
			useNotify, _ := flags.GetBool("notify")

			return cmd.Run(dir, withRepositorySource, withRecipeName, watchAll, useNotify)
		},
	}

	flags := command.Flags()

	flags.StringP("repository", "o", "", "with repository source")
	flags.StringP("recipe", "i", "", "with recipe name")

	flags.BoolP("all", "a", false, "watch recipe too")
	flags.BoolP("notify", "n", false, "use system notifications")

	return command
}

func (cmd *WatchCmd) Run(dir string, withRepositorySource string, withRecipeName string, watchAll bool, useNotify bool) error {
	// Check directory
	if dir != "." {
		if _, err := os.Stat(dir); err != nil {
			return fmt.Errorf("invalid directory: %s", dir)
		}
	}

	// Find project manifest
	prjManifest, err := cmd.ProjectLoader.Find(dir, true)
	if err != nil {
		return err
	}

	if prjManifest == nil {
		return fmt.Errorf("project not found: %s", dir)
	}

	// Sync function
	syncFunc := cmd.getSyncFunc(prjManifest, withRepositorySource, withRecipeName, watchAll)

	// Watcher
	watcher, err := cmd.WatcherManager.NewWatcher()
	if err != nil {
		return fmt.Errorf("error creating watcher: %v", err)
	}
	defer watcher.Close()

	// Sync
	if err := syncFunc(watcher); err != nil {
		return err
	}

	cmd.Log.Info("Start watching...")

	// Watch
	watcher.Watch(func(watcher models.WatcherInterface) {
		if err := syncFunc(watcher); err != nil {
			cmd.Log.Error(err.Error())
			if useNotify {
				_ = beeep.Alert("Manala", strings.Replace(err.Error(), `"`, `\"`, -1), "")
			}
		} else {
			if useNotify {
				_ = beeep.Notify("Manala", "Project synced", "")
			}
		}
	})

	return nil
}

func (cmd *WatchCmd) getSyncFunc(prjManifest *os.File, withRepositorySource string, withRecipeName string, watchAll bool) func(watcher models.WatcherInterface) error {
	return func(watcher models.WatcherInterface) error {
		// Load project
		prj, err := cmd.ProjectLoader.Load(prjManifest, withRepositorySource, withRecipeName)
		if err != nil {
			return err
		}

		// Validate project
		if err := validator.ValidateProject(prj); err != nil {
			return err
		}

		cmd.Log.Info("Project validated")

		// Watch project
		if err := watcher.SetProject(prj); err != nil {
			return fmt.Errorf("error setting project watching: %v", err)
		}

		// Watch recipe
		if watchAll {
			if err := watcher.SetRecipe(prj.Recipe()); err != nil {
				return fmt.Errorf("error setting recipe watching: %v", err)
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
