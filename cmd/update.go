package cmd

import (
	"github.com/apex/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"manala/loaders"
	"manala/syncer"
)

// UpdateCmd represents the update command
func UpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update",
		Aliases: []string{"up"},
		Short:   "Update project",
		Long: `Update (manala update) will update project, based on
recipe and related variables defined in manala.yaml.

Example: manala update -> resulting in an update in current directory`,
		Run:  updateRun,
		Args: cobra.NoArgs,
	}

	return cmd
}

func updateRun(cmd *cobra.Command, args []string) {
	// Loaders
	repoLoader := loaders.NewRepositoryLoader(viper.GetString("cache_dir"))
	recLoader := loaders.NewRecipeLoader()
	prjLoader := loaders.NewProjectLoader(repoLoader, recLoader, viper.GetString("repository"))

	// Load project
	prj, err := prjLoader.Load(viper.GetString("dir"))
	if err != nil {
		log.Fatal(err.Error())
	}

	// Sync project
	if err := syncer.SyncProject(prj); err != nil {
		log.Fatal(err.Error())
	}

	log.Info("Project synced")
}
