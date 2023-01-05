package repository

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"github.com/hashicorp/go-getter/v2"
	"manala/core"
	internalCache "manala/internal/cache"
	internalLog "manala/internal/log"
)

func NewGetterManager(log *internalLog.Logger, cache *internalCache.Cache) *GetterManager {
	return &GetterManager{
		log:   log,
		cache: cache,
	}
}

type GetterManager struct {
	log   *internalLog.Logger
	cache *internalCache.Cache
}

func (manager *GetterManager) LoadRepository(url string) (core.Repository, error) {
	// Log
	manager.log.
		WithField("manager", "getter").
		WithField("url", url).
		Debug("load repository")
	manager.log.IncreasePadding()
	defer manager.log.DecreasePadding()

	// Repository cache directory should be unique
	hash := sha256.New224()
	hash.Write([]byte(url))
	cacheDir, err := manager.cache.Dir("repositories", hex.EncodeToString(hash.Sum(nil)))
	if err != nil {
		return nil, err
	}

	request := &getter.Request{
		Src:     url,
		Dst:     cacheDir,
		GetMode: getter.ModeDir,
	}

	result := NewGetterResult()

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
	if err := result.HandleError(err); err != nil {
		return nil, err.WithField("url", url)
	}

	return NewRepository(
		url,
		response.Dst,
	), nil
}
