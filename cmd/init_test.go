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
	s.Run("Already Existing Project", func() {
		projDir := internalTesting.DataPath(s, "project")

		err := s.executor.execute([]string{
			projDir,
		})

		s.Error(err)
		s.Empty(s.executor.stdout)
		s.Empty(s.executor.stderr)

		report := internalReport.NewErrorReport(err)

		reportAssert := &internalReport.Assert{
			Err: "already existing project",
			Fields: map[string]interface{}{
				"dir": projDir,
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
			Err: "unable to process empty repository url",
		}
		reportAssert.Equal(&s.Suite, report)
	})

	s.Run("Repository Not Found", func() {
		repoUrl := internalTesting.DataPath(s, "repository")

		err := s.executor.execute([]string{
			"--repository", repoUrl,
		})

		s.Error(err)
		s.Empty(s.executor.stdout)
		s.Empty(s.executor.stderr)

		report := internalReport.NewErrorReport(err)

		reportAssert := &internalReport.Assert{
			Err: "unsupported repository url",
			Fields: map[string]interface{}{
				"url": repoUrl,
			},
		}
		reportAssert.Equal(&s.Suite, report)
	})

	s.Run("Wrong Repository", func() {
		repoUrl := internalTesting.DataPath(s, "repository")

		err := s.executor.execute([]string{
			"--repository", repoUrl,
		})

		s.Error(err)
		s.Empty(s.executor.stdout)
		s.Empty(s.executor.stderr)

		report := internalReport.NewErrorReport(err)

		reportAssert := &internalReport.Assert{
			Err: "unsupported repository url",
			Fields: map[string]interface{}{
				"url": repoUrl,
			},
		}
		reportAssert.Equal(&s.Suite, report)
	})

	s.Run("Empty Repository", func() {
		repoUrl := internalTesting.DataPath(s, "repository")

		err := s.executor.execute([]string{
			"--repository", repoUrl,
		})

		s.Error(err)
		s.Empty(s.executor.stdout)
		s.Empty(s.executor.stderr)

		report := internalReport.NewErrorReport(err)

		reportAssert := &internalReport.Assert{
			Err: "empty repository",
			Fields: map[string]interface{}{
				"dir": repoUrl,
			},
		}
		reportAssert.Equal(&s.Suite, report)
	})
}

func (s *InitSuite) TestRecipeError() {

	s.Run("Recipe Not Found", func() {
		repoUrl := internalTesting.DataPath(s, "repository")

		err := s.executor.execute([]string{
			"--repository", repoUrl,
			"--recipe", "not_found",
		})

		s.Error(err)
		s.Empty(s.executor.stdout)
		s.Empty(s.executor.stderr)

		report := internalReport.NewErrorReport(err)

		reportAssert := &internalReport.Assert{
			Err: "recipe manifest not found",
			Fields: map[string]interface{}{
				"file": filepath.Join(repoUrl, "not_found", ".manala.yaml"),
			},
		}
		reportAssert.Equal(&s.Suite, report)
	})

	s.Run("Wrong Recipe Manifest", func() {
		repoUrl := internalTesting.DataPath(s, "repository")

		err := s.executor.execute([]string{
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
		repoUrl := internalTesting.DataPath(s, "repository")

		err := s.executor.execute([]string{
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

func (s *InitSuite) Test() {

	s.Run("Custom Repository", func() {
		projDir := internalTesting.DataPath(s, "project")
		repoUrl := internalTesting.DataPath(s, "repository")

		_ = os.RemoveAll(projDir)

		err := s.executor.execute([]string{
			projDir,
			"--repository", repoUrl,
			"--recipe", "recipe",
		})

		s.NoError(err)
		s.Empty(s.executor.stdout)
		s.goldie.AssertWithTemplate(s.T(), internalTesting.Path(s, "stderr"), map[string]interface{}{
			"dst": projDir,
			"src": filepath.Join(repoUrl, "recipe"),
		}, s.executor.stderr.Bytes())
		s.DirExists(projDir)
		s.FileExists(filepath.Join(projDir, ".manala.yaml"))
		manContent, _ := os.ReadFile(filepath.Join(projDir, ".manala.yaml"))
		s.goldie.AssertWithTemplate(s.T(), internalTesting.Path(s, "manifest"), map[string]interface{}{
			"repository": repoUrl,
		}, manContent)
	})

	s.Run("Default Repository", func() {
		projDir := internalTesting.DataPath(s, "project")
		repoUrl := internalTesting.DataPath(s, "repository")

		_ = os.RemoveAll(projDir)

		s.config.Set("default-repository", repoUrl)

		err := s.executor.execute([]string{
			projDir,
			"--recipe", "recipe",
		})

		s.NoError(err)
		s.Empty(s.executor.stdout)
		s.goldie.AssertWithTemplate(s.T(), internalTesting.Path(s, "stderr"), map[string]interface{}{
			"dst": projDir,
			"src": filepath.Join(repoUrl, "recipe"),
		}, s.executor.stderr.Bytes())
		s.DirExists(projDir)
		s.FileExists(filepath.Join(projDir, ".manala.yaml"))
		manContent, _ := os.ReadFile(filepath.Join(projDir, ".manala.yaml"))
		s.goldie.AssertWithTemplate(s.T(), internalTesting.Path(s, "manifest"), map[string]interface{}{
			"repository": repoUrl,
		}, manContent)
		s.FileExists(filepath.Join(projDir, "file.txt"))
		fileContent, _ := os.ReadFile(filepath.Join(projDir, "file.txt"))
		s.goldie.Assert(s.T(), internalTesting.Path(s, "file.txt"), fileContent)
		s.FileExists(filepath.Join(projDir, "template"))
		templateContent, _ := os.ReadFile(filepath.Join(projDir, "template"))
		s.goldie.Assert(s.T(), internalTesting.Path(s, "template"), templateContent)
	})
}
