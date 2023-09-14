package schema

func (s *Suite) TestFormatGitRepo() {
	checker := &GitRepoFormatChecker{}

	s.True(checker.IsFormat("https://github.com/manala/manala-recipes.git"))
	s.True(checker.IsFormat("git@github.com:manala/manala.git"))

	s.False(checker.IsFormat("foo"))
}

func (s *Suite) TestFormatFilePath() {
	checker := &FilePathFormatChecker{}

	s.True(checker.IsFormat("/"))
	s.True(checker.IsFormat("/foo"))
	s.True(checker.IsFormat("/foo/bar"))
	s.True(checker.IsFormat("/foo-bar"))

	s.False(checker.IsFormat("foo"))
	s.False(checker.IsFormat("/foo/*"))
}

func (s *Suite) TestFormatDomain() {
	checker := &DomainFormatChecker{}

	s.True(checker.IsFormat("foo.bar"))
	s.True(checker.IsFormat("foo.bar.baz"))
	s.True(checker.IsFormat("foo-bar.baz"))
	s.True(checker.IsFormat("fo.ba"))

	s.False(checker.IsFormat("foo"))
	s.False(checker.IsFormat("foo.b"))
	s.False(checker.IsFormat("foo_.bar"))
}
