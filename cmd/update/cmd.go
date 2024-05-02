package update

import (
	"github.com/spf13/cobra"
	"manala/app"
	"manala/app/api"
	"path/filepath"
)

func NewCmd(api *api.Api) *cobra.Command {
	// Flags
	var repositoryUrl, repositoryRef, recipeName string
	var recursive bool

	// Command
	cmd := &cobra.Command{
		Use:               "update [dir]",
		Aliases:           []string{"up"},
		Args:              cobra.MaximumNArgs(1),
		DisableAutoGenTag: true,
		Short:             "Synchronize project(s)",
		Long: `Update (manala update) will synchronize project(s), based on repository's recipe
and related variables defined in manifest (.manala.yaml).

Example: manala update -> resulting in an update in a project dir (default to the
current directory)`,
		RunE: func(_ *cobra.Command, args []string) error {
			// Args
			dir := filepath.Clean(append(args, "")[0])

			return run(api, dir, repositoryUrl, repositoryRef, recipeName, recursive)
		},
	}

	// Set flags
	cmd.Flags().StringVarP(&repositoryUrl, "repository", "o", "", "use repository")
	cmd.Flags().StringVar(&repositoryRef, "ref", "", "use repository ref")
	cmd.Flags().StringVarP(&recipeName, "recipe", "i", "", "use recipe")
	cmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "set recursive mode")

	return cmd
}

func run(api *api.Api, dir, repositoryUrl, repositoryRef, recipeName string, recursive bool) error {
	// Get repository loader
	repositoryLoader := api.NewRepositoryLoader(
		api.WithRepositoryLoaderUrl(repositoryUrl),
		api.WithRepositoryLoaderRef(repositoryRef),
	)

	// Get recipe loader
	recipeLoader := api.NewRecipeLoader(
		api.WithRecipeLoaderName(recipeName),
	)

	if recursive {
		// Get project loader
		projectLoader := api.NewProjectLoader(repositoryLoader, recipeLoader)

		// Recursively load projects
		return projectLoader.LoadRecursive(dir,
			func(project app.Project) error {
				// Sync project
				return api.NewProjectSyncer().Sync(project)
			},
		)
	} else {
		// Get project loader
		projectLoader := api.NewProjectLoader(repositoryLoader, recipeLoader,
			api.WithProjectLoaderFrom(true),
		)

		// Load project
		project, err := projectLoader.Load(dir)
		if err != nil {
			return err
		}

		// Sync project
		return api.NewProjectSyncer().Sync(project)
	}
}
