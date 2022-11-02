package application

type Option func(app *Application)

func WithRepositoryUrl(url string) Option {
	return func(app *Application) {
		app.repositoryManager.WithUppermostUrl(url)
	}
}

func WithRecipeName(name string) Option {
	return func(app *Application) {
		app.recipeManager.WithUppermostName(name)
	}
}
