package cmd

import (
	"bytes"
	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/stretchr/testify/suite"
	"manala/internal/config"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

/******************/
/* Update - Suite */
/******************/

type UpdateTestSuite struct {
	suite.Suite
	wd string
}

func TestUpdateTestSuite(t *testing.T) {
	// Run
	suite.Run(t, new(UpdateTestSuite))
}

func (s *UpdateTestSuite) SetupSuite() {
	// Current working directory
	s.wd, _ = os.Getwd()
}

func (s *UpdateTestSuite) ExecuteCmd(dir string, args []string) (*bytes.Buffer, *bytes.Buffer, error) {
	if dir != "" {
		_ = os.Chdir(dir)
	}

	stdOut := bytes.NewBufferString("")
	stdErr := bytes.NewBufferString("")

	conf := config.New()
	conf.SetDefault("repository", filepath.Join(s.wd, "testdata/update/repository/default"))

	logger := &log.Logger{
		Handler: cli.New(stdErr),
		Level:   log.InfoLevel,
	}

	// Command
	command := (&UpdateCmd{}).Command(conf, logger)
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

/******************/
/* Update - Tests */
/******************/

func (s *UpdateTestSuite) Test() {
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
			test: "Default",
			dir:  "testdata/update/project/default",
			args: []string{},
			stdErr: `   • Project loaded            recipe=foo repository={{ wd }}{{ ps }}testdata{{ ps }}update{{ ps }}repository{{ ps }}default
   • Repository loaded        
   • Recipe loaded            
   • Project validated        
   • Synced file               path=file_default_foo
   • Project synced           
`,
			file: "testdata/update/project/default/file_default_foo",
		},
		{
			test: "Default with repository",
			dir:  "testdata/update/project/default",
			args: []string{"--repository", filepath.Join(s.wd, "testdata/update/repository/custom")},
			stdErr: `   • Project loaded            recipe=foo repository={{ wd }}{{ ps }}testdata{{ ps }}update{{ ps }}repository{{ ps }}custom
   • Repository loaded        
   • Recipe loaded            
   • Project validated        
   • Synced file               path=file_custom_foo
   • Project synced           
`,
			file: "testdata/update/project/default/file_custom_foo",
		},
		{
			test: "Default with invalid repository",
			dir:  "testdata/update/project/default",
			args: []string{"--repository", "testdata/update/repository/invalid"},
			err:  "\"testdata/update/repository/invalid\" directory does not exists",
			stdErr: `   • Project loaded            recipe=foo repository=testdata/update/repository/invalid
`,
		},
		{
			test: "Default with recipe",
			dir:  "testdata/update/project/default",
			args: []string{"--recipe", "bar"},
			stdErr: `   • Project loaded            recipe=bar repository={{ wd }}{{ ps }}testdata{{ ps }}update{{ ps }}repository{{ ps }}default
   • Repository loaded        
   • Recipe loaded            
   • Project validated        
   • Synced file               path=file_default_bar
   • Project synced           
`,
			file: "testdata/update/project/default/file_default_bar",
		},
		{
			test: "Default with invalid recipe",
			dir:  "testdata/update/project/default",
			args: []string{"--recipe", "invalid"},
			err:  "recipe not found",
			stdErr: `   • Project loaded            recipe=invalid repository={{ wd }}{{ ps }}testdata{{ ps }}update{{ ps }}repository{{ ps }}default
   • Repository loaded        
`,
		},
		{
			test:   "Default with repository and recipe",
			dir:    "testdata/update/project/default",
			args:   []string{"--repository", filepath.Join(s.wd, "testdata/update/repository/custom"), "--recipe", "bar"},
			err:    "",
			stdOut: "",
			stdErr: `   • Project loaded            recipe=bar repository={{ wd }}{{ ps }}testdata{{ ps }}update{{ ps }}repository{{ ps }}custom
   • Repository loaded        
   • Recipe loaded            
   • Project validated        
   • Synced file               path=file_custom_bar
   • Project synced           
`,
			file: "testdata/update/project/default/file_custom_bar",
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
			s.Equal(
				strings.NewReplacer("{{ wd }}", s.wd, "{{ ps }}", string(os.PathSeparator)).Replace(t.stdErr),
				stdErr.String(),
			)
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
			s.Equal(
				strings.NewReplacer("{{ wd }}", s.wd, "{{ ps }}", string(os.PathSeparator)).Replace(t.stdErr),
				stdErr.String(),
			)
			// File
			if t.file != "" {
				s.FileExists(t.file)
			}
		})
	}
}

