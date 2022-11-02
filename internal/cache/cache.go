package cache

import (
	internalReport "manala/internal/report"
	"os"
	"path/filepath"
)

func New(dir string, opts ...Option) *Cache {
	cache := &Cache{
		dir: dir,
	}

	// Options
	for _, opt := range opts {
		opt(cache)
	}

	return cache
}

type Cache struct {
	dir     string
	userDir string
}

func (cache *Cache) Dir(dirs ...string) (string, error) {
	if cache.dir != "" {
		return filepath.Join(cache.dir, filepath.Join(dirs...)), nil
	}

	// Fallback to user cache dir
	userDir, err := os.UserCacheDir()
	if err != nil {
		return "", internalReport.NewError(err).
			WithMessage("unable to get user cache dir")
	}

	return filepath.Join(userDir, cache.userDir, filepath.Join(dirs...)), nil
}
