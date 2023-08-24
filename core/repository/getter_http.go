package repository

import (
	"context"
	"github.com/hashicorp/go-getter/v2"
	"log/slog"
	"time"
)

func NewHttpGetter(log *slog.Logger, result *GetterResult) *HttpGetter {
	return &HttpGetter{
		log:    log.With("getter", "http"),
		result: result,
		HttpGetter: &getter.HttpGetter{
			// Will lookup and use auth information found in the user's netrc file if available
			Netrc: true,
			// Disables the client's usage of the "X-Terraform-Get" header value
			XTerraformGetDisabled: true,
			// Enforce a timeout when the server supports HEAD requests
			HeadFirstTimeout: 10 * time.Second,
			// Enforce a timeout when making a request to an HTTP server and reading its response body
			ReadTimeout: 30 * time.Second,
		},
		protocol: "http",
	}
}

type HttpGetter struct {
	log    *slog.Logger
	result *GetterResult
	*getter.HttpGetter
	protocol string
}

func (g *HttpGetter) Detect(req *getter.Request) (bool, error) {
	// Log
	g.log.Debug("try to detect repository",
		"src", req.Src,
	)

	// Detect
	ok, err := g.HttpGetter.Detect(req)

	if err != nil {
		// Log
		g.log.Debug("unable to detect repository",
			"error", err,
		)

		g.result.SetDetectError(err, g.protocol)

		return ok, err
	}

	return ok, nil
}

func (g *HttpGetter) Get(ctx context.Context, req *getter.Request) error {
	// Log
	g.log.Debug("get repository",
		"src", req.Src,
	)

	// Get
	err := g.HttpGetter.Get(ctx, req)

	if err != nil {
		// Log
		g.log.Debug("unable to get repository",
			"error", err,
		)

		g.result.AddGetError(err, g.protocol)

		return err
	}

	return nil
}
