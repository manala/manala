package core

import (
	"github.com/xeipuuv/gojsonschema"
	"io"
	internalReport "manala/internal/report"
	internalSyncer "manala/internal/syncer"
	internalTemplate "manala/internal/template"
	internalWatcher "manala/internal/watcher"
	"path/filepath"
)

/***********/
/* Project */
/***********/

func NewProjectManifestMock() *ProjectManifestMock {
	mock := &ProjectManifestMock{}
	return mock
}

type ProjectManifestMock struct {
	path       string
	recipe     string
	repository string
	vars       map[string]interface{}
}

func (manifest *ProjectManifestMock) WithDir(dir string) *ProjectManifestMock {
	manifest.path = filepath.Join(dir, "manifest")
	return manifest
}

func (manifest *ProjectManifestMock) Path() string {
	return manifest.path
}

func (manifest *ProjectManifestMock) Recipe() string {
	return manifest.recipe
}

func (manifest *ProjectManifestMock) Repository() string {
	return manifest.repository
}

func (manifest *ProjectManifestMock) WithVars(vars map[string]interface{}) *ProjectManifestMock {
	manifest.vars = vars
	return manifest
}

func (manifest *ProjectManifestMock) Vars() map[string]interface{} {
	return manifest.vars
}

func (manifest *ProjectManifestMock) ReadFrom(_ io.Reader) error {
	return nil
}

func (manifest *ProjectManifestMock) Report(_ gojsonschema.ResultError, _ *internalReport.Report) {
}

/**********/
/* Recipe */
/**********/

func NewRecipeManagerMock() *RecipeManagerMock {
	mock := &RecipeManagerMock{}
	return mock
}

type RecipeManagerMock struct {
	recipe Recipe
}

func (manager *RecipeManagerMock) WithLoadRecipe(rec Recipe) *RecipeManagerMock {
	manager.recipe = rec
	return manager
}

func (manager *RecipeManagerMock) LoadRecipe(_ string) (Recipe, error) {
	return manager.recipe, nil
}

func NewRecipeMock() *RecipeMock {
	mock := &RecipeMock{
		schema:                  map[string]interface{}{},
		template:                &internalTemplate.Template{},
		projectManifestTemplate: &internalTemplate.Template{},
	}
	return mock
}

type RecipeMock struct {
	name                    string
	vars                    map[string]interface{}
	schema                  map[string]interface{}
	repository              Repository
	template                *internalTemplate.Template
	projectManifestTemplate *internalTemplate.Template
}

func (rec *RecipeMock) Path() string {
	return ""
}

func (rec *RecipeMock) WithName(name string) *RecipeMock {
	rec.name = name
	return rec
}

func (rec *RecipeMock) Name() string {
	return rec.name
}

func (rec *RecipeMock) Description() string {
	return ""
}

func (rec *RecipeMock) WithVars(vars map[string]interface{}) *RecipeMock {
	rec.vars = vars
	return rec
}

func (rec *RecipeMock) Vars() map[string]interface{} {
	return rec.vars
}

func (rec *RecipeMock) Sync() []internalSyncer.UnitInterface {
	return nil
}

func (rec *RecipeMock) WithSchema(schema map[string]interface{}) *RecipeMock {
	rec.schema = schema
	return rec
}

func (rec *RecipeMock) Schema() map[string]interface{} {
	return rec.schema
}

func (rec *RecipeMock) InitVars(_ func(options []RecipeOption) error) (map[string]interface{}, error) {
	return nil, nil
}

func (rec *RecipeMock) WithRepository(repo Repository) *RecipeMock {
	rec.repository = repo
	return rec
}

func (rec *RecipeMock) Repository() Repository {
	return rec.repository
}

func (rec *RecipeMock) Template() *internalTemplate.Template {
	return rec.template
}

func (rec *RecipeMock) ProjectManifestTemplate() *internalTemplate.Template {
	return rec.projectManifestTemplate
}

func (rec *RecipeMock) Watch(watcher *internalWatcher.Watcher) error {
	return nil
}

func NewRecipeManifestMock() *RecipeManifestMock {
	mock := &RecipeManifestMock{}
	return mock
}

