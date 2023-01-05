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
	manager.AddUrl("url", -10)

	repo, err := manager.LoadRepository("url")

	s.Nil(repo)
	s.ErrorIs(err, cascadingError)
}

func (s *UrlProcessorManagerSuite) TestLoadRepository() {
	log := internalLog.New(io.Discard)

	cascadingRepoMock := core.NewRepositoryMock()

	tests := []struct {
		name                 string
		url                  string
		urls                 map[int]string
		expectedCascadingUrl string
	}{
		{
			name: "Lowermost Url Only",
			url:  "",
			urls: map[int]string{
				-10: "lowermost_url",
				10:  "",
			},
			expectedCascadingUrl: "lowermost_url",
		},
		{
			name: "Lowermost Url And Url",
			url:  "url",
			urls: map[int]string{
				-10: "lowermost_url",
				10:  "",
			},
			expectedCascadingUrl: "url",
		},
		{
			name: "Lowermost Url And Url And Uppermost Url",
			url:  "url",
			urls: map[int]string{
				-10: "lowermost_url",
				10:  "uppermost_url",
			},
			expectedCascadingUrl: "uppermost_url",
		},
		{
			name: "Lowermost Url And Uppermost Url",
			url:  "",
			urls: map[int]string{
				-10: "lowermost_url",
				10:  "uppermost_url",
			},
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

			for priority, url := range test.urls {
				manager.AddUrl(url, priority)
			}

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
		urls                 map[int]string
		expectedCascadingUrl string
	}{
		{
			name: "Lowermost Url",
			urls: map[int]string{
				-10: "lowermost_url",
				10:  "",
			},
			expectedCascadingUrl: "lowermost_url",
		},
		{
			name: "Lowermost Url And Uppermost Url",
			urls: map[int]string{
				-10: "lowermost_url",
				10:  "uppermost_url",
			},
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

			for priority, url := range test.urls {
				manager.AddUrl(url, priority)
			}

			repo, err := manager.LoadPrecedingRepository()

			s.NoError(err)
			s.Equal(cascadingRepoMock, repo)

			cascadingManagerMock.AssertCalled(s.T(), "LoadRepository", test.expectedCascadingUrl)
		})
	}
}

func (s *UrlProcessorManagerSuite) TestProcessUrl() {
	log := internalLog.New(io.Discard)

	cascadingManagerMock := core.NewRepositoryManagerMock()

	tests := []struct {
		url      string
		urls     map[int]string
		expected string
		error    bool
	}{
		{
			url: "",
			urls: map[int]string{
				-10: "",
				10:  "",
			},
			error: true,
		},
		{
			url: "",
			urls: map[int]string{
				-10: "lower",
				10:  "",
			},
			expected: "lower",
		},
		{
			url: "url",
			urls: map[int]string{
				-10: "lower",
				10:  "",
			},
			expected: "url",
		},
		{
			url: "url",
			urls: map[int]string{
				-10: "lower",
				10:  "upper",
			},
			expected: "upper",
		},
		// Windows
		{
			url: "",
			urls: map[int]string{
				-10: "",
				10:  `foo\bar`,
			},
			expected: `foo\bar`,
		},
		// Query
		{
			url: "?url=url",
			urls: map[int]string{
				-10: "",
				10:  "",
			},
			error: true,
		},
		{
			url: "?url=url",
			urls: map[int]string{
				-10: "lower",
				10:  "",
			},
			expected: "lower?url=url",
		},
		{
			url: "url",
			urls: map[int]string{
				-10: "lower",
				10:  "?upper=upper",
			},
			expected: "url?upper=upper",
		},
		{
			url: "url?url=url",
			urls: map[int]string{
				-10: "lower",
				10:  "upper?upper=upper",
			},
			expected: "upper?upper=upper",
		},
		{
			url: "?url=url",
			urls: map[int]string{
				-10: "lower?lower=lower",
				10:  "?upper=upper",
			},
			expected: "lower?lower=lower&upper=upper&url=url",
		},
	}

	for i, test := range tests {
		s.Run(fmt.Sprint(i), func() {
			manager := NewUrlProcessorManager(
				log,
				cascadingManagerMock,
			)

			for priority, url := range test.urls {
				manager.AddUrl(url, priority)
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
