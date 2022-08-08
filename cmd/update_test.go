package cmd

import (
	"bytes"
	"errors"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/suite"
	internalConfig "manala/internal/config"
	internalLog "manala/internal/log"
	"os"
	"path/filepath"
	"testing"
)

type UpdateSuite struct {
	suite.Suite
	config   *internalConfig.Config
	executor *cmdExecutor
}

func TestUpdateSuite(t *testing.T) {
	suite.Run(t, new(UpdateSuite))
}

func (s *UpdateSuite) SetupTest() {
	s.config = internalConfig.New()
	s.executor = newCmdExecutor(func(stderr *bytes.Buffer) *cobra.Command {
		return newUpdateCmd(
			s.config,
			internalLog.New(stderr),
		)
	})
}

var updateTestRepositoryPath = filepath.Join("testdata", "update", "repository")
var updateTestProjectPath = filepath.Join("testdata", "update", "project")

func (s *UpdateSuite) TestProjectError() {

	s.Run("Project Not Found", func() {
		err := s.executor.execute([]string{
			filepath.Join(updateTestProjectPath, "not_found"),
		})

		s.True(errors.As(err, &internalError))
		s.Equal("project manifest not found", internalError.Message)
		s.Empty(s.executor.stdout.String())
		s.Empty(s.executor.stderr.String())
	})

	s.Run("Wrong Project Manifest", func() {
		err := s.executor.execute([]string{
			filepath.Join(updateTestProjectPath, "wrong_manifest"),
		})

		s.True(errors.As(err, &internalError))
		s.Equal("wrong project manifest", internalError.Message)
		s.Empty(s.executor.stdout.String())
		s.Empty(s.executor.stderr.String())
	})

	s.Run("Empty Project Manifest", func() {
		err := s.executor.execute([]string{
			filepath.Join(updateTestProjectPath, "empty_manifest"),
		})

		s.True(errors.As(err, &internalError))
		s.Equal("empty project manifest", internalError.Message)
		s.Empty(s.executor.stdout.String())
		s.Empty(s.executor.stderr.String())
	})

	s.Run("Invalid Project Manifest", func() {
		err := s.executor.execute([]string{
			filepath.Join(updateTestProjectPath, "invalid_manifest"),
		})

		s.True(errors.As(err, &internalError))
		s.Equal("project validation error", internalError.Message)
		s.Empty(s.executor.stdout.String())
		s.Empty(s.executor.stderr.String())
	})
}

func (s *UpdateSuite) TestRecursiveProjectError() {

	s.Run("Project Not Found", func() {
		err := s.executor.execute([]string{
			filepath.Join(updateTestProjectPath, "not_found"),
			"--recursive",
		})

		s.NoError(err)
		s.Empty(s.executor.stdout.String())
		s.Equal(`   • load projects from        path=`+filepath.Join(updateTestProjectPath, "not_found")+`
`, s.executor.stderr.String())
	})

	s.Run("Wrong Project Manifest", func() {
		err := s.executor.execute([]string{
			filepath.Join(updateTestProjectPath, "wrong_manifest"),
			"--recursive",
		})

		s.True(errors.As(err, &internalError))
		s.Equal("wrong project manifest", internalError.Message)
		s.Empty(s.executor.stdout.String())
		s.Equal(`   • load projects from        path=`+filepath.Join(updateTestProjectPath, "wrong_manifest")+`
`, s.executor.stderr.String())
	})

	s.Run("Empty Project Manifest", func() {
		err := s.executor.execute([]string{
			filepath.Join(updateTestProjectPath, "empty_manifest"),
			"--recursive",
		})

		s.True(errors.As(err, &internalError))
		s.Equal("empty project manifest", internalError.Message)
		s.Empty(s.executor.stdout.String())
		s.Equal(`   • load projects from        path=`+filepath.Join(updateTestProjectPath, "empty_manifest")+`
`, s.executor.stderr.String())
	})

	s.Run("Invalid Project Manifest", func() {
		err := s.executor.execute([]string{
			filepath.Join(updateTestProjectPath, "invalid_manifest"),
			"--recursive",
		})

		s.True(errors.As(err, &internalError))
		s.Equal("project validation error", internalError.Message)
		s.Empty(s.executor.stdout.String())
		s.Equal(`   • load projects from        path=`+filepath.Join(updateTestProjectPath, "invalid_manifest")+`
`, s.executor.stderr.String())
	})
}

func (s *UpdateSuite) TestRepositoryError() {

	s.Run("No Repository", func() {
		err := s.executor.execute([]string{
			filepath.Join(updateTestProjectPath, "no_repository"),
		})

		s.True(errors.As(err, &internalError))
		s.Equal("unsupported repository", internalError.Message)
		s.Empty(s.executor.stdout.String())
		s.Empty(s.executor.stderr.String())
	})

	s.Run("Repository Not Found", func() {
		err := s.executor.execute([]string{
			filepath.Join(updateTestProjectPath, "no_repository"),
			"--repository", filepath.Join(updateTestRepositoryPath, "not_found"),
		})

		s.True(errors.As(err, &internalError))
		s.Equal("repository not found", internalError.Message)
		s.Empty(s.executor.stdout.String())
		s.Empty(s.executor.stderr.String())
	})

	s.Run("Wrong Repository", func() {
		err := s.executor.execute([]string{
			filepath.Join(updateTestProjectPath, "wrong_repository"),
			"--repository", filepath.Join(updateTestRepositoryPath, "wrong"),
		})

		s.True(errors.As(err, &internalError))
		s.Equal("wrong repository", internalError.Message)
		s.Empty(s.executor.stdout.String())
		s.Empty(s.executor.stderr.String())
	})
}