func (s *UpdateTestSuite) TestNotFound() {
	s.Run("relative", func() {
		// Execute
		stdOut, stdErr, err := s.ExecuteCmd(
			"testdata/update/project/not_found",
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
			[]string{"testdata/update/project/not_found"},
		)
		s.Error(err)
		s.Equal("project not found: testdata/update/project/not_found", err.Error())
		s.Equal("", stdOut.String())
		s.Equal("", stdErr.String())
	})
}

func (s *UpdateTestSuite) TestInvalid() {
	s.Run("dir", func() {
		// Execute
		stdOut, stdErr, err := s.ExecuteCmd(
			"",
			[]string{"testdata/update/project/invalid"},
		)
		s.Error(err)
		s.Equal("invalid directory: testdata/update/project/invalid", err.Error())
		s.Equal("", stdOut.String())
		s.Equal("", stdErr.String())
	})
}

func (s *UpdateTestSuite) TestTraverse() {
	s.Run("relative", func() {
		// Clean
		_ = os.Remove("testdata/update/project/traverse/file_default_foo")
		// Execute
		stdOut, stdErr, err := s.ExecuteCmd(
			"testdata/update/project/traverse/level",
			[]string{},
		)
		s.NoError(err)
		s.Equal("", stdOut.String())
		s.Equal(
			strings.NewReplacer("{{ wd }}", s.wd, "{{ ps }}", string(os.PathSeparator)).Replace(`   • Project loaded            recipe=foo repository={{ wd }}{{ ps }}testdata{{ ps }}update{{ ps }}repository{{ ps }}default
   • Repository loaded        
   • Recipe loaded            
   • Project validated        
   • Synced file               path=file_default_foo
   • Project synced           
`),
			stdErr.String(),
		)
		s.FileExists("testdata/update/project/traverse/file_default_foo")
	})
	s.Run("dir", func() {
		// Clean
		_ = os.Remove("testdata/update/project/traverse/file_default_foo")
		// Execute
		stdOut, stdErr, err := s.ExecuteCmd(
			"",
			[]string{"testdata/update/project/traverse/level"},
		)
		s.NoError(err)
		s.Equal("", stdOut.String())
		s.Equal(
			strings.NewReplacer("{{ wd }}", s.wd, "{{ ps }}", string(os.PathSeparator)).Replace(`   • Project loaded            recipe=foo repository={{ wd }}{{ ps }}testdata{{ ps }}update{{ ps }}repository{{ ps }}default
   • Repository loaded        
   • Recipe loaded            
   • Project validated        
   • Synced file               path=file_default_foo
   • Project synced           
`),
			stdErr.String(),
		)
		s.FileExists("testdata/update/project/traverse/file_default_foo")
	})
}

