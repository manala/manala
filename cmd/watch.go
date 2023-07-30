package cmd

import (
	"github.com/spf13/cobra"
	"log/slog"
	"manala/app/interfaces"
	"manala/core/application"
	"manala/internal/ui/output"
	"path/filepath"
)

func newWatchCmd(conf interfaces.Config, log *slog.Logger, out output.Output) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "watch [dir]",
		Args:              cobra.MaximumNArgs(1),
		DisableAutoGenTag: true,
		Short:             "Watch project",
		Long: `Watch (manala watch) will watch project files, and launch update on changes.

Example: manala watch -> resulting in a watch in a project dir (default to the current directory)`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Application options
			var appOptions []application.Option

			// Flag - Repository url
			if cmd.Flags().Changed("repository") {
				repoUrl, _ := cmd.Flags().GetString("repository")
				appOptions = append(appOptions, application.WithRepositoryUrl(repoUrl))
			}

			// Flag - Repository ref
			if cmd.Flags().Changed("ref") {
				repoRef, _ := cmd.Flags().GetString("ref")
				appOptions = append(appOptions, application.WithRepositoryRef(repoRef))
			}

			// Flag - Recipe name
			if cmd.Flags().Changed("recipe") {
				recName, _ := cmd.Flags().GetString("recipe")
				appOptions = append(appOptions, application.WithRecipeName(recName))
			}

			// Flag - All
			all, _ := cmd.Flags().GetBool("all")

			// Flag - Notify
			notify, _ := cmd.Flags().GetBool("notify")

			// Application
			app := application.NewApplication(
				conf,
				log,
				out,
				appOptions...,
			)

			// Load project
			proj, err := app.LoadProjectFrom(
				filepath.Clean(append(args, "")[0]),
			)
			if err != nil {
				return err
			}

			// Sync project
			if err := app.SyncProject(proj); err != nil {
				return err
			}

			// Watch project
			return app.WatchProject(
				proj,
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
