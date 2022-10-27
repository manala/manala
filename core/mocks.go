package core

import (
	"github.com/stretchr/testify/mock"
	"github.com/xeipuuv/gojsonschema"
	"io"
	internalReport "manala/internal/report"
	internalSyncer "manala/internal/syncer"
	internalTemplate "manala/internal/template"
	internalWatcher "manala/internal/watcher"
)

/***********/
/* Project */
/***********/

func NewProjectManifestMock() *ProjectManifestMock {
	return &ProjectManifestMock{}
}

type ProjectManifestMock struct {
	mock.Mock
}

func (manifest *ProjectManifestMock) Path() string {
	args := manifest.Called()
	return args.String(0)
}

func (manifest *ProjectManifestMock) Recipe() string {
	args := manifest.Called()
	return args.String(0)
}

func (manifest *ProjectManifestMock) Repository() string {
	args := manifest.Called()
	return args.String(0)
}

func (manifest *ProjectManifestMock) Vars() map[string]interface{} {
	args := manifest.Called()
	return args.Get(0).(map[string]interface{})
}

func (manifest *ProjectManifestMock) ReadFrom(reader io.Reader) error {
	args := manifest.Called(reader)
	return args.Error(0)
}

func (manifest *ProjectManifestMock) Report(result gojsonschema.ResultError, report *internalReport.Report) {
	manifest.Called(result, report)
}

/**********/
/* Recipe */
/**********/

func NewRecipeManagerMock() *RecipeManagerMock {
	return &RecipeManagerMock{}
}

type RecipeManagerMock struct {
	mock.Mock
}

func (manager *RecipeManagerMock) LoadRecipe(name string) (Recipe, error) {
	args := manager.Called(name)
	return args.Get(0).(Recipe), args.Error(1)
}

func NewRecipeMock() *RecipeMock {
	return &RecipeMock{}
}

type RecipeMock struct {
	mock.Mock
}

func (rec *RecipeMock) Path() string {
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

func (rec *RecipeMock) InitVars(callback func(options []RecipeOption) error) (map[string]interface{}, error) {
	args := rec.Called(callback)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (rec *RecipeMock) Repository() Repository {
	args := rec.Called()
	return args.Get(0).(Repository)
}

func (rec *RecipeMock) Template() *internalTemplate.Template {
	args := rec.Called()
	return args.Get(0).(*internalTemplate.Template)
}

func (rec *RecipeMock) ProjectManifestTemplate() *internalTemplate.Template {
	args := rec.Called()
	return args.Get(0).(*internalTemplate.Template)
}

func (rec *RecipeMock) Watch(watcher *internalWatcher.Watcher) error {
	args := rec.Called(watcher)
	return args.Error(0)
}

func NewRecipeManifestMock() *RecipeManifestMock {
	return &RecipeManifestMock{}
}

type RecipeManifestMock struct {
	mock.Mock
}

func (manifest *RecipeManifestMock) Path() string {
	args := manifest.Called()
	return args.String(0)
}

func (manifest *RecipeManifestMock) Description() string {
	args := manifest.Called()
	return args.String(0)
}

func (manifest *RecipeManifestMock) Template() string {
	args := manifest.Called()
	return args.String(0)
}

func (manifest *RecipeManifestMock) Vars() map[string]interface{} {
	args := manifest.Called()
	return args.Get(0).(map[string]interface{})
}

func (manifest *RecipeManifestMock) Sync() []internalSyncer.UnitInterface {
	args := manifest.Called()
	return args.Get(0).([]internalSyncer.UnitInterface)
}

func (manifest *RecipeManifestMock) Schema() map[string]interface{} {
	args := manifest.Called()
	return args.Get(0).(map[string]interface{})
}

func (manifest *RecipeManifestMock) ReadFrom(reader io.Reader) error {
	args := manifest.Called(reader)
	return args.Error(0)
}

func (manifest *RecipeManifestMock) Report(result gojsonschema.ResultError, report *internalReport.Report) {
	manifest.Called(result, report)
}

func (manifest *RecipeManifestMock) InitVars(callback func(options []RecipeOption) error) (map[string]interface{}, error) {
	args := manifest.Called(callback)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func NewRecipeOptionMock() *RecipeOptionMock {
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

/**************/
/* Repository */
/**************/

func NewRepositoryManagerMock() *RepositoryManagerMock {
	return &RepositoryManagerMock{}
}

type RepositoryManagerMock struct {
	mock.Mock
}

func (manager *RepositoryManagerMock) LoadRepository(path string) (Repository, error) {
	args := manager.Called(path)
	return args.Get(0).(Repository), args.Error(1)
}

func NewRepositoryMock() *RepositoryMock {
	return &RepositoryMock{}
}

type RepositoryMock struct {
	mock.Mock
}

func (repo *RepositoryMock) Path() string {
	args := repo.Called()
	return args.String(0)
}

func (repo *RepositoryMock) Source() string {
	args := repo.Called()
	return args.String(0)
}

func (repo *RepositoryMock) Dir() string {
	args := repo.Called()
	return args.String(0)
}

func (repo *RepositoryMock) LoadRecipe(name string) (Recipe, error) {
	args := repo.Called(name)
	return args.Get(0).(Recipe), args.Error(1)
}

func (repo *RepositoryMock) WalkRecipes(walker func(rec Recipe)) error {
	args := repo.Called(walker)
	return args.Error(0)
}
