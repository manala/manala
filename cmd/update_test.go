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

/******************/
/* Update - Suite */
/******************/

type UpdateTestSuite struct {
	suite.Suite
}

func TestUpdateTestSuite(t *testing.T) {
	// Config
	viper.SetDefault("repository", "testdata/repository/default")
	// Run
	suite.Run(t, new(UpdateTestSuite))
}

/******************/
/* Update - Tests */
/******************/

func (s *UpdateTestSuite) Test() {
	for _, t := range []struct {
		test   string
		args   []string
		err    string
		stdErr string
		stdOut string
		file   [2]string
	}{
		{
			test:   "Default project",
			args:   []string{"testdata/project/default"},
			err:    "",
			stdOut: "",
			stdErr: `   • Project loaded            recipe=foo repository=
   • Repository loaded        
   • Recipe loaded            
   • Project validated        
   • Synced file               path=testdata/project/default/file
   • Project synced           
`,
			file: [2]string{
				"testdata/project/default/file",
				`Default foo file
`,
			},
		},
		{
			test:   "Default project force repository",
			args:   []string{"testdata/project/default", "--repository", "testdata/repository/custom"},
			err:    "",
			stdOut: "",
			stdErr: `   • Project loaded            recipe=foo repository=testdata/repository/custom
   • Repository loaded        
   • Recipe loaded            
   • Project validated        
   • Synced file               path=testdata/project/default/file
   • Project synced           
`,
			file: [2]string{
				"testdata/project/default/file",
				`Custom foo file
`,
			},
		},
		{
			test: "Default project force invalid repository",
			args: []string{"testdata/project/default", "--repository", "testdata/repository/invalid"},
			err:  "\"testdata/repository/invalid\" directory does not exists",
		},
		{
			test:   "Default project force recipe",
			args:   []string{"testdata/project/default", "--recipe", "bar"},
			err:    "",
			stdOut: "",
			stdErr: `   • Project loaded            recipe=bar repository=
   • Repository loaded        
   • Recipe loaded            
   • Project validated        
   • Synced file               path=testdata/project/default/file
   • Project synced           
`,
			file: [2]string{
				"testdata/project/default/file",
				`Default bar file
`,
			},
		},
		{
			test: "Default project force invalid recipe",
			args: []string{"testdata/project/default", "--recipe", "invalid"},
			err:  "recipe not found",
		},
		{
			test:   "Default project force repository force recipe",
			args:   []string{"testdata/project/default", "--repository", "testdata/repository/custom", "--recipe", "bar"},
			err:    "",
			stdOut: "",
			stdErr: `   • Project loaded            recipe=bar repository=testdata/repository/custom
   • Repository loaded        
   • Recipe loaded            
   • Project validated        
   • Synced file               path=testdata/project/default/file
   • Project synced           
`,
			file: [2]string{
				"testdata/project/default/file",
				`Custom bar file
`,
			},
		},
		{
			test:   "Custom project",
			args:   []string{"testdata/project/custom"},
			err:    "",
			stdOut: "",
			stdErr: `   • Project loaded            recipe=foo repository=testdata/repository/custom
   • Repository loaded        
   • Recipe loaded            
   • Project validated        
   • Synced file               path=testdata/project/custom/file
   • Project synced           
`,
			file: [2]string{
				"testdata/project/custom/file",
				`Custom foo file
`,
			},
		},
		{
			test:   "Custom project force repository",
			args:   []string{"testdata/project/custom", "--repository", "testdata/repository/force"},
			err:    "",
			stdOut: "",
			stdErr: `   • Project loaded            recipe=foo repository=testdata/repository/force
   • Repository loaded        
   • Recipe loaded            
   • Project validated        
   • Synced file               path=testdata/project/custom/file
   • Project synced           
`,
			file: [2]string{
				"testdata/project/custom/file",
				`Force foo file
`,
			},
		},
		{
			test: "Custom project force invalid repository",
			args: []string{"testdata/project/custom", "--repository", "testdata/repository/invalid"},
			err:  "\"testdata/repository/invalid\" directory does not exists",
		},
		{
			test:   "Custom project force recipe",
			args:   []string{"testdata/project/custom", "--recipe", "bar"},
			err:    "",
			stdOut: "",
			stdErr: `   • Project loaded            recipe=bar repository=testdata/repository/custom
   • Repository loaded        
   • Recipe loaded            
   • Project validated        
   • Synced file               path=testdata/project/custom/file
   • Project synced           
`,
			file: [2]string{
				"testdata/project/custom/file",
				`Custom bar file
`,
			},
		},
		{
			test:   "Custom project force repository force recipe",
			args:   []string{"testdata/project/custom", "--repository", "testdata/repository/force", "--recipe", "bar"},
			err:    "",
			stdOut: "",
			stdErr: `   • Project loaded            recipe=bar repository=testdata/repository/force
   • Repository loaded        
   • Recipe loaded            
   • Project validated        
   • Synced file               path=testdata/project/custom/file
   • Project synced           
`,
			file: [2]string{
				"testdata/project/custom/file",
				`Force bar file
`,
			},
		},
		{
			test: "Custom project force invalid recipe",
			args: []string{"testdata/project/custom", "--recipe", "invalid"},
			err:  "recipe not found",
		},
	} {
		s.Run(t.test, func() {
			// Command
			cmd := UpdateCmd()

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
