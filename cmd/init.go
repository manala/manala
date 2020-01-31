package cmd

import (
	"github.com/apex/log"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"manala/loaders"
	"manala/models"
	"manala/syncer"
	"os"
	"strings"
)

// InitCmd represents the init command
func InitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "init",
		Aliases: []string{"in"},
		Short:   "Init project",
		Long: `Init (manala init) will init project.

A optional dir could be passed as argument.

Example: manala init -> resulting in an init in current directory
Example: manala init /foo/bar -> resulting in an init in /foo/bar directory`,
		Run:  initRun,
		Args: cobra.NoArgs,
	}

	return cmd
}

func initRun(cmd *cobra.Command, args []string) {
	// Loaders
	repoLoader := loaders.NewRepositoryLoader(viper.GetString("cache_dir"))
	recLoader := loaders.NewRecipeLoader()
	prjLoader := loaders.NewProjectLoader(repoLoader, recLoader, viper.GetString("repository"))

	// Ensure project is not yet initialized by checking configuration file existence
	cfgFile, _ := prjLoader.ConfigFile(viper.GetString("dir"))
	if cfgFile != nil {
		log.Fatal("Project already initialized")
	}

	// Load repository
	repo, err := repoLoader.Load(viper.GetString("repository"))
	if err != nil {
		log.Fatal(err.Error())
	}

	var recipes []models.RecipeInterface

	// Walk into recipes
	if err := recLoader.Walk(repo, func(rec models.RecipeInterface) {
		recipes = append(recipes, rec)
	}); err != nil {
		log.Fatal(err.Error())
	}

	prompt := promptui.Select{
		Items: recipes,
		Templates: &promptui.SelectTemplates{
			Label:    "Select recipe:",
			Active:   `{{ "▸" | bold }} {{ .Name | underline }}`,
			Inactive: "  {{ .Name }}",
			Selected: `{{ "✔" | green }} {{ .Name | faint }}`,
			Details: `
		{{ .Description }}`,
		},
		Searcher: func(input string, index int) bool {
			rec := recipes[index]
			name := strings.Replace(strings.ToLower(rec.Name()), " ", "", -1)
			description := strings.Replace(strings.ToLower(rec.Name()), " ", "", -1)
			input = strings.Replace(strings.ToLower(input), " ", "", -1)

			return strings.Contains(name, input) || strings.Contains(description, input)
		},
		Size:              12,
		StartInSearchMode: true,
	}

	index, _, err := prompt.Run()

	if err != nil {
		switch err {
		case promptui.ErrInterrupt:
			os.Exit(130)
		default:
			log.Fatal(err.Error())
		}
	}

	prj := models.NewProject(
		viper.GetString("dir"),
		recipes[index],
	)

	// Sync project
	if err := syncer.SyncProject(prj); err != nil {
		log.Fatal(err.Error())
	}

	log.Info("Project synced")
}
