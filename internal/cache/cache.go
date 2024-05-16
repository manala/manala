package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"manala/internal/serrors"
	"os"
	"path/filepath"
)

func New(dir string) *Cache {
	cache := &Cache{
		dir: dir,
	}

	return cache
}

type Cache struct {
	dir     string
	dirs    []string
	userDir string
}

func (cache *Cache) Dir() (string, error) {
	dirs := filepath.Join(cache.dirs...)

	// Use user dir
	if cache.dir == "" && cache.userDir != "" {
		userDir, err := os.UserCacheDir()
		if err != nil {
			return "", serrors.New("unable to get user cache dir").
				WithErrors(err)
		}

		return filepath.Join(userDir, cache.userDir, dirs), nil
	}

	return filepath.Join(cache.dir, dirs), nil
}

func (cache *Cache) WithDir(dir string) *Cache {
	clone := *cache
	clone.dirs = append(clone.dirs, dir)

	return &clone
}

func (cache *Cache) WithUserDir(dir string) *Cache {
	clone := *cache
	clone.userDir = dir

	return &clone
}

func (cache *Cache) WithHashDir(dir string) *Cache {
	clone := *cache

	hash := sha256.New224()
	hash.Write([]byte(dir))

	clone.dirs = append(clone.dirs,
		hex.EncodeToString(hash.Sum(nil)),
	)

	return &clone
}
