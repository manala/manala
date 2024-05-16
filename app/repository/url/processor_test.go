package url

import (
	"github.com/stretchr/testify/suite"
	"manala/internal/log"
	"manala/internal/serrors"
	"testing"
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
		expectedUrl string
		expectedErr *serrors.Assertion
	}{
		{
			test: "QuerySemicolon",
			url:  "foo?bar;baz",
			expectedErr: &serrors.Assertion{
				Message: "unable to process repository query",
				Arguments: []any{
					"query", "bar;baz",
				},
				Errors: []*serrors.Assertion{
					{
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
			expectedUrl: "",
		},
		{
			test: "TailUrl",
			url:  "",
			urls: map[int]string{
				10:  "",
				-10: "tail_url",
			},
			expectedUrl: "tail_url",
		},
		{
			test: "UrlAndTailUrl",
			url:  "url",
			urls: map[int]string{
				10:  "",
				-10: "tail_url",
			},
			expectedUrl: "url",
		},
		{
			test: "All",
			url:  "url",
			urls: map[int]string{
				10:  "head_url",
				-10: "tail_url",
			},
			expectedUrl: "head_url",
		},
		{
			test: "HeadUrlAndTailUrl",
			url:  "",
			urls: map[int]string{
				10:  "head_url",
				-10: "tail_url",
			},
			expectedUrl: "head_url",
		},
		{
			test:        "Windows",
			url:         `foo\bar`,
			expectedUrl: `foo\bar`,
		},
		{
			test:        "Query",
			url:         "?query=query",
			expectedUrl: "",
		},
		{
			test: "QueryAndTailUrl",
			url:  "?query=query",
			urls: map[int]string{
				10:  "",
				-10: "tail_url",
			},
			expectedUrl: "tail_url?query=query",
		},
		{
			test: "UrlAndHeadQueryAndTailUrl",
			url:  "url",
			urls: map[int]string{
				-10: "tail_url",
			},
			queries: map[int]map[string]string{
				10: {"head_query": "head_query"},
			},
			expectedUrl: "url?head_query=head_query",
		},
		{
			test: "UrlQueryAndHeadUrlQueryAndTailUrl",
			url:  "url?query=query",
			urls: map[int]string{
				10:  "head_url?head_query=head_query",
				-10: "tail_url",
			},
			expectedUrl: "head_url?head_query=head_query",
		},
		{
			test: "QueryAndHeadQueryAndTailUrlQuery",
			url:  "?query=query",
			urls: map[int]string{
				-10: "tail_url?tail_query=tail_query",
			},
			queries: map[int]map[string]string{
				10: {"head_query": "head_query"},
			},
			expectedUrl: "tail_url?head_query=head_query&query=query&tail_query=tail_query",
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			processor := NewProcessor(log.Discard)

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
				s.Empty(url)
				serrors.Equal(s.T(), test.expectedErr, err)
			} else {
				s.Equal(test.expectedUrl, url)
				s.NoError(err)
			}
		})
	}
}
