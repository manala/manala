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
	"os"
	"path/filepath"
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
		projPath := internalTesting.DataPath(s, "project")

		err := s.executor.execute([]string{
			projPath,
		})

		s.Error(err)
		s.Empty(s.executor.stdout)
		s.Empty(s.executor.stderr)

		report := internalReport.NewErrorReport(err)

		reportAssert := &internalReport.Assert{
			Message: "irregular project manifest",
			Err:     "empty yaml file",
			Fields: map[string]interface{}{
				"path": filepath.Join(projPath, ".manala.yaml"),
			},
		}
		reportAssert.Equal(&s.Suite, report)
	})

	s.Run("Already Existing Project", func() {
		projPath := internalTesting.DataPath(s, "project")

		err := s.executor.execute([]string{
			projPath,
		})

		s.Error(err)
		s.Empty(s.executor.stdout)
		s.Empty(s.executor.stderr)

		report := internalReport.NewErrorReport(err)

		reportAssert := &internalReport.Assert{
			Err: "already existing project",
			Fields: map[string]interface{}{
				"path": projPath,
			},
		}
		reportAssert.Equal(&s.Suite, report)
	})
}

func (s *InitSuite) TestRepositoryError() {

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

func (s *InitSuite) TestRecipeError() {

	s.Run("Recipe Not Found", func() {
		repoPath := internalTesting.DataPath(s, "repository")

		err := s.executor.execute([]string{
			"--repository", repoPath,
			"--recipe", "not_found",
		})

		s.Error(err)
		s.Empty(s.executor.stdout)
		s.Empty(s.executor.stderr)

		report := internalReport.NewErrorReport(err)

		reportAssert := &internalReport.Assert{
			Err: "recipe not found",
		}
		reportAssert.Equal(&s.Suite, report)
	})

	s.Run("Wrong Recipe Manifest", func() {
		repoPath := internalTesting.DataPath(s, "repository")

		err := s.executor.execute([]string{
			"--repository", repoPath,
			"--recipe", "recipe",
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
			"--recipe", "recipe",
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

func (s *InitSuite) Test() {

	s.Run("Custom Repository", func() {
		projPath := internalTesting.DataPath(s, "project")
		repoPath := internalTesting.DataPath(s, "repository")

		_ = os.RemoveAll(projPath)

		err := s.executor.execute([]string{
			projPath,
			"--repository", repoPath,
			"--recipe", "recipe",
		})

		s.NoError(err)
		s.Empty(s.executor.stdout)
		s.goldie.AssertWithTemplate(s.T(), internalTesting.Path(s, "stderr"), map[string]interface{}{
			"dst": projPath,
			"src": filepath.Join(repoPath, "recipe"),
		}, s.executor.stderr.Bytes())
		s.DirExists(projPath)
		s.FileExists(filepath.Join(projPath, ".manala.yaml"))
		manifestContent, _ := os.ReadFile(filepath.Join(projPath, ".manala.yaml"))
		s.goldie.AssertWithTemplate(s.T(), internalTesting.Path(s, "manifest"), map[string]interface{}{
			"repository": repoPath,
		}, manifestContent)
	})

	s.Run("Default Repository", func() {
		projPath := internalTesting.DataPath(s, "project")
		repoPath := internalTesting.DataPath(s, "repository")

		_ = os.RemoveAll(projPath)

		s.config.Set("default-repository", repoPath)

		err := s.executor.execute([]string{
			projPath,
			"--recipe", "recipe",
		})

		s.NoError(err)
		s.Empty(s.executor.stdout)
		s.goldie.AssertWithTemplate(s.T(), internalTesting.Path(s, "stderr"), map[string]interface{}{
			"dst": projPath,
			"src": filepath.Join(repoPath, "recipe"),
		}, s.executor.stderr.Bytes())
		s.DirExists(projPath)
		s.FileExists(filepath.Join(projPath, ".manala.yaml"))
		manifestContent, _ := os.ReadFile(filepath.Join(projPath, ".manala.yaml"))
		s.goldie.AssertWithTemplate(s.T(), internalTesting.Path(s, "manifest"), map[string]interface{}{
			"repository": repoPath,
		}, manifestContent)
		s.FileExists(filepath.Join(projPath, "file.txt"))
		fileContent, _ := os.ReadFile(filepath.Join(projPath, "file.txt"))
		s.goldie.Assert(s.T(), internalTesting.Path(s, "file.txt"), fileContent)
		s.FileExists(filepath.Join(projPath, "template"))
		templateContent, _ := os.ReadFile(filepath.Join(projPath, "template"))
		s.goldie.Assert(s.T(), internalTesting.Path(s, "template"), templateContent)
	})
}
