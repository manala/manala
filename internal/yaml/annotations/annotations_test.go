package annotations_test

import (
	"testing"

	"github.com/manala/manala/internal/yaml/annotations"

	"github.com/stretchr/testify/suite"
)

type AnnotationsSuite struct{ suite.Suite }

func TestAnnotationsSuite(t *testing.T) {
	suite.Run(t, new(AnnotationsSuite))
}

func (s *AnnotationsSuite) TestLookup() {
	src := `
		# @foo bar
		# @bar baz
	`
	annots, err := annotations.Parse(src)
	s.Require().NoError(err)
	s.Require().Len(annots, 2)

	var annot *annotations.Annotation
	var ok bool

	annot, ok = annots.Lookup("foo")
	s.True(ok)
	s.Equal("foo", annot.Name())
	s.Equal("bar", annot.Value())

	annot, ok = annots.Lookup("bar")
	s.True(ok)
	s.Equal("bar", annot.Name())
	s.Equal("baz", annot.Value())

	annot, ok = annots.Lookup("baz")
	s.False(ok)
	s.Nil(annot)
}
