package repository

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"github.com/hashicorp/go-getter/v2"
	"log/slog"
	"manala/app"
	"manala/internal/cache"
	"manala/internal/serrors"
)

func NewGetterManager(log *slog.Logger, cache *cache.Cache) *GetterManager {
	return &GetterManager{
		log:   log.With("manager", "getter"),
		cache: cache,
	}
}

type GetterManager struct {
	log   *slog.Logger
	cache *cache.Cache
}

func (manager *GetterManager) LoadRepository(url string) (app.Repository, error) {
	// Log
	manager.log.Debug("load repository",
		"url", url,
	)

	// Repository cache directory should be unique
	hash := sha256.New224()
	hash.Write([]byte(url))
	cacheDir, err := manager.cache.Dir(
		"repositories",
		hex.EncodeToString(hash.Sum(nil)),
	)
	if err != nil {
		return nil, serrors.New("unable to get cache dir").
			WithErrors(err)
	}

	request := &getter.Request{
		Src:     url,
		Dst:     cacheDir,
		GetMode: getter.ModeDir,
	}

	result := NewGetterResult(url)

	client := &getter.Client{
		// Prevent copying or writing files through symlinks
		DisableSymlinks: true,
		Getters: []getter.Getter{
			NewGitGetter(manager.log, result),
			NewS3Getter(manager.log, result),
			NewHttpGetter(manager.log, result),
			NewFileGetter(manager.log, result),
		},
		Decompressors: getter.Decompressors,
	}

	response, err := client.Get(context.Background(), request)
	if _err := result.Error(err); _err != nil {
		return nil, _err
	}

	return NewRepository(
		url,
		response.Dst,
	), nil
}
