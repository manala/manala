package cmd

import (
	"github.com/spf13/cobra"
	"log/slog"
	"manala/app/api"
	"manala/app/config"
	"manala/internal/ui"
	"path/filepath"
)

func newWatchCmd(config config.Config, log *slog.Logger, out ui.Output) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "watch [dir]",
		Args:              cobra.MaximumNArgs(1),
		DisableAutoGenTag: true,
		Short:             "Watch project",
		Long: `Watch (manala watch) will watch project files, and launch update on changes.

Example: manala watch -> resulting in a watch in a project dir (default to the current directory)`,
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

			// Flag - All
			all, _ := cmd.Flags().GetBool("all")

			// Flag - Notify
			notify, _ := cmd.Flags().GetBool("notify")

			// Api
			app := api.New(
				config,
				log,
				out,
				apiOptions...,
			)

			// Load project
			project, err := app.LoadProjectFrom(
				filepath.Clean(append(args, "")[0]),
			)
			if err != nil {
				return err
			}

			// Sync project
			if err := app.SyncProject(project); err != nil {
				return err
			}

			// Watch project
			return app.WatchProject(
				project,
				all,
				notify,
			)
		},
	}

	// Flags
	cmd.Flags().StringP("repository", "o", "", "use repository")
	cmd.Flags().String("ref", "", "use repository ref")
	cmd.Flags().StringP("recipe", "i", "", "use recipe")
	cmd.Flags().BoolP("all", "a", false, "watch recipe too")
	cmd.Flags().BoolP("notify", "n", false, "use system notifications")

	return cmd
}
