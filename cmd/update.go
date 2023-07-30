package cmd

import (
	"github.com/spf13/cobra"
	"log/slog"
	"manala/app/interfaces"
	"manala/core/application"
	"manala/internal/ui/output"
	"path/filepath"
)

func newUpdateCmd(conf interfaces.Config, log *slog.Logger, out output.Output) *cobra.Command {
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

			// Flag - Recursive
			recursive, _ := cmd.Flags().GetBool("recursive")

			// Application
			app := application.NewApplication(
				conf,
				log,
				out,
				appOptions...,
			)

			// Get args
			dir := filepath.Clean(append(args, "")[0])

			if recursive {
				// Recursively load projects
				return app.WalkProjects(
					dir,
					func(proj interfaces.Project) error {
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
	cmd.Flags().String("ref", "", "use repository ref")
	cmd.Flags().StringP("recipe", "i", "", "use recipe")
	cmd.Flags().BoolP("recursive", "r", false, "set recursive mode")

	return cmd
}
