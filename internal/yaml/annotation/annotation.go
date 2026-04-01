package annotation

import (
	"strings"
)

type Annotation struct {
	nameToken   Token
	valueTokens []Token
}

// Name returns the annotation name.
func (a *Annotation) Name() string {
	return a.nameToken.Value
}

// Value returns the annotation value, joining text lines with newlines.
func (a *Annotation) Value() string {
	if len(a.valueTokens) == 0 {
		return ""
	}

	var b strings.Builder
	for i, t := range a.valueTokens {
		if i > 0 {
			b.WriteString("\n")
		}
		b.WriteString(t.Value)
	}

	return b.String()
}

type Annotations []*Annotation

// Lookup returns the annotation with the given name, if any.
func (annotations Annotations) Lookup(name string) (*Annotation, bool) {
	for _, a := range annotations {
		if a.Name() == name {
			return a, true
		}
	}

	return nil, false
}
