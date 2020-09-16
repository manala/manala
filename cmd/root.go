package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// RootCmd represents the base command when called without any subcommands
func RootCmd(version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "manala",
		Short: "Let your project's plumbing up to date",
		Long: `Manala synchronize some boring parts of your projects,
such as makefile targets, virtualization and provisioning files...

Recipes are pulled from a git repository, or a local directory.`,
		SilenceErrors:     true,
		SilenceUsage:      true,
		Version:           version,
		DisableAutoGenTag: true,
	}

	cmd.PersistentFlags().StringP("cache-dir", "c", viper.GetString("cache_dir"), "cache directory")
	_ = viper.BindPFlag("cache_dir", cmd.PersistentFlags().Lookup("cache-dir"))

	cmd.PersistentFlags().BoolP("debug", "d", viper.GetBool("debug"), "debug mode")
	_ = viper.BindPFlag("debug", cmd.PersistentFlags().Lookup("debug"))

	return cmd
}

func addRepositoryFlag(cmd *cobra.Command, usage string) {
	cmd.Flags().StringP("repository", "o", "", usage)
}

func addRecipeFlag(cmd *cobra.Command, usage string) {
	cmd.Flags().StringP("recipe", "i", "", usage)
}
