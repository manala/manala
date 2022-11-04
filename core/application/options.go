package application

type Option func(app *Application)

func WithRepositoryUrl(url string) Option {
	return func(app *Application) {
		app.repositoryManager.AddUrl(url, 10)
	}
}

func WithRecipeName(name string) Option {
	return func(app *Application) {
		app.recipeManager.AddName(name, 10)
	}
}
