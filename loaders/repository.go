package loaders

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/apex/log"
	"github.com/go-git/go-git/v5"
	"github.com/mingrammer/commonregex"
	"manala/models"
	"os"
	"path/filepath"
)

func NewRepositoryLoader(log log.Interface) RepositoryLoaderInterface {
	return &repositoryLoader{
		log:   log,
		cache: make(map[string]models.RepositoryInterface),
	}
}

type RepositoryLoaderInterface interface {
	Load(src string, cacheDir string) (models.RepositoryInterface, error)
}

type repositoryLoader struct {
	log   log.Interface
	cache map[string]models.RepositoryInterface
}

func (ld *repositoryLoader) Load(src string, cacheDir string) (models.RepositoryInterface, error) {
	// Check if repository already in cache
	if repo, ok := ld.cache[src]; ok {
		return repo, nil
	}

	var err error
	var repo models.RepositoryInterface

	// Is source a git repo ?
	if commonregex.GitRepoRegex.MatchString(src) {
		repo, err = ld.loadGit(src, cacheDir)
	} else {
		repo, err = ld.loadDir(src)
	}

	if err != nil {
		return nil, err
	}

	// Cache repository
	ld.cache[src] = repo

	return repo, nil
}

func (ld *repositoryLoader) loadDir(src string) (models.RepositoryInterface, error) {
	ld.log.WithField("source", src).Debug("Loading dir repository...")

	stat, err := os.Stat(src)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("\"%s\" directory does not exists", src)
		}
		return nil, err
	} else if !stat.IsDir() {
		return nil, fmt.Errorf("\"%s\" is not a directory", src)
	}

	return models.NewRepository(src, src), nil
}

func (ld *repositoryLoader) loadGit(src string, cacheDir string) (models.RepositoryInterface, error) {
	hash := md5.New()
	hash.Write([]byte(src))

	ld.log.WithField("source", src).Debug("Loading git repository...")

	// Repository cache directory should be unique
	dir := filepath.Join(cacheDir, "repositories", hex.EncodeToString(hash.Sum(nil)))

	ld.log.WithField("dir", dir).Debug("Opening repository cache...")

Load:
	if err := os.MkdirAll(dir, os.FileMode(0700)); err != nil {
		return nil, err
	}

	gitRepository, err := git.PlainOpen(dir)

	switch err {

	// Repository not in cache, let's clone it
	case git.ErrRepositoryNotExists:
		ld.log.Debug("Cloning git repository cache...")

		_, err = git.PlainClone(dir, false, &git.CloneOptions{
			URL:               src,
			RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		})
		if err != nil {
			return nil, fmt.Errorf("unable to clone repository: %w", err)
		}

	// Repository already in cache, let's pull it
	case nil:
		ld.log.Debug("Getting git repository worktree cache...")

		gitRepositoryWorktree, err := gitRepository.Worktree()
		if err != nil {
			return nil, fmt.Errorf("invalid repository: %w", err)
		}

		ld.log.Debug("Pulling cache git repository worktree...")

		if err := gitRepositoryWorktree.Pull(&git.PullOptions{
			RemoteName: "origin",
		}); err != nil {
			switch err {
			case git.NoErrAlreadyUpToDate:
			case git.ErrNonFastForwardUpdate:
				ld.log.Debug("Fast forward update detected, delete repository cache and retry with cloning...")
				if err := os.RemoveAll(dir); err != nil {
					return nil, fmt.Errorf("unable to delete repository cache: %w", err)
				}
				goto Load
			default:
				return nil, fmt.Errorf("unable to pull repository: %w", err)
			}
		}

	// Unable to open repository...
	default:
		return nil, fmt.Errorf("unable to open repository: %w", err)
	}

	return models.NewRepository(src, dir), nil
}
