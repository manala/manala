package cmd

import (
	"bytes"
	"github.com/stretchr/testify/suite"
	"manala/app"
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
	"testing/fstest"
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

	conf := config.New(
		config.WithMainRepository(filepath.Join(s.wd, "testdata/init/repository/default")),
	)

	log := logger.New(logger.WithWriter(stdErr))

	fsManager := fs.NewManager()
	modelFsManager := models.NewFsManager(fsManager)
	templateManager := template.NewManager()
	modelTemplateManager := models.NewTemplateManager(templateManager, modelFsManager)

	repositoryLoader := loaders.NewRepositoryLoader(log, conf)
	recipeLoader := loaders.NewRecipeLoader(log, modelFsManager)

	cmd := &InitCmd{
		App: &app.App{
			RepositoryLoader: repositoryLoader,
			RecipeLoader:     recipeLoader,
			ProjectLoader:    loaders.NewProjectLoader(log, conf, repositoryLoader, recipeLoader),
			TemplateManager:  modelTemplateManager,
			Sync:             syncer.New(log, modelFsManager, modelTemplateManager),
			Log:              log,
		},
		Conf: conf,
		Assets: fstest.MapFS{
			"assets/.manala.yaml.tmpl": {Data: []byte(`manala:
   recipe: {{ .Recipe.Name }}
   repository: {{ .Recipe.Repository.Source }}
`)},
		},
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
		test     string
		dir      string
		args     []string
		err      string
		stdErr   string
		stdOut   string
		manifest string
		file     string
	}{
		{
			test: "Use recipe",
			dir:  "testdata/init/project/default",
			args: []string{"--recipe", "foo"},
			stdErr: `   • Project loaded            recipe=foo repository={{ wd }}{{ ps }}testdata{{ ps }}init{{ ps }}repository{{ ps }}default
   • Repository loaded        
   • Recipe loaded            
   • Project validated        
   • Synced file               path=file_default_foo
   • Project synced           
`,
			manifest: `manala:
   recipe: foo
   repository: {{ wd }}{{ ps }}testdata{{ ps }}init{{ ps }}repository{{ ps }}default
`,
			file: "file_default_foo",
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
			stdErr: `   • Project loaded            recipe=foo repository={{ wd }}{{ ps }}testdata{{ ps }}init{{ ps }}repository{{ ps }}custom
   • Repository loaded        
   • Recipe loaded            
   • Project validated        
   • Synced file               path=file_custom_foo
   • Project synced           
`,
			manifest: `manala:
   recipe: foo
   repository: {{ wd }}{{ ps }}testdata{{ ps }}init{{ ps }}repository{{ ps }}custom
`,
			file: "file_custom_foo",
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
			// Tests - Error
			if t.err != "" {
				s.Error(err)
				s.Equal(t.err, err.Error())
			} else {
				s.NoError(err)
			}
			s.Equal(t.stdOut, stdOut.String())
			// Tests - Std
			s.Equal(
				strings.NewReplacer("{{ wd }}", s.wd, "{{ ps }}", string(os.PathSeparator)).Replace(t.stdErr),
				stdErr.String(),
			)
			// Tests - Manifest
			if t.manifest != "" {
				s.FileExists(filepath.Join(t.dir, ".manala.yaml"))
				content, _ := os.ReadFile(filepath.Join(t.dir, ".manala.yaml"))
				s.Equal(
					strings.NewReplacer("{{ wd }}", s.wd, "{{ ps }}", string(os.PathSeparator)).Replace(t.manifest),
					string(content),
				)
			}
			// Tests - File
			if t.file != "" {
				s.FileExists(filepath.Join(t.dir, t.file))
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
			// Tests - Error
			if t.err != "" {
				s.Error(err)
				s.Equal(t.err, err.Error())
			} else {
				s.NoError(err)
			}
			// Tests - Std
			s.Equal(t.stdOut, stdOut.String())
			s.Equal(
				strings.NewReplacer("{{ wd }}", s.wd, "{{ ps }}", string(os.PathSeparator)).Replace(t.stdErr),
				stdErr.String(),
			)
			// Tests - Manifest
			if t.manifest != "" {
				s.FileExists(filepath.Join(t.dir, ".manala.yaml"))
				content, _ := os.ReadFile(filepath.Join(t.dir, ".manala.yaml"))
				s.Equal(
					strings.NewReplacer("{{ wd }}", s.wd, "{{ ps }}", string(os.PathSeparator)).Replace(t.manifest),
					string(content),
				)
			}
			// Tests - File
			if t.file != "" {
				s.FileExists(filepath.Join(t.dir, t.file))
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
		// Tests - Error
		s.Error(err)
		s.Equal("project already exists: .", err.Error())
		// Tests - Std
		s.Equal("", stdOut.String())
		s.Equal("", stdErr.String())
	})
	s.Run("dir", func() {
		// Execute
		stdOut, stdErr, err := s.ExecuteCommand(
			"",
			[]string{"testdata/init/project/already_exists"},
		)
		// Tests - Error
		s.Error(err)
		s.Equal("project already exists: testdata/init/project/already_exists", err.Error())
		// Tests - Std
		s.Equal("", stdOut.String())
		s.Equal("", stdErr.String())
	})
}

func (s *InitTestSuite) TestTemplate() {
	s.Run("default", func() {
		// Clean
		_ = os.RemoveAll("testdata/init/project/default")
		_ = os.Mkdir("testdata/init/project/default", 0755)
		// Execute
		stdOut, stdErr, err := s.ExecuteCommand(
			"testdata/init/project/default",
			[]string{"--recipe", "foo", "--repository", filepath.Join(s.wd, "testdata/init/repository/template")},
		)
		// Tests - Error
		s.NoError(err)
		// Tests - Std
		s.Equal("", stdOut.String())
		s.Equal(
			strings.NewReplacer("{{ wd }}", s.wd, "{{ ps }}", string(os.PathSeparator)).Replace(`   • Project loaded            recipe=foo repository={{ wd }}{{ ps }}testdata{{ ps }}init{{ ps }}repository{{ ps }}template
   • Repository loaded        
   • Recipe loaded            
   • Project validated        
   • Project synced           
`),
			stdErr.String(),
		)
		// Tests - Manifest
		s.FileExists("testdata/init/project/default/.manala.yaml")
		content, _ := os.ReadFile("testdata/init/project/default/.manala.yaml")
		s.Equal(
			strings.NewReplacer("{{ wd }}", s.wd, "{{ ps }}", string(os.PathSeparator)).Replace(`manala:
   recipe: foo
   repository: {{ wd }}{{ ps }}testdata{{ ps }}init{{ ps }}repository{{ ps }}template

# Foo
foo:
    bar: baz
`),
			// Ensure windows CRLF conversion
			strings.NewReplacer("\r\n", "\n").Replace(string(content)),
		)
	})
}
