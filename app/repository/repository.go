package repository

func NewRepository(url string, dir string) *Repository {
	return &Repository{
		url: url,
		dir: dir,
	}
}

type Repository struct {
	url string
	dir string
}

func (repository *Repository) Url() string {
	return repository.url
}

func (repository *Repository) Dir() string {
	return repository.dir
}
