package cmd

import (
	"bytes"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/suite"
	internalConfig "manala/internal/config"
	internalLog "manala/internal/log"
	"os"
	"path/filepath"
	"testing"
)

type InitSuite struct {
	suite.Suite
	config   *internalConfig.Config
	executor *cmdExecutor
}

func TestInitSuite(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	suite.Run(t, new(InitSuite))
}

func (s *InitSuite) SetupTest() {
	s.config = internalConfig.New()
	s.executor = newCmdExecutor(func(stderr *bytes.Buffer) *cobra.Command {
		return newInitCmd(
			s.config,
			internalLog.New(stderr),
		)
	})
}

var initTestRepositoryPath = filepath.Join("testdata", "init", "repository")
var initTestProjectPath = filepath.Join("testdata", "init", "project")

func (s *InitSuite) TestProjectError() {

	s.Run("Already Existing Empty Project", func() {
		err := s.executor.execute([]string{
			filepath.Join(initTestProjectPath, "already_existing_empty"),
		})

		s.ErrorAs(err, &internalError)
		s.Equal("empty project manifest", internalError.Message)
		s.Empty(s.executor.stdout.String())
		s.Empty(s.executor.stderr.String())
	})

	s.Run("Already Existing Project", func() {
		err := s.executor.execute([]string{
			filepath.Join(initTestProjectPath, "already_existing"),
		})

		s.ErrorAs(err, &internalError)
		s.Equal("already existing project", internalError.Message)
		s.Empty(s.executor.stdout.String())
		s.Empty(s.executor.stderr.String())
	})
}

func (s *InitSuite) TestRepositoryError() {

	s.Run("No Repository", func() {
		err := s.executor.execute([]string{})

		s.ErrorAs(err, &internalError)
		s.Equal("unsupported repository", internalError.Message)
		s.Empty(s.executor.stdout.String())
		s.Empty(s.executor.stderr.String())
	})

	s.Run("Repository Not Found", func() {
		err := s.executor.execute([]string{
			"--repository", filepath.Join(initTestRepositoryPath, "not_found"),
		})

		s.ErrorAs(err, &internalError)
		s.Equal("repository not found", internalError.Message)
		s.Empty(s.executor.stdout.String())
		s.Empty(s.executor.stderr.String())
	})

	s.Run("Empty Repository", func() {
		err := s.executor.execute([]string{
			"--repository", filepath.Join(initTestRepositoryPath, "empty"),
		})

		s.ErrorAs(err, &internalError)
		s.Equal("empty repository", internalError.Message)
		s.Empty(s.executor.stdout.String())
		s.Empty(s.executor.stderr.String())
	})

	s.Run("Wrong Repository", func() {
		err := s.executor.execute([]string{
			"--repository", filepath.Join(initTestRepositoryPath, "wrong"),
		})

		s.ErrorAs(err, &internalError)
		s.Equal("wrong repository", internalError.Message)
		s.Empty(s.executor.stdout.String())
		s.Empty(s.executor.stderr.String())
	})
}

func (s *InitSuite) TestRecipeError() {

	s.Run("Recipe Not Found", func() {
		err := s.executor.execute([]string{
			"--repository", filepath.Join(initTestRepositoryPath, "default"),
			"--recipe", "not_found",
		})

		s.ErrorAs(err, &internalError)
		s.Equal("recipe manifest not found", internalError.Message)
		s.Empty(s.executor.stdout.String())
		s.Empty(s.executor.stderr.String())
	})

	s.Run("Wrong Recipe Manifest", func() {
		err := s.executor.execute([]string{
			"--repository", filepath.Join(initTestRepositoryPath, "wrong_recipe"),
			"--recipe", "recipe",
		})

		s.ErrorAs(err, &internalError)
		s.Equal("wrong recipe manifest", internalError.Message)
		s.Empty(s.executor.stdout.String())
		s.Empty(s.executor.stderr.String())
	})

	s.Run("Invalid Recipe Manifest", func() {
		err := s.executor.execute([]string{
			"--repository", filepath.Join(initTestRepositoryPath, "invalid_recipe"),
			"--recipe", "recipe",
		})

		s.ErrorAs(err, &internalError)
		s.Equal("recipe validation error", internalError.Message)
		s.Empty(s.executor.stdout.String())
		s.Empty(s.executor.stderr.String())
	})
}

func (s *InitSuite) Test() {

	s.config.Set("default-repository", filepath.Join(initTestRepositoryPath, "default"))

	s.Run("Custom Repository", func() {
		_ = os.RemoveAll(filepath.Join(initTestProjectPath, "custom"))

		err := s.executor.execute([]string{
			filepath.Join(initTestProjectPath, "custom"),
			"--repository", filepath.Join(initTestRepositoryPath, "custom"),
			"--recipe", "recipe",
		})

		s.NoError(err)
		s.Empty(s.executor.stdout.String())
		s.Equal(`    • sync project                                   dst=`+filepath.Join(initTestProjectPath, "custom")+` src=`+filepath.Join(initTestRepositoryPath, "custom", "recipe")+`
      • file synced                                  path=file
`, s.executor.stderr.String())

		s.DirExists(filepath.Join(initTestProjectPath, "custom"))
		s.FileExists(filepath.Join(initTestProjectPath, "custom", ".manala.yaml"))
		s.FileExists(filepath.Join(initTestProjectPath, "custom", "file"))
		manifestContent, _ := os.ReadFile(filepath.Join(initTestProjectPath, "custom", ".manala.yaml"))
		s.Equal(`####################################################################
#                         !!! REMINDER !!!                         #
# Don't forget to run `+"`manala up`"+` each time you update this file ! #
####################################################################

manala:
    recipe: recipe
    repository: `+filepath.Join(initTestRepositoryPath, "custom")+`
`, string(manifestContent))
		fileContent, _ := os.ReadFile(filepath.Join(initTestProjectPath, "custom", "file"))
		s.Equal(`Custom recipe file`, string(fileContent))
	})

	s.Run("Default Repository", func() {
		_ = os.RemoveAll(filepath.Join(initTestProjectPath, "default"))

		err := s.executor.execute([]string{
			filepath.Join(initTestProjectPath, "default"),
			"--recipe", "recipe",
		})

		s.NoError(err)
		s.Empty(s.executor.stdout.String())
		s.Equal(`    • sync project                                   dst=`+filepath.Join(initTestProjectPath, "default")+` src=`+filepath.Join(initTestRepositoryPath, "default", "recipe")+`
      • file synced                                  path=file
      • file synced                                  path=template
`, s.executor.stderr.String())

		s.DirExists(filepath.Join(initTestProjectPath, "default"))
		s.FileExists(filepath.Join(initTestProjectPath, "default", ".manala.yaml"))
		manifestContent, _ := os.ReadFile(filepath.Join(initTestProjectPath, "default", ".manala.yaml"))
		s.Equal(`manala:
    recipe: recipe
    repository: `+filepath.Join(initTestRepositoryPath, "default")+`
foo: bar
bar: foo
`, string(manifestContent))
		s.FileExists(filepath.Join(initTestProjectPath, "default", "file"))
		fileContent, _ := os.ReadFile(filepath.Join(initTestProjectPath, "default", "file"))
		s.Equal(`Default recipe file`, string(fileContent))
		s.FileExists(filepath.Join(initTestProjectPath, "default", "template"))
		templateContent, _ := os.ReadFile(filepath.Join(initTestProjectPath, "default", "template"))
		s.Equal(`Default recipe template
foo: bar
bar: foo`, string(templateContent))
	})
}
