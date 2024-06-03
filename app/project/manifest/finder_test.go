package manifest_test

import (
	"manala/app/project/manifest"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
)

type FinderSuite struct{ suite.Suite }

func TestFinderSuite(t *testing.T) {
	suite.Run(t, new(FinderSuite))
}

func (s *FinderSuite) Test() {
	projectsDir := filepath.FromSlash("testdata/FinderSuite/projects")

	finder := manifest.NewFinder()

	s.Run("Found", func() {
		s.True(finder.Find(
			filepath.Join(projectsDir, "found"),
		))
	})

	s.Run("NotFound", func() {
		s.False(finder.Find(
			filepath.Join(projectsDir, "not_found"),
		))
	})
}
