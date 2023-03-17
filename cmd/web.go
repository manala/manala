package cmd

import (
	"github.com/spf13/cobra"
	"manala/app/interfaces"
	"manala/core/application"
	internalLog "manala/internal/log"
	"manala/web"
	"path/filepath"
)

func newWebCmd(conf interfaces.Config, log *internalLog.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "web",
		Aliases:           []string{"w"},
		Args:              cobra.MaximumNArgs(1),
		DisableAutoGenTag: true,
		Short:             "Web interface",
		Long: `Web (manala web) will launch web interface.

Example: manala web -> resulting in a web interface launch`,
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
				appOptions...,
			)

			// Get args
			dir := filepath.Clean(append(args, "")[0])

			// Web
			return web.New(log, app, conf, dir).ListenAndServe()
		},
	}

	// Flags
	cmd.Flags().StringP("repository", "o", "", "use repository")
	cmd.Flags().String("ref", "", "use repository ref")

	cmd.Flags().Int("port", conf.WebPort(), "server port")
	conf.BindWebPortFlag(cmd.Flags().Lookup("port"))

	return cmd
}
