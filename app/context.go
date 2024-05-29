package app

import "context"

type repositoryUrlKey struct{}

func WithRepositoryUrl(ctx context.Context, repositoryUrl string) context.Context {
	if repositoryUrl == "" {
		return ctx
	}
	return context.WithValue(ctx, repositoryUrlKey{}, repositoryUrl)
}

// RepositoryUrl gets the repository url from the context.
func RepositoryUrl(ctx context.Context) (string, bool) {
	value, ok := ctx.Value(repositoryUrlKey{}).(string)
	return value, ok
}

type repositoryRefKey struct{}

func WithRepositoryRef(ctx context.Context, repositoryUrl string) context.Context {
	if repositoryUrl == "" {
		return ctx
	}
	return context.WithValue(ctx, repositoryRefKey{}, repositoryUrl)
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
