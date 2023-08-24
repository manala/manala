package cmd

import (
	"github.com/spf13/cobra"
	"log/slog"
	"manala/app/interfaces"
	"manala/core/application"
	"manala/internal/ui/components"
	"manala/internal/ui/output"
)

func newListCmd(conf interfaces.Config, log *slog.Logger, out output.Output) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "list",
		Aliases:           []string{"ls"},
		Args:              cobra.NoArgs,
		DisableAutoGenTag: true,
		Short:             "List recipes",
		Long: `List (manala list) will list recipes available on
repository.

Example: manala list -> resulting in a recipes list display`,
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

			// Application
			app := application.NewApplication(
				conf,
				log,
				out,
				appOptions...,
			)

			table := &components.Table{}

			// Walk into recipes
			if err := app.WalkRecipes(func(rec interfaces.Recipe) error {
				table.AddRow(rec.Name(), rec.Description())
				return nil
			}); err != nil {
				return err
			}

			out.Table(table)

			return nil
		},
	}

	// Flags
	cmd.Flags().StringP("repository", "o", "", "use repository")
	cmd.Flags().String("ref", "", "use repository ref")

	return cmd
}
