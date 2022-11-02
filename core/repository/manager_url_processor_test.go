package repository

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"io"
	"manala/core"
	internalLog "manala/internal/log"
	"testing"
)

type UrlProcessorManagerSuite struct{ suite.Suite }

func TestUrlProcessorManagerSuite(t *testing.T) {
	suite.Run(t, new(UrlProcessorManagerSuite))
}

func (s *UrlProcessorManagerSuite) TestLoadRepositoryErrors() {
	log := internalLog.New(io.Discard)

	cascadingRepoMock := core.NewRepositoryMock()
	cascadingError := errors.New("error")

	cascadingManagerMock := core.NewRepositoryManagerMock()
	cascadingManagerMock.
		On("LoadRepository", mock.Anything).Return(cascadingRepoMock, cascadingError)

	manager := NewUrlProcessorManager(
		log,
		cascadingManagerMock,
	)
	manager.WithLowermostUrl("url")

	repo, err := manager.LoadRepository("url")

	s.Nil(repo)
	s.ErrorIs(err, cascadingError)
}

func (s *UrlProcessorManagerSuite) TestLoadRepository() {
	log := internalLog.New(io.Discard)

	cascadingRepoMock := core.NewRepositoryMock()

	tests := []struct {
		name                 string
		lowermostUrl         string
		url                  string
		uppermostUrl         string
		expectedCascadingUrl string
	}{
		{
			name:                 "Lowermost Url Only",
			lowermostUrl:         "lowermost_url",
			url:                  "",
			uppermostUrl:         "",
			expectedCascadingUrl: "lowermost_url",
		},
		{
			name:                 "Lowermost Url And Url",
			lowermostUrl:         "lowermost_url",
			url:                  "url",
			uppermostUrl:         "",
			expectedCascadingUrl: "url",
		},
		{
			name:                 "Lowermost Url And Url And Uppermost Url",
			lowermostUrl:         "lowermost_url",
			url:                  "url",
			uppermostUrl:         "uppermost_url",
			expectedCascadingUrl: "uppermost_url",
		},
		{
			name:                 "Lowermost Url And Uppermost Url",
			lowermostUrl:         "lowermost_url",
			url:                  "",
			uppermostUrl:         "uppermost_url",
			expectedCascadingUrl: "uppermost_url",
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			cascadingManagerMock := core.NewRepositoryManagerMock()
			cascadingManagerMock.
				On("LoadRepository", mock.Anything).Return(cascadingRepoMock, nil)

			manager := NewUrlProcessorManager(
				log,
				cascadingManagerMock,
			)
			manager.WithLowermostUrl(test.lowermostUrl)
			manager.WithUppermostUrl(test.uppermostUrl)

			repo, err := manager.LoadRepository(test.url)

			s.NoError(err)
			s.Equal(cascadingRepoMock, repo)

			cascadingManagerMock.AssertCalled(s.T(), "LoadRepository", test.expectedCascadingUrl)
		})
	}
}

func (s *UrlProcessorManagerSuite) TestLoadPrecedenceRepository() {
	log := internalLog.New(io.Discard)

	cascadingRepoMock := core.NewRepositoryMock()

	tests := []struct {
		name                 string
		lowermostUrl         string
		uppermostUrl         string
		expectedCascadingUrl string
	}{
		{
			name:                 "Lowermost Url",
			lowermostUrl:         "lowermost_url",
			uppermostUrl:         "",
			expectedCascadingUrl: "lowermost_url",
		},
		{
			name:                 "Lowermost Url And Uppermost Url",
			lowermostUrl:         "lowermost_url",
			uppermostUrl:         "uppermost_url",
			expectedCascadingUrl: "uppermost_url",
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			cascadingManagerMock := core.NewRepositoryManagerMock()
			cascadingManagerMock.
				On("LoadRepository", mock.Anything).Return(cascadingRepoMock, nil)

			manager := NewUrlProcessorManager(
				log,
				cascadingManagerMock,
			)
			manager.WithLowermostUrl(test.lowermostUrl)
			manager.WithUppermostUrl(test.uppermostUrl)

			repo, err := manager.LoadPrecedingRepository()

			s.NoError(err)
			s.Equal(cascadingRepoMock, repo)

			cascadingManagerMock.AssertCalled(s.T(), "LoadRepository", test.expectedCascadingUrl)
		})
	}
}

func (s *UrlProcessorManagerSuite) TestProcessUrl() {
	tests := []struct {
		lowermostUrl string
		url          string
		uppermostUrl string
		expected     string
		error        bool
	}{
		{
			lowermostUrl: "",
			url:          "",
			uppermostUrl: "",
			error:        true,
		},
		{
			lowermostUrl: "lower",
			url:          "",
			uppermostUrl: "",
			expected:     "lower",
		},
		{
			lowermostUrl: "",
			url:          "url",
			uppermostUrl: "",
			expected:     "url",
		},
		{
			lowermostUrl: "lower",
			url:          "url",
			uppermostUrl: "",
			expected:     "url",
		},
		{
			lowermostUrl: "",
			url:          "",
			uppermostUrl: "upper",
			expected:     "upper",
		},
		{
			lowermostUrl: "lower",
			url:          "",
			uppermostUrl: "upper",
			expected:     "upper",
		},
		{
			lowermostUrl: "",
			url:          "url",
			uppermostUrl: "upper",
			expected:     "upper",
		},
		{
			lowermostUrl: "lower",
			url:          "url",
			uppermostUrl: "upper",
			expected:     "upper",
		},
	}

	for i, test := range tests {
		s.Run(fmt.Sprint(i), func() {
			manager := &UrlProcessorManager{
				lowermostUrl: test.lowermostUrl,
				uppermostUrl: test.uppermostUrl,
			}

			actual, err := manager.processUrl(test.url)

			if test.error {
				var _error *core.UnprocessableRepositoryUrlError
				s.ErrorAs(err, &_error)
				s.Empty(actual)
			} else {
				s.NoError(err)
				s.Equal(test.expected, actual)
			}
		})
	}
}
