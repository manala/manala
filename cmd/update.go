package cmd

import (
	"github.com/apex/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"manala/loaders"
	"manala/syncer"
	"manala/validator"
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
		RunE: updateRun,
		Args: cobra.MaximumNArgs(1),
	}

	addRepositoryFlag(cmd, "force repository")
	addRecipeFlag(cmd, "force recipe")

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

	// Project directory
	var dir string
	if len(args) != 0 {
		// Get directory from first command arg
		dir = args[0]
	}

	// Find project file
	prjFile, err := prjLoader.Find(dir, true)
	if err != nil {
		return err
	}

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
