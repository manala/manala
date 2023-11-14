package cmd

import (
	"github.com/spf13/cobra"
	"log/slog"
	"manala/app"
	"manala/app/api"
	"manala/app/config"
	"manala/internal/ui"
	"path/filepath"
)

func newUpdateCmd(config config.Config, log *slog.Logger, out ui.Output) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "update [dir]",
		Aliases:           []string{"up"},
		Args:              cobra.MaximumNArgs(1),
		DisableAutoGenTag: true,
		Short:             "Synchronize project(s)",
		Long: `Update (manala update) will synchronize project(s), based on
repository's recipe and related variables defined in manifest (.manala.yaml).

Example: manala update -> resulting in an update in a project dir (default to the current directory)`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Api options
			var apiOptions []api.Option

			// Flag - Repository url
			if cmd.Flags().Changed("repository") {
				repositoryUrl, _ := cmd.Flags().GetString("repository")
				apiOptions = append(apiOptions, api.WithRepositoryUrl(repositoryUrl))
			}

			// Flag - Repository ref
			if cmd.Flags().Changed("ref") {
				repositoryRef, _ := cmd.Flags().GetString("ref")
				apiOptions = append(apiOptions, api.WithRepositoryRef(repositoryRef))
			}

			// Flag - Recipe name
			if cmd.Flags().Changed("recipe") {
				recipeName, _ := cmd.Flags().GetString("recipe")
				apiOptions = append(apiOptions, api.WithRecipeName(recipeName))
			}

			// Flag - Recursive
			recursive, _ := cmd.Flags().GetBool("recursive")

			// Api
			api := api.New(
				config,
				log,
				out,
				apiOptions...,
			)

			// Get args
			dir := filepath.Clean(append(args, "")[0])

			if recursive {
				// Recursively load projects
				return api.WalkProjects(
					dir,
					func(project app.Project) error {
						// Sync project
						return api.SyncProject(project)
					},
				)
			} else {
				// Load project
				project, err := api.LoadProjectFrom(
					dir,
				)
				if err != nil {
					return err
				}

				// Sync project
				return api.SyncProject(project)
			}
		},
	}

	// Flags
	cmd.Flags().StringP("repository", "o", "", "use repository")
	cmd.Flags().String("ref", "", "use repository ref")
	cmd.Flags().StringP("recipe", "i", "", "use recipe")
	cmd.Flags().BoolP("recursive", "r", false, "set recursive mode")

	return cmd
}
