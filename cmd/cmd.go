package cmd

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"io"
	"manala/app/config"
	"manala/app/interfaces"
	internalLog "manala/internal/log"
	internalReport "manala/internal/report"
	"os"
)

// Styles
var styles = struct {
	Primary, Secondary lipgloss.Style
}{
	Primary: lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#1A1A1A", Dark: "#DDDDDD"}).
		Bold(true),
	Secondary: lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#1A1A1A", Dark: "#DDDDDD"}).
		Italic(true),
}

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

	// Log
	log := internalLog.New(stderr)

	// Config
	conf := config.New()

	// Debug
	cobra.OnInitialize(func() {
		if conf.Debug() {
			log.LevelDebug()
		}
	})

	// Root command
	cmd := newCmd(version, conf)
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)

	// Sub commands
	cmd.AddCommand(
		newInitCmd(conf, log),
		newListCmd(conf, log),
		newMascotCmd(),
		newUpdateCmd(conf, log),
		newWatchCmd(conf, log),
		newWebCmd(conf, log),
	)

	// Docs generation command
	if version == "dev" {
		cmd.AddCommand(newDocsCmd(cmd))
	}

	// Execute
	if err := cmd.Execute(); err != nil {
		report := internalReport.NewErrorReport(err)
		log.Report(report)
		os.Exit(1)
	}
}
