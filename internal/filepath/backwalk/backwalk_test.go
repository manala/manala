package backwalk

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
)

type Suite struct{ suite.Suite }

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) Test() {
	dir := filepath.FromSlash("testdata/Test")

	i := 0
	err := WalkDir(
		filepath.Join(dir, "foo", "bar"),
		func(path string, _ os.DirEntry, err error) error {
			s.Require().NoError(err)
			s.Equal(
				[]string{
					filepath.Join(dir, "foo", "bar"),
					filepath.Join(dir, "foo"),
					dir,
				}[i],
				path,
			)

			i++

			if path == dir {
				return filepath.SkipAll
			}

			return nil
		},
	)
	s.NoError(err)
}
