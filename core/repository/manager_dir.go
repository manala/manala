package repository

import (
	"errors"
	"fmt"
	"github.com/caarlos0/log"
	"manala/core"
	internalLog "manala/internal/log"
	internalOs "manala/internal/os"
	internalReport "manala/internal/report"
	"os"
)

/*******/
/* Dir */
/*******/

func NewDirManager(log *internalLog.Logger) *DirManager {
	return &DirManager{
		log: log,
	}
}

type DirManager struct {
	log *internalLog.Logger
}

func (manager *DirManager) LoadRepository(url string) (core.Repository, error) {
	// Is url empty ?
	if url == "" {
		return nil, internalReport.NewError(
			core.NewUnsupportedRepositoryError("unsupported empty repository url"),
		)
	}

	// Log
	manager.log.WithFields(log.Fields{
		"url":     url,
		"manager": "dir",
	}).Debug("try load")

	repo := NewRepository(
		url,
		url,
	)

	stat, err := os.Stat(repo.dir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, internalReport.NewError(core.NewNotFoundRepositoryError("repository not found")).
				WithField("url", repo.url)
		}
		return nil, internalReport.NewError(internalOs.NewError(err)).
			WithMessage("file system error")
	} else if !stat.IsDir() {
		return nil, internalReport.NewError(fmt.Errorf("wrong repository")).
			WithField("dir", repo.dir)
	}

	return repo, nil
}
