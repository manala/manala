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

type ListSuite struct {
	suite.Suite
	config   *internalConfig.Config
	executor *cmdExecutor
}

func TestListSuite(t *testing.T) {
	suite.Run(t, new(ListSuite))
}

func (s *ListSuite) SetupTest() {
	s.config = internalConfig.New()
	s.executor = newCmdExecutor(func(stderr *bytes.Buffer) *cobra.Command {
		return newListCmd(
			s.config,
			internalLog.New(stderr),
		)
	})
}

var listTestRepositoryPath = filepath.Join("testdata", "list", "repository")

func (s *ListSuite) TestRepositoryError() {

	s.Run("No Repository", func() {
		err := s.executor.execute([]string{})

		s.True(errors.As(err, &internalError))
		s.Equal("unsupported repository", internalError.Message)
		s.Empty(s.executor.stdout.String())
		s.Empty(s.executor.stderr.String())
	})

	s.Run("Repository Not Found", func() {
		err := s.executor.execute([]string{
			"--repository", filepath.Join(listTestRepositoryPath, "not_found"),
		})

		s.True(errors.As(err, &internalError))
		s.Equal("repository not found", internalError.Message)
		s.Empty(s.executor.stdout.String())
		s.Empty(s.executor.stderr.String())
	})

	s.Run("Wrong Repository", func() {
		err := s.executor.execute([]string{
			"--repository", filepath.Join(listTestRepositoryPath, "wrong"),
		})

		s.True(errors.As(err, &internalError))
		s.Equal("wrong repository", internalError.Message)
		s.Empty(s.executor.stdout.String())
		s.Empty(s.executor.stderr.String())
	})

	s.Run("Empty Repository", func() {
		err := s.executor.execute([]string{
			"--repository", filepath.Join(listTestRepositoryPath, "empty"),
		})

		s.True(errors.As(err, &internalError))
		s.Equal("empty repository", internalError.Message)
		s.Empty(s.executor.stdout.String())
		s.Empty(s.executor.stderr.String())
	})
}

func (s *ListSuite) TestRecipeError() {

	s.Run("Wrong Recipe Manifest", func() {
		err := s.executor.execute([]string{
			"--repository", filepath.Join(listTestRepositoryPath, "wrong_recipe"),
		})

		s.True(errors.As(err, &internalError))
		s.Equal("wrong recipe manifest", internalError.Message)
		s.Empty(s.executor.stdout.String())
		s.Empty(s.executor.stderr.String())
	})

	s.Run("Invalid Recipe Manifest", func() {
		err := s.executor.execute([]string{
			"--repository", filepath.Join(listTestRepositoryPath, "invalid_recipe"),
		})

		s.True(errors.As(err, &internalError))
		s.Equal("recipe validation error", internalError.Message)
		s.Empty(s.executor.stdout.String())
		s.Empty(s.executor.stderr.String())
	})
}

func (s *ListSuite) Test() {

	s.config.Set("default-repository", filepath.Join(listTestRepositoryPath, "default"))

	s.Run("Custom Repository", func() {
		err := s.executor.execute([]string{
			"--repository", filepath.Join(listTestRepositoryPath, "custom"),
		})

		s.NoError(err)
		s.Contains(s.executor.stdout.String(), "foo: Custom foo recipe")
		s.Contains(s.executor.stdout.String(), "bar: Custom bar recipe")
		s.Empty(s.executor.stderr.String())
	})

	s.Run("Default Repository", func() {
		err := s.executor.execute([]string{})

		s.NoError(err)
		s.Contains(s.executor.stdout.String(), "foo: Default foo recipe")
		s.Contains(s.executor.stdout.String(), "bar: Default bar recipe")
		s.Empty(s.executor.stderr.String())
	})
}
