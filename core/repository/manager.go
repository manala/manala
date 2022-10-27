package repository

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/caarlos0/log"
	"github.com/go-git/go-git/v5"
	"manala/core"
	internalLog "manala/internal/log"
	internalOs "manala/internal/os"
	internalReport "manala/internal/report"
	"os"
	"path/filepath"
)

/***********/
/* Default */
/***********/

func NewDefaultManager(log *internalLog.Logger, path string, manager core.RepositoryManager) *DefaultManager {
	return &DefaultManager{
		log:     log,
		path:    path,
		manager: manager,
	}
}

type DefaultManager struct {
	log     *internalLog.Logger
	path    string
	manager core.RepositoryManager
}

func (manager *DefaultManager) LoadRepository(path string) (core.Repository, error) {
	// Replace path by default if empty
	if path == "" {
		path = manager.path
	}

	return manager.manager.LoadRepository(path)
}

/*********/
/* Cache */
/*********/

func NewCacheManager(log *internalLog.Logger, manager core.RepositoryManager) *CacheManager {
	return &CacheManager{
		log:     log,
		cache:   make(map[string]core.Repository),
		manager: manager,
	}
}

type CacheManager struct {
	log     *internalLog.Logger
	cache   map[string]core.Repository
	manager core.RepositoryManager
}

func (manager *CacheManager) LoadRepository(path string) (core.Repository, error) {
	// Check if repository already in cache
	if repo, ok := manager.cache[path]; ok {
		manager.log.Debug("load from cache")
		return repo, nil
	}

	repo, err := manager.manager.LoadRepository(path)
	if err != nil {
		return nil, err
	}

	// Cache repository
	manager.cache[path] = repo

	return repo, nil
}

/*********/
/* Chain */
/*********/

func NewChainManager(log *internalLog.Logger, managers []core.RepositoryManager) *ChainManager {
	return &ChainManager{
		log:      log,
		managers: managers,
	}
}

type ChainManager struct {
	log      *internalLog.Logger
	managers []core.RepositoryManager
}

func (manager *ChainManager) LoadRepository(path string) (core.Repository, error) {
	var repo core.Repository

	// Try managers
	for _, _manager := range manager.managers {
		_repo, err := _manager.LoadRepository(path)
		if err != nil {
			var _unsupportedRepositoryError *core.UnsupportedRepositoryError
			if errors.As(err, &_unsupportedRepositoryError) {
				continue
			}
			return nil, err
		}
		repo = _repo
		break
	}

	if repo == nil {
		return nil, internalReport.NewError(
			core.NewUnsupportedRepositoryError("unsupported repository"),
		)
	}

	return repo, nil
}

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

func (manager *DirManager) LoadRepository(path string) (core.Repository, error) {
	// Is path empty ?
	if path == "" {
		return nil, internalReport.NewError(
			core.NewUnsupportedRepositoryError("unsupported repository"),
		)
	}

	// Log
	manager.log.WithFields(log.Fields{
		"path":   path,
		"loader": "dir",
	}).Debug("try load")

	repo := NewRepository(
		manager.log,
		path,
		path,
	)

	stat, err := os.Stat(repo.dir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, internalReport.NewError(core.NewNotFoundRepositoryError("repository not found")).
				WithField("path", repo.path)
		}
		return nil, internalReport.NewError(internalOs.NewError(err)).
			WithMessage("file system error")
	} else if !stat.IsDir() {
		return nil, internalReport.NewError(fmt.Errorf("wrong repository")).
			WithField("dir", repo.dir)
	}

	return repo, nil
}

/*******/
/* Git */
/*******/

func NewGitManager(log *internalLog.Logger, cacheDir string) *GitManager {
	return &GitManager{
		log:      log,
		cacheDir: cacheDir,
	}
}

type GitManager struct {
	log      *internalLog.Logger
	cacheDir string
}

func (manager *GitManager) LoadRepository(path string) (core.Repository, error) {
	// Is git repo format ?
	if !(&core.GitRepoFormatChecker{}).IsFormat(path) {
		return nil, internalReport.NewError(
			core.NewUnsupportedRepositoryError("unsupported repository"),
		)
	}

	// Log
	manager.log.WithFields(log.Fields{
		"path":   path,
		"loader": "git",
	}).Debug("try load")

	hash := sha1.New()
	hash.Write([]byte(path))

	// Repository cache directory should be unique
	dir := filepath.Join(manager.cacheDir, "repositories", hex.EncodeToString(hash.Sum(nil)))

	manager.log.WithField("dir", dir).Debug("open git repository cache")

	repo := NewRepository(
		manager.log,
		path,
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
			URL:               repo.path,
			RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		})
		if err != nil {
			return nil, internalReport.NewError(err).
				WithMessage("clone git repository").
				WithField("dir", repo.dir).
				WithField("url", repo.path)
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
