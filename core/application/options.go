package application

type Option func(app *Application)

func WithRepositoryUrl(url string) Option {
	return func(app *Application) {
		priority := 10

		app.log.
			WithField("url", url).
			WithField("priority", priority).
			Debug("option repository")

		app.repositoryManager.AddUrl(url, priority)
	}
}

func WithRepositoryRef(ref string) Option {
	return func(app *Application) {
		priority := 20

		app.log.
			WithField("ref", ref).
			WithField("priority", priority).
			Debug("option repository")

		app.repositoryManager.AddUrlQuery("ref", ref, priority)
	}
}

func WithRecipeName(name string) Option {
	return func(app *Application) {
		priority := 10

		app.log.
			WithField("name", name).
			WithField("priority", priority).
			Debug("option recipe")

		app.recipeManager.AddName(name, priority)
	}
}
