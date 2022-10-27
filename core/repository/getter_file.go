package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/go-getter/v2"
	internalLog "manala/internal/log"
	internalOs "manala/internal/os"
	"os"
	"path/filepath"
)

func NewFileGetter(log *internalLog.Logger, result *GetterResult) *FileGetter {
	return &FileGetter{
		log:        log,
		result:     result,
		FileGetter: &getter.FileGetter{},
		protocol:   "file",
	}
}

type FileGetter struct {
	log    *internalLog.Logger
	result *GetterResult
	*getter.FileGetter
	protocol string
}

func (g *FileGetter) Detect(req *getter.Request) (bool, error) {
	g.log.
		WithField("protocol", g.protocol).
		Debug("detect")

	// Stat
	stat, err := os.Stat(req.Src)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, internalOs.NewError(err)
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
		g.log.
			WithField("protocol", g.protocol).
			WithError(err).
			Debug("unable to detect")

		g.result.SetDetectError(err, g.protocol)

		return ok, err
	}

	return ok, nil
}

func (g *FileGetter) Get(ctx context.Context, req *getter.Request) error {
	g.log.
		WithField("protocol", g.protocol).
		Debug("get")

	// Get
	err := g.FileGetter.Get(ctx, req)

	if err != nil {
		g.log.
			WithField("protocol", g.protocol).
			WithError(err).
			Debug("unable to get ")

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
