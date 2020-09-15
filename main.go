package main

import (
	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"manala/cmd"
	"os"
)

// Default repository
var repository = "https://github.com/manala/manala-recipes.git"

// Set at build time, by goreleaser, via ldflags
var version = "dev"

func main() {
	// Log handler
	log.SetHandler(cli.Default)

	// Config
	viper.SetEnvPrefix("manala")
	viper.AutomaticEnv()

	viper.SetDefault("repository", repository)
	viper.SetDefault("debug", false)

	cacheDir, err := os.UserCacheDir()
	if err != nil {
		log.WithError(err).Fatal("Error getting cache dir")
	}
	viper.SetDefault("cache_dir", cacheDir)

	// Commands
	rootCmd := cmd.RootCmd(version)
	rootCmd.AddCommand(cmd.InitCmd())
	rootCmd.AddCommand(cmd.ListCmd())
	rootCmd.AddCommand(cmd.UpdateCmd())
	rootCmd.AddCommand(cmd.WatchCmd())

	// Log handler
	log.SetHandler(cli.New(rootCmd.ErrOrStderr()))

	cobra.OnInitialize(func() {
		// Debug
		if viper.GetBool("debug") {
			log.SetLevel(log.DebugLevel)
		}
	})

	if err = rootCmd.Execute(); err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}
