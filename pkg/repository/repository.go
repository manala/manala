package repository

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/apex/log"
	"github.com/mingrammer/commonregex"
	"gopkg.in/src-d/go-git.v4"
	"os"
	"path"
)

type Repository struct {
	Src string
	Dir string
}

var (
	ErrUnclonable = errors.New("repository unclonable")
	ErrUnopenable = errors.New("repository unopenable")
	ErrInvalid    = errors.New("repository invalid")
)

// Create a repository
func New(src string) *Repository {
	return &Repository{
		Src: src,
	}
}

// Load a repository
func Load(src string, cacheDir string) (*Repository, error) {
	// Is src a git repo ?
	if commonregex.GitRepoRegex.MatchString(src) {
		return loadGit(src, cacheDir)
	}

	log.WithField("src", src).Debug("Loading repository...")

	repo := New(src)
	repo.Dir = src

	return repo, nil
}

var cache = make(map[string]*Repository)

func loadGit(src string, cacheDir string) (*Repository, error) {

	// Check if repository already in cache
	if repo, ok := cache[src]; ok {
		return repo, nil
	}

	repo := New(src)

	hash := md5.New()
	hash.Write([]byte(repo.Src))

	log.WithField("src", repo.Src).Debug("Loading git repository...")

	// Repository cache directory should be unique
	repo.Dir = path.Join(cacheDir, "repositories", hex.EncodeToString(hash.Sum(nil)))

	log.WithField("dir", repo.Dir).Debug("Opening repository cache...")

	err := os.MkdirAll(repo.Dir, os.FileMode(0700))
	if err != nil {
		return nil, err
	}

	gitRepository, err := git.PlainOpen(repo.Dir)

	if err != nil {
		switch err {
		case git.ErrRepositoryNotExists:
			log.Debug("Cloning git repository cache...")

			gitRepository, err = git.PlainClone(repo.Dir, false, &git.CloneOptions{
				URL:               src,
				RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
			})
			if err != nil {
				return nil, ErrUnclonable
			}
		default:
			return nil, ErrUnopenable
		}
	} else {
		log.Debug("Getting git repository worktree cache...")

		gitRepositoryWorktree, err := gitRepository.Worktree()
		if err != nil {
			return nil, ErrInvalid
		}

		log.Debug("Pulling cache git repository worktree...")

		if err := gitRepositoryWorktree.Pull(&git.PullOptions{
			RemoteName: "origin",
		}); err != nil {
			switch err {
			case git.NoErrAlreadyUpToDate:
			default:
				return nil, err
			}
		}
	}

	// Cache repository
	cache[src] = repo

	return repo, nil
}
