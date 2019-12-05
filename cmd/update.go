package cmd

import (
	"github.com/apex/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"manala/pkg/project"
	"manala/pkg/recipe"
	"manala/pkg/repository"
	"manala/pkg/syncer"
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
	// Load project
	prj, err := project.Load(viper.GetString("dir"), viper.GetString("repository"))
	if err != nil {
		log.Fatal(err.Error())
	}

	log.WithFields(log.Fields{
		"recipe":     prj.Config.Recipe,
		"repository": prj.Config.Repository,
	}).Info("Project loaded")

	// Load repository
	repo, err := repository.Load(prj.Config.Repository, viper.GetString("cache_dir"))
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Info("Repository loaded")

	// Lod recipe
	rec, err := recipe.Load(repo, prj.Config.Recipe)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Info("Recipe loaded")

	// Sync project
	if err := syncer.SyncProject(prj, rec); err != nil {
		log.Fatal(err.Error())
	}

	log.Info("Project synced")
}
