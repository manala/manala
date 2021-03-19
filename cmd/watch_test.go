package cmd

import (
	"bytes"
	"github.com/stretchr/testify/suite"
	"manala/config"
	"manala/fs"
	"manala/loaders"
	"manala/logger"
	"manala/models"
	"manala/syncer"
	"manala/template"
	"os"
	"path/filepath"
	"strings"
	"testing"
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
}

func (s *WatchTestSuite) ExecuteCmd(dir string, args []string) (*bytes.Buffer, *bytes.Buffer, error) {
	if dir != "" {
		_ = os.Chdir(dir)
	}

	stdOut := bytes.NewBufferString("")
	stdErr := bytes.NewBufferString("")

	conf := config.New(
		config.WithMainRepository(filepath.Join(s.wd, "testdata/update/repository/default")),
	)

	log := logger.New(logger.WithWriter(stdErr))

	fsManager := fs.NewManager()
	modelFsManager := models.NewFsManager(fsManager)
	templateManager := template.NewManager()
	modelTemplateManager := models.NewTemplateManager(templateManager, modelFsManager)
	modelWatcherManager := models.NewWatcherManager(log)

	repositoryLoader := loaders.NewRepositoryLoader(log, conf)
	recipeLoader := loaders.NewRecipeLoader(log, modelFsManager)

	cmd := &WatchCmd{
		Log:            log,
		ProjectLoader:  loaders.NewProjectLoader(log, conf, repositoryLoader, recipeLoader),
		WatcherManager: modelWatcherManager,
		Sync:           syncer.New(log, modelFsManager, modelTemplateManager),
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
			test: "Default with invalid repository",
			dir:  "testdata/watch/project/default",
			args: []string{"--repository", "testdata/watch/repository/invalid"},
			err:  "\"testdata/watch/repository/invalid\" directory does not exists",
			stdErr: `   • Project loaded            recipe=foo repository=testdata/watch/repository/invalid
`,
		},
		{
			test: "Default with invalid recipe",
			dir:  "testdata/watch/project/default",
			args: []string{"--recipe", "invalid"},
			err:  "recipe not found",
			stdErr: `   • Project loaded            recipe=invalid repository={{ wd }}testdata/update/repository/default
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
			s.Equal(
				strings.NewReplacer("{{ wd }}", s.wd+"/").Replace(t.stdErr),
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
				strings.NewReplacer("{{ wd }}", s.wd+"/").Replace(t.stdErr),
				stdErr.String(),
			)
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
