package main

import (
	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/spf13/viper"
	"manala/cmd"
	"os"
	"path/filepath"
)

// Set at build time, by goreleaser, via ldflags
var version = "dev"

// Default repository
var defaultRepository = "https://github.com/manala/manala-recipes.git"

func main() {
	// Logger
	logger := &log.Logger{
		Handler: cli.New(os.Stderr),
		Level:   log.InfoLevel,
	}

	// Config
	config := viper.New()
	config.SetEnvPrefix("manala")
	config.AutomaticEnv()
	config.SetDefault("debug", false)
	config.Set("version", version)
	config.SetDefault("repository", defaultRepository)

	// Cache dir
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		logger.Fatal(err.Error())
	}
	config.SetDefault("cache-dir", filepath.Join(cacheDir, "manala"))

	// Commands
	rootCommand := (&cmd.RootCmd{}).Command(config, logger)
	rootCommand.AddCommand(
		(&cmd.InitCmd{}).Command(config, logger),
		(&cmd.ListCmd{}).Command(config, logger),
		(&cmd.UpdateCmd{}).Command(config, logger),
		(&cmd.WatchCmd{}).Command(config, logger),
		(&cmd.MascotCmd{}).Command(),
	)

	// Docs generation command
	if config.GetString("version") == "dev" {
		rootCommand.AddCommand(
			(&cmd.DocsCmd{}).Command(rootCommand, "docs/commands"),
		)
	}

	// Execute
	if err := rootCommand.Execute(); err != nil {
		logger.Fatal(err.Error())
	}
}
