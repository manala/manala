package main

import (
	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"manala/cmd"
	"manala/internal/config"
	"os"
	"path/filepath"
	"strings"
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

	// Conf
	conf := config.New()
	conf.SetEnvPrefix("MANALA")
	conf.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	conf.AutomaticEnv()
	conf.SetDefault("debug", false)
	conf.Set("version", version)
	conf.SetDefault("repository", defaultRepository)

	// Cache dir
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		logger.Fatal(err.Error())
	}
	conf.SetDefault("cache-dir", filepath.Join(cacheDir, "manala"))

	// Commands
	rootCommand := (&cmd.RootCmd{}).Command(conf, logger)
	rootCommand.AddCommand(
		(&cmd.InitCmd{}).Command(conf, logger),
		(&cmd.ListCmd{}).Command(conf, logger),
		(&cmd.UpdateCmd{}).Command(conf, logger),
		(&cmd.WatchCmd{}).Command(conf, logger),
		(&cmd.MascotCmd{}).Command(),
	)

	// Docs generation command
	if conf.GetString("version") == "dev" {
		rootCommand.AddCommand(
			(&cmd.DocsCmd{}).Command(rootCommand, "docs/commands"),
		)
	}

	// Execute
	if err := rootCommand.Execute(); err != nil {
		logger.Fatal(err.Error())
	}
}
