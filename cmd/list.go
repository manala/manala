package cmd

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"manala/core"
	"manala/core/application"
	internalConfig "manala/internal/config"
	internalLog "manala/internal/log"
)

func newListCmd(config *internalConfig.Config, log *internalLog.Logger) *cobra.Command {
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

			// Application
			app := application.NewApplication(
				config,
				log,
				appOptions...,
			)

			var recs []core.Recipe
			maxNameWidth := 0

			// Walk into recipes
			if err := app.WalkRecipes(func(rec core.Recipe) error {
				recs = append(recs, rec)

				nameWidth := lipgloss.Width(rec.Name())
				if nameWidth > maxNameWidth {
					maxNameWidth = nameWidth
				}

				return nil
			}); err != nil {
				return err
			}

			nameStyle := styles.Primary.Copy().
				Width(maxNameWidth).
				MarginRight(2)
			descriptionStyle := styles.Secondary.Copy()

			for _, rec := range recs {
				if _, err := fmt.Fprintf(
					cmd.OutOrStdout(),
					"%s%s\n",
					nameStyle.Render(rec.Name()),
					descriptionStyle.Render(rec.Description()),
				); err != nil {
					return err
				}
			}

			return nil
		},
	}

	// Flags
	cmd.Flags().StringP("repository", "o", "", "use repository")

	return cmd
}
