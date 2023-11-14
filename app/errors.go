package app

/***********/
/* Project */
/***********/

type AlreadyExistingProjectError struct{ Dir string }

func (err *AlreadyExistingProjectError) Error() string         { return "already existing project" }
func (err *AlreadyExistingProjectError) ErrorArguments() []any { return []any{"dir", err.Dir} }

type NotFoundProjectManifestError struct{ File string }

func (err *NotFoundProjectManifestError) Error() string         { return "project manifest not found" }
func (err *NotFoundProjectManifestError) ErrorArguments() []any { return []any{"file", err.File} }

/**********/
/* Recipe */
/**********/

type NotFoundRecipeManifestError struct{ File string }

func (err *NotFoundRecipeManifestError) Error() string         { return "recipe manifest not found" }
func (err *NotFoundRecipeManifestError) ErrorArguments() []any { return []any{"file", err.File} }

type UnprocessableRecipeNameError struct{}

func (err *UnprocessableRecipeNameError) Error() string { return "unable to process recipe name" }

/**************/
/* Repository */
/**************/

type UnsupportedRepositoryError struct{ Url string }

func (err *UnsupportedRepositoryError) Error() string         { return "unsupported repository url" }
func (err *UnsupportedRepositoryError) ErrorArguments() []any { return []any{"url", err.Url} }

type UnprocessableRepositoryUrlError struct{}

func (err *UnprocessableRepositoryUrlError) Error() string { return "unable to process repository url" }

type EmptyRepositoryError struct{ Repository Repository }

func (err *EmptyRepositoryError) Error() string         { return "empty repository" }
func (err *EmptyRepositoryError) ErrorArguments() []any { return []any{"url", err.Repository.Url()} }
