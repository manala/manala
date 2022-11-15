package cmd

import (
	"github.com/spf13/cobra"
	"manala/core/application"
	internalConfig "manala/internal/config"
	internalLog "manala/internal/log"
	"path/filepath"
)

func newWatchCmd(config *internalConfig.Config, log *internalLog.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "watch [dir]",
		Args:              cobra.MaximumNArgs(1),
		DisableAutoGenTag: true,
		Short:             "Watch project",
		Long: `Watch (manala watch) will watch project files, and launch update on changes.

Example: manala watch -> resulting in a watch in a project dir (default to the current directory)`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get flags
			repoUrl, _ := cmd.Flags().GetString("repository")
			recName, _ := cmd.Flags().GetString("recipe")
			all, _ := cmd.Flags().GetBool("all")
			notify, _ := cmd.Flags().GetBool("notify")

			// Application
			app := application.NewApplication(
				config,
				log,
				application.WithRepositoryUrl(repoUrl),
				application.WithRecipeName(recName),
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
	cmd.Flags().StringP("recipe", "i", "", "use recipe")
	cmd.Flags().BoolP("all", "a", false, "watch recipe too")
	cmd.Flags().BoolP("notify", "n", false, "use system notifications")

	return cmd
}
