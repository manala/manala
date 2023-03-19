package mocks

import (
	"github.com/stretchr/testify/mock"
	"github.com/xeipuuv/gojsonschema"
	"io"
	"manala/app/interfaces"
	internalReport "manala/internal/report"
	internalSyncer "manala/internal/syncer"
	internalTemplate "manala/internal/template"
	internalWatcher "manala/internal/watcher"
)

func MockRecipe() *RecipeMock {
	return &RecipeMock{}
}

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

func (rec *RecipeMock) Sync() []internalSyncer.UnitInterface {
	args := rec.Called()
	return args.Get(0).([]internalSyncer.UnitInterface)
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

func (rec *RecipeMock) Template() *internalTemplate.Template {
	args := rec.Called()
	return args.Get(0).(*internalTemplate.Template)
}

func (rec *RecipeMock) ProjectManifestTemplate() *internalTemplate.Template {
	args := rec.Called()
	return args.Get(0).(*internalTemplate.Template)
}

/************/
/* Manifest */
/************/

func MockRecipeManifest() *RecipeManifestMock {
	return &RecipeManifestMock{}
}

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

func (man *RecipeManifestMock) Sync() []internalSyncer.UnitInterface {
	args := man.Called()
	return args.Get(0).([]internalSyncer.UnitInterface)
}

func (man *RecipeManifestMock) Schema() map[string]interface{} {
	args := man.Called()
	return args.Get(0).(map[string]interface{})
}

func (man *RecipeManifestMock) ReadFrom(reader io.Reader) error {
	args := man.Called(reader)
	return args.Error(0)
}

func (man *RecipeManifestMock) Report(result gojsonschema.ResultError, report *internalReport.Report) {
	man.Called(result, report)
}

func (man *RecipeManifestMock) InitVars(callback func(options []interfaces.RecipeOption) error) (map[string]interface{}, error) {
	args := man.Called(callback)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

/**********/
/* Option */
/**********/

func MockRecipeOption() *RecipeOptionMock {
	return &RecipeOptionMock{}
}

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

func MockRecipeManager() *RecipeManagerMock {
	return &RecipeManagerMock{}
}

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

func (manager *RecipeManagerMock) WatchRecipe(rec interfaces.Recipe, watcher *internalWatcher.Watcher) error {
	args := manager.Called(rec, watcher)
	return args.Error(0)
}
