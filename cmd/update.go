package cmd

import (
	"fmt"
	"github.com/apex/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"manala/loaders"
	"manala/syncer"
	"manala/validator"
	"os"
	"path/filepath"
	"strings"
)

// UpdateCmd represents the update command
func UpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update [dir]",
		Aliases: []string{"up"},
		Short:   "Update project",
		Long: `Update (manala update) will update project, based on
recipe and related variables defined in manala.yaml.

Example: manala update -> resulting in an update in a directory (default to the current directory)`,
		Args:              cobra.MaximumNArgs(1),
		DisableAutoGenTag: true,
		RunE:              updateRun,
	}

	addRepositoryFlag(cmd, "force repository")
	addRecipeFlag(cmd, "force recipe")

	cmd.Flags().BoolP("recursive", "r", false, "recursive")

	return cmd
}

func updateRun(cmd *cobra.Command, args []string) error {
	// Loaders
	repoLoader := loaders.NewRepositoryLoader(
		viper.GetString("cache_dir"),
		viper.GetString("repository"),
	)
	recLoader := loaders.NewRecipeLoader()
	repoName, _ := cmd.Flags().GetString("repository")
	recName, _ := cmd.Flags().GetString("recipe")
	prjLoader := loaders.NewProjectLoader(repoLoader, recLoader, repoName, recName)

	// Directory
	var dir string
	if len(args) != 0 {
		// Get directory from first command arg
		dir = args[0]
	}

	recursive, _ := cmd.Flags().GetBool("recursive")

	if recursive == true {
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			// Only directories
			if err != nil || !info.IsDir() {
				return err
			}

			// Only not dotted directories
			if strings.HasPrefix(filepath.Base(path), ".") {
				return filepath.SkipDir
			}

			// Find project file
			prjFile, err := prjLoader.Find(path, false)
			if err != nil {
				return err
			}

			// Update
			if prjFile != nil {
				if err := updateRunFunc(prjLoader, prjFile); err != nil {
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
		prjFile, err := prjLoader.Find(dir, true)
		if err != nil {
			return err
		}

		if prjFile == nil {
			return fmt.Errorf("project not found: %s", dir)
		}

		// Update
		if err = updateRunFunc(prjLoader, prjFile); err != nil {
			return err
		}
	}

	return nil
}

func updateRunFunc(prjLoader loaders.ProjectLoaderInterface, prjFile *os.File) error {
	// Load project
	prj, err := prjLoader.Load(prjFile)
	if err != nil {
		return err
	}

	// Validate project
	if err := validator.ValidateProject(prj); err != nil {
		return err
	}

	log.Info("Project validated")

	// Sync project
	if err := syncer.SyncProject(prj); err != nil {
		return err
	}

	log.Info("Project synced")

	return nil
}
