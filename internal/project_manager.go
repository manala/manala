package internal

import (
	internalLog "manala/internal/log"
)

func NewProjectManager(
	logger *internalLog.Logger,
	repositoryLoader RepositoryLoaderInterface,
	projectLoader ProjectLoaderInterface,
) *ProjectManager {
	return &ProjectManager{
		Log:              logger,
		RepositoryLoader: repositoryLoader,
		ProjectLoader:    projectLoader,
	}
}

type ProjectManager struct {
	Log              *internalLog.Logger
	RepositoryLoader RepositoryLoaderInterface
	ProjectLoader    ProjectLoaderInterface
}

func (manager *ProjectManager) LoadProjectManifest(path string) (*ProjectManifest, error) {
	return manager.ProjectLoader.LoadProjectManifest(path)
}

func (manager *ProjectManager) LoadProject(path string, repositoryPath string, recipeName string) (*Project, error) {
	return manager.ProjectLoader.LoadProject(path, repositoryPath, recipeName)
}
