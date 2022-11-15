package cmd

import (
	"github.com/spf13/cobra"
	"manala/core"
	"manala/core/application"
	internalConfig "manala/internal/config"
	internalLog "manala/internal/log"
	"path/filepath"
)

func newUpdateCmd(config *internalConfig.Config, log *internalLog.Logger) *cobra.Command {
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

			// Get flags
			repoUrl, _ := cmd.Flags().GetString("repository")
			recName, _ := cmd.Flags().GetString("recipe")
			recursive, _ := cmd.Flags().GetBool("recursive")

			// Application
			app := application.NewApplication(
				config,
				log,
				application.WithRepositoryUrl(repoUrl),
				application.WithRecipeName(recName),
			)

			// Get args
			dir := filepath.Clean(append(args, "")[0])

			if recursive {
				// Recursively load projects
				return app.WalkProjects(
					dir,
					func(proj core.Project) error {
						// Sync project
						return app.SyncProject(proj)
					},
				)
			} else {
				// Load project
				proj, err := app.LoadProjectFrom(
					dir,
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
