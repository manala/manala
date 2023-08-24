package application

type Option func(app *Application)

func WithRepositoryUrl(url string) Option {
	return func(app *Application) {
		priority := 10

		// Log
		app.log.Debug("app repository option",
			"url", url,
			"priority", priority,
		)

		app.repositoryManager.AddUrl(url, priority)
	}
}

func WithRepositoryRef(ref string) Option {
	return func(app *Application) {
		priority := 20

		// Log
		app.log.Debug("app repository option",
			"ref", ref,
			"priority", priority,
		)

		app.repositoryManager.AddUrlQuery("ref", ref, priority)
	}
}

func WithRecipeName(name string) Option {
	return func(app *Application) {
		priority := 10

		// Log
		app.log.Debug("app recipe option",
			"name", name,
			"priority", priority,
		)

		app.recipeManager.AddName(name, priority)
	}
}
