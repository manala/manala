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

func (repo *Repository) Url() string {
	return repo.url
}

func (repo *Repository) Dir() string {
	return repo.dir
}
