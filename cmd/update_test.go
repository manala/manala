package cmd

import (
	"bytes"
	"errors"
	"github.com/sebdah/goldie/v2"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/suite"
	internalConfig "manala/internal/config"
	internalLog "manala/internal/log"
	internalTesting "manala/internal/testing"
	"os"
	"testing"
)

type UpdateSuite struct {
	suite.Suite
	goldie   *goldie.Goldie
	config   *internalConfig.Config
	executor *cmdExecutor
}

func TestUpdateSuite(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	suite.Run(t, new(UpdateSuite))
}

func (s *UpdateSuite) SetupTest() {
	s.goldie = goldie.New(s.T())
	s.config = internalConfig.New()
	s.executor = newCmdExecutor(func(stderr *bytes.Buffer) *cobra.Command {
		return newUpdateCmd(
			s.config,
			internalLog.New(stderr),
		)
	})
}

func (s *UpdateSuite) TestProjectError() {

	s.Run("Project Not Found", func() {
		err := s.executor.execute([]string{
			internalTesting.DataPath(s, "project"),
		})

		s.True(errors.As(err, &internalError))
		s.Equal("project manifest not found", internalError.Message)
		s.Empty(s.executor.stdout)
		s.Empty(s.executor.stderr)
	})

	s.Run("Wrong Project Manifest", func() {
		err := s.executor.execute([]string{
			internalTesting.DataPath(s, "project"),
		})

		s.True(errors.As(err, &internalError))
		s.Equal("wrong project manifest", internalError.Message)
		s.Empty(s.executor.stdout)
		s.Empty(s.executor.stderr)
	})

	s.Run("Empty Project Manifest", func() {
		err := s.executor.execute([]string{
			internalTesting.DataPath(s, "project"),
		})

		s.True(errors.As(err, &internalError))
		s.Equal("empty project manifest", internalError.Message)
		s.Empty(s.executor.stdout)
		s.Empty(s.executor.stderr)
	})

	s.Run("Invalid Project Manifest", func() {
		err := s.executor.execute([]string{
			internalTesting.DataPath(s, "project"),
		})

		s.True(errors.As(err, &internalError))
		s.Equal("project validation error", internalError.Message)
		s.Empty(s.executor.stdout)
		s.Empty(s.executor.stderr)
	})
}

func (s *UpdateSuite) TestRecursiveProjectError() {

	s.Run("Project Not Found", func() {
		err := s.executor.execute([]string{
			internalTesting.DataPath(s, "project"),
			"--recursive",
		})

		s.NoError(err)
		s.Empty(s.executor.stdout)
		s.goldie.AssertWithTemplate(s.T(), internalTesting.Path(s, "stderr"), map[string]interface{}{
			"path": internalTesting.DataPath(s, "project"),
		}, s.executor.stderr.Bytes())
	})

	s.Run("Wrong Project Manifest", func() {
		err := s.executor.execute([]string{
			internalTesting.DataPath(s, "project"),
			"--recursive",
		})

		s.True(errors.As(err, &internalError))
		s.Equal("wrong project manifest", internalError.Message)
		s.Empty(s.executor.stdout)
		s.goldie.AssertWithTemplate(s.T(), internalTesting.Path(s, "stderr"), map[string]interface{}{
			"path": internalTesting.DataPath(s, "project"),
		}, s.executor.stderr.Bytes())
	})

	s.Run("Empty Project Manifest", func() {
		err := s.executor.execute([]string{
			internalTesting.DataPath(s, "project"),
			"--recursive",
		})

		s.True(errors.As(err, &internalError))
		s.Equal("empty project manifest", internalError.Message)
		s.Empty(s.executor.stdout)
		s.goldie.AssertWithTemplate(s.T(), internalTesting.Path(s, "stderr"), map[string]interface{}{
			"path": internalTesting.DataPath(s, "project"),
		}, s.executor.stderr.Bytes())
	})

	s.Run("Invalid Project Manifest", func() {
		err := s.executor.execute([]string{
			internalTesting.DataPath(s, "project"),
			"--recursive",
		})

		s.True(errors.As(err, &internalError))
		s.Equal("project validation error", internalError.Message)
		s.Empty(s.executor.stdout)
		s.goldie.AssertWithTemplate(s.T(), internalTesting.Path(s, "stderr"), map[string]interface{}{
			"path": internalTesting.DataPath(s, "project"),
		}, s.executor.stderr.Bytes())
	})
}

func (s *UpdateSuite) TestRepositoryError() {

	s.Run("No Repository", func() {
		err := s.executor.execute([]string{
			internalTesting.DataPath(s, "project"),
		})

		s.True(errors.As(err, &internalError))
		s.Equal("unsupported repository", internalError.Message)
		s.Empty(s.executor.stdout)
		s.Empty(s.executor.stderr)
	})

	s.Run("Repository Not Found", func() {
		err := s.executor.execute([]string{
			internalTesting.DataPath(s, "project"),
			"--repository", internalTesting.DataPath(s, "repository"),
		})

		s.True(errors.As(err, &internalError))
		s.Equal("repository not found", internalError.Message)
		s.Empty(s.executor.stdout)
		s.Empty(s.executor.stderr)
	})

	s.Run("Wrong Repository", func() {
		err := s.executor.execute([]string{
			internalTesting.DataPath(s, "project"),
			"--repository", internalTesting.DataPath(s, "repository"),
		})

		s.True(errors.As(err, &internalError))
		s.Equal("wrong repository", internalError.Message)
		s.Empty(s.executor.stdout)
		s.Empty(s.executor.stderr)
	})
}

