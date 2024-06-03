package schema_test

import (
	"manala/internal/schema"
	"testing"

	"github.com/stretchr/testify/suite"
)

type FormatsSuite struct{ suite.Suite }

func TestFormatsSuite(t *testing.T) {
	suite.Run(t, new(FormatsSuite))
}

func (s *FormatsSuite) TestGitRepo() {
	checker := &schema.GitRepoFormatChecker{}

	s.True(checker.IsFormat("https://github.com/manala/manala-recipes.git"))
	s.True(checker.IsFormat("git@github.com:manala/manala.git"))

	s.False(checker.IsFormat("foo"))
}

func (s *FormatsSuite) TestFilePath() {
	checker := &schema.FilePathFormatChecker{}

	s.True(checker.IsFormat("/"))
	s.True(checker.IsFormat("/foo"))
	s.True(checker.IsFormat("/foo/bar"))
	s.True(checker.IsFormat("/foo-bar"))

	s.False(checker.IsFormat("foo"))
	s.False(checker.IsFormat("/foo/*"))
}

func (s *FormatsSuite) TestDomain() {
	checker := &schema.DomainFormatChecker{}

	s.True(checker.IsFormat("foo.bar"))
	s.True(checker.IsFormat("foo.bar.baz"))
	s.True(checker.IsFormat("foo-bar.baz"))
	s.True(checker.IsFormat("fo.ba"))

	s.False(checker.IsFormat("foo"))
	s.False(checker.IsFormat("foo.b"))
	s.False(checker.IsFormat("foo_.bar"))
}
