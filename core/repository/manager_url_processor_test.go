package repository

import (
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"io"
	"log/slog"
	"manala/app/mocks"
	"manala/core"
	"manala/internal/errors/serrors"
	"testing"
)

type UrlProcessorManagerSuite struct{ suite.Suite }

func TestUrlProcessorManagerSuite(t *testing.T) {
	suite.Run(t, new(UrlProcessorManagerSuite))
}

func (s *UrlProcessorManagerSuite) TestLoadRepositoryErrors() {
	cascadingRepoMock := &mocks.RepositoryMock{}
	cascadingError := serrors.New("error")

	cascadingManagerMock := &mocks.RepositoryManagerMock{}
	cascadingManagerMock.
		On("LoadRepository", mock.Anything).Return(cascadingRepoMock, cascadingError)

	manager := NewUrlProcessorManager(
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		cascadingManagerMock,
	)
	manager.AddUrl("url", -10)

	repo, err := manager.LoadRepository("url")

	s.Nil(repo)

	s.ErrorIs(err, cascadingError)
}

func (s *UrlProcessorManagerSuite) TestLoadRepository() {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))

	cascadingRepoMock := &mocks.RepositoryMock{}

	tests := []struct {
		test     string
		url      string
		urls     map[int]string
		expected string
	}{
		{
			test: "LowermostUrlOnly",
			url:  "",
			urls: map[int]string{
				-10: "lowermost_url",
				10:  "",
			},
			expected: "lowermost_url",
		},
		{
			test: "LowermostUrlAndUrl",
			url:  "url",
			urls: map[int]string{
				-10: "lowermost_url",
				10:  "",
			},
			expected: "url",
		},
		{
			test: "LowermostUrlAndUrlAndUppermostUrl",
			url:  "url",
			urls: map[int]string{
				-10: "lowermost_url",
				10:  "uppermost_url",
			},
			expected: "uppermost_url",
		},
		{
			test: "LowermostUrlAndUppermostUrl",
			url:  "",
			urls: map[int]string{
				-10: "lowermost_url",
				10:  "uppermost_url",
			},
			expected: "uppermost_url",
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			cascadingManagerMock := &mocks.RepositoryManagerMock{}
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

			cascadingManagerMock.AssertCalled(s.T(), "LoadRepository", test.expected)
		})
	}
}

func (s *UrlProcessorManagerSuite) TestLoadPrecedenceRepository() {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))

	cascadingRepoMock := &mocks.RepositoryMock{}

	tests := []struct {
		test     string
		urls     map[int]string
		expected string
	}{
		{
			test: "LowermostUrl",
			urls: map[int]string{
				-10: "lowermost_url",
				10:  "",
			},
			expected: "lowermost_url",
		},
		{
			test: "LowermostUrlAndUppermostUrl",
			urls: map[int]string{
				-10: "lowermost_url",
				10:  "uppermost_url",
			},
			expected: "uppermost_url",
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			cascadingManagerMock := &mocks.RepositoryManagerMock{}
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

			cascadingManagerMock.AssertCalled(s.T(), "LoadRepository", test.expected)
		})
	}
}

func (s *UrlProcessorManagerSuite) TestProcessUrl() {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))

	cascadingManagerMock := &mocks.RepositoryManagerMock{}

	tests := []struct {
		test        string
		url         string
		urls        map[int]string
		expectedUrl string
		expectedErr *serrors.Assert
	}{
		{
			test: "1",
			url:  "",
			urls: map[int]string{
				-10: "",
				10:  "",
			},
			expectedErr: &serrors.Assert{
				Type:    &core.UnprocessableRepositoryUrlError{},
				Message: "unable to process repository url",
			},
		},
		{
			test: "2",
			url:  "",
			urls: map[int]string{
				-10: "lower",
				10:  "",
			},
			expectedUrl: "lower",
		},
		{
			test: "3",
			url:  "url",
			urls: map[int]string{
				-10: "lower",
				10:  "",
			},
			expectedUrl: "url",
		},
		{
			test: "4",
			url:  "url",
			urls: map[int]string{
				-10: "lower",
				10:  "upper",
			},
			expectedUrl: "upper",
		},
		// Windows
		{
			test: "5",
			url:  "",
			urls: map[int]string{
				-10: "",
				10:  `foo\bar`,
			},
			expectedUrl: `foo\bar`,
		},
		// Query
		{
			test: "6",
			url:  "?url=url",
			urls: map[int]string{
				-10: "",
				10:  "",
			},
			expectedErr: &serrors.Assert{
				Type:    &core.UnprocessableRepositoryUrlError{},
				Message: "unable to process repository url",
			},
		},
		{
			test: "7",
			url:  "?url=url",
			urls: map[int]string{
				-10: "lower",
				10:  "",
			},
			expectedUrl: "lower?url=url",
		},
		{
			test: "8",
			url:  "url",
			urls: map[int]string{
				-10: "lower",
				10:  "?upper=upper",
			},
			expectedUrl: "url?upper=upper",
		},
		{
			test: "9",
			url:  "url?url=url",
			urls: map[int]string{
				-10: "lower",
				10:  "upper?upper=upper",
			},
			expectedUrl: "upper?upper=upper",
		},
		{
			test: "10",
			url:  "?url=url",
			urls: map[int]string{
				-10: "lower?lower=lower",
				10:  "?upper=upper",
			},
			expectedUrl: "lower?lower=lower&upper=upper&url=url",
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			manager := NewUrlProcessorManager(
				log,
				cascadingManagerMock,
			)

			for priority, url := range test.urls {
				manager.AddUrl(url, priority)
			}

			actual, err := manager.processUrl(test.url)

			if test.expectedErr != nil {
				s.Empty(actual)
				serrors.Equal(s.Assert(), test.expectedErr, err)

			} else {
				s.Equal(test.expectedUrl, actual)
				s.NoError(err)
			}
		})
	}
}
