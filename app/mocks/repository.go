package mocks

import (
	"github.com/stretchr/testify/mock"
	"manala/app/interfaces"
)

func MockRepository() *RepositoryMock {
	return &RepositoryMock{}
}

type RepositoryMock struct {
	mock.Mock
}

func (repo *RepositoryMock) Url() string {
	args := repo.Called()
	return args.String(0)
}

func (repo *RepositoryMock) Dir() string {
	args := repo.Called()
	return args.String(0)
}

/***********/
/* Manager */
/***********/

func MockRepositoryManager() *RepositoryManagerMock {
	return &RepositoryManagerMock{}
}

type RepositoryManagerMock struct {
	mock.Mock
}

func (manager *RepositoryManagerMock) LoadRepository(url string) (interfaces.Repository, error) {
	args := manager.Called(url)
	return args.Get(0).(interfaces.Repository), args.Error(1)
}
