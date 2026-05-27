package annotation

import (
	"fmt"
	"slices"
)

type SetValue interface {
	Set(annotation *Annotation) error
}

type setBinding struct {
	name  string
	value SetValue
}

type Set struct {
	bindings []*setBinding
}

func NewSet() *Set {
	return &Set{}
}

func (s *Set) Var(value SetValue, name string) {
	b := &setBinding{name: name, value: value}
	if i := slices.IndexFunc(s.bindings, func(b *setBinding) bool {
		return b.name == name
	}); i >= 0 {
		s.bindings[i] = b
		return
	}
	s.bindings = append(s.bindings, b)
}

func (s *Set) Func(name string, fn func(*Annotation) error) {
	s.Var(setFuncValue(fn), name)
}

func (s *Set) BodyFunc(name string, fn func(*Body) error) {
	s.Var(setBodyFuncValue(fn), name)
}

func (s *Set) Parse(src string) error {
	annotations, err := Parse(src)
	if err != nil {
		return err
	}

	// Reject undeclared annotations
	for _, annot := range annotations {
		if s.lookup(annot.Name.String()) == nil {
			return NewError(
				fmt.Errorf("annotation @%s not defined", annot.Name),
				annot.Start(),
			)
		}
	}

	// Dispatch in declaration order
	for _, b := range s.bindings {
		i := slices.IndexFunc(annotations, func(a *Annotation) bool {
			return a.Name.String() == b.name
		})
		if i == -1 {
			continue
		}

		annot := annotations[i]
		if err := b.value.Set(annot); err != nil {
			return err
		}
	}

	return nil
}

func (s *Set) lookup(name string) *setBinding {
	for _, b := range s.bindings {
		if b.name == name {
			return b
		}
	}
	return nil
}

type setFuncValue func(*Annotation) error

func (f setFuncValue) Set(annotation *Annotation) error {
	return f(annotation)
}

type setBodyFuncValue func(*Body) error

func (f setBodyFuncValue) Set(annotation *Annotation) error {
	if annotation.Body == nil {
		return NewError(
			fmt.Errorf("annotation @%s requires a value", annotation.Name),
			annotation.Start(),
		)
	}
	return f(annotation.Body)
}
