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
	"text/template"
)

/****************/
/* Init - Suite */
/****************/

type InitTestSuite struct {
	suite.Suite
	wd string
}

func TestInitTestSuite(t *testing.T) {
	// Run
	suite.Run(t, new(InitTestSuite))
}

func (s *InitTestSuite) SetupSuite() {
	// Current working directory
	s.wd, _ = os.Getwd()
	// Default repository
	viper.SetDefault(
		"repository",
		filepath.Join(s.wd, "testdata/init/repository/default"),
	)
}

func (s *InitTestSuite) ExecuteCmd(dir string, args []string) (*bytes.Buffer, *bytes.Buffer, error) {
	if dir != "" {
		_ = os.Chdir(dir)
	}

	// Command
	cmd := InitCmd()
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
/* Init - Tests */
/****************/

func (s *InitTestSuite) Test() {
	for _, t := range []struct {
		test   string
		dir    string
		args   []string
		err    string
		stdErr string
		stdOut string
		file   string
	}{
		{
			test: "Use recipe",
			dir:  "testdata/init/project/default",
			args: []string{"--recipe", "foo"},
			stdErr: `   • Synced file               path={{ .Dir }}file_default_foo
   • Project synced           
`,
			file: "testdata/init/project/default/file_default_foo",
		},
		{
			test: "Use invalid recipe",
			dir:  "testdata/init/project/default",
			args: []string{"--recipe", "invalid"},
			err:  "recipe not found",
		},
		{
			test: "Use recipe and repository",
			dir:  "testdata/init/project/default",
			args: []string{"--recipe", "foo", "--repository", filepath.Join(s.wd, "testdata/init/repository/custom")},
			stdErr: `   • Synced file               path={{ .Dir }}file_custom_foo
   • Project synced           
`,
			file: "testdata/init/project/default/file_custom_foo",
		},
		{
			test: "Use recipe and invalid repository",
			dir:  "testdata/init/project/default",
			args: []string{"--recipe", "foo", "--repository", "testdata/init/repository/invalid"},
			err:  "\"testdata/init/repository/invalid\" directory does not exists",
		},
	} {
		s.Run(t.test+"/relative", func() {
			// Clean
			_ = os.Remove(t.file)
			// Execute
			stdOut, stdErr, err := s.ExecuteCmd(
				t.dir,
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
			// Stderr
			var stdErrContent bytes.Buffer
			_ = template.Must(template.New("stdErr").Parse(t.stdErr)).Execute(&stdErrContent, map[string]string{
				"Wd":  s.wd + "/",
				"Dir": "",
			})
			s.Equal(stdErrContent.String(), stdErr.String())
			// File
			if t.file != "" {
				s.FileExists(t.file)
			}
		})
		s.Run(t.test+"/dir", func() {
			// Clean
			_ = os.Remove(t.file)
			// Execute
			stdOut, stdErr, err := s.ExecuteCmd(
				"",
				append([]string{t.dir}, t.args...),
			)
			// Test
			if t.err != "" {
				s.Error(err)
				s.Equal(t.err, err.Error())
			} else {
				s.NoError(err)
			}
			s.Equal(t.stdOut, stdOut.String())
			// Stderr
			var stdErrContent bytes.Buffer
			_ = template.Must(template.New("stdErr").Parse(t.stdErr)).Execute(&stdErrContent, map[string]string{
				"Wd":  s.wd + "/",
				"Dir": t.dir + "/",
			})
			s.Equal(stdErrContent.String(), stdErr.String())
			// File
			if t.file != "" {
				s.FileExists(t.file)
			}
		})
	}
}

func (s *InitTestSuite) TestProjectAlreadyExists() {
	s.Run("relative", func() {
		// Execute
		stdOut, stdErr, err := s.ExecuteCmd(
			"testdata/init/project/already_exists",
			[]string{},
		)
		s.Error(err)
		s.Equal("project already exists: .", err.Error())
		s.Equal("", stdOut.String())
		s.Equal("", stdErr.String())
	})
	s.Run("dir", func() {
		// Execute
		stdOut, stdErr, err := s.ExecuteCmd(
			"",
			[]string{"testdata/init/project/already_exists"},
		)
		s.Error(err)
		s.Equal("project already exists: testdata/init/project/already_exists", err.Error())
		s.Equal("", stdOut.String())
		s.Equal("", stdErr.String())
	})
}
