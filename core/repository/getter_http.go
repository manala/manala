package repository

import (
	"context"
	"github.com/hashicorp/go-getter/v2"
	internalLog "manala/internal/log"
	"time"
)

func NewHttpGetter(log *internalLog.Logger, result *GetterResult) *HttpGetter {
	return &HttpGetter{
		log:    log,
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
	log    *internalLog.Logger
	result *GetterResult
	*getter.HttpGetter
	protocol string
}

func (g *HttpGetter) Detect(req *getter.Request) (bool, error) {
	g.log.
		WithField("protocol", g.protocol).
		Debug("detect")

	// Detect
	ok, err := g.HttpGetter.Detect(req)

	if err != nil {
		g.log.
			WithField("protocol", g.protocol).
			WithError(err).
			Debug("unable to detect")

		g.result.SetDetectError(err, g.protocol)

		return ok, err
	}

	return ok, nil
}

func (g *HttpGetter) Get(ctx context.Context, req *getter.Request) error {
	g.log.
		WithField("protocol", g.protocol).
		Debug("get")

	// Get
	err := g.HttpGetter.Get(ctx, req)

	if err != nil {
		g.log.
			WithField("protocol", g.protocol).
			WithError(err).
			Debug("unable to get")

		g.result.AddGetError(err, g.protocol)

		return err
	}

	return nil
}