func (s *UpdateSuite) TestRecipeError() {
	s.config.Set("default-repository", filepath.Join(updateTestRepositoryPath, "default"))

	s.Run("Recipe Not Found", func() {
		err := s.executor.execute([]string{
			filepath.Join(updateTestProjectPath, "recipe_not_found"),
			"--recipe", "not_found",
		})

		s.True(errors.As(err, &internalError))
		s.Equal("recipe manifest not found", internalError.Message)
		s.Empty(s.executor.stdout.String())
		s.Empty(s.executor.stderr.String())
	})

	s.Run("Wrong Recipe Manifest", func() {
		err := s.executor.execute([]string{
			filepath.Join(updateTestProjectPath, "wrong_recipe_manifest"),
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
			filepath.Join(updateTestProjectPath, "invalid_recipe_manifest"),
			"--repository", filepath.Join(initTestRepositoryPath, "invalid_recipe"),
			"--recipe", "recipe",
		})

		s.ErrorAs(err, &internalError)
		s.Equal("recipe validation error", internalError.Message)
		s.Empty(s.executor.stdout.String())
		s.Empty(s.executor.stderr.String())
	})
}

func (s *UpdateSuite) Test() {
	s.config.Set("default-repository", filepath.Join(updateTestRepositoryPath, "default"))

	s.Run("Custom Repository", func() {
		_ = os.Remove(filepath.Join(updateTestProjectPath, "custom_repository", "file"))

		err := s.executor.execute([]string{
			filepath.Join(updateTestProjectPath, "custom_repository"),
			"--repository", filepath.Join(updateTestRepositoryPath, "custom"),
		})

		s.NoError(err)
		s.Empty(s.executor.stdout.String())
		s.Equal(`   • project loaded            path=`+filepath.Join(updateTestProjectPath, "custom_repository")+` recipe=recipe repository=`+filepath.Join(updateTestRepositoryPath, "custom")+`
      • sync project              dst=`+filepath.Join(updateTestProjectPath, "custom_repository")+` src=`+filepath.Join(updateTestRepositoryPath, "custom", "recipe")+`
         • file synced               path=file
`, s.executor.stderr.String())
		s.FileExists(filepath.Join(updateTestProjectPath, "custom_repository", "file"))
		fileContent, _ := os.ReadFile(filepath.Join(updateTestProjectPath, "custom_repository", "file"))
		s.Equal(`Custom recipe file`, string(fileContent))
	})

	s.Run("Custom Recipe", func() {
		_ = os.Remove(filepath.Join(updateTestProjectPath, "custom_recipe", "file"))

		err := s.executor.execute([]string{
			filepath.Join(updateTestProjectPath, "custom_recipe"),
			"--recipe", "custom",
		})

		s.NoError(err)
		s.Empty(s.executor.stdout.String())
		s.Equal(`   • project loaded            path=`+filepath.Join(updateTestProjectPath, "custom_recipe")+` recipe=custom repository=`+filepath.Join(updateTestRepositoryPath, "default")+`
      • sync project              dst=`+filepath.Join(updateTestProjectPath, "custom_recipe")+` src=`+filepath.Join(updateTestRepositoryPath, "default", "custom")+`
         • file synced               path=file
`, s.executor.stderr.String())
		s.FileExists(filepath.Join(updateTestProjectPath, "custom_repository", "file"))
		fileContent, _ := os.ReadFile(filepath.Join(updateTestProjectPath, "custom_recipe", "file"))
		s.Equal(`Default custom file`, string(fileContent))
	})

	s.Run("Default Repository", func() {
		_ = os.Remove(filepath.Join(updateTestProjectPath, "default_repository", "file"))
		_ = os.Remove(filepath.Join(updateTestProjectPath, "default_repository", "template"))

		err := s.executor.execute([]string{
			filepath.Join(updateTestProjectPath, "default_repository"),
		})

		s.NoError(err)
		s.Empty(s.executor.stdout.String())
		s.Equal(`   • project loaded            path=`+filepath.Join(updateTestProjectPath, "default_repository")+` recipe=recipe repository=`+filepath.Join(updateTestRepositoryPath, "default")+`
      • sync project              dst=`+filepath.Join(updateTestProjectPath, "default_repository")+` src=`+filepath.Join(updateTestRepositoryPath, "default", "recipe")+`
         • file synced               path=file
         • file synced               path=template
`, s.executor.stderr.String())
		s.FileExists(filepath.Join(updateTestProjectPath, "default_repository", "file"))
		fileContent, _ := os.ReadFile(filepath.Join(updateTestProjectPath, "default_repository", "file"))
		s.Equal(`Default recipe file`, string(fileContent))
		s.FileExists(filepath.Join(updateTestProjectPath, "default_repository", "template"))
		templateContent, _ := os.ReadFile(filepath.Join(updateTestProjectPath, "default_repository", "template"))
		s.Equal(`Default recipe template
foo: foo
bar: bar`, string(templateContent))
	})
}
