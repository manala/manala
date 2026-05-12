package mocks

import (
	"github.com/stretchr/testify/mock"
)

// Repository mock a Repository.
type Repository struct {
	mock.Mock
}

func (r *Repository) URL() string {
	args := r.Called()

	return args.String(0)
}

func (r *Repository) Dir() string {
	args := r.Called()

	return args.String(0)
}
