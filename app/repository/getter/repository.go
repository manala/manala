package getter

type Repository struct {
	url string
	dir string
}

func NewRepository(url string, dir string) *Repository {
	return &Repository{
		url: url,
		dir: dir,
	}
}

func (repository *Repository) URL() string {
	return repository.url
}

func (repository *Repository) Dir() string {
	return repository.dir
}
