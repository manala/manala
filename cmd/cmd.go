package cmd

import (
	"github.com/spf13/cobra"
	"io"
	internalConfig "manala/internal/config"
	internalLog "manala/internal/log"
	internalReport "manala/internal/report"
	"os"
	"path/filepath"
	"strings"
)

func newCmd(version string) *cobra.Command {
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
	cmd.PersistentFlags().StringP("cache-dir", "c", "", "use cache directory")
	cmd.PersistentFlags().BoolP("debug", "d", false, "set debug mode")

	return cmd
}

func Execute(version string, defaultRepository string, stdout io.Writer, stderr io.Writer) {

	// Logger
	logger := internalLog.New(stderr)

	// Config
	config := internalConfig.New()
	config.SetEnvPrefix("MANALA")
	config.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	config.AutomaticEnv()
	config.SetDefault("debug", false)
	config.Set("default-repository", defaultRepository)

	// Cache dir
	if dir, err := os.UserCacheDir(); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	} else {
		config.SetDefault("cache-dir", filepath.Join(dir, "manala"))
	}

	// Debug
	cobra.OnInitialize(func() {
		if config.GetBool("debug") {
			logger.LevelDebug()
		}
	})

	// Root command
	cmd := newCmd(version)
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)

	// Bind config to persistent flags
	_ = config.BindPFlags(cmd.PersistentFlags())

	// Sub commands
	cmd.AddCommand(
		newInitCmd(config, logger),
		newListCmd(config, logger),
		newMascotCmd(),
		newUpdateCmd(config, logger),
		newWatchCmd(config, logger),
	)

	// Docs generation command
	if version == "dev" {
		cmd.AddCommand(newDocsCmd(cmd))
	}

	// Execute
	if err := cmd.Execute(); err != nil {
		report := internalReport.NewErrorReport(err)
		logger.Report(report)
		os.Exit(1)
	}
}
