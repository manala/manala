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

func newListCmd(config *internalConfig.Config, logger *internalLog.Logger) *cobra.Command {
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
			// Application
			app := application.NewApplication(config, logger)

			// Get flags
			repoPath, _ := cmd.Flags().GetString("repository")

			// Load repository
			repo, err := app.Repository(repoPath)
			if err != nil {
				return err
			}

			var nameStyle = lipgloss.NewStyle().Bold(true)
			var descriptionStyle = lipgloss.NewStyle().Italic(true)

			// Walk into repository recipes
			return repo.WalkRecipes(func(rec core.Recipe) {
				_, _ = fmt.Fprintf(
					cmd.OutOrStdout(),
					"%s: %s\n",
					nameStyle.Render(rec.Name()),
					descriptionStyle.Render(rec.Description()),
				)
			})
		},
	}

	// Flags
	cmd.Flags().StringP("repository", "o", "", "use repository")

	return cmd
}
