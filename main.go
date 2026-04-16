package main

import (
	"errors"
	"os"

	"github.com/manala/manala/app/api"
	"github.com/manala/manala/cmd"
	cmdDocs "github.com/manala/manala/cmd/docs"
	cmdInit "github.com/manala/manala/cmd/init"
	cmdList "github.com/manala/manala/cmd/list"
	cmdMascot "github.com/manala/manala/cmd/mascot"
	cmdUpdate "github.com/manala/manala/cmd/update"
	cmdWatch "github.com/manala/manala/cmd/watch"
	"github.com/manala/manala/internal/caching"
	"github.com/manala/manala/internal/log"
	"github.com/manala/manala/internal/notify"
	"github.com/manala/manala/internal/output"

	"charm.land/lipgloss/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Set at build time, by goreleaser, via ldflags.
var version = "dev"

// Default repository url.
const defaultRepositoryURL = "https://github.com/manala/manala-recipes.git"

func main() {
	// Os
	stdin, stdout, stderr := os.Stdin, os.Stdout, os.Stderr
	env := os.Environ()

	// Notifier
	notifier := notify.New(notify.NewBeeepHandler("Manala"))

	// Logger
	logger := log.New(output.New(stdin, stderr, env))

	// Output
	out := output.New(stdin, stdout, env)

	// Api
	appApi := new(api.API)

	// Commands
	command := cmd.NewCommand(version, stdin, stdout, stderr)
	command.AddCommand(
		cmdInit.NewCommand(logger, appApi, out),
		cmdList.NewCommand(logger, appApi, out),
		cmdMascot.NewCommand(stdin, stdout),
		cmdUpdate.NewCommand(logger, appApi, out),
		cmdWatch.NewCommand(logger, appApi, out, notifier),
	)

	// Commands persistent flags
	command.PersistentFlags().StringP("cache-dir", "c", "", "use cache directory")
	command.PersistentFlags().CountP("verbose", "v", "more verbose output (repeatable)")

	// Docs command only available in dev
	if version == "dev" {
		command.AddCommand(cmdDocs.NewCommand(command))
	}

	cobra.OnInitialize(func() {
		// Viper
		v := viper.New()

		_ = v.BindPFlag("cache_dir", command.PersistentFlags().Lookup("cache-dir"))
		_ = v.BindPFlag("verbose", command.PersistentFlags().Lookup("verbose"))
		v.SetDefault("default_repository", defaultRepositoryURL)

		// Viper - Env
		v.AutomaticEnv()
		v.SetEnvPrefix("MANALA")

		// Cache
		cache := caching.NewCache(v.GetString("cache_dir")).
			WithUserDir("manala")

		// Logger verbose mode
		logger.Verbose(v.GetInt("verbose"))

		// Deferred app api instantiation
		*appApi = *api.New(logger, cache,
			api.WithDefaultRepositoryURL(v.GetString("default_repository")),
		)

		// Log config
		logger.Debug("config",
			"default_repository", v.GetString("default_repository"),
			"cache_dir", v.GetString("cache_dir"),
			"verbose", v.GetInt("verbose"),
		)
	})

	// Execute command
	if err := command.Execute(); err != nil {
		if _, ok := errors.AsType[*cmd.CancelError](err); ok {
			lipgloss.Fprintln(stdout, err.Error())
			os.Exit(0)
		}
		logger.Error(err)
		os.Exit(1)
	}
}
