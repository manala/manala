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
			// Get flags
			repoUrl, _ := cmd.Flags().GetString("repository")

			// Application
			app := application.NewApplication(
				config,
				log,
				application.WithRepositoryUrl(repoUrl),
			)

			// Styles
			var nameStyle = lipgloss.NewStyle().Bold(true)
			var descriptionStyle = lipgloss.NewStyle().Italic(true)

			// Walk into recipes
			return app.WalkRecipes(func(rec core.Recipe) error {
				if _, err := fmt.Fprintf(
					cmd.OutOrStdout(),
					"%s: %s\n",
					nameStyle.Render(rec.Name()),
					descriptionStyle.Render(rec.Description()),
				); err != nil {
					return err
				}

				return nil
			})
		},
	}

	// Flags
	cmd.Flags().StringP("repository", "o", "", "use repository")

	return cmd
}
