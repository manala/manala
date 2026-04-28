package mocks

import (
	"github.com/manala/manala/app"
	"github.com/manala/manala/internal/schema"
	"github.com/manala/manala/internal/sync"

	"github.com/stretchr/testify/mock"
)

/***********/
/* Project */
/***********/

// ProjectMock mock a Project.
type ProjectMock struct {
	mock.Mock
}

func (mock *ProjectMock) Dir() string {
	args := mock.Called()

	return args.String(0)
}

func (mock *ProjectMock) Recipe() app.Recipe {
	args := mock.Called()

	return args.Get(0).(app.Recipe)
}

func (mock *ProjectMock) Vars() map[string]any {
	args := mock.Called()

	return args.Get(0).(map[string]any)
}

func (mock *ProjectMock) Watches() ([]string, error) {
	args := mock.Called()

	return args.Get(0).([]string), args.Error(1)
}

/**********/
/* Recipe */
/**********/

// RecipeMock mock a Recipe.
type RecipeMock struct {
	mock.Mock
}

func (mock *RecipeMock) Dir() string {
	args := mock.Called()

	return args.String(0)
}

func (mock *RecipeMock) Name() string {
	args := mock.Called()

	return args.String(0)
}

func (mock *RecipeMock) Description() string {
	args := mock.Called()

	return args.String(0)
}

func (mock *RecipeMock) Icon() string {
	args := mock.Called()

	return args.String(0)
}

func (mock *RecipeMock) Vars() map[string]any {
	args := mock.Called()

	return args.Get(0).(map[string]any)
}

func (mock *RecipeMock) Sync() []sync.UnitInterface {
	args := mock.Called()

	return args.Get(0).([]sync.UnitInterface)
}

func (mock *RecipeMock) Schema() schema.Schema {
	args := mock.Called()

	return args.Get(0).(schema.Schema)
}

func (mock *RecipeMock) Options() []app.RecipeOption {
	args := mock.Called()

	return args.Get(0).([]app.RecipeOption)
}

func (mock *RecipeMock) Repository() app.Repository {
	args := mock.Called()

	return args.Get(0).(app.Repository)
}

func (mock *RecipeMock) Template() string {
	args := mock.Called()

	return args.String(0)
}

func (mock *RecipeMock) Partials() []string {
	args := mock.Called()

	return args.Get(0).([]string)
}

func (mock *RecipeMock) ProjectValidator() *schema.Validator {
	args := mock.Called()

	return args.Get(0).(*schema.Validator)
}

func (mock *RecipeMock) Watches() ([]string, error) {
	args := mock.Called()

	return args.Get(0).([]string), args.Error(1)
}

/**************/
/* Repository */
/**************/

// RepositoryMock mock a Repository.
type RepositoryMock struct {
	mock.Mock
}

func (mock *RepositoryMock) URL() string {
	args := mock.Called()

	return args.String(0)
}

func (mock *RepositoryMock) Dir() string {
	args := mock.Called()

	return args.String(0)
}
