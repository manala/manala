package repository

import (
	"crypto/sha1"
	"encoding/hex"
	"github.com/caarlos0/log"
	"github.com/go-git/go-git/v5"
	"manala/core"
	internalCache "manala/internal/cache"
	internalLog "manala/internal/log"
	internalOs "manala/internal/os"
	internalReport "manala/internal/report"
	"os"
)

func NewGitManager(log *internalLog.Logger, cache *internalCache.Cache) *GitManager {
	return &GitManager{
		log:   log,
		cache: cache,
	}
}

type GitManager struct {
	log   *internalLog.Logger
	cache *internalCache.Cache
}

func (manager *GitManager) LoadRepository(url string) (core.Repository, error) {
	// Is git repo format ?
	if !(&core.GitRepoFormatChecker{}).IsFormat(url) {
		return nil, internalReport.NewError(
			core.NewUnsupportedRepositoryError("unsupported repository url"),
		).WithField("url", url)
	}

	// Log
	manager.log.WithFields(log.Fields{
		"url":     url,
		"manager": "git",
	}).Debug("try load")

	hash := sha1.New()
	hash.Write([]byte(url))

	// Repository cache directory should be unique
	dir, err := manager.cache.Dir("repositories", hex.EncodeToString(hash.Sum(nil)))
	if err != nil {
		return nil, err
	}

	manager.log.WithField("dir", dir).Debug("open git repository cache")

	repo := NewRepository(
		url,
		dir,
	)

Load:
	if err := os.MkdirAll(repo.dir, os.FileMode(0700)); err != nil {
		return nil, internalReport.NewError(internalOs.NewError(err)).
			WithMessage("file system error")
	}

	gitRepository, err := git.PlainOpen(repo.dir)

	switch err {

	// Repository not in cache, let's clone it
	case git.ErrRepositoryNotExists:
		manager.log.Debug("clone git repository cache")

		_, err = git.PlainClone(repo.dir, false, &git.CloneOptions{
			URL:               repo.url,
			RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		})
		if err != nil {
			return nil, internalReport.NewError(err).
				WithMessage("clone git repository").
				WithField("dir", repo.dir).
				WithField("url", repo.url)
		}

	// Repository already in cache, let's pull it
	case nil:
		manager.log.Debug("get git repository cache worktree")

		gitRepositoryWorktree, err := gitRepository.Worktree()
		if err != nil {
			return nil, internalReport.NewError(err).
				WithMessage("invalid git repository").
				WithField("dir", repo.dir)
		}

		manager.log.Debug("pull git repository cache worktree")

		if err := gitRepositoryWorktree.Pull(&git.PullOptions{
			RemoteName: "origin",
		}); err != nil {
			switch err {
			case git.NoErrAlreadyUpToDate:
			case git.ErrNonFastForwardUpdate:
				manager.log.Debug("fast forward update detected, delete git repository cache and retry with cloning")
				if err := os.RemoveAll(repo.dir); err != nil {
					return nil, internalReport.NewError(err).
						WithMessage("delete git repository").
						WithField("dir", repo.dir)
				}
				goto Load
			default:
				return nil, internalReport.NewError(err).
					WithMessage("pull git repository").
					WithField("dir", repo.dir)
			}
		}

	// Unable to open repository...
	default:
		return nil, internalReport.NewError(err).
			WithMessage("open git repository").
			WithField("dir", repo.dir)
	}

	return repo, nil
}
