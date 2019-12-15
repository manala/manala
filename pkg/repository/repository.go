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

// Create a repository
func New(src string) Interface {
	return &repository{
		src: src,
	}
}

type Interface interface {
	GetSrc() string
	GetDir() string
	Load(cacheDir string) error
}

type repository struct {
	src string
	dir string
}

var (
	ErrUnclonable = errors.New("repository unclonable")
	ErrUnopenable = errors.New("repository unopenable")
	ErrInvalid    = errors.New("repository invalid")
)

func (repo *repository) GetSrc() string {
	return repo.src
}

func (repo *repository) GetDir() string {
	return repo.dir
}

// Load repository
func (repo *repository) Load(cacheDir string) error {
	// Is src a git repo ?
	if commonregex.GitRepoRegex.MatchString(repo.src) {
		return repo.loadGit(cacheDir)
	}

	log.WithField("src", repo.src).Debug("Loading repository...")

	repo.dir = repo.src

	return nil
}

var cache = make(map[string]string)

func (repo *repository) loadGit(cacheDir string) error {

	// Check if repository dir already in cache
	if dir, ok := cache[repo.src]; ok {
		repo.dir = dir
		return nil
	}

	hash := md5.New()
	hash.Write([]byte(repo.src))

	log.WithField("src", repo.src).Debug("Loading git repository...")

	// Repository cache directory should be unique
	repo.dir = path.Join(cacheDir, "repositories", hex.EncodeToString(hash.Sum(nil)))

	log.WithField("dir", repo.dir).Debug("Opening repository cache...")

	err := os.MkdirAll(repo.dir, os.FileMode(0700))
	if err != nil {
		return err
	}

	gitRepository, err := git.PlainOpen(repo.dir)

	if err != nil {
		switch err {
		case git.ErrRepositoryNotExists:
			log.Debug("Cloning git repository cache...")

			gitRepository, err = git.PlainClone(repo.dir, false, &git.CloneOptions{
				URL:               repo.src,
				RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
			})
			if err != nil {
				return ErrUnclonable
			}
		default:
			return ErrUnopenable
		}
	} else {
		log.Debug("Getting git repository worktree cache...")

		gitRepositoryWorktree, err := gitRepository.Worktree()
		if err != nil {
			return ErrInvalid
		}

		log.Debug("Pulling cache git repository worktree...")

		if err := gitRepositoryWorktree.Pull(&git.PullOptions{
			RemoteName: "origin",
		}); err != nil {
			switch err {
			case git.NoErrAlreadyUpToDate:
			default:
				return err
			}
		}
	}

	// Cache repository
	cache[repo.src] = repo.dir

	return nil
}
