package main

import (
	"log/slog"
	"github.com/manala/manala/app/api"
	"github.com/manala/manala/cmd"
	cmdDocs "github.com/manala/manala/cmd/docs"
	cmdInit "github.com/manala/manala/cmd/init"
	cmdList "github.com/manala/manala/cmd/list"
	cmdMascot "github.com/manala/manala/cmd/mascot"
	cmdUpdate "github.com/manala/manala/cmd/update"
	cmdWatch "github.com/manala/manala/cmd/watch"
	"github.com/manala/manala/internal/caching"
	"github.com/manala/manala/internal/notifier"
	"github.com/manala/manala/internal/ui/adapters/charm"
	"github.com/manala/manala/internal/ui/log"
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
		v := viper.New()

		_ = v.BindPFlag("cache_dir", appCmd.PersistentFlags().Lookup("cache-dir"))
		_ = v.BindPFlag("debug", appCmd.PersistentFlags().Lookup("debug"))
		v.SetDefault("default_repository", defaultRepositoryURL)

		// Viper - Env
		v.AutomaticEnv()
		v.SetEnvPrefix("MANALA")

		// App cache
		appCache := caching.NewCache(v.GetString("cache_dir")).
			WithUserDir("manala")

		// Deferred app log instantiation
		appLogHandler := log.NewSlogHandler(ui,
			log.WithSlogHandlerDebug(v.GetBool("debug")),
		)
		*appLog = *slog.New(appLogHandler)

		// Deferred app api instantiation
		*appAPI = *api.New(appLog, appCache,
			api.WithDefaultRepositoryURL(v.GetString("default_repository")),
		)

		// Log config
		appLog.Debug("config",
			"default_repository", v.GetString("default_repository"),
			"cache_dir", v.GetString("cache_dir"),
			"debug", v.GetBool("debug"),
		)
	})

	// Execute app command
	if err := appCmd.Execute(); err != nil {
		ui.Error(err)
		os.Exit(1)
	}
}
