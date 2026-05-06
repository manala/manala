package validation_test

import (
	"testing"

	"github.com/manala/manala/internal/validation"

	"github.com/stretchr/testify/suite"
)

type FormatsSuite struct{ suite.Suite }

func TestFormatsSuite(t *testing.T) {
	suite.Run(t, new(FormatsSuite))
}

func (s *FormatsSuite) TestGitRepo() {
	s.NoError(validation.GitRepoFormat.Validate("https://github.com/manala/manala-recipes.git"))
	s.NoError(validation.GitRepoFormat.Validate("git@github.com:manala/manala.git"))

	s.Error(validation.GitRepoFormat.Validate("foo"))
}

func (s *FormatsSuite) TestFilePath() {
	s.NoError(validation.FilePathFormat.Validate("/"))
	s.NoError(validation.FilePathFormat.Validate("/foo"))
	s.NoError(validation.FilePathFormat.Validate("/foo/bar"))
	s.NoError(validation.FilePathFormat.Validate("/foo-bar"))

	s.Require().Error(validation.FilePathFormat.Validate("foo"))
	s.Require().Error(validation.FilePathFormat.Validate("/foo/*"))
}

func (s *FormatsSuite) TestDomain() {
	s.NoError(validation.DomainFormat.Validate("foo.bar"))
	s.NoError(validation.DomainFormat.Validate("foo.bar.baz"))
	s.NoError(validation.DomainFormat.Validate("foo-bar.baz"))
	s.NoError(validation.DomainFormat.Validate("fo.ba"))

	s.Require().Error(validation.DomainFormat.Validate("foo"))
	s.Require().Error(validation.DomainFormat.Validate("foo.b"))
	s.Require().Error(validation.DomainFormat.Validate("foo_.bar"))
}
