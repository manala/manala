package manifest

import (
	"github.com/stretchr/testify/suite"
	"manala/internal/log"
	"path/filepath"
	"testing"
)

type FinderSuite struct{ suite.Suite }

func TestFinderSuite(t *testing.T) {
	suite.Run(t, new(FinderSuite))
}

func (s *FinderSuite) Test() {
	projectsDir := filepath.FromSlash("testdata/FinderSuite/projects")

	finder := NewFinder(log.Discard)

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
