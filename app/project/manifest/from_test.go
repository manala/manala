package manifest_test

import (
	"path/filepath"
	"testing"

	"github.com/manala/manala/app"
	"github.com/manala/manala/app/project"
	"github.com/manala/manala/app/project/manifest"
	"github.com/manala/manala/app/testing/mocks"
	"github.com/manala/manala/internal/errors/serror/serrortest"
	"github.com/manala/manala/internal/log"
	"github.com/manala/manala/internal/testing/expectation"

	"github.com/stretchr/testify/suite"
)

type FromSuite struct{ suite.Suite }

func TestFromSuite(t *testing.T) {
	suite.Run(t, new(FromSuite))
}

func (s *FromSuite) TestHandleErrors() {
	dir := filepath.FromSlash("testdata/FromSuite/TestHandleErrors")

	tests := []struct {
		test     string
		expected expectation.ErrorExpectation
	}{
		{
			test: "NotExists",
			expected: serrortest.Expectation{
				Msg: "project from dir does not exist",
				Attrs: [][2]any{
					{"dir", filepath.Join(dir, "NotExists", "project")},
				},
			},
		},
		{
			test: "File",
			expected: serrortest.Expectation{
				Msg: "project from dir is not a dir",
				Attrs: [][2]any{
					{"dir", filepath.Join(dir, "File", "project")},
				},
			},
		},
	}

	for _, test := range tests {
		s.Run(test.test, func() {
			_, err := s.handle(filepath.Join(dir, test.test, "project"))

			expectation.ExpectError(s.T(), test.expected, err)
		})
	}
}

func (s *FromSuite) handle(dir string) (app.Project, error) {
	query := &project.LoaderQuery{Dir: dir}

	projectMock := &mocks.Project{}

	chainMock := &project.LoaderHandlerChainMock{}
	chainMock.
		On("Next", query).Return(projectMock, nil)

	handler := manifest.NewFromLoaderHandler(log.Discard)

	return handler.Handle(query, chainMock)
}