func (s *UpdateTestSuite) TestRecursive() {
	s.Run("relative", func() {
		// Clean
		_ = os.Remove("testdata/update/project/recursive/foo/file_default_foo")
		_ = os.Remove("testdata/update/project/recursive/foo/embedded/file_default_bar")
		_ = os.Remove("testdata/update/project/recursive/bar/file_default_bar")
		_ = os.Remove("testdata/update/project/recursive/level/foo/file_default_foo")
		// Execute
		stdOut, stdErr, err := s.ExecuteCmd(
			"testdata/update/project/recursive",
			[]string{"--recursive"},
		)
		s.NoError(err)
		s.Equal("", stdOut.String())
		s.Equal(
			strings.NewReplacer("{{ wd }}", s.wd, "{{ ps }}", string(os.PathSeparator)).Replace(
				`   • Project loaded            recipe=bar repository={{ wd }}{{ ps }}testdata{{ ps }}update{{ ps }}repository{{ ps }}default
   • Repository loaded        
   • Recipe loaded            
   • Project validated        
   • Synced file               path=file_default_bar
   • Project synced           
   • Project loaded            recipe=foo repository={{ wd }}{{ ps }}testdata{{ ps }}update{{ ps }}repository{{ ps }}default
   • Repository loaded        
   • Recipe loaded            
   • Project validated        
   • Synced file               path=file_default_foo
   • Project synced           
   • Project loaded            recipe=bar repository={{ wd }}{{ ps }}testdata{{ ps }}update{{ ps }}repository{{ ps }}default
   • Repository loaded        
   • Recipe loaded            
   • Project validated        
   • Synced file               path=file_default_bar
   • Project synced           
   • Project loaded            recipe=foo repository={{ wd }}{{ ps }}testdata{{ ps }}update{{ ps }}repository{{ ps }}default
   • Repository loaded        
   • Recipe loaded            
   • Project validated        
   • Synced file               path=file_default_foo
   • Project synced           
`), stdErr.String(),
		)
		s.FileExists("testdata/update/project/recursive/foo/file_default_foo")
		s.FileExists("testdata/update/project/recursive/foo/embedded/file_default_bar")
		s.FileExists("testdata/update/project/recursive/bar/file_default_bar")
		s.FileExists("testdata/update/project/recursive/level/foo/file_default_foo")
	})
	s.Run("dir", func() {
		// Clean
		_ = os.Remove("testdata/update/project/recursive/foo/file_default_foo")
		_ = os.Remove("testdata/update/project/recursive/foo/embedded/file_default_bar")
		_ = os.Remove("testdata/update/project/recursive/bar/file_default_bar")
		_ = os.Remove("testdata/update/project/recursive/level/foo/file_default_foo")
		// Execute
		stdOut, stdErr, err := s.ExecuteCmd(
			"",
			[]string{"testdata/update/project/recursive", "--recursive"},
		)
		s.NoError(err)
		s.Equal("", stdOut.String())
		s.Equal(
			strings.NewReplacer("{{ wd }}", s.wd, "{{ ps }}", string(os.PathSeparator)).Replace(`   • Project loaded            recipe=bar repository={{ wd }}{{ ps }}testdata{{ ps }}update{{ ps }}repository{{ ps }}default
   • Repository loaded        
   • Recipe loaded            
   • Project validated        
   • Synced file               path=file_default_bar
   • Project synced           
   • Project loaded            recipe=foo repository={{ wd }}{{ ps }}testdata{{ ps }}update{{ ps }}repository{{ ps }}default
   • Repository loaded        
   • Recipe loaded            
   • Project validated        
   • Synced file               path=file_default_foo
   • Project synced           
   • Project loaded            recipe=bar repository={{ wd }}{{ ps }}testdata{{ ps }}update{{ ps }}repository{{ ps }}default
   • Repository loaded        
   • Recipe loaded            
   • Project validated        
   • Synced file               path=file_default_bar
   • Project synced           
   • Project loaded            recipe=foo repository={{ wd }}{{ ps }}testdata{{ ps }}update{{ ps }}repository{{ ps }}default
   • Repository loaded        
   • Recipe loaded            
   • Project validated        
   • Synced file               path=file_default_foo
   • Project synced           
`),
			stdErr.String(),
		)
		s.FileExists("testdata/update/project/recursive/foo/file_default_foo")
		s.FileExists("testdata/update/project/recursive/foo/embedded/file_default_bar")
		s.FileExists("testdata/update/project/recursive/bar/file_default_bar")
		s.FileExists("testdata/update/project/recursive/level/foo/file_default_foo")
	})
}

func (s *UpdateTestSuite) TestRecursiveNotFound() {
	s.Run("relative", func() {
		// Execute
		stdOut, stdErr, err := s.ExecuteCmd(
			"testdata/update/project/not_found",
			[]string{"--recursive"},
		)
		s.NoError(err)
		s.Equal("", stdOut.String())
		s.Equal("", stdErr.String())
	})
	s.Run("dir", func() {
		// Execute
		stdOut, stdErr, err := s.ExecuteCmd(
			"",
			[]string{"testdata/update/project/not_found", "--recursive"},
		)
		s.NoError(err)
		s.Equal("", stdOut.String())
		s.Equal("", stdErr.String())
	})
}

func (s *UpdateTestSuite) TestRecursiveInvalid() {
	s.Run("dir", func() {
		// Execute
		stdOut, stdErr, err := s.ExecuteCmd(
			"",
			[]string{"testdata/update/project/invalid", "--recursive"},
		)
		s.Error(err)
		s.Equal("invalid directory: testdata/update/project/invalid", err.Error())
		s.Equal("", stdOut.String())
		s.Equal("", stdErr.String())
	})
}
