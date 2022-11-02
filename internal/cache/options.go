package cache

type Option func(cache *Cache)

func WithUserDir(dir string) Option {
	return func(app *Cache) {
		app.userDir = dir
	}
}
