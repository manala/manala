package manifest

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"
)

func NewFinder(log *slog.Logger) *Finder {
	return &Finder{
		log: log,
	}
}

type Finder struct {
	log *slog.Logger
}

func (finder *Finder) Find(dir string) bool {
	finder.log.Info("finding project…")

	manifestFile := filepath.Join(dir, filename)

	if _, err := os.Stat(manifestFile); errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}
