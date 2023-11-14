package repository

import (
	"context"
	"github.com/hashicorp/go-getter/v2"
	"log/slog"
	"manala/internal/schema"
	"time"
)

func NewGitGetter(log *slog.Logger, result *GetterResult) *GitGetter {
	return &GitGetter{
		log:    log.With("getter", "git"),
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
	log    *slog.Logger
	result *GetterResult
	*getter.GitGetter
	protocol string
}

func (g *GitGetter) Detect(req *getter.Request) (bool, error) {
	// Log
	g.log.Debug("try to detect repository",
		"src", req.Src,
	)

	// Force git repo format (ensure backward compatibility)
	if (&schema.GitRepoFormatChecker{}).IsFormat(req.Src) {
		req.Forced = "git"
	}

	// Detect
	ok, err := g.GitGetter.Detect(req)

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

func (g *GitGetter) Get(ctx context.Context, req *getter.Request) error {
	// Log
	g.log.Debug("get repository",
		"src", req.Src,
	)

	// Get
	err := g.GitGetter.Get(ctx, req)

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
