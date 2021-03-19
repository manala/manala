package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"manala/loaders"
	"manala/logger"
	"manala/syncer"
	"manala/validator"
	"os"
	"path/filepath"
	"strings"
)

type UpdateCmd struct {
	Log           logger.Logger
	ProjectLoader loaders.ProjectLoaderInterface
	Sync          *syncer.Syncer
}

func (cmd *UpdateCmd) Command() *cobra.Command {
	command := &cobra.Command{
		Use:     "update [dir]",
		Aliases: []string{"up"},
		Short:   "Update project",
		Long: `Update (manala update) will update project, based on
recipe and related variables defined in manala.yaml.

Example: manala update -> resulting in an update in a directory (default to the current directory)`,
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

			recursive, _ := flags.GetBool("recursive")

			return cmd.Run(dir, withRepositorySource, withRecipeName, recursive)
		},
	}

	flags := command.Flags()

	flags.StringP("repository", "o", "", "with repository source")
	flags.StringP("recipe", "i", "", "with recipe name")

	flags.BoolP("recursive", "r", false, "set recursive mode")

	return command
}

func (cmd *UpdateCmd) Run(dir string, withRepositorySource string, withRecipeName string, recursive bool) error {
	// Check directory
	if dir != "." {
		if _, err := os.Stat(dir); err != nil {
			return fmt.Errorf("invalid directory: %s", dir)
		}
	}

	if recursive {
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			// Only directories
			if err != nil || !info.IsDir() {
				return err
			}

			// Only not dotted directories
			// (except - of course - current directory)
			if strings.HasPrefix(filepath.Base(path), ".") && (path != ".") {
				return filepath.SkipDir
			}

			// Find project manifest
			prjManifest, err := cmd.ProjectLoader.Find(path, false)
			if err != nil {
				return err
			}

			// Sync
			if prjManifest != nil {
				if err := cmd.runProjectSync(prjManifest, withRepositorySource, withRecipeName); err != nil {
					return err
				}
			}

			return nil
		})
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}
	} else {
		// Find project manifest
		prjManifest, err := cmd.ProjectLoader.Find(dir, true)
		if err != nil {
			return err
		}

		if prjManifest == nil {
			return fmt.Errorf("project not found: %s", dir)
		}

		// Sync
		if err = cmd.runProjectSync(prjManifest, withRepositorySource, withRecipeName); err != nil {
			return err
		}
	}

	return nil
}

func (cmd *UpdateCmd) runProjectSync(prjManifest *os.File, withRepositorySource string, withRecipeName string) error {
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

	// Sync project
	if err := cmd.Sync.SyncProject(prj); err != nil {
		return err
	}

	cmd.Log.Info("Project synced")

	return nil
}
