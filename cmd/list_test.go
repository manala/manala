package cmd

import (
	"bytes"
	"github.com/stretchr/testify/suite"
	"manala/app"
	"manala/config"
	"manala/fs"
	"manala/loaders"
	"manala/logger"
	"manala/models"
	"os"
	"path/filepath"
	"testing"
)

/****************/
/* List - Suite */
/****************/

type ListTestSuite struct {
	suite.Suite
	wd string
}

func TestListTestSuite(t *testing.T) {
	// Run
	suite.Run(t, new(ListTestSuite))
}

func (s *ListTestSuite) SetupSuite() {
	// Current working directory
	s.wd, _ = os.Getwd()
}

func (s *ListTestSuite) ExecuteCommand(dir string, args []string) (*bytes.Buffer, *bytes.Buffer, error) {
	if dir != "" {
		_ = os.Chdir(dir)
	}

	stdOut := bytes.NewBufferString("")
	stdErr := bytes.NewBufferString("")

	conf := config.New(
		config.WithMainRepository(filepath.Join(s.wd, "testdata/list/repository/default")),
	)

	log := logger.New(logger.WithWriter(stdErr))

	fsManager := fs.NewManager()
	modelFsManager := models.NewFsManager(fsManager)

	repositoryLoader := loaders.NewRepositoryLoader(log, conf)
	recipeLoader := loaders.NewRecipeLoader(log, modelFsManager)

	cmd := &ListCmd{
		App: &app.App{
			RepositoryLoader: repositoryLoader,
			RecipeLoader:     recipeLoader,
		},
		Conf: conf,
		Out:  stdOut,
	}

	// Command
	command := cmd.Command()
	command.SetArgs(args)
	command.SilenceErrors = true
	command.SilenceUsage = true
	command.SetOut(stdOut)
	command.SetErr(stdErr)

	err := command.Execute()

	if dir != "" {
		_ = os.Chdir(s.wd)
	}

	return stdOut, stdErr, err
}

/****************/
/* List - Tests */
/****************/

func (s *ListTestSuite) Test() {
	for _, t := range []struct {
		test   string
		args   []string
		err    string
		stdErr string
		stdOut string
	}{
		{
			test: "Default repository",
			args: []string{},
			stdOut: `bar: Default bar recipe
foo: Default foo recipe
`,
		},
		{
			test: "Custom repository",
			args: []string{"--repository", "testdata/list/repository/custom"},
			stdOut: `bar: Custom bar recipe
foo: Custom foo recipe
`,
		},
		{
			test: "Nonexistent repository",
			args: []string{"--repository", "testdata/list/repository/nonexistent"},
			err:  "\"testdata/list/repository/nonexistent\" directory does not exists",
		},
	} {
		s.Run(t.test, func() {
			// Execute
			stdOut, stdErr, err := s.ExecuteCommand(
				"",
				t.args,
			)
			// Tests
			if t.err != "" {
				s.Error(err)
				s.Equal(t.err, err.Error())
			} else {
				s.NoError(err)
			}
			s.Equal(t.stdOut, stdOut.String())
			s.Equal(t.stdErr, stdErr.String())
		})
	}
}
