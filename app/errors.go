package app

/***********/
/* Project */
/***********/

type AlreadyExistingProjectError struct{ Dir string }

func (err *AlreadyExistingProjectError) Error() string   { return "already existing project" }
func (err *AlreadyExistingProjectError) Attrs() [][2]any { return [][2]any{{"dir", err.Dir}} }

type NotFoundProjectError struct{ Dir string }

func (err *NotFoundProjectError) Error() string   { return "project not found" }
func (err *NotFoundProjectError) Attrs() [][2]any { return [][2]any{{"dir", err.Dir}} }

/**********/
/* Recipe */
/**********/

type NotFoundRecipeError struct {
	Repository Repository
	Name       string
}

func (err *NotFoundRecipeError) Error() string { return "recipe not found" }
func (err *NotFoundRecipeError) Attrs() [][2]any {
	return [][2]any{{"repository", err.Repository.URL()}, {"name", err.Name}}
}

/**************/
/* Repository */
/**************/

type NotFoundRepositoryError struct{ URL string }

func (err *NotFoundRepositoryError) Error() string   { return "repository not found" }
func (err *NotFoundRepositoryError) Attrs() [][2]any { return [][2]any{{"url", err.URL}} }

type EmptyRepositoryError struct{ Repository Repository }

func (err *EmptyRepositoryError) Error() string   { return "empty repository" }
func (err *EmptyRepositoryError) Attrs() [][2]any { return [][2]any{{"url", err.Repository.URL()}} }
