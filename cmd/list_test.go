package cmd

import (
	"bytes"
	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
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
	// Config
	viper.SetDefault("repository", "testdata/list/repository/default")
	// Run
	suite.Run(t, new(ListTestSuite))
}

func (s *ListTestSuite) SetupSuite() {
	// Current working directory
	s.wd, _ = os.Getwd()
	// Default repository
	viper.SetDefault(
		"repository",
		filepath.Join(s.wd, "testdata/list/repository/default"),
	)
}

func (s *ListTestSuite) ExecuteCmd(dir string, args []string) (*bytes.Buffer, *bytes.Buffer, error) {
	if dir != "" {
		_ = os.Chdir(dir)
	}

	// Command
	cmd := ListCmd()
	cmd.SetArgs(args)
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	stdOut := bytes.NewBufferString("")
	cmd.SetOut(stdOut)
	stdErr := bytes.NewBufferString("")
	cmd.SetErr(stdErr)

	log.SetHandler(cli.New(cmd.ErrOrStderr()))

	err := cmd.Execute()

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
			test: "Use repository",
			args: []string{"--repository", "testdata/list/repository/custom"},
			stdOut: `bar: Custom bar recipe
foo: Custom foo recipe
`,
		},
		{
			test: "Use invalid repository",
			args: []string{"--repository", "testdata/list/repository/invalid"},
			err:  "\"testdata/list/repository/invalid\" directory does not exists",
		},
	} {
		s.Run(t.test, func() {
			// Execute
			stdOut, stdErr, err := s.ExecuteCmd(
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
