package app

import (
	"manala/internal/schema"
	"manala/internal/syncer"
	"manala/internal/template"
	"manala/internal/validator"

	"github.com/stretchr/testify/mock"
)

/***********/
/* Project */
/***********/

// ProjectMock mock a project.
type ProjectMock struct {
	mock.Mock
}

func (mock *ProjectMock) Dir() string {
	args := mock.Called()
	return args.String(0)
}

func (mock *ProjectMock) Recipe() Recipe {
	args := mock.Called()
	return args.Get(0).(Recipe)
}

func (mock *ProjectMock) Vars() map[string]any {
	args := mock.Called()
	return args.Get(0).(map[string]any)
}

func (mock *ProjectMock) Template() *template.Template {
	args := mock.Called()
	return args.Get(0).(*template.Template)
}

func (mock *ProjectMock) Watches() ([]string, error) {
	args := mock.Called()
	return args.Get(0).([]string), args.Error(1)
}

/**********/
/* Recipe */
/**********/

// RecipeMock mock a recipe.
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

func (mock *RecipeMock) Sync() []syncer.UnitInterface {
	args := mock.Called()
	return args.Get(0).([]syncer.UnitInterface)
}

func (mock *RecipeMock) Schema() schema.Schema {
	args := mock.Called()
	return args.Get(0).(schema.Schema)
}

func (mock *RecipeMock) Options() []RecipeOption {
	args := mock.Called()
	return args.Get(0).([]RecipeOption)
}

func (mock *RecipeMock) Repository() Repository {
	args := mock.Called()
	return args.Get(0).(Repository)
}

func (mock *RecipeMock) Template() *template.Template {
	args := mock.Called()
	return args.Get(0).(*template.Template)
}

func (mock *RecipeMock) ProjectManifestTemplate() *template.Template {
	args := mock.Called()
	return args.Get(0).(*template.Template)
}

func (mock *RecipeMock) ProjectValidator() validator.Validator {
	args := mock.Called()
	return args.Get(0).(validator.Validator)
}

func (mock *RecipeMock) Watches() ([]string, error) {
	args := mock.Called()
	return args.Get(0).([]string), args.Error(1)
}

/**************/
/* Repository */
/**************/

// RepositoryMock mock a repository.
type RepositoryMock struct {
	mock.Mock
}

func (mock *RepositoryMock) Url() string {
	args := mock.Called()
	return args.String(0)
}

func (mock *RepositoryMock) Dir() string {
	args := mock.Called()
	return args.String(0)
}
