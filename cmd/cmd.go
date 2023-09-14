package cmd

import (
	"github.com/spf13/cobra"
	"io"
	"log/slog"
	"manala/app/config"
	"manala/internal/ui/adapters/charm"
	"manala/internal/ui/log"
	"os"
)

func newCmd(version string, config config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "manala",
		Version:           version,
		DisableAutoGenTag: true,
		SilenceErrors:     true,
		SilenceUsage:      true,
		Short:             "Let your project's plumbing up to date",
		Long: `Manala synchronize some boring parts of your projects,
such as makefile targets, virtualization and provisioning files...

Recipes are pulled from a git repository, or a local directory.`,
	}

	// Persistent flags
	cmd.PersistentFlags().StringP("cache-dir", "c", config.CacheDir(), "use cache directory")
	config.BindCacheDirFlag(cmd.PersistentFlags().Lookup("cache-dir"))

	cmd.PersistentFlags().BoolP("debug", "d", config.Debug(), "set debug mode")
	config.BindDebugFlag(cmd.PersistentFlags().Lookup("debug"))

	return cmd
}

func Execute(version string, stdin io.Reader, stdout io.Writer, stderr io.Writer) {
	// Config
	config := config.NewViperConfig()

	// Ui Adapter
	ui := charm.New(stdin, stdout, stderr)

	// Log handler
	logHandler := log.NewSlogHandler(ui)

	// Debug
	cobra.OnInitialize(func() {
		if config.Debug() {
			logHandler.LevelDebug()
		}
	})

	// Logger
	logger := slog.New(logHandler)

	// Root command
	cmd := newCmd(version, config)
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)

	// Sub commands
	cmd.AddCommand(
		newInitCmd(config, logger, ui, ui),
		newListCmd(config, logger, ui),
		newMascotCmd(ui),
		newUpdateCmd(config, logger, ui),
		newWatchCmd(config, logger, ui),
	)

	// Docs generation command
	if version == "dev" {
		cmd.AddCommand(newDocsCmd(cmd))
	}

	// Execute
	if err := cmd.Execute(); err != nil {
		ui.Error(err)
		os.Exit(1)
	}
}
