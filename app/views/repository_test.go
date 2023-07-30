package views

import (
	"github.com/stretchr/testify/suite"
	"manala/app/mocks"
	"testing"
)

type RepositorySuite struct{ suite.Suite }

func TestRepositorySuite(t *testing.T) {
	suite.Run(t, new(RepositorySuite))
}

func (s *RepositorySuite) TestNormalize() {
	repoUrl := "url"

	repoMock := &mocks.RepositoryMock{}
	repoMock.
		On("Url").Return(repoUrl)

	repoView := NormalizeRepository(repoMock)

	s.Equal(repoUrl, repoView.Url)
	s.Equal(repoUrl, repoView.Path)
	s.Equal(repoUrl, repoView.Source)
}
