package cmd

import (
	"github.com/spf13/cobra"
	"manala/app"
	"manala/internal"
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
			// App
			manala := app.New(config, logger)

			// Get flags
			repositoryPath, _ := cmd.Flags().GetString("repository")
			recipeName, _ := cmd.Flags().GetString("recipe")
			recursive, _ := cmd.Flags().GetBool("recursive")

			// Get args
			path := filepath.Clean(append(args, "")[0])

			if recursive {
				// Recursively load projects
				return manala.WalkProjects(
					path,
					repositoryPath,
					recipeName,
					func(project *internal.Project) error {
						// Sync project
						return manala.SyncProject(project)
					},
				)
			} else {
				// Load project
				project, err := manala.ProjectFrom(
					path,
					repositoryPath,
					recipeName,
				)
				if err != nil {
					return err
				}

				// Sync project
				return manala.SyncProject(project)
			}
		},
	}

	// Flags
	cmd.Flags().StringP("repository", "o", "", "use repository")
	cmd.Flags().StringP("recipe", "i", "", "use recipe")
	cmd.Flags().BoolP("recursive", "r", false, "set recursive mode")

	return cmd
}
