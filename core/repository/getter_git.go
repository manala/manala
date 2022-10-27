package repository

import (
	"context"
	"github.com/hashicorp/go-getter/v2"
	"manala/core"
	internalLog "manala/internal/log"
	"time"
)

func NewGitGetter(log *internalLog.Logger, result *GetterResult) *GitGetter {
	return &GitGetter{
		log:    log,
		result: result,
		GitGetter: &getter.GitGetter{
			Detectors: []getter.Detector{
				&getter.GitHubDetector{},
				&getter.GitDetector{},
				&getter.BitBucketDetector{},
				&getter.GitLabDetector{},
			},
			Timeout: 30 * time.Second,
		},
		protocol: "git",
	}
}

type GitGetter struct {
	log    *internalLog.Logger
	result *GetterResult
	*getter.GitGetter
	protocol string
}

func (g *GitGetter) Detect(req *getter.Request) (bool, error) {
	g.log.
		WithField("protocol", g.protocol).
		Debug("detect")

	// Force git repo format (ensure backward compatibility)
	if (&core.GitRepoFormatChecker{}).IsFormat(req.Src) {
		req.Forced = "git"
	}

	// Detect
	ok, err := g.GitGetter.Detect(req)

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

func (g *GitGetter) Get(ctx context.Context, req *getter.Request) error {
	g.log.
		WithField("protocol", g.protocol).
		Debug("get")

	// Get
	err := g.GitGetter.Get(ctx, req)

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
