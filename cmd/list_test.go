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
	"testing"
)

type ListSuite struct {
	suite.Suite
	goldie   *goldie.Goldie
	config   *internalConfig.Config
	executor *cmdExecutor
}

func TestListSuite(t *testing.T) {
	suite.Run(t, new(ListSuite))
}

func (s *ListSuite) SetupTest() {
	s.goldie = goldie.New(s.T())
	s.config = internalConfig.New()
	s.executor = newCmdExecutor(func(stderr *bytes.Buffer) *cobra.Command {
		return newListCmd(
			s.config,
			internalLog.New(stderr),
		)
	})
}

func (s *ListSuite) TestRepositoryError() {

	s.Run("No Repository", func() {
		err := s.executor.execute([]string{})

		s.True(errors.As(err, &internalError))
		s.Equal("unsupported repository", internalError.Message)
		s.Empty(s.executor.stdout)
		s.Empty(s.executor.stderr)
	})

	s.Run("Repository Not Found", func() {
		err := s.executor.execute([]string{
			"--repository", internalTesting.DataPath(s, "repository"),
		})

		s.True(errors.As(err, &internalError))
		s.Equal("repository not found", internalError.Message)
		s.Empty(s.executor.stdout)
		s.Empty(s.executor.stderr)
	})

	s.Run("Wrong Repository", func() {
		err := s.executor.execute([]string{
			"--repository", internalTesting.DataPath(s, "repository"),
		})

		s.True(errors.As(err, &internalError))
		s.Equal("wrong repository", internalError.Message)
		s.Empty(s.executor.stdout)
		s.Empty(s.executor.stderr)
	})

	s.Run("Empty Repository", func() {
		err := s.executor.execute([]string{
			"--repository", internalTesting.DataPath(s, "repository"),
		})

		s.True(errors.As(err, &internalError))
		s.Equal("empty repository", internalError.Message)
		s.Empty(s.executor.stdout)
		s.Empty(s.executor.stderr)
	})
}

func (s *ListSuite) TestRecipeError() {

	s.Run("Wrong Recipe Manifest", func() {
		err := s.executor.execute([]string{
			"--repository", internalTesting.DataPath(s, "repository"),
		})

		s.True(errors.As(err, &internalError))
		s.Equal("wrong recipe manifest", internalError.Message)
		s.Empty(s.executor.stdout)
		s.Empty(s.executor.stderr)
	})

	s.Run("Invalid Recipe Manifest", func() {
		err := s.executor.execute([]string{
			"--repository", internalTesting.DataPath(s, "repository"),
		})

		s.True(errors.As(err, &internalError))
		s.Equal("recipe validation error", internalError.Message)
		s.Empty(s.executor.stdout)
		s.Empty(s.executor.stderr)
	})
}

func (s *ListSuite) Test() {

	s.Run("Custom Repository", func() {
		err := s.executor.execute([]string{
			"--repository", internalTesting.DataPath(s, "repository"),
		})

		s.NoError(err)
		s.goldie.Assert(s.T(), internalTesting.Path(s, "stdout"), s.executor.stdout.Bytes())
		s.Empty(s.executor.stderr)
	})

	s.Run("Default Repository", func() {
		s.config.Set("default-repository", internalTesting.DataPath(s, "repository"))

		err := s.executor.execute([]string{})

		s.NoError(err)
		s.goldie.Assert(s.T(), internalTesting.Path(s, "stdout"), s.executor.stdout.Bytes())
		s.Empty(s.executor.stderr)
	})
}
