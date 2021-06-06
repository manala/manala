package models

// NewRepository creates a repository
func NewRepository(source string, dir string, main bool) RepositoryInterface {
	return &repository{
		source: source,
		dir:    dir,
		main:   main,
	}
}

type RepositoryInterface interface {
	model
	Source() string
	Main() bool
}

type repository struct {
	source string
	dir    string
	main   bool
}

func (repo *repository) Source() string {
	return repo.source
}

func (repo *repository) getDir() string {
	return repo.dir
}

func (repo *repository) Main() bool {
	return repo.main
}