func (s *UpdateSuite) TestRecipeError() {

	s.Run("Recipe Not Found", func() {
		s.config.Set("default-repository", internalTesting.DataPath(s, "repository"))

		err := s.executor.execute([]string{
			internalTesting.DataPath(s, "project"),
			"--recipe", "recipe",
		})

		s.True(errors.As(err, &internalError))
		s.Equal("recipe manifest not found", internalError.Message)
		s.Empty(s.executor.stdout)
		s.Empty(s.executor.stderr)
	})

	s.Run("Wrong Recipe Manifest", func() {
		err := s.executor.execute([]string{
			internalTesting.DataPath(s, "project"),
			"--repository", internalTesting.DataPath(s, "repository"),
			"--recipe", "recipe",
		})

		s.ErrorAs(err, &internalError)
		s.Equal("wrong recipe manifest", internalError.Message)
		s.Empty(s.executor.stdout)
		s.Empty(s.executor.stderr)
	})

	s.Run("Invalid Recipe Manifest", func() {
		err := s.executor.execute([]string{
			internalTesting.DataPath(s, "project"),
			"--repository", internalTesting.DataPath(s, "repository"),
			"--recipe", "recipe",
		})

		s.ErrorAs(err, &internalError)
		s.Equal("recipe validation error", internalError.Message)
		s.Empty(s.executor.stdout)
		s.Empty(s.executor.stderr)
	})
}

func (s *UpdateSuite) Test() {

	s.Run("Custom Repository", func() {
		_ = os.Remove(internalTesting.DataPath(s, "project", "file.txt"))

		err := s.executor.execute([]string{
			internalTesting.DataPath(s, "project"),
			"--repository", internalTesting.DataPath(s, "repository"),
		})

		s.NoError(err)
		s.Empty(s.executor.stdout)
		s.goldie.AssertWithTemplate(s.T(), internalTesting.Path(s, "stderr"), map[string]interface{}{
			"path":       internalTesting.DataPath(s, "project"),
			"repository": internalTesting.DataPath(s, "repository"),
			"dst":        internalTesting.DataPath(s, "project"),
			"src":        internalTesting.DataPath(s, "repository", "recipe"),
		}, s.executor.stderr.Bytes())
		s.FileExists(internalTesting.DataPath(s, "project", "file.txt"))
		fileContent, _ := os.ReadFile(internalTesting.DataPath(s, "project", "file.txt"))
		s.goldie.Assert(s.T(), internalTesting.Path(s, "file.txt"), fileContent)
	})

	s.Run("Custom Recipe", func() {
		_ = os.Remove(internalTesting.DataPath(s, "project", "file.txt"))

		s.config.Set("default-repository", internalTesting.DataPath(s, "repository"))

		err := s.executor.execute([]string{
			internalTesting.DataPath(s, "project"),
			"--recipe", "recipe",
		})

		s.NoError(err)
		s.Empty(s.executor.stdout)
		s.goldie.AssertWithTemplate(s.T(), internalTesting.Path(s, "stderr"), map[string]interface{}{
			"path":       internalTesting.DataPath(s, "project"),
			"repository": internalTesting.DataPath(s, "repository"),
			"dst":        internalTesting.DataPath(s, "project"),
			"src":        internalTesting.DataPath(s, "repository", "recipe"),
		}, s.executor.stderr.Bytes())
		s.FileExists(internalTesting.DataPath(s, "project", "file.txt"))
		fileContent, _ := os.ReadFile(internalTesting.DataPath(s, "project", "file.txt"))
		s.goldie.Assert(s.T(), internalTesting.Path(s, "file.txt"), fileContent)
	})

	s.Run("Default Repository", func() {
		_ = os.Remove(internalTesting.DataPath(s, "project", "file.txt"))
		_ = os.Remove(internalTesting.DataPath(s, "project", "template.txt"))

		s.config.Set("default-repository", internalTesting.DataPath(s, "repository"))

		err := s.executor.execute([]string{
			internalTesting.DataPath(s, "project"),
		})

		s.NoError(err)
		s.Empty(s.executor.stdout)
		s.goldie.AssertWithTemplate(s.T(), internalTesting.Path(s, "stderr"), map[string]interface{}{
			"path":       internalTesting.DataPath(s, "project"),
			"repository": internalTesting.DataPath(s, "repository"),
			"dst":        internalTesting.DataPath(s, "project"),
			"src":        internalTesting.DataPath(s, "repository", "recipe"),
		}, s.executor.stderr.Bytes())
		s.FileExists(internalTesting.DataPath(s, "project", "file.txt"))
		fileContent, _ := os.ReadFile(internalTesting.DataPath(s, "project", "file.txt"))
		s.goldie.Assert(s.T(), internalTesting.Path(s, "file.txt"), fileContent)
		s.FileExists(internalTesting.DataPath(s, "project", "template.txt"))
		templateContent, _ := os.ReadFile(internalTesting.DataPath(s, "project", "template.txt"))
		s.goldie.Assert(s.T(), internalTesting.Path(s, "template.txt"), templateContent)
	})
}
