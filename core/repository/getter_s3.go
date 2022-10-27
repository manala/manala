package repository

import (
	"context"
	"github.com/hashicorp/go-getter/s3/v2"
	"github.com/hashicorp/go-getter/v2"
	internalLog "manala/internal/log"
	"time"
)

func NewS3Getter(log *internalLog.Logger, result *GetterResult) *S3Getter {
	return &S3Getter{
		log:    log,
		result: result,
		Getter: &s3.Getter{
			Timeout: 30 * time.Second,
		},
		protocol: "s3",
	}
}

type S3Getter struct {
	log    *internalLog.Logger
	result *GetterResult
	*s3.Getter
	protocol string
}

func (g *S3Getter) Detect(req *getter.Request) (bool, error) {
	g.log.
		WithField("protocol", "s3").
		Debug("detect")

	// Detect
	ok, err := g.Getter.Detect(req)

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

func (g *S3Getter) Get(ctx context.Context, req *getter.Request) error {
	g.log.
		WithField("protocol", g.protocol).
		Debug("get")

	// Get
	err := g.Getter.Get(ctx, req)

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
