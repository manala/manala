package annotation_test

import (
	"errors"
	"testing"

	"github.com/manala/manala/internal/testing/expectation"
	"github.com/manala/manala/internal/testing/heredoc"
	yamlannotation "github.com/manala/manala/internal/yaml/annotation"

	"github.com/stretchr/testify/suite"
)

type SetSuite struct{ suite.Suite }

func TestSetSuite(t *testing.T) {
	suite.Run(t, new(SetSuite))
}

func (s *SetSuite) TestFunc() {
	src := heredoc.Doc(`
		# @foo bar
	`)

	var annot *yamlannotation.Annotation

	set := yamlannotation.NewSet()
	set.Func("foo", func(a *yamlannotation.Annotation) error {
		annot = a
		return nil
	})

	err := set.Parse(src)
	s.Require().NoError(err)

	s.Require().NotNil(annot)
	s.Equal("foo", annot.Name.String())
	s.Require().NotNil(annot.Body)
	s.Equal("bar", annot.Body.String())
}

func (s *SetSuite) TestBodyFunc() {
	src := heredoc.Doc(`
		# @foo bar
	`)

	var body *yamlannotation.Body

	set := yamlannotation.NewSet()
	set.BodyFunc("foo", func(b *yamlannotation.Body) error {
		body = b
		return nil
	})

	err := set.Parse(src)
	s.Require().NoError(err)

	s.Require().NotNil(body)
	s.Equal("bar", body.String())
}

func (s *SetSuite) TestBodyFuncNoBody() {
	src := heredoc.Doc(`
		# @foo
	`)

	called := false

	set := yamlannotation.NewSet()
	set.BodyFunc("foo", func(_ *yamlannotation.Body) error {
		called = true
		return nil
	})

	err := set.Parse(src)

	s.False(called)
	expectation.ExpectError(s.T(), yamlannotation.ErrorExpectation{
		Position: [2]int{1, 3},
		Err:      expectation.ErrorMessage("annotation @foo requires a value"),
	}, err)
}

func (s *SetSuite) TestAbsent() {
	src := heredoc.Doc(`
		# @foo bar
	`)

	called := false

	set := yamlannotation.NewSet()
	set.Func("foo", func(_ *yamlannotation.Annotation) error {
		return nil
	})
	set.Func("bar", func(_ *yamlannotation.Annotation) error {
		called = true
		return nil
	})

	err := set.Parse(src)
	s.Require().NoError(err)

	s.False(called)
}

func (s *SetSuite) TestUndeclared() {
	src := heredoc.Doc(`
		# @foo bar
	`)

	set := yamlannotation.NewSet()

	err := set.Parse(src)

	expectation.ExpectError(s.T(), yamlannotation.ErrorExpectation{
		Position: [2]int{1, 3},
		Err:      expectation.ErrorMessage("annotation @foo not defined"),
	}, err)
}

func (s *SetSuite) TestDeclarationOrder() {
	// Annotations are dispatched in declaration order, regardless of
	// their order in the source.
	src := heredoc.Doc(`
		# @foo bar
		# @bar baz
	`)

	var order []string

	set := yamlannotation.NewSet()
	set.Func("bar", func(_ *yamlannotation.Annotation) error {
		order = append(order, "bar")
		return nil
	})
	set.Func("foo", func(_ *yamlannotation.Annotation) error {
		order = append(order, "foo")
		return nil
	})

	err := set.Parse(src)
	s.Require().NoError(err)

	s.Equal([]string{"bar", "foo"}, order)
}

func (s *SetSuite) TestVarReplace() {
	src := heredoc.Doc(`
		# @foo bar
	`)

	var got string

	set := yamlannotation.NewSet()
	set.Func("foo", func(_ *yamlannotation.Annotation) error {
		got = "first"
		return nil
	})
	// Re-registering the same name replaces the previous binding.
	set.Func("foo", func(_ *yamlannotation.Annotation) error {
		got = "second"
		return nil
	})

	err := set.Parse(src)
	s.Require().NoError(err)

	s.Equal("second", got)
}

func (s *SetSuite) TestFuncError() {
	src := heredoc.Doc(`
		# @foo bar
	`)

	set := yamlannotation.NewSet()
	set.Func("foo", func(_ *yamlannotation.Annotation) error {
		return errors.New("boom")
	})

	err := set.Parse(src)

	s.Require().Error(err)
	s.Equal("boom", err.Error())
}

func (s *SetSuite) TestParseError() {
	src := heredoc.Doc(`
		# @foo bar
		# @foo baz
	`)

	set := yamlannotation.NewSet()
	set.Func("foo", func(_ *yamlannotation.Annotation) error {
		return nil
	})

	err := set.Parse(src)

	expectation.ExpectError(s.T(), yamlannotation.ErrorExpectation{
		Position: [2]int{2, 3},
		Err:      expectation.ErrorMessage("duplicate @foo annotation"),
	}, err)
}
