package cmd

import (
	"github.com/spf13/cobra"
	"manala/core"
	"manala/core/application"
	internalConfig "manala/internal/config"
	internalLog "manala/internal/log"
	"path/filepath"
)

func newUpdateCmd(config *internalConfig.Config, logger *internalLog.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "update [path]",
		Aliases:           []string{"up"},
		Args:              cobra.MaximumNArgs(1),
		DisableAutoGenTag: true,
		Short:             "Synchronize project(s)",
		Long: `Update (manala update) will synchronize project(s), based on
repository's recipe and related variables defined in manifest (.manala.yaml).

Example: manala update -> resulting in an update in a path (default to the current directory)`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Application
			app := application.NewApplication(config, logger)

			// Get flags
			repoPath, _ := cmd.Flags().GetString("repository")
			recName, _ := cmd.Flags().GetString("recipe")
			recursive, _ := cmd.Flags().GetBool("recursive")

			// Get args
			path := filepath.Clean(append(args, "")[0])

			if recursive {
				// Recursively load projects
				return app.WalkProjects(
					path,
					repoPath,
					recName,
					func(proj core.Project) error {
						// Sync project
						return app.SyncProject(proj)
					},
				)
			} else {
				// Load project
				proj, err := app.ProjectFrom(
					path,
					repoPath,
					recName,
				)
				if err != nil {
					return err
				}

				// Sync project
				return app.SyncProject(proj)
			}
		},
	}

	// Flags
	cmd.Flags().StringP("repository", "o", "", "use repository")
	cmd.Flags().StringP("recipe", "i", "", "use recipe")
	cmd.Flags().BoolP("recursive", "r", false, "set recursive mode")

	return cmd
}
