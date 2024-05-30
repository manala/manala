package main

import (
	"log/slog"
	"manala/app/api"
	"manala/cmd"
	cmdDocs "manala/cmd/docs"
	cmdInit "manala/cmd/init"
	cmdList "manala/cmd/list"
	cmdMascot "manala/cmd/mascot"
	cmdUpdate "manala/cmd/update"
	cmdWatch "manala/cmd/watch"
	"manala/internal/cache"
	"manala/internal/notifier"
	"manala/internal/ui/adapters/charm"
	"manala/internal/ui/log"
	"os"

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

	// User interface
	ui := charm.New(in, out, err)

	// Notify
	notify := notifier.NewBeeep("Manala")

	var (
		appLog = new(slog.Logger)
		appAPI = new(api.API)
	)

	// App commands
	appCmd := cmd.NewCmd(version, in, out, err)
	appCmd.AddCommand(
		cmdInit.NewCmd(appLog, appAPI, ui),
		cmdList.NewCmd(appLog, appAPI, ui),
		cmdMascot.NewCmd(),
		cmdUpdate.NewCmd(appLog, appAPI),
		cmdWatch.NewCmd(appLog, appAPI, ui, notify),
	)

	// App commands persistent flags
	appCmd.PersistentFlags().StringP("cache-dir", "c", "", "use cache directory")
	appCmd.PersistentFlags().BoolP("debug", "d", false, "set debug mode")

	// Docs app command only available in dev
	if version == "dev" {
		appCmd.AddCommand(cmdDocs.NewCmd(appCmd))
	}

	cobra.OnInitialize(func() {
		// Viper
		_ = viper.BindPFlag("cache_dir", appCmd.PersistentFlags().Lookup("cache-dir"))
		_ = viper.BindPFlag("debug", appCmd.PersistentFlags().Lookup("debug"))
		viper.SetDefault("default_repository", defaultRepositoryURL)

		// Viper - Env
		viper.AutomaticEnv()
		viper.SetEnvPrefix("MANALA")

		// App cache
		appCache := cache.New(viper.GetString("cache_dir")).
			WithUserDir("manala")

		// Deferred app log instantiation
		appLogHandler := log.NewSlogHandler(ui,
			log.WithSlogHandlerDebug(viper.GetBool("debug")),
		)
		*appLog = *slog.New(appLogHandler)

		// Deferred app api instantiation
		*appAPI = *api.New(appLog, appCache,
			api.WithDefaultRepositoryURL(viper.GetString("default_repository")),
		)

		// Log config
		appLog.Debug("config",
			"default_repository", viper.GetString("default_repository"),
			"cache_dir", viper.GetString("cache_dir"),
			"debug", viper.GetBool("debug"),
		)
	})

	// Execute app command
	if err := appCmd.Execute(); err != nil {
		ui.Error(err)
		os.Exit(1)
	}
}
