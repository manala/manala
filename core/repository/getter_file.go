package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/go-getter/v2"
	"log/slog"
	"manala/internal/errors/serrors"
	"os"
	"path/filepath"
)

func NewFileGetter(log *slog.Logger, result *GetterResult) *FileGetter {
	return &FileGetter{
		log:        log.With("getter", "file"),
		result:     result,
		FileGetter: &getter.FileGetter{},
		protocol:   "file",
	}
}

type FileGetter struct {
	log    *slog.Logger
	result *GetterResult
	*getter.FileGetter
	protocol string
}

func (g *FileGetter) Detect(req *getter.Request) (bool, error) {
	// Log
	g.log.Debug("try to detect repository",
		"src", req.Src,
	)

	// Stat
	stat, err := os.Stat(req.Src)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, serrors.WrapOs("file system error", err).
			WithArguments("path", req.Src)
	} else if !stat.IsDir() {
		return false, nil
	}

	// In local file mode, the returned operation will simply contain the source file path
	req.Inplace = true

	// Set pwd if relative
	if !filepath.IsAbs(req.Src) {
		var err error
		req.Pwd, err = os.Getwd()
		if err != nil {
			return false, fmt.Errorf("unable to get current directory")
		}
	}

	// Detect
	ok, err := g.FileGetter.Detect(req)

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

func (g *FileGetter) Get(ctx context.Context, req *getter.Request) error {
	// Log
	g.log.Debug("get repository",
		"src", req.Src,
	)

	// Get
	err := g.FileGetter.Get(ctx, req)

	if err != nil {
		// Log
		g.log.Debug("unable to get repository",
			"error", err,
		)

		g.result.AddGetError(err, g.protocol)

		return err
	}

	// Switch back to relative dst
	if req.Pwd != "" {
		var err error
		req.Dst, err = filepath.Rel(req.Pwd, req.Dst)
		if err != nil {
			return fmt.Errorf("unable to get relative path")
		}
	}

	return nil
}
