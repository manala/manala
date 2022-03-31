package cmd

import (
	"bytes"
	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/stretchr/testify/suite"
	"manala/internal/config"
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

	conf := config.New()
	conf.SetDefault("repository", filepath.Join(s.wd, "testdata/list/repository/default"))

	logger := &log.Logger{
		Handler: cli.New(stdErr),
		Level:   log.InfoLevel,
	}

	// Command
	command := (&ListCmd{}).Command(conf, logger)
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
	}{
		{
			test: "Default repository",
			args: []string{},
		},
		{
			test: "Custom repository",
			args: []string{"--repository", "testdata/list/repository/custom"},
		},
		{
			test: "Nonexistent repository",
			args: []string{"--repository", "testdata/list/repository/nonexistent"},
			err:  "\"testdata/list/repository/nonexistent\" directory does not exists",
		},
	} {
		s.Run(t.test, func() {
			// Execute
			_, stdErr, err := s.ExecuteCommand(
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
			s.Equal(t.stdErr, stdErr.String())
		})
	}
}
