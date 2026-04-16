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

	"charm.land/lipgloss/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Set at build time, by goreleaser, via ldflags.
var version = "dev"

// Default repository url.
const defaultRepositoryURL = "https://github.com/manala/manala-recipes.git"

func main() {
	// Streams
	in, out, err := os.Stdin, os.Stdout, os.Stderr

	// Notifier
	notifier := notify.New(notify.NewBeeepHandler("Manala"))

	appLog := log.New(err)
	appAPI := new(api.API)

	// App commands
	appCommand := cmd.NewCommand(version, in, out, err)
	appCommand.AddCommand(
		cmdInit.NewCommand(appLog, appAPI, out),
		cmdList.NewCommand(appLog, appAPI, out),
		cmdMascot.NewCommand(in, out),
		cmdUpdate.NewCommand(appLog, appAPI, out),
		cmdWatch.NewCommand(appLog, appAPI, out, notifier),
	)

	// App commands persistent flags
	appCommand.PersistentFlags().StringP("cache-dir", "c", "", "use cache directory")
	appCommand.PersistentFlags().CountP("verbose", "v", "more verbose output (repeatable)")

	// Docs app command only available in dev
	if version == "dev" {
		appCommand.AddCommand(cmdDocs.NewCommand(appCommand))
	}

	cobra.OnInitialize(func() {
		// Viper
		v := viper.New()

		_ = v.BindPFlag("cache_dir", appCommand.PersistentFlags().Lookup("cache-dir"))
		_ = v.BindPFlag("verbose", appCommand.PersistentFlags().Lookup("verbose"))
		v.SetDefault("default_repository", defaultRepositoryURL)

		// Viper - Env
		v.AutomaticEnv()
		v.SetEnvPrefix("MANALA")

		// App cache
		appCache := caching.NewCache(v.GetString("cache_dir")).
			WithUserDir("manala")

		// App log verbose mode
		appLog.Verbose(v.GetInt("verbose"))

		// Deferred app api instantiation
		*appAPI = *api.New(appLog, appCache,
			api.WithDefaultRepositoryURL(v.GetString("default_repository")),
		)

		// Log config
		appLog.Debug("config",
			"default_repository", v.GetString("default_repository"),
			"cache_dir", v.GetString("cache_dir"),
			"verbose", v.GetInt("verbose"),
		)
	})

	// Execute app command
	if err := appCommand.Execute(); err != nil {
		if _, ok := errors.AsType[*cmd.CancelError](err); ok {
			lipgloss.Fprintln(out, err.Error())
			os.Exit(0)
		}
		appLog.Error(err)
		os.Exit(1)
	}
}
