package views

import (
	"github.com/stretchr/testify/suite"
	"manala/app/mocks"
	"testing"
)

type RecipeSuite struct{ suite.Suite }

func TestRecipeSuite(t *testing.T) {
	suite.Run(t, new(RecipeSuite))
}

func (s *RecipeSuite) TestNormalize() {
	repoUrl := "url"

	repoMock := &mocks.RepositoryMock{}
	repoMock.
		On("Url").Return(repoUrl)

	recName := "name"

	recMock := &mocks.RecipeMock{}
	recMock.
		On("Name").Return(recName).
		On("Repository").Return(repoMock)

	recView := NormalizeRecipe(recMock)

	s.Equal(recName, recView.Name)
	s.Equal(repoUrl, recView.Repository.Url)
}
