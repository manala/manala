package backwalk

import (
	"github.com/stretchr/testify/suite"
	"os"
	"path/filepath"
	"testing"
)

type BackwalkSuite struct{ suite.Suite }

func TestBackwalkSuite(t *testing.T) {
	suite.Run(t, new(BackwalkSuite))
}

func (s *BackwalkSuite) Test() {
	dir := filepath.FromSlash("testdata/BackwalkSuite/Test")

	i := 0
	err := Backwalk(
		filepath.Join(dir, "foo", "bar"),
		func(path string, entry os.DirEntry, err error) error {
			s.NoError(err)
			s.Equal(
				[]string{
					filepath.Join(dir, "foo", "bar"),
					filepath.Join(dir, "foo"),
					dir,
				}[i],
				path,
			)
			i = i + 1
			if path == dir {
				return filepath.SkipDir
			}
			return nil
		},
	)
	s.NoError(err)
}
