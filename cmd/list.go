package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"manala/app"
	"manala/internal"
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
			// App
			manala := app.New(config, logger)

			// Get flags
			repositoryPath, _ := cmd.Flags().GetString("repository")

			// Load repository
			repository, err := manala.Repository(
				repositoryPath,
			)
			if err != nil {
				return err
			}

			// Walk repository recipes
			return repository.WalkRecipes(func(recipe *internal.Recipe) {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s: %s\n", recipe.Name(), recipe.Description())
			})
		},
	}

	// Flags
	cmd.Flags().StringP("repository", "o", "", "use repository")

	return cmd
}
