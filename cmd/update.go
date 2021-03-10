package cmd

import (
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
	Log           *logger.Logger
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

			repoSrc, _ := flags.GetString("repository")
			recName, _ := flags.GetString("recipe")

			recursive, _ := flags.GetBool("recursive")

			return cmd.Run(dir, repoSrc, recName, recursive)
		},
	}

	flags := command.Flags()

	flags.StringP("repository", "o", "", "with repository source")
	flags.StringP("recipe", "i", "", "with recipe name")

	flags.BoolP("recursive", "r", false, "set recursive mode")

	return command
}

func (cmd *UpdateCmd) Run(dir string, repoSrc string, recName string, recursive bool) error {
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

			// Find project file
			prjFile, err := cmd.ProjectLoader.Find(path, false)
			if err != nil {
				return err
			}

			// Sync
			if prjFile != nil {
				if err := cmd.runProjectSync(prjFile, repoSrc, recName); err != nil {
					return err
				}
			}

			return nil
		})
		if err != nil && !os.IsNotExist(err) {
			return err
		}
	} else {
		// Find project file
		prjFile, err := cmd.ProjectLoader.Find(dir, true)
		if err != nil {
			return err
		}

		if prjFile == nil {
			return fmt.Errorf("project not found: %s", dir)
		}

		// Sync
		if err = cmd.runProjectSync(prjFile, repoSrc, recName); err != nil {
			return err
		}
	}

	return nil
}

func (cmd *UpdateCmd) runProjectSync(prjFile *os.File, repoSrc string, recName string) error {
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

	// Sync project
	if err := cmd.Sync.SyncProject(prj); err != nil {
		return err
	}

	cmd.Log.Info("Project synced")

	return nil
}
