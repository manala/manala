package app

import "context"

type repositoryURLKey struct{}

func WithRepositoryURL(ctx context.Context, repositoryURL string) context.Context {
	if repositoryURL == "" {
		return ctx
	}

	return context.WithValue(ctx, repositoryURLKey{}, repositoryURL)
}

// RepositoryURL gets the repository URL from the context.
func RepositoryURL(ctx context.Context) (string, bool) {
	value, ok := ctx.Value(repositoryURLKey{}).(string)

	return value, ok
}

type repositoryRefKey struct{}

func WithRepositoryRef(ctx context.Context, repositoryURL string) context.Context {
	if repositoryURL == "" {
		return ctx
	}

	return context.WithValue(ctx, repositoryRefKey{}, repositoryURL)
}

// RepositoryRef gets the repository ref from the context.
func RepositoryRef(ctx context.Context) (string, bool) {
	value, ok := ctx.Value(repositoryRefKey{}).(string)

	return value, ok
}

type recipeNameKey struct{}

func WithRecipeName(ctx context.Context, recipeName string) context.Context {
	if recipeName == "" {
		return ctx
	}

	return context.WithValue(ctx, recipeNameKey{}, recipeName)
}

// RecipeName gets the recipe name from the context.
func RecipeName(ctx context.Context) (string, bool) {
	value, ok := ctx.Value(recipeNameKey{}).(string)

	return value, ok
}
