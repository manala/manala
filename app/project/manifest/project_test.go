package manifest //nolint:testpackage

import (
	_ "embed"
	"path/filepath"
	"testing"

	"github.com/manala/manala/app/testing/mocks"

	"github.com/stretchr/testify/suite"
)

type ProjectSuite struct {
	suite.Suite
}

func TestProjectSuite(t *testing.T) {
	suite.Run(t, new(ProjectSuite))
}

func (s *ProjectSuite) Test() {
	dir := "dir"
	recipeMock := &mocks.Recipe{}
	vars := map[string]any{"foo": "bar"}

	project := &Project{
		dir:    dir,
		recipe: recipeMock,
		vars:   vars,
	}

	s.Equal(dir, project.Dir())
	s.Equal(recipeMock, project.Recipe())
	s.Equal(vars, project.Vars())

	watches, err := project.Watches()
	s.Require().NoError(err)
	s.Equal([]string{filepath.Join(dir, ".manala.yaml")}, watches)
}
