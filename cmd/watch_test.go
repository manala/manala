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

/*****************/
/* Watch - Suite */
/*****************/

type WatchTestSuite struct {
	suite.Suite
	wd string
}

func TestWatchTestSuite(t *testing.T) {
	// Run
	suite.Run(t, new(WatchTestSuite))
}

func (s *WatchTestSuite) SetupSuite() {
	// Current working directory
	s.wd, _ = os.Getwd()
	// Default repository
	viper.SetDefault(
		"repository",
		filepath.Join(s.wd, "testdata/watch/repository/default"),
	)
}

func (s *WatchTestSuite) ExecuteCmd(dir string, args []string) (*bytes.Buffer, *bytes.Buffer, error) {
	if dir != "" {
		_ = os.Chdir(dir)
	}

	// Command
	cmd := WatchCmd()
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

/*****************/
/* Watch - Tests */
/*****************/

func (s *WatchTestSuite) Test() {
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
			test: "Default force invalid repository",
			dir:  "testdata/watch/project/default",
			args: []string{"--repository", "testdata/watch/repository/invalid"},
			err:  "\"testdata/watch/repository/invalid\" directory does not exists",
			stdErr: `   • Project loaded            recipe=foo repository=testdata/watch/repository/invalid
`,
		},
		{
			test: "Default force invalid recipe",
			dir:  "testdata/watch/project/default",
			args: []string{"--recipe", "invalid"},
			err:  "recipe not found",
			stdErr: `   • Project loaded            recipe=invalid repository=
   • Repository loaded        
`,
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

func (s *WatchTestSuite) TestNotFound() {
	s.Run("relative", func() {
		// Execute
		stdOut, stdErr, err := s.ExecuteCmd(
			"testdata/watch/project/not_found",
			[]string{},
		)
		s.Error(err)
		s.Equal("project not found: .", err.Error())
		s.Equal("", stdOut.String())
		s.Equal("", stdErr.String())
	})
	s.Run("dir", func() {
		// Execute
		stdOut, stdErr, err := s.ExecuteCmd(
			"",
			[]string{"testdata/watch/project/not_found"},
		)
		s.Error(err)
		s.Equal("project not found: testdata/watch/project/not_found", err.Error())
		s.Equal("", stdOut.String())
		s.Equal("", stdErr.String())
	})
}

func (s *WatchTestSuite) TestInvalid() {
	s.Run("dir", func() {
		// Execute
		stdOut, stdErr, err := s.ExecuteCmd(
			"",
			[]string{"testdata/watch/project/invalid"},
		)
		s.Error(err)
		s.Equal("invalid directory: testdata/watch/project/invalid", err.Error())
		s.Equal("", stdOut.String())
		s.Equal("", stdErr.String())
	})
}
