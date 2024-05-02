package app

/***********/
/* Project */
/***********/

type AlreadyExistingProjectError struct{ Dir string }

func (err *AlreadyExistingProjectError) Error() string         { return "already existing project" }
func (err *AlreadyExistingProjectError) ErrorArguments() []any { return []any{"dir", err.Dir} }

type NotFoundProjectError struct{ Dir string }

func (err *NotFoundProjectError) Error() string         { return "project not found" }
func (err *NotFoundProjectError) ErrorArguments() []any { return []any{"dir", err.Dir} }

/**********/
/* Recipe */
/**********/

type NotFoundRecipeError struct {
	Repository Repository
	Name       string
}

func (err *NotFoundRecipeError) Error() string { return "recipe not found" }
func (err *NotFoundRecipeError) ErrorArguments() []any {
	return []any{"repository", err.Repository.Url(), "name", err.Name}
}

/**************/
/* Repository */
/**************/

type NotFoundRepositoryError struct{ Url string }

func (err *NotFoundRepositoryError) Error() string         { return "repository not found" }
func (err *NotFoundRepositoryError) ErrorArguments() []any { return []any{"url", err.Url} }

type UnsupportedRepositoryError struct{ Url string }

func (err *UnsupportedRepositoryError) Error() string         { return "unsupported repository url" }
func (err *UnsupportedRepositoryError) ErrorArguments() []any { return []any{"url", err.Url} }

type EmptyRepositoryError struct{ Repository Repository }

func (err *EmptyRepositoryError) Error() string         { return "empty repository" }
func (err *EmptyRepositoryError) ErrorArguments() []any { return []any{"url", err.Repository.Url()} }
