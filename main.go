package main

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
)

// Set at build time, by goreleaser, via ldflags
var version = "dev"

// Default repository url
const defaultRepositoryUrl = "https://github.com/manala/manala-recipes.git"

func main() {
	// Standard streams
	stdIn := os.Stdin
	stdOut := os.Stdout
	stdErr := os.Stderr

	// User interface
	ui := charm.New(stdIn, stdOut, stdErr)

	// Notify
	notify := notifier.NewBeeep("Manala")

	var (
		appLog = new(slog.Logger)
		appApi = new(api.Api)
	)

	// App commands
	appCmd := cmd.NewCmd(version, stdOut, stdErr)
	appCmd.AddCommand(
		cmdInit.NewCmd(appLog, appApi, ui),
		cmdList.NewCmd(appLog, appApi, ui),
		cmdMascot.NewCmd(stdOut),
		cmdUpdate.NewCmd(appLog, appApi),
		cmdWatch.NewCmd(appLog, appApi, ui, notify),
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
		viper.SetDefault("default_repository", defaultRepositoryUrl)

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
		*appApi = *api.New(appLog, appCache,
			api.WithDefaultRepositoryUrl(viper.GetString("default_repository")),
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