type RecipeManifestMock struct {
	path        string
	description string
	template    string
	vars        map[string]interface{}
	initVars    map[string]interface{}
}

func (manifest *RecipeManifestMock) WithDir(dir string) *RecipeManifestMock {
	manifest.path = filepath.Join(dir, "manifest")
	return manifest
}

func (manifest *RecipeManifestMock) Path() string {
	return manifest.path
}

func (manifest *RecipeManifestMock) WithDescription(description string) *RecipeManifestMock {
	manifest.description = description
	return manifest
}

func (manifest *RecipeManifestMock) Description() string {
	return manifest.description
}

func (manifest *RecipeManifestMock) WithTemplate(template string) *RecipeManifestMock {
	manifest.template = template
	return manifest
}

func (manifest *RecipeManifestMock) Template() string {
	return manifest.template
}

func (manifest *RecipeManifestMock) WithVars(vars map[string]interface{}) *RecipeManifestMock {
	manifest.vars = vars
	return manifest
}

func (manifest *RecipeManifestMock) Vars() map[string]interface{} {
	return manifest.vars
}

func (manifest *RecipeManifestMock) Sync() []internalSyncer.UnitInterface {
	return nil
}

func (manifest *RecipeManifestMock) Schema() map[string]interface{} {
	return map[string]interface{}{}
}

func (manifest *RecipeManifestMock) ReadFrom(_ io.Reader) error {
	return nil
}

func (manifest *RecipeManifestMock) Report(_ gojsonschema.ResultError, _ *internalReport.Report) {
}

func (manifest *RecipeManifestMock) WithInitVars(vars map[string]interface{}) *RecipeManifestMock {
	manifest.initVars = vars
	return manifest
}

func (manifest *RecipeManifestMock) InitVars(_ func(options []RecipeOption) error) (map[string]interface{}, error) {
	return manifest.initVars, nil
}

func NewRecipeOptionMock() *RecipeOptionMock {
	mock := &RecipeOptionMock{}
	return mock
}

type RecipeOptionMock struct {
	label  string
	schema map[string]interface{}
	value  *interface{}
}

func (option *RecipeOptionMock) Label() string {
	return option.label
}

func (option *RecipeOptionMock) WithSchema(schema map[string]interface{}) *RecipeOptionMock {
	option.schema = schema
	return option
}

func (option *RecipeOptionMock) Schema() map[string]interface{} {
	return option.schema
}

func (option *RecipeOptionMock) WithSetValue(value *interface{}) *RecipeOptionMock {
	option.value = value
	return option
}

func (option *RecipeOptionMock) Set(value interface{}) error {
	*option.value = value
	return nil
}

/**************/
/* Repository */
/**************/

func NewRepositoryManagerMock() *RepositoryManagerMock {
	mock := &RepositoryManagerMock{}
	return mock
}

type RepositoryManagerMock struct {
	repository Repository
}

func (manager *RepositoryManagerMock) WithLoadRepository(repo Repository) *RepositoryManagerMock {
	manager.repository = repo
	return manager
}

func (manager *RepositoryManagerMock) LoadRepository(_ []string) (Repository, error) {
	return manager.repository, nil
}

func NewRepositoryMock() *RepositoryMock {
	mock := &RepositoryMock{}
	return mock
}

type RepositoryMock struct {
	path   string
	dir    string
	recipe Recipe
}

func (repo *RepositoryMock) WithPath(path string) *RepositoryMock {
	repo.path = path
	return repo
}

func (repo *RepositoryMock) Path() string {
	return repo.path
}

func (repo *RepositoryMock) Source() string {
	return repo.Path()
}

func (repo *RepositoryMock) WithDir(dir string) *RepositoryMock {
	repo.dir = dir
	return repo
}

func (repo *RepositoryMock) Dir() string {
	return repo.dir
}

func (repo *RepositoryMock) WithLoadRecipe(rec Recipe) *RepositoryMock {
	repo.recipe = rec
	return repo
}

func (repo *RepositoryMock) LoadRecipe(_ string) (Recipe, error) {
	return repo.recipe, nil
}

func (repo *RepositoryMock) WalkRecipes(_ func(rec Recipe)) error {
	return nil
}
