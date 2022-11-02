package cmd

import (
	"bytes"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/suite"
	internalConfig "manala/internal/config"
	internalLog "manala/internal/log"
	internalReport "manala/internal/report"
	internalTesting "manala/internal/testing"
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

func (s *WatchSuite) TestProjectError() {

	s.Run("Project Not Found", func() {
		projDir := internalTesting.DataPath(s, "project")

		err := s.executor.execute([]string{
			projDir,
		})

		s.Error(err)
		s.Empty(s.executor.stdout)
		s.Empty(s.executor.stderr)

		report := internalReport.NewErrorReport(err)

		reportAssert := &internalReport.Assert{
			Err: "project manifest not found",
			Fields: map[string]interface{}{
				"dir": projDir,
			},
		}
		reportAssert.Equal(&s.Suite, report)
	})

	s.Run("Wrong Project Manifest", func() {
		projDir := internalTesting.DataPath(s, "project")

		err := s.executor.execute([]string{
			projDir,
		})

		s.Error(err)
		s.Empty(s.executor.stdout)
		s.Empty(s.executor.stderr)

		report := internalReport.NewErrorReport(err)

		reportAssert := &internalReport.Assert{
			Err: "project manifest is a directory",
			Fields: map[string]interface{}{
				"dir": filepath.Join(projDir, ".manala.yaml"),
			},
		}
		reportAssert.Equal(&s.Suite, report)
	})

	s.Run("Empty Project Manifest", func() {
		projDir := internalTesting.DataPath(s, "project")

		err := s.executor.execute([]string{
			projDir,
		})

		s.Error(err)
		s.Empty(s.executor.stdout)
		s.Empty(s.executor.stderr)

		report := internalReport.NewErrorReport(err)

		reportAssert := &internalReport.Assert{
			Message: "irregular project manifest",
			Err:     "empty yaml file",
			Fields: map[string]interface{}{
				"file": filepath.Join(projDir, ".manala.yaml"),
			},
		}
		reportAssert.Equal(&s.Suite, report)
	})

	s.Run("Invalid Project Manifest", func() {
		projDir := internalTesting.DataPath(s, "project")

		err := s.executor.execute([]string{
			projDir,
		})

		s.Error(err)
		s.Empty(s.executor.stdout)
		s.Empty(s.executor.stderr)

		report := internalReport.NewErrorReport(err)

		reportAssert := &internalReport.Assert{
			Err: "invalid project manifest",
			Fields: map[string]interface{}{
				"file": filepath.Join(projDir, ".manala.yaml"),
			},
			Reports: []internalReport.Assert{
				{
					Message: "missing manala recipe field",
					Fields: map[string]interface{}{
						"line":     1,
						"column":   9,
						"property": "recipe",
					},
				},
			},
		}
		reportAssert.Equal(&s.Suite, report)
	})
}

func (s *WatchSuite) TestRepositoryError() {

	s.Run("No Repository", func() {
		projDir := internalTesting.DataPath(s, "project")

		err := s.executor.execute([]string{
			projDir,
		})

		s.Error(err)
		s.Empty(s.executor.stdout)
		s.Empty(s.executor.stderr)

		report := internalReport.NewErrorReport(err)

		reportAssert := &internalReport.Assert{
			Err: "unable to process empty repository url",
		}
		reportAssert.Equal(&s.Suite, report)
	})

	s.Run("Repository Not Found", func() {
		projDir := internalTesting.DataPath(s, "project")
		repoUrl := internalTesting.DataPath(s, "repository")

		err := s.executor.execute([]string{
			projDir,
			"--repository", repoUrl,
		})

		s.Error(err)
		s.Empty(s.executor.stdout)
		s.Empty(s.executor.stderr)

		report := internalReport.NewErrorReport(err)

		reportAssert := &internalReport.Assert{
			Err: "repository not found",
			Fields: map[string]interface{}{
				"url": repoUrl,
			},
		}
		reportAssert.Equal(&s.Suite, report)
	})

	s.Run("Wrong Repository", func() {
		projDir := internalTesting.DataPath(s, "project")
		repoUrl := internalTesting.DataPath(s, "repository")

		err := s.executor.execute([]string{
			projDir,
			"--repository", repoUrl,
		})

		s.Error(err)
		s.Empty(s.executor.stdout)
		s.Empty(s.executor.stderr)

		report := internalReport.NewErrorReport(err)

		reportAssert := &internalReport.Assert{
			Err: "wrong repository",
			Fields: map[string]interface{}{
				"dir": repoUrl,
			},
		}
		reportAssert.Equal(&s.Suite, report)
	})
}

func (s *WatchSuite) TestRecipeError() {

	s.Run("Recipe Not Found", func() {
		projDir := internalTesting.DataPath(s, "project")
		repoUrl := internalTesting.DataPath(s, "repository")

		s.config.Set("default-repository", repoUrl)

		err := s.executor.execute([]string{
			projDir,
			"--recipe", "recipe",
		})

		s.Error(err)
		s.Empty(s.executor.stdout)
		s.Empty(s.executor.stderr)

		report := internalReport.NewErrorReport(err)

		reportAssert := &internalReport.Assert{
			Err: "recipe manifest not found",
			Fields: map[string]interface{}{
				"file": filepath.Join(repoUrl, "recipe", ".manala.yaml"),
			},
		}
		reportAssert.Equal(&s.Suite, report)
	})

	s.Run("Wrong Recipe Manifest", func() {
		projDir := internalTesting.DataPath(s, "project")
		repoUrl := internalTesting.DataPath(s, "repository")

		err := s.executor.execute([]string{
			projDir,
			"--repository", repoUrl,
			"--recipe", "recipe",
		})

		s.Error(err)
		s.Empty(s.executor.stdout)
		s.Empty(s.executor.stderr)

		report := internalReport.NewErrorReport(err)

		reportAssert := &internalReport.Assert{
			Err: "recipe manifest is a directory",
			Fields: map[string]interface{}{
				"dir": filepath.Join(repoUrl, "recipe", ".manala.yaml"),
			},
		}
		reportAssert.Equal(&s.Suite, report)
	})

	s.Run("Invalid Recipe Manifest", func() {
		projDir := internalTesting.DataPath(s, "project")
		repoUrl := internalTesting.DataPath(s, "repository")

		err := s.executor.execute([]string{
			projDir,
			"--repository", repoUrl,
			"--recipe", "recipe",
		})

		s.Error(err)
		s.Empty(s.executor.stdout)
		s.Empty(s.executor.stderr)

		report := internalReport.NewErrorReport(err)

		reportAssert := &internalReport.Assert{
			Err: "invalid recipe manifest",
			Fields: map[string]interface{}{
				"file": filepath.Join(repoUrl, "recipe", ".manala.yaml"),
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
