package filepath

import (
	"github.com/stretchr/testify/suite"
	internalTesting "manala/internal/testing"
	"os"
	"path/filepath"
	"testing"
)

type BackwalkSuite struct{ suite.Suite }

func TestBackwalkSuite(t *testing.T) {
	suite.Run(t, new(BackwalkSuite))
}

func (s *BackwalkSuite) Test() {
	i := 0
	err := Backwalk(
		internalTesting.DataPath(s, "foo", "bar"),
		func(path string, file os.DirEntry, err error) error {
			s.NoError(err)
			s.Equal(
				[]string{
					internalTesting.DataPath(s, "foo", "bar"),
					internalTesting.DataPath(s, "foo"),
					internalTesting.DataPath(s),
				}[i],
				path,
			)
			i = i + 1
			if path == internalTesting.DataPath(s) {
				return filepath.SkipDir
			}
			return nil
		},
	)
	s.NoError(err)
}
