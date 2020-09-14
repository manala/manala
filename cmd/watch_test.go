package cmd

import (
	"bytes"
	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

/*****************/
/* Watch - Suite */
/*****************/

type WatchTestSuite struct {
	suite.Suite
}

func TestWatchTestSuite(t *testing.T) {
	// Config
	viper.SetDefault("repository", "testdata/repository/default")
	// Run
	suite.Run(t, new(WatchTestSuite))
}

/*****************/
/* Watch - Tests */
/*****************/

func (s *WatchTestSuite) Test() {
	for _, t := range []struct {
		test   string
		args   []string
		err    string
		stdErr string
		stdOut string
		file   [2]string
	}{
		{
			test: "Default project force invalid repository",
			args: []string{"testdata/project/default", "--repository", "testdata/repository/invalid"},
			err:  "\"testdata/repository/invalid\" directory does not exists",
		},
		{
			test: "Default project force invalid recipe",
			args: []string{"testdata/project/default", "--recipe", "invalid"},
			err:  "recipe not found",
		},
		{
			test: "Custom project force invalid repository",
			args: []string{"testdata/project/custom", "--repository", "testdata/repository/invalid"},
			err:  "\"testdata/repository/invalid\" directory does not exists",
		},
		{
			test: "Custom project force invalid recipe",
			args: []string{"testdata/project/custom", "--recipe", "invalid"},
			err:  "recipe not found",
		},
	} {
		s.Run(t.test, func() {
			// Command
			cmd := WatchCmd()

			// Io
			stdOut := bytes.NewBufferString("")
			cmd.SetOut(stdOut)
			stdErr := bytes.NewBufferString("")
			cmd.SetErr(stdErr)
			log.SetHandler(cli.New(cmd.ErrOrStderr()))

			// Clean
			_ = os.Remove("testdata/project/default/file")
			_ = os.Remove("testdata/project/custom/file")

			// Execute
			cmd.SetArgs(t.args)
			err := cmd.Execute()

			// Test error
			s.Error(err)
			s.Equal(t.err, err.Error())
		})
	}
}
