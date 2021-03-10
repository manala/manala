package cmd

import (
	"bytes"
	"github.com/stretchr/testify/suite"
	"manala/config"
	"manala/loaders"
	"manala/logger"
	"manala/syncer"
	"manala/template"
	"os"
	"path/filepath"
	"strings"
	"testing"
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
}

func (s *InitTestSuite) ExecuteCommand(dir string, args []string) (*bytes.Buffer, *bytes.Buffer, error) {
	if dir != "" {
		_ = os.Chdir(dir)
	}

	stdOut := bytes.NewBufferString("")
	stdErr := bytes.NewBufferString("")

	conf := config.New("test", filepath.Join(s.wd, "testdata/init/repository/default"))

	log := logger.New(conf)
	log.SetOut(stdErr)

	tmpl := template.New()

	repositoryLoader := loaders.NewRepositoryLoader(log, conf)
	recipeLoader := loaders.NewRecipeLoader(log)

	cmd := &InitCmd{
		Log:              log,
		RepositoryLoader: repositoryLoader,
		RecipeLoader:     recipeLoader,
		ProjectLoader:    loaders.NewProjectLoader(log, repositoryLoader, recipeLoader),
		Sync:             syncer.New(log, tmpl),
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
			stdErr: `   • Synced file               path={{ dir }}file_default_foo
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
			stdErr: `   • Synced file               path={{ dir }}file_custom_foo
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
			_ = os.RemoveAll(t.dir)
			_ = os.Mkdir(t.dir, 0755)
			// Execute
			stdOut, stdErr, err := s.ExecuteCommand(
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
				strings.NewReplacer("{{ wd }}", s.wd+"/", "{{ dir }}", "").Replace(t.stdErr),
				stdErr.String(),
			)
			// File
			if t.file != "" {
				s.FileExists(t.file)
			}
		})
		s.Run(t.test+"/dir", func() {
			// Clean
			_ = os.RemoveAll(t.dir)
			_ = os.Mkdir(t.dir, 0755)
			// Execute
			stdOut, stdErr, err := s.ExecuteCommand(
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
				strings.NewReplacer("{{ wd }}", s.wd+"/", "{{ dir }}", t.dir+"/").Replace(t.stdErr),
				stdErr.String(),
			)
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
		stdOut, stdErr, err := s.ExecuteCommand(
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
		stdOut, stdErr, err := s.ExecuteCommand(
			"",
			[]string{"testdata/init/project/already_exists"},
		)
		s.Error(err)
		s.Equal("project already exists: testdata/init/project/already_exists", err.Error())
		s.Equal("", stdOut.String())
		s.Equal("", stdErr.String())
	})
}
