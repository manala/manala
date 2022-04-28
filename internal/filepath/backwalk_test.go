package filepath

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

var backwalkTestPath = filepath.Join("testdata", "backwalk")

func (s *BackwalkSuite) Test() {
	i := 0
	err := Backwalk(
		filepath.Join(backwalkTestPath, "foo", "bar"),
		func(path string, file os.DirEntry, err error) error {
			s.NoError(err)
			s.Equal(
				[]string{
					filepath.Join(backwalkTestPath, "foo", "bar"),
					filepath.Join(backwalkTestPath, "foo"),
					backwalkTestPath,
				}[i],
				path,
			)
			i = i + 1
			if path == backwalkTestPath {
				return filepath.SkipDir
			}
			return nil
		},
	)
	s.NoError(err)
}
