package models

// NewRepository creates a repository
func NewRepository(source string, dir string) RepositoryInterface {
	return &repository{
		source: source,
		dir:    dir,
	}
}

type RepositoryInterface interface {
	model
	Source() string
}

type repository struct {
	source string
	dir    string
}

func (repo *repository) Source() string {
	return repo.source
}

func (repo *repository) getDir() string {
	return repo.dir
}
