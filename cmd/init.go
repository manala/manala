package cmd

import (
	"github.com/apex/log"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"manala/pkg/project"
	"manala/pkg/recipe"
	"manala/pkg/repository"
	"manala/pkg/syncer"
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
	// Create project
	prj := project.New(viper.GetString("dir"))

	if prj.IsExist() {
		log.Fatal("Project already initialized")
	}

	// Load repository
	repo, err := repository.Load(viper.GetString("repository"), viper.GetString("cache_dir"))
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Info("Repository loaded")

	var recipes []recipe.Recipe

	// Walk into recipes
	err = recipe.Walk(repo, func(rec *recipe.Recipe) {
		recipes = append(recipes, *rec)
	})
	if err != nil {
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
{{ .Config.Description }}`,
		},
		Searcher: func(input string, index int) bool {
			rec := recipes[index]
			name := strings.Replace(strings.ToLower(rec.Name), " ", "", -1)
			description := strings.Replace(strings.ToLower(rec.Name), " ", "", -1)
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

	rec := recipes[index]

	// Sync project
	if err := syncer.SyncProject(prj, &rec); err != nil {
		log.Fatal(err.Error())
	}

	log.Info("Project synced")
}
