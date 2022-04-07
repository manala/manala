package cmd

import (
	"github.com/spf13/cobra"
	"manala/app"
	internalConfig "manala/internal/config"
	internalLog "manala/internal/log"
	"path/filepath"
)

func newWatchCmd(config *internalConfig.Config, logger *internalLog.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "watch [path]",
		Args:              cobra.MaximumNArgs(1),
		DisableAutoGenTag: true,
		Short:             "List recipes",
		Long: `Watch (manala watch) will watch project, and launch update on changes.

Example: manala watch -> resulting in a watch in a path (default to the current directory)`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// App
			manala := app.New(config, logger)

			// Get flags
			repositoryPath, _ := cmd.Flags().GetString("repository")
			recipeName, _ := cmd.Flags().GetString("recipe")
			all, _ := cmd.Flags().GetBool("all")
			notify, _ := cmd.Flags().GetBool("notify")

			// Load project
			project, err := manala.ProjectFrom(
				filepath.Clean(append(args, "")[0]),
				repositoryPath,
				recipeName,
			)
			if err != nil {
				return err
			}

			// Sync project
			if err := manala.SyncProject(project); err != nil {
				return err
			}

			// Watch project
			return manala.WatchProject(
				project,
				repositoryPath,
				recipeName,
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
