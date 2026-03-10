package manifest

import (
	"errors"
	"os"
	"path/filepath"
)

type Finder struct{}

func NewFinder() *Finder {
	return &Finder{}
}

func (finder *Finder) Find(dir string) bool {
	manifestFile := filepath.Join(dir, filename)

	if _, err := os.Stat(manifestFile); errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}
