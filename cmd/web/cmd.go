package web

import (
	"context"
	"log/slog"
	"manala/app"
	"manala/app/api"
	"manala/internal/ui"
	"manala/web"
	"net/url"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/charmbracelet/x/exp/open"
	"github.com/spf13/cobra"
)

func NewCmd(log *slog.Logger, api *api.API, out ui.Output) *cobra.Command {
	// Flags
	var (
		repositoryURL, repositoryRef string
		address                      string
		noBrowser                    bool
	)

	// Command
	cmd := &cobra.Command{
		Use:               "web",
		Aliases:           []string{"w"},
		Args:              cobra.MaximumNArgs(1),
		DisableAutoGenTag: true,
		Short:             "Web interface",
		Long: `Web (manala web) will launch web interface.

Example: manala web -> resulting in a web interface launch`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Args
			dir := filepath.Clean(append(args, "")[0])

			// Context
			ctx := cmd.Context()
			ctx = app.WithRepositoryURL(ctx, repositoryURL)
			ctx = app.WithRepositoryRef(ctx, repositoryRef)

			return run(ctx, log, api, out, dir, address, !noBrowser)
		},
	}

	// Set flags
	cmd.Flags().StringVarP(&repositoryURL, "repository", "o", "", "use repository")
	cmd.Flags().StringVar(&repositoryRef, "ref", "", "use repository ref")

	cmd.Flags().StringVar(&address, "address", "127.0.0.1:9400", "address")
	cmd.Flags().BoolVar(&noBrowser, "no-browser", false, "no browser")

	return cmd
}

func run(ctx context.Context, log *slog.Logger, api *api.API, out ui.Output, dir, address string, browser bool) error {
	if browser {
		// Handle context as query values
		query := url.Values{}
		if url, ok := app.RepositoryURL(ctx); ok {
			query.Add("repository", url)
		}

		if ref, ok := app.RepositoryRef(ctx); ok {
			query.Add("ref", ref)
		}

		// Compose url
		url := strings.Builder{}
		url.WriteString("http://" + address)

		if len(query) > 0 {
			url.WriteString("?" + query.Encode())
		}

		// Open url in browser
		go func() {
			<-time.After(100 * time.Millisecond)

			if err := open.Open(context.Background(), url.String()); err != nil {
				log.Warn("unable to open browser", "error", err)
			}
		}()
	}

	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Start web server
	log.Info("starting web server…", "address", address)

	return web.NewServer(log, api, out, dir).Serve(ctx, address)
}
