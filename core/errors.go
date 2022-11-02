package core

/***********/
/* Project */
/***********/

func NewNotFoundProjectManifestError(message string) *NotFoundProjectManifestError {
	return &NotFoundProjectManifestError{
		message: message,
	}
}

type NotFoundProjectManifestError struct {
	message string
}

func (err *NotFoundProjectManifestError) Error() string {
	return err.message
}

/**********/
/* Recipe */
/**********/

func NewNotFoundRecipeManifestError(message string) *NotFoundRecipeManifestError {
	return &NotFoundRecipeManifestError{
		message: message,
	}
}

type NotFoundRecipeManifestError struct {
	message string
}

func (err *NotFoundRecipeManifestError) Error() string {
	return err.message
}

func NewUnprocessableRecipeNameError(message string) *UnprocessableRecipeNameError {
	return &UnprocessableRecipeNameError{
		message: message,
	}
}

type UnprocessableRecipeNameError struct {
	message string
}

func (err *UnprocessableRecipeNameError) Error() string {
	return err.message
}

/**************/
/* Repository */
/**************/

func NewNotFoundRepositoryError(message string) *NotFoundRepositoryError {
	return &NotFoundRepositoryError{
		message: message,
	}
}

type NotFoundRepositoryError struct {
	message string
}

func (err *NotFoundRepositoryError) Error() string {
	return err.message
}

func NewUnsupportedRepositoryError(message string) *UnsupportedRepositoryError {
	return &UnsupportedRepositoryError{
		message: message,
	}
}

type UnsupportedRepositoryError struct {
	message string
}

func (err *UnsupportedRepositoryError) Error() string {
	return err.message
}

func NewUnprocessableRepositoryUrlError(message string) *UnprocessableRepositoryUrlError {
	return &UnprocessableRepositoryUrlError{
		message: message,
	}
}

type UnprocessableRepositoryUrlError struct {
	message string
}

func (err *UnprocessableRepositoryUrlError) Error() string {
	return err.message
}
