package repository

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/apex/log"
	"github.com/mingrammer/commonregex"
	"gopkg.in/src-d/go-git.v4"
	"io/ioutil"
	"manala/pkg/recipe"
	"os"
	"path"
	"path/filepath"
	"strings"
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
	WalkRecipes(fn walkRecipesFunc) error
	LoadRecipe(name string) (recipe.Interface, error)
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

// Walk into recipes
func (repo *repository) WalkRecipes(fn walkRecipesFunc) error {
	files, err := ioutil.ReadDir(repo.GetDir())
	if err != nil {
		return err
	}

	for _, file := range files {
		// Exclude dot files
		if strings.HasPrefix(file.Name(), ".") {
			continue
		}
		if file.IsDir() {
			rec := recipe.New(filepath.Join(repo.GetDir(), file.Name()))
			if err := rec.Load(recipe.Config{}); err != nil {
				return err
			}
			fn(rec)
		}
	}

	return nil
}

type walkRecipesFunc func(rec recipe.Interface)

// Load recipe by its name
func (repo *repository) LoadRecipe(name string) (recipe.Interface, error) {
	var baseRec recipe.Interface

	if err := repo.WalkRecipes(func(rec recipe.Interface) {
		fmt.Printf("%s: %s\n", rec.GetName(), rec.GetConfig().Description)
		if rec.GetName() == name {
			baseRec = rec
		}
	}); err != nil {
		return nil, err
	}

	if baseRec != nil {
		return baseRec, nil
	}

	return nil, fmt.Errorf("recipe not found")
}
