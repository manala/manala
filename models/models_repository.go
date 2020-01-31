package models

// Create a repository
func NewRepository(src string, dir string) RepositoryInterface {
	return &repository{
		src: src,
		dir: dir,
	}
}

type RepositoryInterface interface {
	Src() string
	Dir() string
}

type repository struct {
	src string
	dir string
}

func (repo *repository) Src() string {
	return repo.src
}

func (repo *repository) Dir() string {
	return repo.dir
}
