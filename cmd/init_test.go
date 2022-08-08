package cmd

import (
	"bytes"
	"github.com/sebdah/goldie/v2"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/suite"
	internalConfig "manala/internal/config"
	internalLog "manala/internal/log"
	internalTesting "manala/internal/testing"
	"os"
	"testing"
)

type InitSuite struct {
	suite.Suite
	goldie   *goldie.Goldie
	config   *internalConfig.Config
	executor *cmdExecutor
}

func TestInitSuite(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	suite.Run(t, new(InitSuite))
}

func (s *InitSuite) SetupTest() {
	s.goldie = goldie.New(s.T())
	s.config = internalConfig.New()
	s.executor = newCmdExecutor(func(stderr *bytes.Buffer) *cobra.Command {
		return newInitCmd(
			s.config,
			internalLog.New(stderr),
		)
	})
}

func (s *InitSuite) TestProjectError() {

	s.Run("Already Existing Empty Project", func() {
		err := s.executor.execute([]string{
			internalTesting.DataPath(s, "project"),
		})

		s.ErrorAs(err, &internalError)
		s.Equal("empty project manifest", internalError.Message)
		s.Empty(s.executor.stdout)
		s.Empty(s.executor.stderr)
	})

	s.Run("Already Existing Project", func() {
		err := s.executor.execute([]string{
			internalTesting.DataPath(s, "project"),
		})

		s.ErrorAs(err, &internalError)
		s.Equal("already existing project", internalError.Message)
		s.Empty(s.executor.stdout)
		s.Empty(s.executor.stderr)
	})
}

func (s *InitSuite) TestRepositoryError() {

	s.Run("No Repository", func() {
		err := s.executor.execute([]string{})

		s.ErrorAs(err, &internalError)
		s.Equal("unsupported repository", internalError.Message)
		s.Empty(s.executor.stdout)
		s.Empty(s.executor.stderr)
	})

	s.Run("Repository Not Found", func() {
		err := s.executor.execute([]string{
			"--repository", internalTesting.DataPath(s, "repository"),
		})

		s.ErrorAs(err, &internalError)
		s.Equal("repository not found", internalError.Message)
		s.Empty(s.executor.stdout)
		s.Empty(s.executor.stderr)
	})

	s.Run("Empty Repository", func() {
		err := s.executor.execute([]string{
			"--repository", internalTesting.DataPath(s, "repository"),
		})

		s.ErrorAs(err, &internalError)
		s.Equal("empty repository", internalError.Message)
		s.Empty(s.executor.stdout)
		s.Empty(s.executor.stderr)
	})

	s.Run("Wrong Repository", func() {
		err := s.executor.execute([]string{
			"--repository", internalTesting.DataPath(s, "repository"),
		})

		s.ErrorAs(err, &internalError)
		s.Equal("wrong repository", internalError.Message)
		s.Empty(s.executor.stdout)
		s.Empty(s.executor.stderr)
	})
}

func (s *InitSuite) TestRecipeError() {

	s.Run("Recipe Not Found", func() {
		err := s.executor.execute([]string{
			"--repository", internalTesting.DataPath(s, "repository"),
			"--recipe", "not_found",
		})

		s.ErrorAs(err, &internalError)
		s.Equal("recipe manifest not found", internalError.Message)
		s.Empty(s.executor.stdout)
		s.Empty(s.executor.stderr)
	})

	s.Run("Wrong Recipe Manifest", func() {
		err := s.executor.execute([]string{
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
			"--repository", internalTesting.DataPath(s, "repository"),
			"--recipe", "recipe",
		})

		s.ErrorAs(err, &internalError)
		s.Equal("recipe validation error", internalError.Message)
		s.Empty(s.executor.stdout)
		s.Empty(s.executor.stderr)
	})
}

func (s *InitSuite) Test() {

	s.Run("Custom Repository", func() {
		_ = os.RemoveAll(internalTesting.DataPath(s, "project"))

		err := s.executor.execute([]string{
			internalTesting.DataPath(s, "project"),
			"--repository", internalTesting.DataPath(s, "repository"),
			"--recipe", "recipe",
		})

		s.NoError(err)
		s.Empty(s.executor.stdout)
		s.goldie.AssertWithTemplate(s.T(), internalTesting.Path(s, "stderr"), map[string]interface{}{
			"dst": internalTesting.DataPath(s, "project"),
			"src": internalTesting.DataPath(s, "repository", "recipe"),
		}, s.executor.stderr.Bytes())
		s.DirExists(internalTesting.DataPath(s, "project"))
		s.FileExists(internalTesting.DataPath(s, "project", ".manala.yaml"))
		manifestContent, _ := os.ReadFile(internalTesting.DataPath(s, "project", ".manala.yaml"))
		s.goldie.AssertWithTemplate(s.T(), internalTesting.Path(s, "manifest"), map[string]interface{}{
			"repository": internalTesting.DataPath(s, "repository"),
		}, manifestContent)
	})

	s.Run("Default Repository", func() {
		_ = os.RemoveAll(internalTesting.DataPath(s, "project"))

		s.config.Set("default-repository", internalTesting.DataPath(s, "repository"))

		err := s.executor.execute([]string{
			internalTesting.DataPath(s, "project"),
			"--recipe", "recipe",
		})

		s.NoError(err)
		s.Empty(s.executor.stdout)
		s.goldie.AssertWithTemplate(s.T(), internalTesting.Path(s, "stderr"), map[string]interface{}{
			"dst": internalTesting.DataPath(s, "project"),
			"src": internalTesting.DataPath(s, "repository", "recipe"),
		}, s.executor.stderr.Bytes())
		s.DirExists(internalTesting.DataPath(s, "project"))
		s.FileExists(internalTesting.DataPath(s, "project", ".manala.yaml"))
		manifestContent, _ := os.ReadFile(internalTesting.DataPath(s, "project", ".manala.yaml"))
		s.goldie.AssertWithTemplate(s.T(), internalTesting.Path(s, "manifest"), map[string]interface{}{
			"repository": internalTesting.DataPath(s, "repository"),
		}, manifestContent)
		s.FileExists(internalTesting.DataPath(s, "project", "file.txt"))
		fileContent, _ := os.ReadFile(internalTesting.DataPath(s, "project", "file.txt"))
		s.goldie.Assert(s.T(), internalTesting.Path(s, "file.txt"), fileContent)
		s.FileExists(internalTesting.DataPath(s, "project", "template.txt"))
		templateContent, _ := os.ReadFile(internalTesting.DataPath(s, "project", "template.txt"))
		s.goldie.Assert(s.T(), internalTesting.Path(s, "template.txt"), templateContent)
	})
}
