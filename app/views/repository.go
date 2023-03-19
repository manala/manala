package views

import (
	"manala/app/interfaces"
)

func NormalizeRepository(repo interfaces.Repository) *RepositoryView {
	url := repo.Url()

	return &RepositoryView{
		Url:    url,
		Path:   url,
		Source: url,
	}
}

type RepositoryView struct {
	Url string
	// Path ensure backward compatibility, when "path" was used instead of "url" to define repository origin
	Path string
	// Source ensure backward compatibility, when "source" was used instead of "path" to define repository origin
	Source string
}
