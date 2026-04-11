package annotation

import (
	"github.com/manala/manala/internal/json/unmarshaler"
)

// Set is a set of parsed annotations, returned by Parse.
type Set struct {
	annotations []*Annotation
}

// Lookup returns the annotation with the given name, if any.
func (s *Set) Lookup(name string) (*Annotation, bool) {
	for _, a := range s.annotations {
		if a.Name.String() == name {
			return a, true
		}
	}

	return nil, false
}

// Len returns the number of annotations in the set.
func (s *Set) Len() int {
	return len(s.annotations)
}

// JSONVar unmarshals the value of the named annotation as JSON into p.
// If the annotation is not present, p is left unchanged and no error is returned.
func (s *Set) JSONVar(p any, name string) error {
	annot, ok := s.Lookup(name)
	if !ok {
		return nil
	}

	// Stencil preserves source positions for accurate error reporting
	value := annot.Value.Stencil()

	return unmarshaler.Unmarshal([]byte(value), p)
}

// Func calls fn with the named annotation, if present.
func (s *Set) Func(name string, fn func(*Annotation) error) error {
	annot, ok := s.Lookup(name)
	if !ok {
		return nil
	}

	return fn(annot)
}
