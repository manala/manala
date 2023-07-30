package repository

import (
	"context"
	"github.com/hashicorp/go-getter/s3/v2"
	"github.com/hashicorp/go-getter/v2"
	"log/slog"
	"time"
)

func NewS3Getter(log *slog.Logger, result *GetterResult) *S3Getter {
	return &S3Getter{
		log:    log.With("getter", "s3"),
		result: result,
		Getter: &s3.Getter{
			Timeout: 30 * time.Second,
		},
		protocol: "s3",
	}
}

type S3Getter struct {
	log    *slog.Logger
	result *GetterResult
	*s3.Getter
	protocol string
}

func (g *S3Getter) Detect(req *getter.Request) (bool, error) {
	// Log
	g.log.Debug("try to detect repository",
		"src", req.Src,
	)

	// Detect
	ok, err := g.Getter.Detect(req)

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

func (g *S3Getter) Get(ctx context.Context, req *getter.Request) error {
	// Log
	g.log.Debug("get repository",
		"src", req.Src,
	)

	// Get
	err := g.Getter.Get(ctx, req)

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
