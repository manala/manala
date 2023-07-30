package mocks

import (
	"github.com/stretchr/testify/mock"
	"io"
	"manala/app/interfaces"
	"manala/internal/syncer"
	"manala/internal/template"
	"manala/internal/validation"
	"manala/internal/watcher"
)

type RecipeMock struct {
	mock.Mock
}

func (rec *RecipeMock) Dir() string {
	args := rec.Called()
	return args.String(0)
}

func (rec *RecipeMock) Name() string {
	args := rec.Called()
	return args.String(0)
}

func (rec *RecipeMock) Description() string {
	args := rec.Called()
	return args.String(0)
}

func (rec *RecipeMock) Vars() map[string]interface{} {
	args := rec.Called()
	return args.Get(0).(map[string]interface{})
}

func (rec *RecipeMock) Sync() []syncer.UnitInterface {
	args := rec.Called()
	return args.Get(0).([]syncer.UnitInterface)
}

func (rec *RecipeMock) Schema() map[string]interface{} {
	args := rec.Called()
	return args.Get(0).(map[string]interface{})
}

func (rec *RecipeMock) InitVars(callback func(options []interfaces.RecipeOption) error) (map[string]interface{}, error) {
	args := rec.Called(callback)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (rec *RecipeMock) Repository() interfaces.Repository {
	args := rec.Called()
	return args.Get(0).(interfaces.Repository)
}

func (rec *RecipeMock) Template() *template.Template {
	args := rec.Called()
	return args.Get(0).(*template.Template)
}

func (rec *RecipeMock) ProjectManifestTemplate() *template.Template {
	args := rec.Called()
	return args.Get(0).(*template.Template)
}

/************/
/* Manifest */
/************/

type RecipeManifestMock struct {
	mock.Mock
}

func (man *RecipeManifestMock) Description() string {
	args := man.Called()
	return args.String(0)
}

func (man *RecipeManifestMock) Template() string {
	args := man.Called()
	return args.String(0)
}

func (man *RecipeManifestMock) Vars() map[string]interface{} {
	args := man.Called()
	return args.Get(0).(map[string]interface{})
}

func (man *RecipeManifestMock) Sync() []syncer.UnitInterface {
	args := man.Called()
	return args.Get(0).([]syncer.UnitInterface)
}

func (man *RecipeManifestMock) Schema() map[string]interface{} {
	args := man.Called()
	return args.Get(0).(map[string]interface{})
}

func (man *RecipeManifestMock) ReadFrom(reader io.Reader) error {
	args := man.Called(reader)
	return args.Error(0)
}

func (man *RecipeManifestMock) InitVars(callback func(options []interfaces.RecipeOption) error) (map[string]interface{}, error) {
	args := man.Called(callback)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (man *RecipeManifestMock) ValidationResultErrorDecorator() validation.ResultErrorDecorator {
	args := man.Called()
	return args.Get(0).(validation.ResultErrorDecorator)
}

/**********/
/* Option */
/**********/

type RecipeOptionMock struct {
	mock.Mock
}

func (option *RecipeOptionMock) Label() string {
	args := option.Called()
	return args.String(0)
}

func (option *RecipeOptionMock) Schema() map[string]interface{} {
	args := option.Called()
	return args.Get(0).(map[string]interface{})
}

func (option *RecipeOptionMock) Set(value interface{}) error {
	args := option.Called(value)
	return args.Error(0)
}

/***********/
/* Manager */
/***********/

type RecipeManagerMock struct {
	mock.Mock
}

func (manager *RecipeManagerMock) LoadRecipe(repo interfaces.Repository, name string) (interfaces.Recipe, error) {
	args := manager.Called(repo, name)
	return args.Get(0).(interfaces.Recipe), args.Error(1)
}

func (manager *RecipeManagerMock) WalkRecipes(repo interfaces.Repository, walker func(rec interfaces.Recipe) error) error {
	args := manager.Called(repo, walker)
	return args.Error(0)
}

func (manager *RecipeManagerMock) WatchRecipe(rec interfaces.Recipe, watcher *watcher.Watcher) error {
	args := manager.Called(rec, watcher)
	return args.Error(0)
}
