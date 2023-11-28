package app

import (
	"github.com/stretchr/testify/mock"
	"manala/internal/path"
	"manala/internal/schema"
	"manala/internal/syncer"
	"manala/internal/template"
	"manala/internal/validator"
	"manala/internal/watcher"
)

/***********/
/* Project */
/***********/

// Project describe a project interface
type Project interface {
	Dir() string
	Recipe() Recipe
	Vars() map[string]any
	Template() *template.Template
}

// ProjectMock mock a project
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

// ProjectManifest describe a project manifest interface
type ProjectManifest interface {
	Recipe() string
	Repository() string
	Vars() map[string]any
}

// ProjectManifestMock mock a project manifest
type ProjectManifestMock struct {
	mock.Mock
}

func (mock *ProjectManifestMock) Recipe() string {
	args := mock.Called()
	return args.String(0)
}

func (mock *ProjectManifestMock) Repository() string {
	args := mock.Called()
	return args.String(0)
}

func (mock *ProjectManifestMock) Vars() map[string]any {
	args := mock.Called()
	return args.Get(0).(map[string]any)
}

// ProjectManager describe a project manager interface
type ProjectManager interface {
	IsProject(dir string) bool
	CreateProject(dir string, recipe Recipe, vars map[string]any) (Project, error)
	LoadProject(dir string) (Project, error)
	WatchProject(project Project, watcher *watcher.Watcher) error
}

/**********/
/* Recipe */
/**********/

// Recipe describe a recipe interface
type Recipe interface {
	Dir() string
	Name() string
	Description() string
	Icon() string
	Vars() map[string]any
	Sync() []syncer.UnitInterface
	Schema() schema.Schema
	Options() []RecipeOption
	Repository() Repository
	Template() *template.Template
	ProjectManifestTemplate() *template.Template
	ProjectValidator() validator.Validator
}

// RecipeMock mock a recipe
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

// RecipeManifest describe a recipe manifest interface
type RecipeManifest interface {
	Description() string
	Icon() string
	Template() string
	Vars() map[string]any
	Sync() []syncer.UnitInterface
	Schema() schema.Schema
	Options() []RecipeOption
}

// RecipeManifestMock mock a recipe manifest
type RecipeManifestMock struct {
	mock.Mock
}

func (mock *RecipeManifestMock) Description() string {
	args := mock.Called()
	return args.String(0)
}

func (mock *RecipeManifestMock) Icon() string {
	args := mock.Called()
	return args.String(0)
}

func (mock *RecipeManifestMock) Template() string {
	args := mock.Called()
	return args.String(0)
}

func (mock *RecipeManifestMock) Vars() map[string]any {
	args := mock.Called()
	return args.Get(0).(map[string]any)
}

func (mock *RecipeManifestMock) Sync() []syncer.UnitInterface {
	args := mock.Called()
	return args.Get(0).([]syncer.UnitInterface)
}

func (mock *RecipeManifestMock) Schema() schema.Schema {
	args := mock.Called()
	return args.Get(0).(schema.Schema)
}

func (mock *RecipeManifestMock) Options() []RecipeOption {
	args := mock.Called()
	return args.Get(0).([]RecipeOption)
}

// RecipeOption describe a recipe option interface
type RecipeOption interface {
	Name() string
	Label() string
	Help() string
	Path() path.Path
	Schema() schema.Schema
	Validate(value any) (validator.Violations, error)
}

// RecipeOptionMock mock a recipe option
type RecipeOptionMock struct {
	mock.Mock
}

func (mock *RecipeOptionMock) Name() string {
	args := mock.Called()
	return args.String(0)
}

func (mock *RecipeOptionMock) Label() string {
	args := mock.Called()
	return args.String(0)
}

func (mock *RecipeOptionMock) Help() string {
	args := mock.Called()
	return args.String(0)
}

func (mock *RecipeOptionMock) Path() path.Path {
	args := mock.Called()
	return args.Get(0).(path.Path)
}

func (mock *RecipeOptionMock) Schema() schema.Schema {
	args := mock.Called()
	return args.Get(0).(schema.Schema)
}

func (mock *RecipeOptionMock) Validate(value any) (validator.Violations, error) {
	args := mock.Called(value)
	return args.Get(0).(validator.Violations), args.Error(1)
}

// RecipeManager describe a recipe manager interface
type RecipeManager interface {
	LoadRecipe(repository Repository, name string) (Recipe, error)
	RepositoryRecipes(repository Repository) ([]Recipe, error)
	WatchRecipe(recipe Recipe, watcher *watcher.Watcher) error
}

// RecipeManagerMock mock a recipe manager
type RecipeManagerMock struct {
	mock.Mock
}

func (manager *RecipeManagerMock) LoadRecipe(repository Repository, name string) (Recipe, error) {
	args := manager.Called(repository, name)
	return args.Get(0).(Recipe), args.Error(1)
}

func (manager *RecipeManagerMock) RepositoryRecipes(repository Repository) ([]Recipe, error) {
	args := manager.Called(repository)
	return args.Get(0).([]Recipe), args.Error(1)
}

func (manager *RecipeManagerMock) WatchRecipe(recipe Recipe, watcher *watcher.Watcher) error {
	args := manager.Called(recipe, watcher)
	return args.Error(0)
}

/**************/
/* Repository */
/**************/

// Repository describe a repository interface
type Repository interface {
	Url() string
	Dir() string
}

// RepositoryMock mock a repository
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

// RepositoryManager describe a repository manager interface
type RepositoryManager interface {
	LoadRepository(url string) (Repository, error)
}

// RepositoryManagerMock mock a repository manager
type RepositoryManagerMock struct {
	mock.Mock
}

func (mock *RepositoryManagerMock) LoadRepository(url string) (Repository, error) {
	args := mock.Called(url)
	return args.Get(0).(Repository), args.Error(1)
}
