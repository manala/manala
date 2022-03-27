package main

import (
	"embed"
	"manala/app"
	"manala/cmd"
	"manala/config"
	"manala/fs"
	"manala/loaders"
	"manala/logger"
	"manala/models"
	"manala/syncer"
	"manala/template"
	"os"
)

// Main repository source
var mainRepository = "https://github.com/manala/manala-recipes.git"

// Set at build time, by goreleaser, via ldflags
var version = "dev"

//go:embed assets/*
var assets embed.FS

func main() {
	// Config
	conf := config.New(
		config.WithVersion(version),
		config.WithMainRepository(mainRepository),
	)

	// Logger
	log := logger.New(
		logger.WithConfig(conf),
		logger.WithWriter(os.Stderr),
	)

	// Managers
	fsManager := fs.NewManager()
	modelFsManager := models.NewFsManager(fsManager)
	templateManager := template.NewManager()
	modelTemplateManager := models.NewTemplateManager(templateManager, modelFsManager)
	modelWatcherManager := models.NewWatcherManager(log)

	// Syncer
	sync := syncer.New(log, modelFsManager, modelTemplateManager)

	// Loaders
	repositoryLoader := loaders.NewRepositoryLoader(log, conf)
	recipeLoader := loaders.NewRecipeLoader(log, modelFsManager)
	projectLoader := loaders.NewProjectLoader(log, conf, repositoryLoader, recipeLoader)

	// App
	manala := &app.App{
		RepositoryLoader: repositoryLoader,
		RecipeLoader:     recipeLoader,
		ProjectLoader:    projectLoader,
		TemplateManager:  modelTemplateManager,
		WatcherManager:   modelWatcherManager,
		Sync:             sync,
		Log:              log,
	}

	// Commands
	rootCommand := (&cmd.RootCmd{Conf: conf}).Command()
	rootCommand.AddCommand(
		(&cmd.InitCmd{App: manala, Conf: conf, Assets: assets}).Command(),
		(&cmd.ListCmd{App: manala, Conf: conf, Out: rootCommand.OutOrStdout()}).Command(),
		(&cmd.UpdateCmd{App: manala}).Command(),
		(&cmd.WatchCmd{App: manala}).Command(),
		(&cmd.MascotCmd{Assets: assets}).Command(),
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
