package cmd

import (
	"github.com/spf13/cobra"
	"io"
	"log/slog"
	"manala/app/config"
	"manala/app/interfaces"
	"manala/internal/ui/log"
	"manala/internal/ui/output/lipgloss"
	"os"
)

func newCmd(version string, conf interfaces.Config) *cobra.Command {
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
	cmd.PersistentFlags().StringP("cache-dir", "c", conf.CacheDir(), "use cache directory")
	conf.BindCacheDirFlag(cmd.PersistentFlags().Lookup("cache-dir"))

	cmd.PersistentFlags().BoolP("debug", "d", conf.Debug(), "set debug mode")
	conf.BindDebugFlag(cmd.PersistentFlags().Lookup("debug"))

	return cmd
}

func Execute(version string, stdout io.Writer, stderr io.Writer) {
	// Config
	conf := config.New()

	// Ui Output
	out := lipgloss.New(stdout, stderr)

	// Log handler
	logHandler := log.NewSlogHandler(out)

	// Debug
	cobra.OnInitialize(func() {
		if conf.Debug() {
			logHandler.LevelDebug()
		}
	})

	// Logger
	logger := slog.New(logHandler)

	// Root command
	cmd := newCmd(version, conf)
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)

	// Sub commands
	cmd.AddCommand(
		newInitCmd(conf, logger, out),
		newListCmd(conf, logger, out),
		newMascotCmd(),
		newUpdateCmd(conf, logger, out),
		newWatchCmd(conf, logger, out),
	)

	// Docs generation command
	if version == "dev" {
		cmd.AddCommand(newDocsCmd(cmd))
	}

	// Execute
	if err := cmd.Execute(); err != nil {
		out.Error(err)
		os.Exit(1)
	}
}
