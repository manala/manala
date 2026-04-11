package annotation

import (
	"strings"

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

	// Build a JSON value from annotation tokens, padded with empty lines
	// and leading spaces to match each token's line/column in the source.
	// Newlines and spaces are valid JSON whitespace, so padding is transparent.
	// This makes the unmarshaler's error position directly correct
	// relative to the source, with no resolution needed.
	//
	//   # comment         →  ``
	//   # @foo {          →  `       {`
	//   #   "bar": 123    →  `    "bar": 123`
	//   # }               →  `  }`
	var b strings.Builder
	line := 1
	for _, token := range annot.Value.Tokens {
		for line < token.Line {
			b.WriteString("\n")
			line++
		}
		b.WriteString(strings.Repeat(" ", token.Column-1))
		b.WriteString(token.Value)
	}

	return unmarshaler.Unmarshal([]byte(b.String()), p)
}

// Func calls fn with the named annotation, if present.
func (s *Set) Func(name string, fn func(*Annotation) error) error {
	annot, ok := s.Lookup(name)
	if !ok {
		return nil
	}

	return fn(annot)
}
