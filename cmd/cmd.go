package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"log/slog"
	"manala/app/config"
	"manala/internal/ui/adapters/charm"
	"manala/internal/ui/log"
	"strings"
)

func Execute(version string, conf *config.Config, stdin io.Reader, stdout io.Writer, stderr io.Writer) int {
	// Define root command
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

	// Set standard streams
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)

	// Set persistent flags
	cmd.PersistentFlags().StringP("cache-dir", "c", conf.CacheDir, "use cache directory")
	cmd.PersistentFlags().BoolP("debug", "d", conf.Debug, "set debug mode")

	// Ui Adapter
	ui := charm.New(stdin, stdout, stderr)

	// Viper
	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvPrefix("MANALA")
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	// Bind command persistent flags
	if err := v.BindPFlags(cmd.PersistentFlags()); err != nil {
		ui.Error(err)
		return 1
	}

	// Logger
	loggerHandler := log.NewSlogHandler(ui)
	logger := slog.New(loggerHandler)

	cobra.OnInitialize(func() {
		// Unmarshall config
		_ = v.Unmarshal(&conf)

		// Debug
		if conf.Debug {
			loggerHandler.Level(slog.LevelDebug)
		}

		// Log config
		logger.Debug("config",
			"repository", conf.Repository,
			"cache-dir", conf.CacheDir,
			"debug", conf.Debug,
		)
	})

	// Sub commands
	cmd.AddCommand(
		newInitCmd(conf, logger, ui, ui),
		newListCmd(conf, logger, ui),
		newMascotCmd(ui),
		newUpdateCmd(conf, logger, ui),
		newWatchCmd(conf, logger, ui),
	)

	// Docs generation command
	if version == "dev" {
		cmd.AddCommand(newDocsCmd(cmd))
	}

	// Execute
	if err := cmd.Execute(); err != nil {
		ui.Error(err)
		return 1
	}

	return 0
}
