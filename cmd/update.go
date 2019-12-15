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
	// Create project
	prj := project.New(viper.GetString("dir"))

	// Load project
	if err := prj.Load(project.Config{
		Repository: viper.GetString("repository"),
	}); err != nil {
		log.Fatal(err.Error())
	}

	log.WithFields(log.Fields{
		"recipe":     prj.GetConfig().Recipe,
		"repository": prj.GetConfig().Repository,
	}).Info("Project loaded")

	// Load repository
	repo := repository.New(prj.GetConfig().Repository)
	if err := repo.Load(viper.GetString("cache_dir")); err != nil {
		log.Fatal(err.Error())
	}

	log.Info("Repository loaded")

	// Load recipe
	rec := recipe.New(prj.GetConfig().Recipe)
	if err := rec.Load(repo); err != nil {
		log.Fatal(err.Error())
	}

	log.Info("Recipe loaded")

	// Sync project
	if err := syncer.SyncProject(prj, rec); err != nil {
		log.Fatal(err.Error())
	}

	log.Info("Project synced")
}
