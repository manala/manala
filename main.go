package main

import (
	"manala/cmd"
	"manala/config"
	"manala/loaders"
	"manala/logger"
	"manala/syncer"
	"manala/template"
	"os"
)

// Main repository source
var repository = "https://github.com/manala/manala-recipes.git"

// Set at build time, by goreleaser, via ldflags
var version = "dev"

func main() {
	// Config
	conf := config.New(version, repository)

	// Logger
	log := logger.New(conf)

	// Template
	tmpl := template.New()

	// Syncer
	sync := syncer.New(log, tmpl)

	// Loaders
	repositoryLoader := loaders.NewRepositoryLoader(log, conf)
	recipeLoader := loaders.NewRecipeLoader(log)
	projectLoader := loaders.NewProjectLoader(log, repositoryLoader, recipeLoader)

	// Commands
	rootCommand := (&cmd.RootCmd{Conf: conf}).Command()
	rootCommand.AddCommand(
		(&cmd.InitCmd{Log: log, RepositoryLoader: repositoryLoader, RecipeLoader: recipeLoader, ProjectLoader: projectLoader, Sync: sync}).Command(),
		(&cmd.ListCmd{RepositoryLoader: repositoryLoader, RecipeLoader: recipeLoader, Out: rootCommand.OutOrStdout()}).Command(),
		(&cmd.UpdateCmd{Log: log, ProjectLoader: projectLoader, Sync: sync}).Command(),
		(&cmd.WatchCmd{Log: log, ProjectLoader: projectLoader, Sync: sync}).Command(),
	)

	// Docs generation command
	if conf.Version() == "dev" {
		rootCommand.AddCommand(
			(&cmd.DocsCmd{RootCommand: rootCommand, Dir: "docs/commands"}).Command(),
		)
	}

	// Execute command
	if err := rootCommand.Execute(); err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}
