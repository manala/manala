package mocks

import (
	"github.com/manala/manala/app"
	"github.com/manala/manala/app/sync"

	"github.com/stretchr/testify/mock"
)

// Recipe mock a Recipe.
type Recipe struct {
	mock.Mock
}

func (r *Recipe) Dir() string {
	args := r.Called()

	return args.String(0)
}

func (r *Recipe) Name() string {
	args := r.Called()

	return args.String(0)
}

func (r *Recipe) Description() string {
	args := r.Called()

	return args.String(0)
}

func (r *Recipe) Icon() string {
	args := r.Called()

	return args.String(0)
}

func (r *Recipe) Template() string {
	args := r.Called()

	return args.String(0)
}

func (r *Recipe) Partials() []string {
	args := r.Called()

	return args.Get(0).([]string)
}

func (r *Recipe) Sync() []sync.Unit {
	args := r.Called()

	return args.Get(0).([]sync.Unit)
}

func (r *Recipe) Repository() app.Repository {
	args := r.Called()

	return args.Get(0).(app.Repository)
}

func (r *Recipe) Vars() map[string]any {
	args := r.Called()

	return args.Get(0).(map[string]any)
}

func (r *Recipe) Schema() map[string]any {
	args := r.Called()

	return args.Get(0).(map[string]any)
}

func (r *Recipe) Options() []app.RecipeOption {
	args := r.Called()

	return args.Get(0).([]app.RecipeOption)
}

func (r *Recipe) Watches() ([]string, error) {
	args := r.Called()

	return args.Get(0).([]string), args.Error(1)
}
