package cmd

import (
	"bytes"
	"errors"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/suite"
	internalConfig "manala/internal/config"
	internalLog "manala/internal/log"
	"path/filepath"
	"testing"
)

type WatchSuite struct {
	suite.Suite
	config   *internalConfig.Config
	executor *cmdExecutor
}

func TestWatchSuite(t *testing.T) {
	suite.Run(t, new(WatchSuite))
}

func (s *WatchSuite) SetupTest() {
	s.config = internalConfig.New()
	s.executor = newCmdExecutor(func(stderr *bytes.Buffer) *cobra.Command {
		return newWatchCmd(
			s.config,
			internalLog.New(stderr),
		)
	})
}

var watchTestRepositoryPath = filepath.Join("testdata", "watch", "repository")
var watchTestProjectPath = filepath.Join("testdata", "watch", "project")

func (s *WatchSuite) TestProjectError() {

	s.Run("Project Not Found", func() {
		err := s.executor.execute([]string{
			filepath.Join(watchTestProjectPath, "not_found"),
		})

		s.True(errors.As(err, &internalError))
		s.Equal("project manifest not found", internalError.Message)
		s.Empty(s.executor.stdout.String())
		s.Empty(s.executor.stderr.String())
	})

	s.Run("Wrong Project Manifest", func() {
		err := s.executor.execute([]string{
			filepath.Join(watchTestProjectPath, "wrong_manifest"),
		})

		s.True(errors.As(err, &internalError))
		s.Equal("wrong project manifest", internalError.Message)
		s.Empty(s.executor.stdout.String())
		s.Empty(s.executor.stderr.String())
	})

	s.Run("Empty Project Manifest", func() {
		err := s.executor.execute([]string{
			filepath.Join(watchTestProjectPath, "empty_manifest"),
		})

		s.True(errors.As(err, &internalError))
		s.Equal("empty project manifest", internalError.Message)
		s.Empty(s.executor.stdout.String())
		s.Empty(s.executor.stderr.String())
	})

	s.Run("Invalid Project Manifest", func() {
		err := s.executor.execute([]string{
			filepath.Join(watchTestProjectPath, "invalid_manifest"),
		})

		s.True(errors.As(err, &internalError))
		s.Equal("project validation error", internalError.Message)
		s.Empty(s.executor.stdout.String())
		s.Empty(s.executor.stderr.String())
	})
}

func (s *WatchSuite) TestRepositoryError() {

	s.Run("No Repository", func() {
		err := s.executor.execute([]string{
			filepath.Join(watchTestProjectPath, "no_repository"),
		})

		s.True(errors.As(err, &internalError))
		s.Equal("unsupported repository", internalError.Message)
		s.Empty(s.executor.stdout.String())
		s.Empty(s.executor.stderr.String())
	})

	s.Run("Repository Not Found", func() {
		err := s.executor.execute([]string{
			filepath.Join(watchTestProjectPath, "no_repository"),
			"--repository", filepath.Join(watchTestRepositoryPath, "not_found"),
		})

		s.True(errors.As(err, &internalError))
		s.Equal("repository not found", internalError.Message)
		s.Empty(s.executor.stdout.String())
		s.Empty(s.executor.stderr.String())
	})

	s.Run("Wrong Repository", func() {
		err := s.executor.execute([]string{
			filepath.Join(watchTestProjectPath, "wrong_repository"),
			"--repository", filepath.Join(watchTestRepositoryPath, "wrong"),
		})

		s.True(errors.As(err, &internalError))
		s.Equal("wrong repository", internalError.Message)
		s.Empty(s.executor.stdout.String())
		s.Empty(s.executor.stderr.String())
	})
}

func (s *WatchSuite) TestRecipeError() {
	s.config.Set("default-repository", filepath.Join(watchTestRepositoryPath, "default"))

	s.Run("Recipe Not Found", func() {
		err := s.executor.execute([]string{
			filepath.Join(watchTestProjectPath, "recipe_not_found"),
			"--recipe", "not_found",
		})

		s.True(errors.As(err, &internalError))
		s.Equal("recipe manifest not found", internalError.Message)
		s.Empty(s.executor.stdout.String())
		s.Empty(s.executor.stderr.String())
	})

	s.Run("Wrong Recipe Manifest", func() {
		err := s.executor.execute([]string{
			filepath.Join(watchTestProjectPath, "wrong_recipe_manifest"),
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
			filepath.Join(watchTestProjectPath, "invalid_recipe_manifest"),
			"--repository", filepath.Join(initTestRepositoryPath, "invalid_recipe"),
			"--recipe", "recipe",
		})

		s.ErrorAs(err, &internalError)
		s.Equal("recipe validation error", internalError.Message)
		s.Empty(s.executor.stdout.String())
		s.Empty(s.executor.stderr.String())
	})
}
