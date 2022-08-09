package cmd

import (
	"bytes"
	"errors"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/suite"
	internalConfig "manala/internal/config"
	internalLog "manala/internal/log"
	internalTesting "manala/internal/testing"
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

func (s *WatchSuite) TestProjectError() {

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

func (s *WatchSuite) TestRepositoryError() {

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

func (s *WatchSuite) TestRecipeError() {

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
