package repository

import (
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"io"
	"log/slog"
	"manala/app"
	"manala/internal/serrors"
	"testing"
)

type UrlProcessorManagerSuite struct{ suite.Suite }

func TestUrlProcessorManagerSuite(t *testing.T) {
	suite.Run(t, new(UrlProcessorManagerSuite))
}

func (s *UrlProcessorManagerSuite) TestLoadRepositoryErrors() {
	cascadingMock := &app.RepositoryMock{}
	cascadingError := serrors.New("error")

	cascadingManagerMock := &app.RepositoryManagerMock{}
	cascadingManagerMock.
		On("LoadRepository", mock.Anything).Return(cascadingMock, cascadingError)

	manager := NewUrlProcessorManager(
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		cascadingManagerMock,
	)
	manager.AddUrl("url", -10)

	repository, err := manager.LoadRepository("url")

	s.Equal(cascadingMock, repository)
	s.Equal(cascadingError, err)
}

func (s *UrlProcessorManagerSuite) TestLoadRepository() {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))

	cascadingMock := &app.RepositoryMock{}

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
			cascadingManagerMock := &app.RepositoryManagerMock{}
			cascadingManagerMock.
				On("LoadRepository", mock.Anything).Return(cascadingMock, nil)

			manager := NewUrlProcessorManager(
				log,
				cascadingManagerMock,
			)

			for priority, url := range test.urls {
				manager.AddUrl(url, priority)
			}

			repository, err := manager.LoadRepository(test.url)

			s.NoError(err)
			s.Equal(cascadingMock, repository)

			cascadingManagerMock.AssertCalled(s.T(), "LoadRepository", test.expected)
		})
	}
}

func (s *UrlProcessorManagerSuite) TestLoadPrecedenceRepository() {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))

	cascadingMock := &app.RepositoryMock{}

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
			cascadingManagerMock := &app.RepositoryManagerMock{}
			cascadingManagerMock.
				On("LoadRepository", mock.Anything).Return(cascadingMock, nil)

			manager := NewUrlProcessorManager(
				log,
				cascadingManagerMock,
			)

			for priority, url := range test.urls {
				manager.AddUrl(url, priority)
			}

			repository, err := manager.LoadPrecedingRepository()

			s.NoError(err)

			s.Equal(cascadingMock, repository)

			cascadingManagerMock.AssertCalled(s.T(), "LoadRepository", test.expected)
		})
	}
}

func (s *UrlProcessorManagerSuite) TestProcessUrl() {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))

	cascadingManagerMock := &app.RepositoryManagerMock{}

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
				Type:    &app.UnprocessableRepositoryUrlError{},
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
				Type:    &app.UnprocessableRepositoryUrlError{},
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

			url, err := manager.processUrl(test.url)

			if test.expectedErr != nil {
				s.Empty(url)
				serrors.Equal(s.Assert(), test.expectedErr, err)

			} else {
				s.Equal(test.expectedUrl, url)
				s.NoError(err)
			}
		})
	}
}
