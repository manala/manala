package cmd

import (
	"bytes"
	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"os"
	"testing"
)

/****************/
/* Init - Suite */
/****************/

type InitTestSuite struct {
	suite.Suite
}

func TestInitTestSuite(t *testing.T) {
	// Config
	viper.SetDefault("repository", "testdata/init/repository/default")
	// Run
	suite.Run(t, new(InitTestSuite))
}

/****************/
/* Init - Tests */
/****************/

func (s *InitTestSuite) Test() {
	for _, t := range []struct {
		test   string
		args   []string
		err    string
		stdErr string
		stdOut string
		file   [2]string
	}{
		{
			test:   "Use recipe",
			args:   []string{"testdata/init/project/init", "--recipe", "foo"},
			err:    "",
			stdOut: "",
			stdErr: `   • Project directory created dir=testdata/init/project/init
   • Synced file               path=testdata/init/project/init/file
   • Project synced           
`,
			file: [2]string{
				"testdata/init/project/init/file",
				`Default foo file
`,
			},
		},
		{
			test: "Use invalid recipe",
			args: []string{"testdata/init/project/init", "--recipe", "invalid"},
			err:  "recipe not found",
		},
		{
			test:   "Use recipe use repository",
			args:   []string{"testdata/init/project/init", "--recipe", "foo", "--repository", "testdata/init/repository/custom"},
			err:    "",
			stdOut: "",
			stdErr: `   • Project directory created dir=testdata/init/project/init
   • Synced file               path=testdata/init/project/init/file
   • Project synced           
`,
			file: [2]string{
				"testdata/init/project/init/file",
				`Custom foo file
`,
			},
		},
		{
			test: "Use recipe use invalid repository",
			args: []string{"testdata/init/project/init", "--recipe", "foo", "--repository", "testdata/init/repository/invalid"},
			err:  "\"testdata/init/repository/invalid\" directory does not exists",
		},
	} {
		s.Run(t.test, func() {
			// Command
			cmd := InitCmd()

			// Io
			stdOut := bytes.NewBufferString("")
			cmd.SetOut(stdOut)
			stdErr := bytes.NewBufferString("")
			cmd.SetErr(stdErr)
			log.SetHandler(cli.New(stdErr))

			// Clean
			_ = os.RemoveAll("testdata/init/project/init")

			// Execute
			cmd.SetArgs(t.args)
			err := cmd.Execute()

			// Test error
			if t.err != "" {
				s.Error(err)
				s.Equal(t.err, err.Error())
			} else {
				s.NoError(err)

				// Test stdout
				if t.stdOut == "" {
					s.Zero(stdOut.Len())
				} else {
					s.Equal(t.stdOut, stdOut.String())
				}

				// Test stderr
				if t.stdErr == "" {
					s.Zero(stdErr.Len())
				} else {
					s.Equal(t.stdErr, stdErr.String())
				}

				// Test file
				if t.file[0] != "" {
					s.FileExists(t.file[0])
					content, _ := ioutil.ReadFile(t.file[0])
					s.Equal(t.file[1], string(content))
				}
			}
		})
	}
}

func (s *InitTestSuite) TestProjectAlreadyExists() {
	// Command
	cmd := InitCmd()

	// Io
	stdOut := bytes.NewBufferString("")
	cmd.SetOut(stdOut)
	stdErr := bytes.NewBufferString("")
	cmd.SetErr(stdErr)
	log.SetHandler(cli.New(stdErr))

	// Execute
	cmd.SetArgs([]string{"testdata/init/project/already_exists"})
	err := cmd.Execute()

	// Test error
	s.Error(err)
	s.Equal("project already exists: testdata/init/project/already_exists", err.Error())
}
