package cmd

import (
	"bytes"
	"github.com/sebdah/goldie/v2"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/suite"
	internalConfig "manala/internal/config"
	internalLog "manala/internal/log"
	internalReport "manala/internal/report"
	internalTesting "manala/internal/testing"
	"path/filepath"
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

		s.Error(err)
		s.Empty(s.executor.stdout)
		s.Empty(s.executor.stderr)

		report := internalReport.NewErrorReport(err)

		reportAssert := &internalReport.Assert{
			Err: "unsupported repository",
		}
		reportAssert.Equal(&s.Suite, report)
	})

	s.Run("Repository Not Found", func() {
		repoPath := internalTesting.DataPath(s, "repository")

		err := s.executor.execute([]string{
			"--repository", repoPath,
		})

		s.Error(err)
		s.Empty(s.executor.stdout)
		s.Empty(s.executor.stderr)

		report := internalReport.NewErrorReport(err)

		reportAssert := &internalReport.Assert{
			Err: "repository not found",
			Fields: map[string]interface{}{
				"path": repoPath,
			},
		}
		reportAssert.Equal(&s.Suite, report)
	})

	s.Run("Wrong Repository", func() {
		repoPath := internalTesting.DataPath(s, "repository")

		err := s.executor.execute([]string{
			"--repository", repoPath,
		})

		s.Error(err)
		s.Empty(s.executor.stdout)
		s.Empty(s.executor.stderr)

		report := internalReport.NewErrorReport(err)

		reportAssert := &internalReport.Assert{
			Err: "wrong repository",
			Fields: map[string]interface{}{
				"dir": repoPath,
			},
		}
		reportAssert.Equal(&s.Suite, report)
	})

	s.Run("Empty Repository", func() {
		repoPath := internalTesting.DataPath(s, "repository")

		err := s.executor.execute([]string{
			"--repository", repoPath,
		})

		s.Error(err)
		s.Empty(s.executor.stdout)
		s.Empty(s.executor.stderr)

		report := internalReport.NewErrorReport(err)

		reportAssert := &internalReport.Assert{
			Err: "empty repository",
			Fields: map[string]interface{}{
				"dir": repoPath,
			},
		}
		reportAssert.Equal(&s.Suite, report)
	})
}

func (s *ListSuite) TestRecipeError() {

	s.Run("Wrong Recipe Manifest", func() {
		repoPath := internalTesting.DataPath(s, "repository")

		err := s.executor.execute([]string{
			"--repository", repoPath,
		})

		s.Error(err)
		s.Empty(s.executor.stdout)
		s.Empty(s.executor.stderr)

		report := internalReport.NewErrorReport(err)

		reportAssert := &internalReport.Assert{
			Err: "recipe manifest is a directory",
			Fields: map[string]interface{}{
				"path": filepath.Join(repoPath, "recipe", ".manala.yaml"),
			},
		}
		reportAssert.Equal(&s.Suite, report)
	})

	s.Run("Invalid Recipe Manifest", func() {
		repoPath := internalTesting.DataPath(s, "repository")

		err := s.executor.execute([]string{
			"--repository", repoPath,
		})

		s.Error(err)
		s.Empty(s.executor.stdout)
		s.Empty(s.executor.stderr)

		report := internalReport.NewErrorReport(err)

		reportAssert := &internalReport.Assert{
			Err: "invalid recipe manifest",
			Fields: map[string]interface{}{
				"path": filepath.Join(repoPath, "recipe", ".manala.yaml"),
			},
			Reports: []internalReport.Assert{
				{
					Message: "missing manala description field",
					Fields: map[string]interface{}{
						"line":     1,
						"column":   9,
						"property": "description",
					},
				},
			},
		}
		reportAssert.Equal(&s.Suite, report)
	})
}

func (s *ListSuite) Test() {

	s.Run("Custom Repository", func() {
		repoPath := internalTesting.DataPath(s, "repository")

		err := s.executor.execute([]string{
			"--repository", repoPath,
		})

		s.NoError(err)
		s.goldie.Assert(s.T(), internalTesting.Path(s, "stdout"), s.executor.stdout.Bytes())
		s.Empty(s.executor.stderr)
	})

	s.Run("Default Repository", func() {
		repoPath := internalTesting.DataPath(s, "repository")

		s.config.Set("default-repository", repoPath)

		err := s.executor.execute([]string{})

		s.NoError(err)
		s.goldie.Assert(s.T(), internalTesting.Path(s, "stdout"), s.executor.stdout.Bytes())
		s.Empty(s.executor.stderr)
	})
}
