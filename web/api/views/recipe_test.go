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
	recName := "name"
	recDescription := "description"

	recMock := mocks.MockRecipe()
	recMock.
		On("Name").Return(recName).
		On("Description").Return(recDescription)

	recView := NormalizeRecipe(recMock)

	s.Equal(recName, recView.Name)
	s.Equal(recDescription, recView.Description)
}
