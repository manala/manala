package main

import (
	"manala/app"
	"manala/cmd"
	"os"
	"path/filepath"
)

// Default repository
var defaultRepository = "https://github.com/manala/manala-recipes.git"

// Set at build time, by goreleaser, via ldflags
var version = "dev"

func main() {
	// App
	manala := app.New(
		app.WithVersion(version),
		app.WithDefaultRepository(defaultRepository),
		app.WithLogWriter(os.Stderr),
	)

	// Config
	manala.Config.SetEnvPrefix("manala")
	manala.Config.AutomaticEnv()

	// Cache dir
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		manala.Log.Fatal(err.Error())
	}
	manala.Config.SetDefault("cache-dir", filepath.Join(cacheDir, "manala"))

	// Root command
	rootCommand := (&cmd.RootCmd{
		App:          manala,
		OnInitialize: manala.ApplyConfig,
	}).Command()

	// Commands
	rootCommand.AddCommand(
		(&cmd.InitCmd{App: manala}).Command(),
		(&cmd.ListCmd{App: manala, Out: rootCommand.OutOrStdout()}).Command(),
		(&cmd.UpdateCmd{App: manala}).Command(),
		(&cmd.WatchCmd{App: manala}).Command(),
		(&cmd.MascotCmd{}).Command(),
	)

	// Docs generation command
	if manala.Config.GetString("version") == "dev" {
		rootCommand.AddCommand(
			(&cmd.DocsCmd{RootCommand: rootCommand, Dir: "docs/commands"}).Command(),
		)
	}

	// Execute
	if err := rootCommand.Execute(); err != nil {
		manala.Log.Fatal(err.Error())
	}
}
