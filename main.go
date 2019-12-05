package main

import (
	"github.com/OpenPeeDeeP/xdg"
	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"manala/cmd"
	"os"
	"path"
)

// Default repository
var repository = "https://github.com/manala/manala-recipes.git"

// Set at build time, by goreleaser, via ldflags
var version = "dev"

func main() {
	// Set log handler
	log.SetHandler(cli.Default)

	// Config
	viper.SetEnvPrefix("manala")
	viper.AutomaticEnv()

	dir, err := os.Getwd()
	if err != nil {
		log.WithError(err).Fatal("Error getting dir")
	}

	viper.SetDefault("dir", dir)
	viper.SetDefault("cache_dir", path.Join(xdg.CacheHome(), "manala"))
	viper.SetDefault("repository", repository)
	viper.SetDefault("debug", false)

	// Commands
	rootCmd := cmd.RootCmd(version)
	rootCmd.AddCommand(cmd.InitCmd())
	rootCmd.AddCommand(cmd.ListCmd())
	rootCmd.AddCommand(cmd.UpdateCmd())
	rootCmd.AddCommand(cmd.WatchCmd())

	cobra.OnInitialize(func() {
		// Debug
		if viper.GetBool("debug") {
			log.SetLevel(log.DebugLevel)
		}
	})

	if err = rootCmd.Execute(); err != nil {
		log.WithError(err).Fatal("Error executing command")
	}
}
