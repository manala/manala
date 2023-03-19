package interfaces

type Repository interface {
	Url() string
	Dir() string
}

type RepositoryManager interface {
	LoadRepository(url string) (Repository, error)
}
