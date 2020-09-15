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
	viper.SetDefault("repository", "testdata/update/repository/default")
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
			args:   []string{"testdata/update/project/default"},
			err:    "",
			stdOut: "",
			stdErr: `   • Project loaded            recipe=foo repository=
   • Repository loaded        
   • Recipe loaded            
   • Project validated        
   • Synced file               path=testdata/update/project/default/file
   • Project synced           
`,
			file: [2]string{
				"testdata/update/project/default/file",
				`Default foo file
`,
			},
		},
		{
			test:   "Default project force repository",
			args:   []string{"testdata/update/project/default", "--repository", "testdata/update/repository/custom"},
			err:    "",
			stdOut: "",
			stdErr: `   • Project loaded            recipe=foo repository=testdata/update/repository/custom
   • Repository loaded        
   • Recipe loaded            
   • Project validated        
   • Synced file               path=testdata/update/project/default/file
   • Project synced           
`,
			file: [2]string{
				"testdata/update/project/default/file",
				`Custom foo file
`,
			},
		},
		{
			test: "Default project force invalid repository",
			args: []string{"testdata/update/project/default", "--repository", "testdata/update/repository/invalid"},
			err:  "\"testdata/update/repository/invalid\" directory does not exists",
		},
		{
			test:   "Default project force recipe",
			args:   []string{"testdata/update/project/default", "--recipe", "bar"},
			err:    "",
			stdOut: "",
			stdErr: `   • Project loaded            recipe=bar repository=
   • Repository loaded        
   • Recipe loaded            
   • Project validated        
   • Synced file               path=testdata/update/project/default/file
   • Project synced           
`,
			file: [2]string{
				"testdata/update/project/default/file",
				`Default bar file
`,
			},
		},
		{
			test: "Default project force invalid recipe",
			args: []string{"testdata/update/project/default", "--recipe", "invalid"},
			err:  "recipe not found",
		},
		{
			test:   "Default project force repository force recipe",
			args:   []string{"testdata/update/project/default", "--repository", "testdata/update/repository/custom", "--recipe", "bar"},
			err:    "",
			stdOut: "",
			stdErr: `   • Project loaded            recipe=bar repository=testdata/update/repository/custom
   • Repository loaded        
   • Recipe loaded            
   • Project validated        
   • Synced file               path=testdata/update/project/default/file
   • Project synced           
`,
			file: [2]string{
				"testdata/update/project/default/file",
				`Custom bar file
`,
			},
		},
		{
			test:   "Custom project",
			args:   []string{"testdata/update/project/custom"},
			err:    "",
			stdOut: "",
			stdErr: `   • Project loaded            recipe=foo repository=testdata/update/repository/custom
   • Repository loaded        
   • Recipe loaded            
   • Project validated        
   • Synced file               path=testdata/update/project/custom/file
   • Project synced           
`,
			file: [2]string{
				"testdata/update/project/custom/file",
				`Custom foo file
`,
			},
		},
		{
			test:   "Custom project force repository",
			args:   []string{"testdata/update/project/custom", "--repository", "testdata/update/repository/force"},
			err:    "",
			stdOut: "",
			stdErr: `   • Project loaded            recipe=foo repository=testdata/update/repository/force
   • Repository loaded        
   • Recipe loaded            
   • Project validated        
   • Synced file               path=testdata/update/project/custom/file
   • Project synced           
`,
			file: [2]string{
				"testdata/update/project/custom/file",
				`Force foo file
`,
			},
		},
		{
			test: "Custom project force invalid repository",
			args: []string{"testdata/update/project/custom", "--repository", "testdata/update/repository/invalid"},
			err:  "\"testdata/update/repository/invalid\" directory does not exists",
		},
		{
			test:   "Custom project force recipe",
			args:   []string{"testdata/update/project/custom", "--recipe", "bar"},
			err:    "",
			stdOut: "",
			stdErr: `   • Project loaded            recipe=bar repository=testdata/update/repository/custom
   • Repository loaded        
   • Recipe loaded            
   • Project validated        
   • Synced file               path=testdata/update/project/custom/file
   • Project synced           
`,
			file: [2]string{
				"testdata/update/project/custom/file",
				`Custom bar file
`,
			},
		},
		{
			test:   "Custom project force repository force recipe",
			args:   []string{"testdata/update/project/custom", "--repository", "testdata/update/repository/force", "--recipe", "bar"},
			err:    "",
			stdOut: "",
			stdErr: `   • Project loaded            recipe=bar repository=testdata/update/repository/force
   • Repository loaded        
   • Recipe loaded            
   • Project validated        
   • Synced file               path=testdata/update/project/custom/file
   • Project synced           
`,
			file: [2]string{
				"testdata/update/project/custom/file",
				`Force bar file
`,
			},
		},
		{
			test: "Custom project force invalid recipe",
			args: []string{"testdata/update/project/custom", "--recipe", "invalid"},
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
			log.SetHandler(cli.New(stdErr))

			// Clean
			_ = os.Remove("testdata/update/project/default/file")
			_ = os.Remove("testdata/update/project/custom/file")

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

func (s *UpdateTestSuite) TestNotFound() {
	// Command
	cmd := UpdateCmd()

	// Io
	stdOut := bytes.NewBufferString("")
	cmd.SetOut(stdOut)
	stdErr := bytes.NewBufferString("")
	cmd.SetErr(stdErr)
	log.SetHandler(cli.New(stdErr))

	// Execute
	cmd.SetArgs([]string{"testdata/update/project/not_found"})
	err := cmd.Execute()

	s.Error(err)
	s.Equal("project not found: testdata/update/project/not_found", err.Error())
}

func (s *UpdateTestSuite) TestInvalid() {
	// Command
	cmd := UpdateCmd()

	// Io
	stdOut := bytes.NewBufferString("")
	cmd.SetOut(stdOut)
	stdErr := bytes.NewBufferString("")
	cmd.SetErr(stdErr)
	log.SetHandler(cli.New(stdErr))

	// Execute
	cmd.SetArgs([]string{"testdata/update/project/invalid"})
	err := cmd.Execute()

	s.Error(err)
	s.Equal("project not found: testdata/update/project/invalid", err.Error())
}

func (s *UpdateTestSuite) TestTraverse() {
	// Command
	cmd := UpdateCmd()

	// Io
	stdOut := bytes.NewBufferString("")
	cmd.SetOut(stdOut)
	stdErr := bytes.NewBufferString("")
	cmd.SetErr(stdErr)
	log.SetHandler(cli.New(stdErr))

	// Clean
	_ = os.Remove("testdata/update/project/traverse/file")

	// Execute
	cmd.SetArgs([]string{"testdata/update/project/traverse/level"})
	err := cmd.Execute()

	s.NoError(err)

	// Test stdout
	s.Zero(stdOut.Len())

	// Test stderr
	s.Equal(`   • Project loaded            recipe=foo repository=
   • Repository loaded        
   • Recipe loaded            
   • Project validated        
   • Synced file               path=testdata/update/project/traverse/file
   • Project synced           
`, stdErr.String())

	// Test file
	s.FileExists("testdata/update/project/traverse/file")
	content, _ := ioutil.ReadFile("testdata/update/project/traverse/file")
	s.Equal(`Default foo file
`, string(content))
}

func (s *UpdateTestSuite) TestRecursive() {
	// Command
	cmd := UpdateCmd()

	// Io
	stdOut := bytes.NewBufferString("")
	cmd.SetOut(stdOut)
	stdErr := bytes.NewBufferString("")
	cmd.SetErr(stdErr)
	log.SetHandler(cli.New(stdErr))

	// Clean
	_ = os.Remove("testdata/update/project/recursive/foo/file")
	_ = os.Remove("testdata/update/project/recursive/bar/file")
	_ = os.Remove("testdata/update/project/recursive/custom/foo/file")

	// Execute
	cmd.SetArgs([]string{"testdata/update/project/recursive", "--recursive"})
	err := cmd.Execute()

	s.NoError(err)

	// Test stdout
	s.Zero(stdOut.Len())

	// Test stderr
	s.Equal(`   • Project loaded            recipe=bar repository=
   • Repository loaded        
   • Recipe loaded            
   • Project validated        
   • Synced file               path=testdata/update/project/recursive/bar/file
   • Project synced           
   • Project loaded            recipe=foo repository=testdata/update/repository/custom
   • Repository loaded        
   • Recipe loaded            
   • Project validated        
   • Synced file               path=testdata/update/project/recursive/custom/foo/file
   • Project synced           
   • Project loaded            recipe=foo repository=
   • Repository loaded        
   • Recipe loaded            
   • Project validated        
   • Synced file               path=testdata/update/project/recursive/foo/file
   • Project synced           
`, stdErr.String())

	// Test file - Foo
	s.FileExists("testdata/update/project/recursive/foo/file")
	content, _ := ioutil.ReadFile("testdata/update/project/recursive/foo/file")
	s.Equal(`Default foo file
`, string(content))

	// Test file - Bar
	s.FileExists("testdata/update/project/recursive/bar/file")
	content, _ = ioutil.ReadFile("testdata/update/project/recursive/bar/file")
	s.Equal(`Default bar file
`, string(content))

	// Test file - Custom Foo
	s.FileExists("testdata/update/project/recursive/custom/foo/file")
	content, _ = ioutil.ReadFile("testdata/update/project/recursive/custom/foo/file")
	s.Equal(`Custom foo file
`, string(content))
}

func (s *UpdateTestSuite) TestRecursiveNotFound() {
	// Command
	cmd := UpdateCmd()

	// Io
	stdOut := bytes.NewBufferString("")
	cmd.SetOut(stdOut)
	stdErr := bytes.NewBufferString("")
	cmd.SetErr(stdErr)
	log.SetHandler(cli.New(stdErr))

	// Execute
	cmd.SetArgs([]string{"testdata/update/project/not_found", "--recursive"})
	err := cmd.Execute()

	s.NoError(err)

	// Test stdout
	s.Zero(stdOut.Len())

	// Test stderr
	s.Zero(stdErr.Len())
}

func (s *UpdateTestSuite) TestRecursiveInvalid() {
	// Command
	cmd := UpdateCmd()

	// Io
	stdOut := bytes.NewBufferString("")
	cmd.SetOut(stdOut)
	stdErr := bytes.NewBufferString("")
	cmd.SetErr(stdErr)
	log.SetHandler(cli.New(stdErr))

	// Execute
	cmd.SetArgs([]string{"testdata/update/project/invalid", "--recursive"})
	err := cmd.Execute()

	s.NoError(err)

	// Test stdout
	s.Zero(stdOut.Len())

	// Test stderr
	s.Zero(stdErr.Len())
}
