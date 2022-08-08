package internal

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"github.com/caarlos0/log"
	"github.com/go-git/go-git/v5"
	internalGit "manala/internal/git"
	internalLog "manala/internal/log"
	internalOs "manala/internal/os"
	internalValidator "manala/internal/validator"
	"os"
	"path/filepath"
)

type RepositoryLoaderInterface interface {
	LoadRepository(paths []string) (*Repository, error)
}

/*******/
/* Dir */
/*******/

type RepositoryDirLoader struct {
	Log *internalLog.Logger
}

func (loader *RepositoryDirLoader) LoadRepository(paths []string) (*Repository, error) {
	// Get path
	var path string
	for _, _path := range paths {
		if _path != "" {
			path = _path
			break
		}
	}

	// Is path empty ?
	if path == "" {
		return nil, UnsupportedRepositoryPathError()
	}

	// Log
	loader.Log.WithFields(log.Fields{
		"path":   path,
		"loader": "dir",
	}).Debug("try load")

	repository := &Repository{
		log:  loader.Log,
		path: path,
		dir:  path,
	}

	stat, err := os.Stat(repository.dir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, NotFoundRepositoryDirError(repository.path)
		}
		return nil, internalOs.FileSystemError(err)
	} else if !stat.IsDir() {
		return nil, WrongRepositoryDirError(repository.dir)
	}

	repository.recipeLoader = &RecipeRepositoryDirLoader{
		Log:        loader.Log,
		Repository: repository,
	}

	return repository, nil
}

/*******/
/* Git */
/*******/

type RepositoryGitLoader struct {
	Log      *internalLog.Logger
	CacheDir string
}

func (loader *RepositoryGitLoader) LoadRepository(paths []string) (*Repository, error) {
	// Get path
	var path string
	for _, _path := range paths {
		if _path != "" {
			path = _path
			break
		}
	}

	// Is git regex match ?
	if !internalValidator.ValidateFormat("git-repo", path) {
		return nil, UnsupportedRepositoryPathError()
	}

	// Log
	loader.Log.WithFields(log.Fields{
		"path":   path,
		"loader": "git",
	}).Debug("try load")

	repository := &Repository{
		log:  loader.Log,
		path: path,
	}

	hash := sha1.New()
	hash.Write([]byte(repository.path))

	// Repository cache directory should be unique
	repository.dir = filepath.Join(loader.CacheDir, "repositories", hex.EncodeToString(hash.Sum(nil)))

	loader.Log.WithField("dir", repository.dir).Debug("open git repository cache")

Load:
	if err := os.MkdirAll(repository.dir, os.FileMode(0700)); err != nil {
		return nil, internalOs.FileSystemError(err)
	}

	gitRepository, err := git.PlainOpen(repository.dir)

	switch err {

	// Repository not in cache, let's clone it
	case git.ErrRepositoryNotExists:
		loader.Log.Debug("clone git repository cache")

		_, err = git.PlainClone(repository.dir, false, &git.CloneOptions{
			URL:               repository.path,
			RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		})
		if err != nil {
			return nil, internalGit.CloneRepositoryUrlError(repository.dir, repository.path, err)
		}

	// Repository already in cache, let's pull it
	case nil:
		loader.Log.Debug("get git repository cache worktree")

		gitRepositoryWorktree, err := gitRepository.Worktree()
		if err != nil {
			return nil, internalGit.InvalidRepositoryError(repository.dir, err)
		}

		loader.Log.Debug("pull git repository cache worktree")

		if err := gitRepositoryWorktree.Pull(&git.PullOptions{
			RemoteName: "origin",
		}); err != nil {
			switch err {
			case git.NoErrAlreadyUpToDate:
			case git.ErrNonFastForwardUpdate:
				loader.Log.Debug("fast forward update detected, delete git repository cache and retry with cloning")
				if err := os.RemoveAll(repository.dir); err != nil {
					return nil, internalGit.DeleteRepositoryError(repository.dir, err)
				}
				goto Load
			default:
				return nil, internalGit.PullRepositoryError(repository.dir, err)
			}
		}

	// Unable to open repository...
	default:
		return nil, internalGit.OpenRepositoryError(repository.dir, err)
	}

	repository.recipeLoader = &RecipeRepositoryDirLoader{
		Log:        loader.Log,
		Repository: repository,
	}

	return repository, nil
}
