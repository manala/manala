package views

import (
	"github.com/stretchr/testify/suite"
	"manala/app/mocks"
	"testing"
)

type ProjectSuite struct{ suite.Suite }

func TestProjectSuite(t *testing.T) {
	suite.Run(t, new(ProjectSuite))
}

func (s *ProjectSuite) TestNormalize() {
	repoUrl := "url"

	repoMock := mocks.MockRepository()
	repoMock.
		On("Url").Return(repoUrl)

	recName := "name"

	recMock := mocks.MockRecipe()
	recMock.
		On("Name").Return(recName).
		On("Repository").Return(repoMock)

	projVars := map[string]interface{}{
		"foo": "bar",
	}

	projMock := mocks.MockProject()
	projMock.
		On("Vars").Return(projVars).
		On("Recipe").Return(recMock)

	projView := NormalizeProject(projMock)

	s.Equal(projVars, projView.Vars)
	s.Equal(recName, projView.Recipe.Name)
	s.Equal(repoUrl, projView.Recipe.Repository.Url)
}
