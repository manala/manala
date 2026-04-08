package url_test

import (
	stderrors "errors"
	"log/slog"
	"testing"

	"github.com/manala/manala/app/repository/url"
	"github.com/manala/manala/internal/serrors"
	"github.com/manala/manala/internal/testing/errors"

	"github.com/stretchr/testify/suite"
)

type ProcessorSuite struct{ suite.Suite }

func TestProcessorSuite(t *testing.T) {
	suite.Run(t, new(ProcessorSuite))
}

func (s *ProcessorSuite) TestProcess() {
	tests := []struct {
		test        string
		url         string
		urls        map[int]string
		queries     map[int]map[string]string
		expectedURL string
		expectedErr errors.Assertion
	}{
		{
			test: "QuerySemicolon",
			url:  "foo?bar;baz",
			expectedErr: &serrors.Assertion{
				Message: "unable to process repository query",
				Arguments: []any{
					"query", "bar;baz",
				},
				Errors: []errors.Assertion{
					&serrors.Assertion{
						Type:    stderrors.New(""),
						Message: "invalid semicolon separator in query",
					},
				},
			},
		},
		{
			test: "Empty",
			url:  "",
			urls: map[int]string{
				10:  "",
				-10: "",
			},
			expectedURL: "",
		},
		{
			test: "TailURL",
			url:  "",
			urls: map[int]string{
				10:  "",
				-10: "tail_url",
			},
			expectedURL: "tail_url",
		},
		{
			test: "URLAndTailURL",
			url:  "url",
			urls: map[int]string{
				10:  "",
				-10: "tail_url",
			},
			expectedURL: "url",
		},
		{
			test: "All",
			url:  "url",
			urls: map[int]string{
				10:  "head_url",
				-10: "tail_url",
			},
			expectedURL: "head_url",
		},
		{
			test: "HeadURLAndTailURL",
			url:  "",
			urls: map[int]string{
				10:  "head_url",
				-10: "tail_url",
			},
			expectedURL: "head_url",
		},
		{
			test:        "Windows",
			url:         `foo\bar`,
			expectedURL: `foo\bar`,
		},
		{
			test:        "Query",
			url:         "?query=query",
			expectedURL: "",
		},
		{
			test: "QueryAndTailURL",
			url:  "?query=query",
			urls: map[int]string{
				10:  "",
				-10: "tail_url",
			},
			expectedURL: "tail_url?query=query",
		},
		{
			test: "URLAndHeadQueryAndTailURL",
			url:  "url",
			urls: map[int]string{
				-10: "tail_url",
			},
			queries: map[int]map[string]string{
				10: {"head_query": "head_query"},
			},
			expectedURL: "url?head_query=head_query",
		},
		{
			test: "URLQueryAndHeadURLQueryAndTailURL",
			url:  "url?query=query",
			urls: map[int]string{
				10:  "head_url?head_query=head_query",
				-10: "tail_url",
			},
			expectedURL: "head_url?head_query=head_query",
		},
		{
			test: "QueryAndHeadQueryAndTailURLQuery",
			url:  "?query=query",
			urls: map[int]string{
				-10: "tail_url?tail_query=tail_query",
			},
			queries: map[int]map[string]string{
				10: {"head_query": "head_query"},
			},
			expectedURL: "tail_url?head_query=head_query&query=query&tail_query=tail_query",
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			processor := url.NewProcessor(slog.New(slog.DiscardHandler))

			for weight, url := range test.urls {
				processor.Add(url, weight)
			}

			for weight, queries := range test.queries {
				for key, value := range queries {
					processor.AddQuery(key, value, weight)
				}
			}

			url, err := processor.Process(test.url)

			if test.expectedErr != nil {
				errors.Equal(s.T(), test.expectedErr, err)
				s.Empty(url)
			} else {
				s.Require().NoError(err)
				s.Equal(test.expectedURL, url)
			}
		})
	}
}
